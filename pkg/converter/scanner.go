package converter

import (
	"os"
	"path/filepath"
	"strings"
)

func ScanArchives(workDir string) (ScanResult, error) {
	entries, err := os.ReadDir(workDir)
	if err != nil {
		return ScanResult{}, err
	}

	var result ScanResult

	for _, e := range entries {
		if e.IsDir() {
			continue
		}

		name := strings.ToLower(e.Name())

		switch {
		case strings.HasSuffix(name, ".cbr"):
			result.CBRFiles = append(result.CBRFiles, filepath.Join(workDir, e.Name()))
		case strings.HasSuffix(name, ".cbz"):
			result.CBZFiles = append(result.CBZFiles, filepath.Join(workDir, e.Name()))
		}
	}

	return result, nil
}
