package output

import (
	"debtbomb/internal/model"
	"encoding/json"
	"fmt"
	"os"
	"text/tabwriter"
)

// PrintHuman prints a human-readable report
func PrintHuman(bombs []model.DebtBomb, showAll bool) {
	expiredCount := 0
	for _, b := range bombs {
		if b.IsExpired {
			expiredCount++
		}
	}

	if expiredCount > 0 {
		fmt.Printf("💣 %d DebtBombs exploded\n\n", expiredCount)
	} else if !showAll {
		fmt.Println("✨ No expired DebtBombs found")
		return
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 2, ' ', 0)
	for _, b := range bombs {
		if !showAll && !b.IsExpired {
			continue
		}

		status := ""
		if b.IsExpired {
			status = "[EXPIRED]"
		}

		fmt.Fprintf(w, "%s:%d\t%s %s\n", b.File, b.Line, status, b.Expire.Format("2006-01-02"))
		if b.Owner != "" {
			fmt.Fprintf(w, "\tOwner: %s\n", b.Owner)
		}
		if b.Ticket != "" {
			fmt.Fprintf(w, "\tTicket: %s\n", b.Ticket)
		}
		if b.Reason != "" {
			fmt.Fprintf(w, "\tReason: %s\n", b.Reason)
		}
		fmt.Fprintln(w, "")
	}
	w.Flush()
}

// PrintJSON prints the report in JSON format
func PrintJSON(bombs []model.DebtBomb) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(bombs)
}
