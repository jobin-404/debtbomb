package notify

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jobin-404/debtbomb/internal/model"
)

func SendSlack(url, message string) error {
	payload := map[string]string{"text": message}
	return sendWebhook(url, payload)
}

func SendDiscord(url, message string) error {
	payload := map[string]string{"content": message}
	return sendWebhook(url, payload)
}

func SendTeams(url, message string) error {
	payload := map[string]string{"text": message}
	return sendWebhook(url, payload)
}

func sendWebhook(url string, payload interface{}) error {
	data, err := json.Marshal(payload)
	if err != nil {
		return err
	}
	resp, err := http.Post(url, "application/json", bytes.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("webhook failed with status %d", resp.StatusCode)
	}
	return nil
}

func FormatExpiredMessage(bomb model.DebtBomb, ticketKey string) string {
	msg := fmt.Sprintf("🚨 DebtBomb exploded\n%s\n%s\nOwner: %s\nExpires: %s",
		bomb.File, bomb.Reason, bomb.Owner, bomb.Expire.Format("2006-01-02"))
	
	if bomb.Severity != "" {
		msg += fmt.Sprintf("\nSeverity: %s", bomb.Severity)
	}
	
	if ticketKey != "" {
		msg += fmt.Sprintf("\nJira: %s", ticketKey)
	}
	
	return msg
}

func FormatWarningMessage(bomb model.DebtBomb, daysLeft int) string {
	return fmt.Sprintf("⏳ DebtBomb warning (%d days left)\n%s\n%s\nOwner: %s\nExpires: %s",
		daysLeft, bomb.File, bomb.Reason, bomb.Owner, bomb.Expire.Format("2006-01-02"))
}
