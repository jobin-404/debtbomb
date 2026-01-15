package parser

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"io"
	"regexp"
	"strings"
	"time"

	"github.com/jobin-404/debtbomb/internal/model"
)

var singleLineRegex = regexp.MustCompile(`(?:\/\/|#|--|\/\*)\s*@debtbomb\((.*?)\)`)

// Pattern for multi line start: // @debtbomb
var multiLineStartRegex = regexp.MustCompile(`(?:\/\/|#|--|\/\*)\s*@debtbomb\s*$`)

// Regex to find key: value pairs where key is one of our expected fields
var kvRegex = regexp.MustCompile(`(expire|owner|ticket|reason|severity)\s*:\s*([^\/\#\*]+)`)

// Parse scans the content and returns a list of DebtBombs
func Parse(filename string, reader io.Reader) ([]model.DebtBomb, error) {
	var bombs []model.DebtBomb
	var pendingBombs []*model.DebtBomb

	scanner := bufio.NewScanner(reader)
	lineNum := 0

	for scanner.Scan() {
		lineNum++
		lineText := scanner.Text()
		trimmedLine := strings.TrimSpace(lineText)

		if trimmedLine == "" {
			continue
		}

		// Check if this line is a debtbomb
		isBomb := false
		var currentBombs []model.DebtBomb

		// Check single line
		if matches := singleLineRegex.FindStringSubmatch(lineText); len(matches) > 1 {
			bomb, err := parseAttributes(matches[1])
			if err == nil {
				bomb.File = filename
				bomb.Line = lineNum
				bomb.RawText = trimmedLine
				currentBombs = append(currentBombs, bomb)
				isBomb = true

				// Check for inline
				loc := singleLineRegex.FindStringIndex(lineText)
				if loc != nil && loc[0] > 0 {
					// Check if there is code before
					pre := strings.TrimSpace(lineText[:loc[0]])
					if pre != "" {
						// Inline bomb
						for i := range currentBombs {
							currentBombs[i].Snippet = pre
							currentBombs[i].ID = generateID(filename, currentBombs[i].Reason, pre)
							bombs = append(bombs, currentBombs[i])
						}
						// Clear currentBombs so they don't get added to pending
						currentBombs = nil
					}
				}
			}
		} else if idx := strings.Index(lineText, "@debtbomb"); idx != -1 {
			// Check if it was already handled by singleLineRegex (parentheses format)
			if !singleLineRegex.MatchString(lineText) {
				content := lineText[idx+len("@debtbomb"):]
				bomb, err := parseKeyValueStyle(content)
				if err == nil {
					bomb.File = filename
					bomb.Line = lineNum
					bomb.RawText = trimmedLine
					currentBombs = append(currentBombs, bomb)
					isBomb = true

					// Check for inline
					commentStart := regexp.MustCompile(`^\s*(?:\/\/|#|--|\/\*)`)
					if !commentStart.MatchString(lineText) {
						pre := lineText[:idx]
						commentIdx := -1
						// This is naive, but might work for supported languages
						for _, c := range []string{"//", "#", "--", "/*"} {
							if i := strings.LastIndex(pre, c); i != -1 {
								if commentIdx == -1 || i > commentIdx {
									commentIdx = i
								}
							}
						}

						if commentIdx != -1 {
							code := strings.TrimSpace(pre[:commentIdx])
							if code != "" {
								for i := range currentBombs {
									currentBombs[i].Snippet = code
									currentBombs[i].ID = generateID(filename, currentBombs[i].Reason, code)
									bombs = append(bombs, currentBombs[i])
								}
								currentBombs = nil
							}
						}
					}
				}
			}
		}

		if isBomb {
			for i := range currentBombs {
				// Copy to avoid pointer issues with loop var
				b := currentBombs[i]
				pendingBombs = append(pendingBombs, &b)
			}
		} else {
			// Not a bomb line. If we have pending bombs, this line is the snippet.
			if len(pendingBombs) > 0 {
				snippet := trimmedLine
				for _, b := range pendingBombs {
					b.Snippet = snippet
					b.ID = generateID(filename, b.Reason, snippet)
					bombs = append(bombs, *b)
				}
				pendingBombs = nil
			}
		}
	}

	// Flush pending
	if len(pendingBombs) > 0 {
		for _, b := range pendingBombs {
			b.Snippet = "EOF"
			b.ID = generateID(filename, b.Reason, "EOF")
			bombs = append(bombs, *b)
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
		part = strings.TrimSpace(part)
		var kv []string
		if strings.Contains(part, "=") {
			kv = strings.SplitN(part, "=", 2)
		} else if strings.Contains(part, ":") {
			kv = strings.SplitN(part, ":", 2)
		}

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
		case "severity":
			bomb.Severity = val
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
		case "severity":
			bomb.Severity = val
		}

	}

	if !foundExpire {
		return bomb, fmt.Errorf("missing expire date")
	}
	return bomb, nil
}

func generateID(file, reason, snippet string) string {
	h := sha1.New()
	h.Write([]byte(file))
	h.Write([]byte(reason))
	cleanSnippet := strings.TrimSpace(snippet)
	if len(cleanSnippet) > 80 {
		cleanSnippet = cleanSnippet[:80]
	}
	h.Write([]byte(cleanSnippet))
	return fmt.Sprintf("%x", h.Sum(nil))
}
