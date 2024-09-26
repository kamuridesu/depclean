package main

import (
	"fmt"
	"io/fs"
	"math"
	"os"
	"path/filepath"
	"reflect"
	"strings"
)

type PruneData struct {
	Paths []string
	Size  string
}

func (p *PruneData) Equals(other *PruneData) bool {
	return p.Size == other.Size && reflect.DeepEqual(p.Paths, other.Paths)
}

func formatSize(number int) string {
	units := [8]string{"", "Ki", "Mi", "Gi", "Ti", "Pi", "Ei", "Zi"}
	n := float64(number)
	for _, unit := range units {
		if math.Abs(n) < 1024.0 {
			return fmt.Sprintf("%.1f%sB", n, unit)
		}
		n /= 1024.0
	}
	return fmt.Sprintf("%.1fYiB", n)
}

func GetTotalSize(root string) int {
	size := 0
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			panic(err)
		}
		info, err := d.Info()
		if err != nil {
			panic(err)
		}
		if info.Mode()&os.ModeSymlink == 0 {
			if !d.IsDir() {
				size += int(info.Size())
			}
		}
		return nil
	})
	return size
}

func MapAllPaths(root string) *PruneData {
	bar := NewIndefiniteLoadingBar()
	go bar.start()
	paths := []string{}
	size := 0
	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if strings.Contains(path, "node_modules") || strings.Contains(path, "venv") || strings.Contains(path, ".venv") {
			if err != nil {
				panic(err)
			}
			info, err := d.Info()
			if err != nil {
				panic(err)
			}
			if info.Mode()&os.ModeSymlink == 0 && d.IsDir() {
				paths = append(paths, path)
				size += GetTotalSize(path)
				return nil
			}
		}
		return nil
	})
	bar.end()
	return &PruneData{
		Size:  formatSize(size),
		Paths: paths,
	}
}

func DeleteAll(paths []string) {
	progress := startProgress("Deleting ALL files")
	totalLength := len(paths)
	for idx, path := range paths {
		progress.progress(totalLength, idx)
		if err := os.RemoveAll(path); err != nil {
			fmt.Println(err)
			continue
		}
	}
	progress.end()
}
