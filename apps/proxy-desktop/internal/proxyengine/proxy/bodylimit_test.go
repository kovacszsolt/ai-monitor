package proxy

import (
	"bytes"
	"errors"
	"strings"
	"testing"
)

func TestReadLimitedWithinMax(t *testing.T) {
	data, err := readLimited(bytes.NewReader([]byte(`{"model":"m"}`)), 1024)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"model":"m"}` {
		t.Fatalf("got %q", data)
	}
}

func TestReadLimitedRejectsOversized(t *testing.T) {
	_, err := readLimited(bytes.NewReader([]byte("0123456789")), 5)
	if err == nil {
		t.Fatal("expected error")
	}
	if !errors.Is(err, errRequestBodyTooLarge) {
		t.Fatalf("expected errRequestBodyTooLarge, got %v", err)
	}
	if !strings.Contains(err.Error(), "5") {
		t.Fatalf("expected max bytes in message: %v", err)
	}
}
