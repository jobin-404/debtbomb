package model

import "time"

// DebtBomb represents a technical debt item found in the codebase
type DebtBomb struct {
	File      string    `json:"file"`
	Line      int       `json:"line"`
	Expire    time.Time `json:"expire"`
	Owner     string    `json:"owner,omitempty"`
	Ticket    string    `json:"ticket,omitempty"`
	Reason    string    `json:"reason,omitempty"`
	RawText   string    `json:"rawText"`
	IsExpired bool      `json:"isExpired"`
}
