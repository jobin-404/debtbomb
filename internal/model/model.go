package model

import "time"

// DebtBomb represents a technical debt item found in the codebase
type DebtBomb struct {
	ID       string    `json:"id"`
	File     string    `json:"file"`
	Line     int       `json:"line"`
	Expire   time.Time `json:"expire"`
	Owner    string    `json:"owner,omitempty"`
	Ticket   string    `json:"ticket,omitempty"`
	Reason   string    `json:"reason,omitempty"`
	Severity string    `json:"severity,omitempty"`
	RawText  string    `json:"rawText"`
	Snippet  string    `json:"snippet"`

	IsExpired bool `json:"isExpired"`

	GitAuthor  string    `json:"gitAuthor,omitempty"`
	CommitHash string    `json:"commitHash,omitempty"`
	CommitDate time.Time `json:"commitDate,omitempty"`
}

// DebtEvent represents a change in state of a DebtBomb
type DebtEvent struct {
	ID       string
	Type     string // expired | expiring_soon
	File     string
	Line     int
	Expires  time.Time
	DaysLeft int
	Owner    string
	Reason   string
	Snippet  string
}