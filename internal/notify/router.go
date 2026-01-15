package notify

import (
	"fmt"
	"time"

	"github.com/jobin-404/debtbomb/internal/config"
	"github.com/jobin-404/debtbomb/internal/jira"
	"github.com/jobin-404/debtbomb/internal/model"
	"github.com/jobin-404/debtbomb/internal/state"
)

type Router struct {
	Config *config.Config
	Jira   *jira.Client
	State  *state.State
}

func (r *Router) SyncAndNotify(bombs []model.DebtBomb, checkDays int, expiredOnly bool) error {
	today := time.Now().Truncate(24 * time.Hour)

	// Separate bombs
	var expired []model.DebtBomb
	var expiring []model.DebtBomb

	// Map to track current expired bombs (by ID) to handle cleanup later
	currentExpiredIDs := make(map[string]bool)

	for _, b := range bombs {
		if b.IsExpired {
			expired = append(expired, b)
			currentExpiredIDs[b.ID] = true
		} else {
			if expiredOnly {
				continue
			}
			// Check if expiring soon
			daysLeft := int(b.Expire.Sub(today).Hours() / 24)
			if daysLeft >= 0 {
				if checkDays > 0 && daysLeft > checkDays {
					continue
				}
				expiring = append(expiring, b)
			}
		}
	}

	// 1. Process Expired Bombs (Jira Sync + Notify)
	for _, b := range expired {
		ticketKey := r.State.GetTicket(b.ID)

		if ticketKey == "" {
			// New Expiration
			// Create Jira Ticket if configured
			if r.shouldCreateJira() {
				key, err := r.createJiraTicket(b)
				if err != nil {
					fmt.Printf("Failed to create Jira ticket for %s: %v\n", b.ID, err)
					// Continue?
				} else {
					fmt.Printf("Created Jira ticket %s for %s\n", key, b.ID)
					r.State.SetTicket(b.ID, key)
					ticketKey = key

					// Notify "Expired" (State Transition)
					r.notifyExpired(b, ticketKey)
				}
			} else {
				// No Jira configured, but maybe we still notify?
				r.notifyExpired(b, "")
			}
		} else {
			if b.Severity != "" && r.Jira != nil {
				err := r.Jira.UpdatePriority(ticketKey, b.Severity)
				if err != nil {
				}
			}
		}
	}

	// 2. Process Expiring Bombs (Notify only)
	for _, b := range expiring {
		daysLeft := int(b.Expire.Sub(today).Hours() / 24)
		r.notifyExpiringSoon(b, daysLeft)
	}

	// 3. Cleanup Old Tickets
	storedMap := r.State.Snapshot()
	for bombID, ticketKey := range storedMap {
		if !currentExpiredIDs[bombID] {
			// Bomb is no longer expired (or deleted)
			// Close Jira ticket
			if r.Jira != nil {
				if err := r.Jira.CloseTicket(ticketKey); err != nil {
					fmt.Printf("Failed to close ticket %s: %v\n", ticketKey, err)
				} else {
					fmt.Printf("Closed ticket %s for %s\n", ticketKey, bombID)
				}
			}
			// Remove from state
			r.State.RemoveTicket(bombID)
		}
	}

	return r.State.Save()
}

func (r *Router) shouldCreateJira() bool {
	if r.Config == nil || r.Jira == nil {
		return false
	}
	// Check if any notify rule uses jira
	for _, n := range r.Config.Notify {
		if n.Via == "jira" && n.On == "expired" {
			return true
		}
	}
	return false
}

func (r *Router) createJiraTicket(b model.DebtBomb) (string, error) {
	// Find project and issue type from config
	project := r.Config.Jira.DefaultProject
	issueType := r.Config.Jira.IssueType

	summary := fmt.Sprintf("Expired tech debt: %s", b.Reason)
	description := fmt.Sprintf("File: %s\nExpires: %s\nOwner: %s\nSeverity: %s\n\nSnippet:\n%s",
		b.File, b.Expire.Format("2006-01-02"), b.Owner, b.Severity, b.Snippet)

	return r.Jira.CreateTicket(project, summary, description, issueType, b.Severity)
}

func (r *Router) notifyExpired(b model.DebtBomb, ticketKey string) {
	msg := FormatExpiredMessage(b, ticketKey)
	r.sendNotifications("expired", 0, msg)
}

func (r *Router) notifyExpiringSoon(b model.DebtBomb, daysLeft int) {
	msg := FormatWarningMessage(b, daysLeft)
	r.sendNotifications("expiring_soon", daysLeft, msg)
}

func (r *Router) sendNotifications(on string, days int, msg string) {
	for _, n := range r.Config.Notify {
		if n.On != on {
			continue
		}
		if on == "expiring_soon" && n.Days != days {
			continue
		}

		var err error
		switch n.Via {
		case "slack":
			url := r.Config.GetSlackWebhook()
			if url != "" {
				err = SendSlack(url, msg)
			}
		case "discord":
			url := r.Config.GetDiscordWebhook()
			if url != "" {
				err = SendDiscord(url, msg)
			}
		case "teams":
			url := r.Config.GetTeamsWebhook()
			if url != "" {
				err = SendTeams(url, msg)
			}
		}

		if err != nil {
			fmt.Printf("Failed to send notification via %s: %v\n", n.Via, err)
		}
	}
}
