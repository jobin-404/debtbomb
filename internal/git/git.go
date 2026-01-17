package git

import (
	"bufio"
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type BlameInfo struct {
	Hash   string
	Author string
	Date   time.Time
}

func GetBlame(filePath string) (map[int]BlameInfo, error) {

	if !isGitRepo(filepath.Dir(filePath)) {
		return nil, os.ErrNotExist
	}

	cmd := exec.Command("git", "blame", "--line-porcelain", filePath)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}

	result := make(map[int]BlameInfo)
	scanner := bufio.NewScanner(&out)

	var currentHash string
	var currentAuthor string
	var currentDate time.Time
	var currentLine int

	for scanner.Scan() {
		line := scanner.Text()

		// In git blame --line-porcelain, the actual source code line
		// always starts with a literal tab character.
		if strings.HasPrefix(line, "\t") {
			result[currentLine] = BlameInfo{
				Hash:   currentHash,
				Author: currentAuthor,
				Date:   currentDate,
			}
			continue
		}

		parts := strings.SplitN(line, " ", 2)
		if len(parts) < 2 {
			continue
		}

		key := parts[0]

		// Standard SHA-1 hashes are 40 characters.
		if len(key) == 40 {
			currentHash = key
			fields := strings.Fields(parts[1])
			if len(fields) >= 2 {
				if l, err := strconv.Atoi(fields[1]); err == nil {
					currentLine = l
				}
			}
			continue
		}

		switch key {
		case "author":
			currentAuthor = parts[1]
		case "author-time":
			if ts, err := strconv.ParseInt(parts[1], 10, 64); err == nil {
				currentDate = time.Unix(ts, 0)
			}
		}
	}

	return result, nil
}

func isGitRepo(dir string) bool {
	for {
		if _, err := os.Stat(filepath.Join(dir, ".git")); err == nil {
			return true
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return false
		}
		dir = parent
	}
}
