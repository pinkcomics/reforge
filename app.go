package main

import (
	"context"
	"path/filepath"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"converter/pkg/converter"
)

type App struct {
	ctx context.Context
}

func NewApp() *App { return &App{} }

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

type ScanResult struct {
	CBRCount int      `json:"cbrCount"`
	CBZCount int      `json:"cbzCount"`
	CBRFiles []string `json:"cbrFiles"`
	CBZFiles []string `json:"cbzFiles"`
	Error    string   `json:"error,omitempty"`
}

func (a *App) SelectFolder() string {
	dir, err := runtime.OpenDirectoryDialog(a.ctx, runtime.OpenDialogOptions{
		Title: "Selecionar pasta com arquivos CBR/CBZ",
	})
	if err != nil {
		return ""
	}
	return dir
}

func (a *App) ScanFolder(dir string) ScanResult {
	scan, err := converter.ScanArchives(dir)
	if err != nil {
		return ScanResult{Error: err.Error()}
	}

	return ScanResult{
		CBRCount: len(scan.CBRFiles),
		CBZCount: len(scan.CBZFiles),
		CBRFiles: baseNames(scan.CBRFiles),
		CBZFiles: baseNames(scan.CBZFiles),
	}
}

func (a *App) CheckExtractor() bool {
	return converter.CheckExtractor()
}

func (a *App) StartConversion(workDir string, skip, replace bool) string {
	if !converter.CheckExtractor() {
		return "Nenhum extrator RAR encontrado (instale unrar ou 7z)"
	}

	cfg := converter.Config{
		WorkDir: workDir,
		Skip:    skip,
		Replace: replace,
	}

	progress := make(chan converter.ProgressEvent, 256)

	go func() {
		for ev := range progress {
			switch ev.Status {
			case "done":
				runtime.EventsEmit(a.ctx, "progress:done", map[string]string{
					"file": ev.File,
					"dest": ev.Dest,
				})
			case "kept":
				runtime.EventsEmit(a.ctx, "progress:kept", map[string]string{
					"file": ev.File,
				})
			case "error":
				errMsg := ""
				if ev.Err != nil {
					errMsg = ev.Err.Error()
				}
				runtime.EventsEmit(a.ctx, "progress:error", map[string]string{
					"file":  ev.File,
					"error": errMsg,
				})
			}
		}
		runtime.EventsEmit(a.ctx, "progress:finished", map[string]string{})
	}()

	go func() {
		if err := converter.Run(cfg, progress); err != nil {
			runtime.EventsEmit(a.ctx, "progress:error", map[string]string{
				"file":  "",
				"error": err.Error(),
			})
		}
	}()

	return ""
}

func baseNames(paths []string) []string {
	names := make([]string, len(paths))
	for i, p := range paths {
		names[i] = filepath.Base(p)
	}
	return names
}
