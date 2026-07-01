package main

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"

	"reforge/pkg/reforge"
)

type App struct {
	ctx context.Context

	mu         sync.Mutex
	cancelFunc context.CancelFunc
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
	scan, err := reforge.ScanArchives(dir)
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
	return reforge.CheckExtractor()
}

func (a *App) StartConversion(workDir string, skip, replace bool) string {
	if !reforge.CheckExtractor() {
		return "Nenhum extrator RAR encontrado (instale unrar ou 7z)"
	}

	a.mu.Lock()
	if a.cancelFunc != nil {
		a.mu.Unlock()
		return "Já existe uma conversão em andamento"
	}
	runCtx, cancel := context.WithCancel(a.ctx)
	a.cancelFunc = cancel
	a.mu.Unlock()

	cfg := reforge.Config{
		WorkDir: workDir,
		Skip:    skip,
		Replace: replace,
	}

	progress := make(chan reforge.ProgressEvent, 256)

	var (
		total, converted, kept, failed int
		start                          = time.Now()
	)

	scan, _ := reforge.ScanArchives(workDir)
	total = len(scan.CBRFiles) + len(scan.CBZFiles)

	go func() {
		for ev := range progress {
			switch ev.Status {
			case "converting":
				runtime.EventsEmit(a.ctx, "progress:converting", map[string]string{
					"file": ev.File,
				})
			case "done":
				converted++
				runtime.EventsEmit(a.ctx, "progress:done", map[string]string{
					"file": ev.File,
					"dest": ev.Dest,
				})
			case "kept":
				kept++
				runtime.EventsEmit(a.ctx, "progress:kept", map[string]string{
					"file": ev.File,
				})
			case "error":
				failed++
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

		a.mu.Lock()
		cancelled := runCtx.Err() != nil
		a.cancelFunc = nil
		a.mu.Unlock()

		elapsed := time.Since(start)
		summary := reforge.Summary{
			Total:       total,
			Converted:   converted,
			Kept:        kept,
			Failed:      failed,
			Cancelled:   cancelled,
			ElapsedMs:   elapsed.Milliseconds(),
			ElapsedText: formatElapsed(elapsed),
		}

		runtime.EventsEmit(a.ctx, "progress:finished", summary)
	}()

	go func() {
		if err := reforge.Run(runCtx, cfg, progress); err != nil {
			runtime.EventsEmit(a.ctx, "progress:error", map[string]string{
				"file":  "",
				"error": err.Error(),
			})
		}
	}()

	return ""
}

// CancelConversion interrompe uma conversão em andamento, se houver.
func (a *App) CancelConversion() {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.cancelFunc != nil {
		a.cancelFunc()
	}
}

// GetSettings retorna as preferências salvas do usuário.
func (a *App) GetSettings() reforge.Settings {
	s, err := loadSettings()
	if err != nil {
		return reforge.Settings{}
	}
	return s
}

// SaveSettings persiste as preferências do usuário em disco.
func (a *App) SaveSettings(settings reforge.Settings) string {
	if err := saveSettings(settings); err != nil {
		return err.Error()
	}
	return ""
}

func baseNames(paths []string) []string {
	names := make([]string, len(paths))
	for i, p := range paths {
		names[i] = filepath.Base(p)
	}
	return names
}

func formatElapsed(d time.Duration) string {
	d = d.Round(time.Second)
	if d < time.Minute {
		return d.String()
	}
	m := d / time.Minute
	s := (d % time.Minute) / time.Second
	return time.Duration(m*time.Minute + s*time.Second).String()
}

func settingsPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	appDir := filepath.Join(dir, "pInk-reforge")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return "", err
	}
	return filepath.Join(appDir, "settings.json"), nil
}

func loadSettings() (reforge.Settings, error) {
	path, err := settingsPath()
	if err != nil {
		return reforge.Settings{}, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return reforge.Settings{}, nil
		}
		return reforge.Settings{}, err
	}

	var s reforge.Settings
	if err := json.Unmarshal(data, &s); err != nil {
		return reforge.Settings{}, err
	}
	return s, nil
}

func saveSettings(s reforge.Settings) error {
	path, err := settingsPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}
