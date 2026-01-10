package scanner

import (
	"bufio"
	"os"
	"path/filepath"
	"strings"
)

// Config holds configuration for the scanner
type Config struct {
	RootPath string
	Excluded []string
}

// DefaultExcluded returns the default list of excluded directories
func DefaultExcluded() []string {
	return []string{
		"node_modules", ".git", ".svn", ".hg", "vendor", "dist", "build", "out", "target", "coverage",
		".tmp", ".temp", ".cache", ".next", ".nuxt", ".turbo", ".parcel-cache", ".esbuild",
		".gradle", ".mvn",
		"__pycache__", ".venv", "venv", "env", ".mypy_cache", ".pytest_cache",
		"bin", "pkg",
		"obj",
		".storybook", ".vite",
		"third_party", "third-party", "external", "deps", "bower_components",
		".terraform", ".terragrunt-cache", ".cdk.out", "pulumi",
		".idea", ".vscode",
	}
}

// isIgnoredExt checks if the file extension is one we should ignore
func isIgnoredExt(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	switch ext {
	case ".png", ".jpg", ".jpeg", ".gif", ".svg", ".ico", ".webp", ".bmp", ".tiff",
		".mp4", ".mov", ".avi", ".mkv", ".mp3", ".wav", ".flac", ".ogg":
		return true
	case ".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx":
		return true
	case ".zip", ".tar", ".gz", ".7z", ".rar", ".jar", ".war":
		return true
	case ".exe", ".dll", ".so", ".dylib", ".bin", ".ds_store",
		".o", ".a", ".test", ".class", ".pyc":
		return true
	case ".log":
		return true
	case ".eot", ".ttf", ".woff", ".woff2":
		return true
	case ".min.js", ".min.css", ".lock":
		return true
	}
	return false
}

// loadIgnoreFile reads .debtbombignore if it exists
func loadIgnoreFile(root string) ([]string, error) {
	path := filepath.Join(root, ".debtbombignore")
	file, err := os.Open(path)
	if os.IsNotExist(err) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var patterns []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		patterns = append(patterns, line)
	}
	return patterns, scanner.Err()
}

// Scan walks the directory tree and returns a list of files to process
func Scan(config Config) ([]string, error) {
	var files []string
	excludedMap := make(map[string]bool)
	for _, dir := range config.Excluded {
		excludedMap[dir] = true
	}

	ignorePatterns, _ := loadIgnoreFile(config.RootPath)

	err := filepath.Walk(config.RootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(config.RootPath, path)
		if err != nil {
			relPath = path
		}

		if info.IsDir() {
			if excludedMap[info.Name()] {
				return filepath.SkipDir
			}
			for _, p := range ignorePatterns {
				pClean := strings.TrimSuffix(p, "/")
				if matched, _ := filepath.Match(pClean, info.Name()); matched {
					return filepath.SkipDir
				}
				if matched, _ := filepath.Match(pClean, relPath); matched {
					return filepath.SkipDir
				}
			}
			return nil
		}
		if info.Size() > 1024*1024 {
			return nil
		}
		if isIgnoredExt(path) {
			return nil
		}
		for _, p := range ignorePatterns {
			pClean := strings.TrimSuffix(p, "/")
			// Match filename
			if matched, _ := filepath.Match(pClean, info.Name()); matched {
				return nil
			}
			// Match relative path
			if matched, _ := filepath.Match(pClean, relPath); matched {
				return nil
			}
		}

		files = append(files, path)
		return nil
	})

	return files, err
}