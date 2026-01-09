package scanner

import (
	"os"
	"path/filepath"
)

// Config holds configuration for the scanner
type Config struct {
	RootPath string
	Excluded []string
}

// DefaultExcluded returns the default list of excluded directories
func DefaultExcluded() []string {
	return []string{
		"node_modules",
		".git",
		"dist",
		"build",
		"vendor",
	}
}

// Scan walks the directory tree and returns a list of files to process
func Scan(config Config) ([]string, error) {
	var files []string
	excludedMap := make(map[string]bool)
	for _, dir := range config.Excluded {
		excludedMap[dir] = true
	}

	err := filepath.Walk(config.RootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			if excludedMap[info.Name()] {
				return filepath.SkipDir
			}
			return nil
		}

		files = append(files, path)
		return nil
	})

	return files, err
}