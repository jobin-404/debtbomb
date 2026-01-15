package parser

import (
	"strings"
	"testing"
	"time"
)

func TestParseSingleLine(t *testing.T) {
	content := `
	// @debtbomb(expire=2026-01-14, owner:jobin)
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
	expectedDate, _ := time.Parse("2006-01-02", "2026-01-14")
	if !b.Expire.Equal(expectedDate) {
		t.Errorf("Expected date %v, got %v", expectedDate, b.Expire)
	}
	if b.Owner != "test" {
		t.Errorf("Expected owner test, got %s", b.Owner)
	}
}

func TestParseMultiLine(t *testing.T) {
	content := `
	// @debtbomb // expire: 2026-01-12 // owner: test
	`
	bombs, err := Parse("test.go", strings.NewReader(content))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(bombs) != 1 {
		t.Fatalf("Expected 1 bomb, got %d", len(bombs))
	}

	b := bombs[0]
	expectedDate, _ := time.Parse("2006-01-02", "2026-01-12")
	if !b.Expire.Equal(expectedDate) {
		t.Errorf("Expected date %v, got %v", expectedDate, b.Expire)
	}
	if b.Owner != "test" {
		t.Errorf("Expected owner test, got %s", b.Owner)
	}
}

func TestParseMixedSyntax(t *testing.T) {
	content := `
	// @debtbomb(expire=2026-01-15, owner:test-mixed)
	`
	bombs, err := Parse("test.go", strings.NewReader(content))
	if err != nil {
		t.Fatalf("Parse failed: %v", err)
	}

	if len(bombs) != 1 {
		t.Fatalf("Expected 1 bomb, got %d", len(bombs))
	}

	b := bombs[0]
	expectedDate, _ := time.Parse("2006-01-02", "2026-01-15")
	if !b.Expire.Equal(expectedDate) {
		t.Errorf("Expected date %v, got %v", expectedDate, b.Expire)
	}
	if b.Owner != "test-mixed" {
		t.Errorf("Expected owner test-mixed, got %s", b.Owner)
	}
}
