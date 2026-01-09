package parser

import (
	"strings"
	"testing"
	"time"
)

func TestParseSingleLine(t *testing.T) {
	content := `
	// @debtbomb(expire=2023-01-01, owner=test)
	code()
	`
	bombs, err := Parse("test.go", strings.NewReader(content))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(bombs) != 1 {
		t.Fatalf("Expected 1 bomb, got %d", len(bombs))
	}

	b := bombs[0]
	expectedDate, _ := time.Parse("2006-01-02", "2023-01-01")
	if !b.Expire.Equal(expectedDate) {
		t.Errorf("Expected date %v, got %v", expectedDate, b.Expire)
	}
	if b.Owner != "test" {
		t.Errorf("Expected owner test, got %s", b.Owner)
	}
}

func TestParseMultiLine(t *testing.T) {
	content := `
	// @debtbomb // expire: 2023-01-01 // owner: test
	`
	bombs, err := Parse("test.go", strings.NewReader(content))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(bombs) != 1 {
		t.Fatalf("Expected 1 bomb, got %d", len(bombs))
	}

	b := bombs[0]
	expectedDate, _ := time.Parse("2006-01-02", "2023-01-01")
	if !b.Expire.Equal(expectedDate) {
		t.Errorf("Expected date %v, got %v", expectedDate, b.Expire)
	}
	if b.Owner != "test" {
		t.Errorf("Expected owner test, got %s", b.Owner)
	}
}