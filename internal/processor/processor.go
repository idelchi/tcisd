package processor

import (
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

// Processor handles the processing of files.
type Processor struct {
	cfg *config.Config
	// Number of workers
	workers int
	// File types to process
	types []string
	// Files to process
	files []string
	// Results from processing
	results map[string][]string
	// Processing time
	processingTime time.Duration
}

// New creates a new processor.
func New(cfg *config.Config, workers int, types []string) *Processor {
	return &Processor{
		cfg:     cfg,
		workers: workers,
		types:   types,
		results: make(map[string][]string),
	}
}

// Process handles the processing of files.
func (p *Processor) Process() error {
	// Find matching files
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
		return fmt.Errorf("no files found")
	}

	log.Printf("Found %d files to process", len(p.files))

	// Process files with workers
	start := time.Now()

	// Create a buffered channel to hold jobs
	jobs := make(chan string, len(p.files))

	// Create a channel to receive results
	results := make(chan struct {
		file   string
		issues []string
	}, len(p.files))

	// Create a wait group to wait for all workers to finish
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < p.workers; i++ {
		wg.Add(1)
		go p.worker(i, jobs, results, &wg)
	}

	// Send jobs to workers
	for _, file := range p.files {
		jobs <- file
	}
	close(jobs)

	// Wait for all workers to finish
	wg.Wait()
	close(results)

	// Collect results
	for result := range results {
		if len(result.issues) > 0 {
			p.results[result.file] = result.issues
		}
	}

	p.processingTime = time.Since(start)

	return nil
}

// worker processes files.
func (p *Processor) worker(_ int, jobs <-chan string, results chan<- struct {
	file   string
	issues []string
}, wg *sync.WaitGroup,
) {
	defer wg.Done()

	for file := range jobs {
		log.Printf("Processing file: %s", file)

		// Determine file type based on extension
		fileType := detectFileType(file)

		// Skip file if its type is not in the list of types to process
		if !contains(p.types, fileType) {
			continue
		}

		// Read file
		content, err := os.ReadFile(file)
		if err != nil {
			log.Printf("Error reading file %s: %v", file, err)
			continue
		}

		// Get the appropriate remover
		r := remover.ForType(fileType)
		if r == nil {
			log.Printf("No remover found for file type %s", fileType)
			continue
		}

		// Process content
		lines := strings.Split(string(content), "\n")
		processedLines, issues := r.Process(lines)

		// If in format mode and there are issues, write the processed content back to the file
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

// Summary prints a summary of the processing results.
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

// detectFileType determines the file type based on extension.
func detectFileType(file string) string {
	ext := path.Ext(file)

	switch ext {
	case ".go":
		return "go"
	case ".sh", "bash":
		return "bash"
	case ".py", "python":
		return "python"
	default:
		return ext
	}
}

// contains checks if a string is in a slice.
func contains(slice []string, str string) bool {
	for _, s := range slice {
		if s == str {
			return true
		}
	}
	return false
}
