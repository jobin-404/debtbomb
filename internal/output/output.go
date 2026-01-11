package output

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/jobin-404/debtbomb/internal/model"
	"github.com/jobin-404/debtbomb/internal/report"
)

type jsonOutput struct {
	HasExpired bool       `json:"hasExpired"`
	Bombs      []jsonBomb `json:"bombs"`
}

type jsonBomb struct {
	File   string `json:"file"`
	Line   int    `json:"line"`
	Expire string `json:"expire"`
	Owner  string `json:"owner,omitempty"`
	Ticket string `json:"ticket,omitempty"`
	Reason string `json:"reason,omitempty"`
}

// PrintJSON prints the report in JSON format
func PrintJSON(bombs []model.DebtBomb) {
	hasExpired := false
	outputBombs := make([]jsonBomb, 0, len(bombs))

	for _, b := range bombs {
		if b.IsExpired {
			hasExpired = true
		}
		outputBombs = append(outputBombs, jsonBomb{
			File:   b.File,
			Line:   b.Line,
			Expire: b.Expire.Format("2006-01-02"),
			Owner:  b.Owner,
			Ticket: b.Ticket,
			Reason: b.Reason,
		})
	}

	out := jsonOutput{
		HasExpired: hasExpired,
		Bombs:      outputBombs,
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(out); err != nil {
		fmt.Fprintf(os.Stderr, "failed to encode json: %v\n", err)
	}
}

// PrintTable prints a clean ASCII table for the list command
func PrintTable(bombs []model.DebtBomb) {
	fmt.Printf("Found %d DebtBombs\n", len(bombs))
	if len(bombs) == 0 {
		return
	}

	// Headers
	hExpires := "Expires"
	hOwner := "Owner"
	hTicket := "Ticket"
	hLocation := "Location"

	// Initial widths based on headers
	wExpires := len(hExpires)
	wOwner := len(hOwner)
	wTicket := len(hTicket)
	wLocation := len(hLocation)

	// Calculate max widths
	for _, b := range bombs {
		dateStr := b.Expire.Format("2006-01-02")
		timeLeftStr := timeLeft(b.Expire)
		fullExpiresStr := fmt.Sprintf("%s %s", dateStr, timeLeftStr)
		if len(fullExpiresStr) > wExpires {
			wExpires = len(fullExpiresStr)
		}
		if len(b.Owner) > wOwner {
			wOwner = len(b.Owner)
		}
		if len(b.Ticket) > wTicket {
			wTicket = len(b.Ticket)
		}
		loc := fmt.Sprintf("%s:%d", b.File, b.Line)
		if len(loc) > wLocation {
			wLocation = len(loc)
		}
	}

	// Helper to print separator
	printSeparator := func() {
		fmt.Printf("+-%s-+-%s-+-%s-+-%s-+\n",
			strings.Repeat("-", wExpires),
			strings.Repeat("-", wOwner),
			strings.Repeat("-", wTicket),
			strings.Repeat("-", wLocation),
		)
	}

	printSeparator()
	fmt.Printf("| %-*s | %-*s | %-*s | %-*s |\n",
		wExpires, hExpires,
		wOwner, hOwner,
		wTicket, hTicket,
		wLocation, hLocation,
	)
	printSeparator()

	for _, b := range bombs {
		expiresWithTime := fmt.Sprintf("%s %s", b.Expire.Format("2006-01-02"), timeLeft(b.Expire))
		fmt.Printf("| %-*s | %-*s | %-*s | %-*s |\n",
			wExpires, expiresWithTime,
			wOwner, b.Owner,
			wTicket, b.Ticket,
			wLocation, fmt.Sprintf("%s:%d", b.File, b.Line),
		)
	}
	printSeparator()
}

// PrintCheckReport prints the failure report for the check command
func PrintCheckReport(expiredBombs []model.DebtBomb) {
	if len(expiredBombs) == 0 {
		return
	}

	fmt.Printf("DebtBomb exploded: %d expired\n\n", len(expiredBombs))

	for i, b := range expiredBombs {
		fmt.Printf("%s:%d\n", b.File, b.Line)
		fmt.Printf("Expired: %s\n", b.Expire.Format("2006-01-02"))
		if b.Owner != "" {
			fmt.Printf("Owner: %s\n", b.Owner)
		}
		if b.Ticket != "" {
			fmt.Printf("Ticket: %s\n", b.Ticket)
		}
		if b.Reason != "" {
			fmt.Printf("Reason: %s\n", b.Reason)
		}

		if i < len(expiredBombs)-1 {
			fmt.Println("")
		}
	}
}

func timeLeft(deadline time.Time) string {
	tLeft := time.Until(deadline)

	if tLeft <= 0 {
		return "(expired)"
	}

	h := uint16(tLeft.Hours())   // total amount of hours INTEGER
	d := h / 24                  // total amount of days INTEGER
	m := uint16(tLeft.Minutes()) // total amount of minutes INTEGER

	if h > 24 {
		return fmt.Sprintf("(%dd%dh)", d, h%24)
	} else {
		return fmt.Sprintf("(%dh%dm)", h, m%60)
	}
}

func PrintReportJSON(r report.Report) {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(r); err != nil {
		fmt.Fprintf(os.Stderr, "failed to encode json: %v\n", err)
	}
}

func PrintReport(r report.Report) {
	printSection("Debt by owner", r.ByOwner, 5)
	printSection("Debt by folder", r.ByFolder, 5)
	printSection("Debt by reason", r.ByReason, 5)

	fmt.Println("By urgency")
	fmt.Printf("  %-15s %d\n", "Expired", r.ByUrgency.Expired)
	fmt.Printf("  %-15s %d\n", "< 30 days", r.ByUrgency.Within30Days)
	fmt.Printf("  %-15s %d\n", "< 90 days", r.ByUrgency.Within90Days)
	fmt.Printf("  %-15s %d\n", "> 90 days", r.ByUrgency.MoreThan90Days)
	fmt.Println()

	fmt.Println("Extremes")
	if r.Oldest != nil {
		fmt.Printf("  Oldest: %s (%s) in %s:%d\n",
			r.Oldest.Expire.Format("2006-01-02"),
			timeLeft(r.Oldest.Expire),
			r.Oldest.File,
			r.Oldest.Line)
	}
	if r.Newest != nil {
		fmt.Printf("  Newest: %s (%s) in %s:%d\n",
			r.Newest.Expire.Format("2006-01-02"),
			timeLeft(r.Newest.Expire),
			r.Newest.File,
			r.Newest.Line)
	}
}

func printSection(title string, items []report.CountItem, limit int) {
	fmt.Println(title)
	count := 0
	for _, item := range items {
		if count >= limit {
			break
		}
		fmt.Printf("  %-20s %d\n", item.Key, item.Count)
		count++
	}
	if len(items) > limit {
		fmt.Printf("  ... and %d more\n", len(items)-limit)
	}
	fmt.Println()
}
