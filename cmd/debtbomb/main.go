package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/jobin-404/debtbomb/internal/config"
	"github.com/jobin-404/debtbomb/internal/engine"
	"github.com/jobin-404/debtbomb/internal/jira"
	"github.com/jobin-404/debtbomb/internal/model"
	"github.com/jobin-404/debtbomb/internal/notify"
	"github.com/jobin-404/debtbomb/internal/output"
	"github.com/jobin-404/debtbomb/internal/report"
	"github.com/jobin-404/debtbomb/internal/state"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file if it exists
	_ = godotenv.Load()

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
	case "report":
		runReport()
	case "notify":
		runNotify()
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
	fmt.Println("  report  Show aggregated statistics about technical debt")
	fmt.Println("  notify  Notify about expired or expiring debtbombs")
}

func runCheck() {
	checkCmd := flag.NewFlagSet("check", flag.ExitOnError)
	jsonOutput := checkCmd.Bool("json", false, "Output in JSON format")
	warnDays := checkCmd.Int("warn-in-days", 0, "Warn about bombs expiring within N days")
	checkCmd.Parse(os.Args[2:])

	bombs, err := engine.Run(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning: %v\n", err)
		os.Exit(1)
	}

	var expired []model.DebtBomb
	var warning []model.DebtBomb

	for _, b := range bombs {
		if b.IsExpired {
			expired = append(expired, b)
		}
	}

	// Check for warning window
	if *warnDays > 0 {
		today := time.Now().Truncate(24 * time.Hour)
		warningDate := today.AddDate(0, 0, *warnDays)

		for _, b := range bombs {
			if !b.IsExpired {
				// If expire date is before or equal to warning date
				if !b.Expire.After(warningDate) {
					warning = append(warning, b)
				}
			}
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

	if hasExpired || len(warning) > 0 {
		output.PrintCheckReport(expired, warning, *warnDays)
		if hasExpired {
			os.Exit(1)
		}
		// If only warnings, exit 0
		os.Exit(0)
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

func runReport() {
	reportCmd := flag.NewFlagSet("report", flag.ExitOnError)
	jsonOutput := reportCmd.Bool("json", false, "Output in JSON format")
	reportCmd.Parse(os.Args[2:])

	bombs, err := engine.Run(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning: %v\n", err)
		os.Exit(1)
	}

	r := report.Generate(bombs)

	if *jsonOutput {
		output.PrintReportJSON(r)
	} else {
		output.PrintReport(r)
	}
}

func runNotify() {
	notifyCmd := flag.NewFlagSet("notify", flag.ExitOnError)
	expired := notifyCmd.Bool("expired", false, "Process expired bombs")
	expireInDays := notifyCmd.Int("expire-in-days", 0, "Process bombs expiring in N days")
	notifyCmd.Parse(os.Args[2:])

	// Load Config
	cfg, err := config.Load(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Warning: failed to load config: %v\n", err)
		cfg = &config.Config{}
	}

	// Load State
	st, err := state.Load(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error loading state: %v\n", err)
		os.Exit(0)
	}

	// Init Jira
	var jClient *jira.Client
	if cfg.Jira.BaseURL != "" && cfg.Jira.Email != "" && cfg.Jira.APIToken != "" {
		jClient = jira.NewClient(cfg.Jira.BaseURL, cfg.Jira.Email, cfg.Jira.APIToken)
	}

	router := &notify.Router{
		Config: cfg,
		Jira:   jClient,
		State:  st,
	}

	bombs, err := engine.Run(".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error scanning: %v\n", err)
		os.Exit(0)
	}

	if err := router.SyncAndNotify(bombs, *expireInDays, *expired); err != nil {

		fmt.Fprintf(os.Stderr, "Error during sync/notify: %v\n", err)
		os.Exit(0)
	}
}
