package repository

import (
	"database/sql"
	"testing"
	"time"
)

func TestPromptWindowMS(t *testing.T) {
	start := time.Date(2026, 1, 1, 10, 0, 0, 0, time.UTC).Format(time.RFC3339Nano)
	end := time.Date(2026, 1, 1, 10, 0, 5, 500_000_000, time.UTC).Format(time.RFC3339Nano)

	ms, ok := promptWindowMS(sql.NullString{String: start, Valid: true}, sql.NullString{String: end, Valid: true})
	if !ok {
		t.Fatal("expected ok")
	}
	if ms != 5500 {
		t.Fatalf("got %d ms", ms)
	}
}

func TestPromptWindowMS_Invalid(t *testing.T) {
	_, ok := promptWindowMS(sql.NullString{}, sql.NullString{String: "x", Valid: true})
	if ok {
		t.Fatal("expected not ok")
	}
}
