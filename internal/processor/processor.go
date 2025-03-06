package processor

import (
	"errors"
	"fmt"
	"log"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/natefinch/atomic"

	"github.com/idelchi/tcisd/internal/config"
	"github.com/idelchi/tcisd/pkg/remover"
	"github.com/idelchi/wslint/pkg/matcher"
)

type Processor struct {
	cfg            *config.Config
	workers        int
	types          []string
	files          []string
	results        map[string][]string
	processingTime time.Duration
}

func New(cfg *config.Config, workers int, types []string) *Processor {
	return &Processor{
		cfg:     cfg,
		workers: workers,
		types:   types,
		results: make(map[string][]string),
	}
}

func (p *Processor) Process() error {
	m := matcher.New(p.cfg.Hidden, p.cfg.Exclude, log.New(os.Stderr, "", 0))

	for _, path := range p.cfg.Paths {
		for _, pattern := range p.cfg.Patterns {
			fullPattern := fmt.Sprintf("%s/%s", path, pattern)
			if err := m.Match(fullPattern); err != nil {
				return fmt.Errorf("matching pattern %s: %w", fullPattern, err)
			}
		}
	}

	p.files = m.ListFiles()

	if len(p.files) == 0 {
		return errors.New("no files found")
	}

	log.Printf("Found %d files to process", len(p.files))

	start := time.Now()

	jobs := make(chan string, len(p.files))

	results := make(chan struct {
		file   string
		issues []string
	}, len(p.files))

	var wg sync.WaitGroup

	for i := range p.workers {
		wg.Add(1)

		go p.worker(i, jobs, results, &wg)
	}

	for _, file := range p.files {
		jobs <- file
	}

	close(jobs)

	wg.Wait()
	close(results)

	for result := range results {
		if len(result.issues) > 0 {
			p.results[result.file] = result.issues
		}
	}

	p.processingTime = time.Since(start)

	return nil
}

func (p *Processor) worker(_ int, jobs <-chan string, results chan<- struct {
	file   string
	issues []string
}, wg *sync.WaitGroup,
) {
	defer wg.Done()

	for file := range jobs {
		log.Printf("Processing file: %s", file)

		fileType := detectFileType(file)

		if !contains(p.types, fileType) {
			continue
		}

		content, err := os.ReadFile(file)
		if err != nil {
			log.Printf("Error reading file %s: %v", file, err)

			continue
		}

		r := remover.ForType(fileType)
		if r == nil {
			log.Printf("No remover found for file type %s", fileType)

			continue
		}

		lines := strings.Split(string(content), "\n")
		processedLines, issues := r.Process(lines)

		if p.cfg.Mode == config.FormatMode && len(issues) > 0 && !p.cfg.DryRun {
			if err := atomic.WriteFile(file, strings.NewReader(strings.Join(processedLines, "\n"))); err != nil {
				log.Printf("Error writing file %s: %v", file, err)
			}
		}

		results <- struct {
			file   string
			issues []string
		}{file, issues}
	}
}

func (p *Processor) Summary() bool {
	log.Printf("Processed %d files in %s", len(p.files), p.processingTime)

	hasIssues := len(p.results) > 0

	if hasIssues {
		log.Printf("Found comments in %d files:", len(p.results))

		for file, issues := range p.results {
			log.Printf("  %s:", color.YellowString(file))

			for _, issue := range issues {
				log.Printf("    - %s", issue)
			}
		}

		if p.cfg.Mode == config.FormatMode {
			if p.cfg.DryRun {
				log.Println(color.YellowString("Dry run: no files were modified"))
			} else {
				log.Println(color.GreenString("Files were modified successfully"))
			}
		}
	}

	return hasIssues
}

func detectFileType(file string) string {
	// Check for Dockerfile (which doesn't have a standard extension)
	baseName := path.Base(file)
	if baseName == "Dockerfile" || strings.HasPrefix(baseName, "Dockerfile.") {
		return "dockerfile"
	}

	ext := path.Ext(file)

	switch ext {
	case ".go":
		return "go"
	case ".py", "python":
		return "python"
	default:
		return ext
	}
}

func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}

	return false
}
