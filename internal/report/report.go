package report

import (
	"path/filepath"
	"sort"
	"time"

	"github.com/jobin-404/debtbomb/internal/model"
)

type Report struct {
	TotalCount int            `json:"totalCount"`
	ByOwner    []CountItem    `json:"byOwner"`
	ByFolder   []CountItem    `json:"byFolder"`
	ByReason   []CountItem    `json:"byReason"`
	ByUrgency  UrgencyStats   `json:"byUrgency"`
	Oldest     *model.DebtBomb `json:"oldest,omitempty"`
	Newest     *model.DebtBomb `json:"newest,omitempty"`
}

type CountItem struct {
	Key   string `json:"key"`
	Count int    `json:"count"`
}

type UrgencyStats struct {
	Expired      int `json:"expired"`
	Within30Days int `json:"within30Days"`
	Within90Days int `json:"within90Days"`
	MoreThan90Days int `json:"moreThan90Days"`
}

func Generate(bombs []model.DebtBomb) Report {
	report := Report{
		TotalCount: len(bombs),
		ByOwner:    make([]CountItem, 0),
		ByFolder:   make([]CountItem, 0),
		ByReason:   make([]CountItem, 0),
	}

	if len(bombs) == 0 {
		return report
	}

	ownerCounts := make(map[string]int)
	folderCounts := make(map[string]int)
	reasonCounts := make(map[string]int)

	today := time.Now().Truncate(24 * time.Hour)
	day30 := today.AddDate(0, 0, 30)
	day90 := today.AddDate(0, 0, 90)

	report.Oldest = &bombs[0]
	report.Newest = &bombs[0]

	for _, b := range bombs {
		owner := b.Owner
		if owner == "" {
			owner = "(no owner)"
		}
		ownerCounts[owner]++

		dir := filepath.Dir(b.File)
		if dir == "." {
			dir = "(root)"
		}
		folderCounts[dir]++

		reason := b.Reason
		if reason == "" {
			reason = "(no reason)"
		}
		reasonCounts[reason]++
		
		if b.IsExpired {
			report.ByUrgency.Expired++
		} else {
			if b.Expire.Before(day30) {
				report.ByUrgency.Within30Days++
			}
			if b.Expire.Before(day90) {
				report.ByUrgency.Within90Days++
			} else {
				report.ByUrgency.MoreThan90Days++
			}
		}

		// Oldest/Newest
		if b.Expire.Before(report.Oldest.Expire) {
			report.Oldest = &b
		}
		if b.Expire.After(report.Newest.Expire) {
			report.Newest = &b
		}
	}

	// Convert maps to slices and sort
	report.ByOwner = mapToSortedSlice(ownerCounts)
	report.ByFolder = mapToSortedSlice(folderCounts)
	report.ByReason = mapToSortedSlice(reasonCounts)

	return report
}

func mapToSortedSlice(m map[string]int) []CountItem {
	var s []CountItem
	for k, v := range m {
		s = append(s, CountItem{Key: k, Count: v})
	}
	sort.Slice(s, func(i, j int) bool {
		if s[i].Count != s[j].Count {
			return s[i].Count > s[j].Count // Descending count
		}
		return s[i].Key < s[j].Key // Ascending key
	})
	return s
}
