package engine

import (
	"debtbomb/internal/model"
	"debtbomb/internal/parser"
	"debtbomb/internal/scanner"
	"os"
	"sort"
	"sync"
	"time"
)

// Run executes the debtbomb scan and returns all found items
func Run(rootPath string) ([]model.DebtBomb, error) {
	// 1. Scan for files
	files, err := scanner.Scan(scanner.Config{
		RootPath: rootPath,
		Excluded: scanner.DefaultExcluded(),
	})
	if err != nil {
		return nil, err
	}

	// 2. Parse files concurrently
	var allBombs []model.DebtBomb
	var mu sync.Mutex
	var wg sync.WaitGroup

	// Limit concurrency to avoid too many open files
	semaphore := make(chan struct{}, 100)

	for _, file := range files {
		wg.Add(1)
		go func(f string) {
			defer wg.Done()
			semaphore <- struct{}{}
			defer func() { <-semaphore }()

			fileHandle, err := os.Open(f)
			if err != nil {
				return // Skip unreadable files
			}
			defer fileHandle.Close()

			bombs, err := parser.Parse(f, fileHandle)
			if err == nil && len(bombs) > 0 {
				mu.Lock()
				allBombs = append(allBombs, bombs...)
				mu.Unlock()
			}
		}(file)
	}

	wg.Wait()

	// 3. Check expiration status
	today := time.Now().Truncate(24 * time.Hour)
	for i := range allBombs {
		if today.After(allBombs[i].Expire) {
			allBombs[i].IsExpired = true
		}
	}

	// 4. Sort by expire date
	sort.Slice(allBombs, func(i, j int) bool {
		return allBombs[i].Expire.Before(allBombs[j].Expire)
	})

	return allBombs, nil
}
