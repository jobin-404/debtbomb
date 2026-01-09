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

// Pattern for single line: // @debtbomb(expire=2026-02-10, owner=pricing, ticket=JIRA-123)
// Also supports other comment starters like # or --
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

		// Check for multi-line format start
		// Note: The prompt example for multi-line is actually on a single physical line but with multiple comment markers or separators?
		// Example: // @debtbomb // expire: 2026-02-10 // owner: pricing // ticket: JIRA-123 // reason: Temporary Diwali surge override
		// This looks like it's still one line of text, just structured differently.
		// Let's check if the line contains @debtbomb and then parse key-values from the rest of the line.
		if multiLineStartRegex.MatchString(lineText) {
			// This strictly matches "@debtbomb" at end of line or before newline, which might imply the attributes are on following lines?
			// BUT the example shows: // @debtbomb // expire: 2026-02-10 ...
			// This is actually all on one line in the example provided in the prompt.
			// Let's treat it as a line containing "@debtbomb" followed by attributes.
		}

		// Let's try a more general approach.
		// If line contains "@debtbomb", we try to extract attributes.
		if idx := strings.Index(lineText, "@debtbomb"); idx != -1 {
			// Check if it was already handled by singleLineRegex (parentheses format)
			if singleLineRegex.MatchString(lineText) {
				continue
			}

			// Handle "Multi-line" style which is actually space/comment separated on one line or potentially multiple lines?
			// The prompt says "Multi-line: // @debtbomb // expire: 2026-02-10 ..."
			// This looks like one line to me. Let's parse the text after @debtbomb.
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
	// Format: // expire: 2026-02-10 // owner: pricing ...
	// We can split by //, #, -- etc or just look for key: value patterns.
	// A simple regex might be best here to find "key: value"
	
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