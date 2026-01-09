package main

import (
	"github.com/jobin-404/debtbomb/internal/engine"
	"github.com/jobin-404/debtbomb/internal/model"
	"github.com/jobin-404/debtbomb/internal/output"
	"flag"
	"fmt"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {
	case "check":
		runCheck()
	case "list":
		runList()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage: debtbomb <command> [flags]")
	fmt.Println("Commands:")
	fmt.Println("  check   Scan for expired debtbombs and exit 1 if found")
	fmt.Println("  list    List all debtbombs")
}

func runCheck() {
	checkCmd := flag.NewFlagSet("check", flag.ExitOnError)
	warnDays := checkCmd.Int("warn-in-days", 0, "Warn about bombs expiring within N days")
	checkCmd.Parse(os.Args[2:])

	bombs, err := engine.Run(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning: %v\n", err)
		os.Exit(1)
	}

	hasExpired := false
	for _, b := range bombs {
		if b.IsExpired {
			hasExpired = true
		}
	}

	if hasExpired {
		output.PrintHuman(bombs, false)
		os.Exit(1)
	}

	if *warnDays > 0 {
		warningDate := time.Now().AddDate(0, 0, *warnDays)
		foundWarning := false
		for _, b := range bombs {
			if !b.IsExpired && b.Expire.Before(warningDate) {
				if !foundWarning {
					fmt.Printf("⚠️  Warning: DebtBombs expiring within %d days:\n\n", *warnDays)
					foundWarning = true
				}
				fmt.Printf("%s:%d expires on %s\n", b.File, b.Line, b.Expire.Format("2006-01-02"))
			}
		}
	}

	os.Exit(0)
}

func runList() {
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	expiredOnly := listCmd.Bool("expired", false, "Show only expired bombs")
	jsonOutput := listCmd.Bool("json", false, "Output in JSON format")
	listCmd.Parse(os.Args[2:])

	bombs, err := engine.Run(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning: %v\n", err)
		os.Exit(1)
	}

	// Filter if needed
	if *expiredOnly {
		var expired []model.DebtBomb
		for _, b := range bombs {
			if b.IsExpired {
				expired = append(expired, b)
			}
		}
		bombs = expired
	}

	if *jsonOutput {
		output.PrintJSON(bombs)
	} else {
		output.PrintHuman(bombs, true)
	}
}