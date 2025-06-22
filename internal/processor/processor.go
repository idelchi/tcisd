package processor

import (
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path"
	"slices"
	"strings"
	"sync"

	"github.com/fatih/color"
	"github.com/natefinch/atomic"

	"github.com/idelchi/tcisd/internal/config"
	"github.com/idelchi/tcisd/pkg/matcher"
	"github.com/idelchi/tcisd/pkg/remover"
)

type Processor struct {
	cfg     *config.Config
	files   []string
	results map[string][]string
}

func New(cfg *config.Config) *Processor {
	return &Processor{
		cfg:     cfg,
		results: make(map[string][]string),
	}
}

func (p *Processor) Process() error {
	logger := log.New(os.Stderr, "", 0)
	logger.SetOutput(io.Discard)

	m := matcher.New(p.cfg.Hidden, p.cfg.Exclude)

	for _, path := range p.cfg.Paths {
		if err := m.Match(path); err != nil {
			return fmt.Errorf("matching pattern %s: %w", path, err)
		}
	}

	p.files = m.List()

	if len(p.files) == 0 {
		return errors.New("no files found")
	}

	jobs := make(chan string, len(p.files))

	results := make(chan struct {
		file   string
		issues []string
	}, len(p.files))

	var wg sync.WaitGroup

	for i := range p.cfg.Parallel {
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

	return nil
}

func (p *Processor) worker(_ int, jobs <-chan string, results chan<- struct {
	file   string
	issues []string
}, wg *sync.WaitGroup,
) {
	defer wg.Done()

	for file := range jobs {
		fileType := detectFileType(file)

		if !slices.Contains(p.cfg.Types, fileType) {
			continue
		}

		r := remover.ForType(fileType)
		if r == nil {
			log.Printf("No remover found for file type %s", fileType)

			continue
		}

		content, err := os.ReadFile(file)
		if err != nil {
			log.Printf("Error reading file %s: %v", file, err)

			continue
		}

		lines := strings.Split(string(content), "\n")
		processedLines, issues := r.Process(lines)

		if p.cfg.Mode == config.FormatMode && len(issues) > 0 {
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
			log.Println(color.GreenString("Files were modified successfully"))

			hasIssues = false
		}
	}

	return hasIssues
}

func detectFileType(file string) string {
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
