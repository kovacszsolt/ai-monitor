package service

import "testing"

func TestDerivePinnedPromptTitle(t *testing.T) {
	if derivePinnedPromptTitle("  hello world  ") != "hello world" {
		t.Fatal("expected trimmed single line")
	}
	multi := "first line\nsecond line"
	if derivePinnedPromptTitle(multi) != "first line" {
		t.Fatal("expected first line only")
	}
	long := ""
	for i := 0; i < 100; i++ {
		long += "x"
	}
	title := derivePinnedPromptTitle(long)
	if len([]rune(title)) > maxChatPinnedTitleLen+1 {
		t.Fatalf("title too long: %d", len([]rune(title)))
	}
	if !stringsHasSuffix(title, "…") {
		t.Fatal("expected ellipsis on long title")
	}
}

func stringsHasSuffix(s, suffix string) bool {
	return len(s) >= len(suffix) && s[len(s)-len(suffix):] == suffix
}
