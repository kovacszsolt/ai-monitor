package content

import (
	"errors"
	"testing"
)

func TestPromptErrorCode(t *testing.T) {
	if PromptErrorCode(errors.New("invalid json body")) != "invalid_json" {
		t.Fatal("expected invalid_json")
	}
	if PromptErrorCode(errors.New("messages missing")) != "messages_missing" {
		t.Fatal("expected messages_missing")
	}
}
