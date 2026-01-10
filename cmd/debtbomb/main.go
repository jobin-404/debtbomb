package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/jobin-404/debtbomb/internal/engine"
	"github.com/jobin-404/debtbomb/internal/model"
	"github.com/jobin-404/debtbomb/internal/output"
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
	jsonOutput := checkCmd.Bool("json", false, "Output in JSON format")
	checkCmd.Parse(os.Args[2:])

	bombs, err := engine.Run(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning: %v\n", err)
		os.Exit(1)
	}

	var expired []model.DebtBomb
	for _, b := range bombs {
		if b.IsExpired {
			expired = append(expired, b)
		}
	}
	hasExpired := len(expired) > 0

	if *jsonOutput {
		output.PrintJSON(bombs)
		if hasExpired {
			os.Exit(1)
		}
		os.Exit(0)
	}

	if hasExpired {
		output.PrintCheckReport(expired)
		os.Exit(1)
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
		output.PrintTable(bombs)
	}
}