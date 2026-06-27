package converter

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
)

func CheckExtractor() bool {
	return hasCommand("unrar") || hasCommand("7z")
}

func Convert(source, dest string) error {
	tmpDir, err := os.MkdirTemp("", "cbr2cbz-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmpDir)

	if err := extractRAR(source, tmpDir); err != nil {
		return err
	}

	files, err := collectFiles(tmpDir)
	if err != nil || len(files) == 0 {
		return fmt.Errorf("nenhum arquivo extraído de %s", filepath.Base(source))
	}

	return createCBZ(dest, tmpDir, files)
}

func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err = io.Copy(out, in); err != nil {
		return err
	}

	return out.Sync()
}

func extractRAR(source, dest string) error {
	var cmd *exec.Cmd

	if hasCommand("unrar") {
		cmd = exec.Command("unrar", "x", "-inul", source, dest)
	} else {
		cmd = exec.Command("7z", "x", source, "-o"+dest, "-y")
	}

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		if stderr.Len() > 0 {
			return fmt.Errorf(strings.TrimSpace(stderr.String()))
		}
		return err
	}

	return nil
}

func createCBZ(dest, base string, files []string) error {
	out, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer out.Close()

	zw := zip.NewWriter(out)
	defer zw.Close()

	for _, file := range files {
		rel, _ := filepath.Rel(base, file)

		w, err := zw.Create(rel)
		if err != nil {
			return err
		}

		src, err := os.Open(file)
		if err != nil {
			return err
		}

		if _, err = io.Copy(w, src); err != nil {
			src.Close()
			return err
		}
		src.Close()
	}

	return nil
}

func collectFiles(dir string) ([]string, error) {
	var files []string

	err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		files = append(files, path)
		return nil
	})

	sort.Strings(files)
	return files, err
}

func hasCommand(cmd string) bool {
	c := exec.Command(cmd)
	var b bytes.Buffer
	c.Stdout = &b
	c.Stderr = &b
	return c.Run() == nil
}
