package reforge

import (
	"context"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

const outputDir = "zipped"

func Run(ctx context.Context, cfg Config, progress chan<- ProgressEvent) error {
	scan, err := ScanArchives(cfg.WorkDir)
	if err != nil {
		return err
	}

	outputPath := cfg.WorkDir
	if !cfg.Replace {
		outputPath = filepath.Join(cfg.WorkDir, outputDir)
		if err := os.MkdirAll(outputPath, 0755); err != nil {
			return err
		}
	}

	total := len(scan.CBRFiles) + len(scan.CBZFiles)

	jobs := make(chan Job, total)
	results := make(chan Result, total)

	var wg sync.WaitGroup
	workers := runtime.NumCPU()

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go runWorker(ctx, jobs, results, progress, &wg)
	}

	for _, cbr := range scan.CBRFiles {
		name := strings.TrimSuffix(filepath.Base(cbr), ".cbr") + ".cbz"
		jobs <- Job{
			Source: cbr,
			Dest:   filepath.Join(outputPath, name),
			Type:   JobConvert,
		}
	}

	for _, cbz := range scan.CBZFiles {
		jobs <- Job{
			Source: cbz,
			Dest:   filepath.Join(outputPath, filepath.Base(cbz)),
			Type:   JobCopy,
		}
	}

	close(jobs)

	go func() {
		wg.Wait()
		close(results)
	}()

	for res := range results {
		if !res.Success {
			if progress != nil {
				progress <- ProgressEvent{
					Status: "error",
					File:   filepath.Base(res.Source),
					Err:    res.Error,
				}
			}
			continue
		}

		switch res.Type {
		case JobConvert:
			if cfg.Replace {
				_ = os.Remove(res.Source)
			}
			if progress != nil {
				progress <- ProgressEvent{
					Status: "done",
					File:   filepath.Base(res.Source),
					Dest:   filepath.Base(res.Dest),
				}
			}

		case JobCopy:
			if progress != nil {
				progress <- ProgressEvent{
					Status: "kept",
					File:   filepath.Base(res.Source),
				}
			}
		}
	}

	if progress != nil {
		close(progress)
	}

	return nil
}

func runWorker(ctx context.Context, jobs <-chan Job, results chan<- Result, progress chan<- ProgressEvent, wg *sync.WaitGroup) {
	defer wg.Done()

	for job := range jobs {
		select {
		case <-ctx.Done():
			results <- Result{
				Type:    job.Type,
				Source:  job.Source,
				Dest:    job.Dest,
				Success: false,
				Error:   ctx.Err(),
			}
			continue
		default:
		}

		if progress != nil {
			progress <- ProgressEvent{
				Status: "converting",
				File:   filepath.Base(job.Source),
			}
		}

		var err error

		switch job.Type {
		case JobConvert:
			err = Convert(ctx, job.Source, job.Dest)
		case JobCopy:
			err = CopyFile(job.Source, job.Dest)
		}

		results <- Result{
			Type:    job.Type,
			Source:  job.Source,
			Dest:    job.Dest,
			Success: err == nil,
			Error:   err,
		}
	}
}
