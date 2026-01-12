package engine

import (
	"github.com/jobin-404/debtbomb/internal/model"
	"github.com/jobin-404/debtbomb/internal/parser"
	"github.com/jobin-404/debtbomb/internal/scanner"
	"os"
	"runtime"
	"sort"
	"sync"
	"time"
)

// Run executes the debtbomb scan and returns all found items
func Run(rootPath string) ([]model.DebtBomb, error) {
	filesChan := make(chan string, 100)
	resultsChan := make(chan []model.DebtBomb, 100)
	errChan := make(chan error, 1)

	go func() {
		err := scanner.Scan(scanner.Config{
			RootPath: rootPath,
			Excluded: scanner.DefaultExcluded(),
		}, filesChan)
		if err != nil {
			errChan <- err
		}
	}()

	var wg sync.WaitGroup
	numWorkers := runtime.NumCPU() * 2

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for file := range filesChan {
				fileHandle, err := os.Open(file)
				if err != nil {
					continue
				}

				bombs, err := parser.Parse(file, fileHandle)
				fileHandle.Close()

				if err == nil && len(bombs) > 0 {
					resultsChan <- bombs
				}
			}
		}()
	}

	go func() {
		wg.Wait()
		close(resultsChan)
	}()
	var allBombs []model.DebtBomb
	for bombs := range resultsChan {
		allBombs = append(allBombs, bombs...)
	}

	select {
	case err := <-errChan:
		return nil, err
	default:
	}

	today := time.Now().Truncate(24 * time.Hour)
	for i := range allBombs {
		if today.After(allBombs[i].Expire) {
			allBombs[i].IsExpired = true
		}
	}
	sort.Slice(allBombs, func(i, j int) bool {
		return allBombs[i].Expire.Before(allBombs[j].Expire)
	})

	return allBombs, nil
}