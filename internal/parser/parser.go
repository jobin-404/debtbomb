package parser

import (
	"bufio"
	"github.com/jobin-404/debtbomb/internal/model"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"
)

var singleLineRegex = regexp.MustCompile(`(?:\/\/|#|--|\/\*)\s*@debtbomb\((.*?)\)`)

// Pattern for multi line start: // @debtbomb
var multiLineStartRegex = regexp.MustCompile(`(?:\/\/|#|--|\/\*)\s*@debtbomb\s*$`)

// Parse scans the content and returns a list of DebtBombs
func Parse(filename string, reader io.Reader) ([]model.DebtBomb, error) {
	var bombs []model.DebtBomb
	scanner := bufio.NewScanner(reader)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		lineText := scanner.Text()

		// Check for single line format
		if matches := singleLineRegex.FindStringSubmatch(lineText); len(matches) > 1 {
			bomb, err := parseAttributes(matches[1])
			if err == nil {
				bomb.File = filename
				bomb.Line = lineNum
				bomb.RawText = strings.TrimSpace(lineText)
				bombs = append(bombs, bomb)
			}
			continue
		}

		if idx := strings.Index(lineText, "@debtbomb"); idx != -1 {
			// Check if it was already handled by singleLineRegex (parentheses format)
			if singleLineRegex.MatchString(lineText) {
				continue
			}
			content := lineText[idx+len("@debtbomb"):]
			bomb, err := parseKeyValueStyle(content)
			if err == nil {
				bomb.File = filename
				bomb.Line = lineNum
				bomb.RawText = strings.TrimSpace(lineText)
				bombs = append(bombs, bomb)
			}
		}
	}

	return bombs, scanner.Err()
}

func parseAttributes(attrString string) (model.DebtBomb, error) {
	// Format: expire=2026-02-10, owner=pricing, ticket=JIRA-123
	bomb := model.DebtBomb{}
	parts := strings.Split(attrString, ",")
	foundExpire := false

	for _, part := range parts {
		kv := strings.SplitN(strings.TrimSpace(part), "=", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.TrimSpace(kv[0])
		val := strings.TrimSpace(kv[1])

		switch key {
		case "expire":
			t, err := time.Parse("2006-01-02", val)
			if err == nil {
				bomb.Expire = t
				foundExpire = true
			}
		case "owner":
			bomb.Owner = val
		case "ticket":
			bomb.Ticket = val
		case "reason":
			bomb.Reason = val
		}
	}

	if !foundExpire {
		return bomb, fmt.Errorf("missing expire date")
	}
	return bomb, nil
}

func parseKeyValueStyle(content string) (model.DebtBomb, error) {
	
	bomb := model.DebtBomb{}
	foundExpire := false

	// Regex to find key: value pairs where key is one of our expected fields
	kvRegex := regexp.MustCompile(`(expire|owner|ticket|reason)\s*:\s*([^\/\#\*]+)`)
	matches := kvRegex.FindAllStringSubmatch(content, -1)

	for _, match := range matches {
		key := match[1]
		val := strings.TrimSpace(match[2])

		switch key {
		case "expire":
			t, err := time.Parse("2006-01-02", val)
			if err == nil {
				bomb.Expire = t
				foundExpire = true
			}
		case "owner":
			bomb.Owner = val
		case "ticket":
			bomb.Ticket = val
		case "reason":
			bomb.Reason = val
		}
	}

	if !foundExpire {
		return bomb, fmt.Errorf("missing expire date")
	}
	return bomb, nil
}