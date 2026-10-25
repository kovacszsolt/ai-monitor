package content

import (
	"strings"
	"testing"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/routing"
)

func TestExtractRequestOpenAIResponses(t *testing.T) {
	body := []byte(`{"model":"m","instructions":"sys","input":[{"role":"user","content":"Group files"}]}`)
	p, err := ExtractRequest(routing.ProtocolOpenAIResponses, body, 1<<20)
	if err != nil {
		t.Fatal(err)
	}
	if p.LastUserText != "Group files" {
		t.Fatalf("last user: %q", p.LastUserText)
	}
	if !contains(p.MessagesJSON, "sys") || !contains(p.MessagesJSON, "Group files") {
		t.Fatalf("messages json: %s", p.MessagesJSON)
	}
}

func TestExtractRequestOpenAIResponsesStringInput(t *testing.T) {
	body := []byte(`{"model":"m","input":"hello"}`)
	p, err := ExtractRequest(routing.ProtocolOpenAIResponses, body, 1<<20)
	if err != nil {
		t.Fatal(err)
	}
	if p.LastUserText != "hello" {
		t.Fatalf("last user: %q", p.LastUserText)
	}
}

func TestExtractResponseOpenAIResponses(t *testing.T) {
	body := []byte(`{"output":[{"type":"reasoning","content":[]},{"type":"message","content":[{"type":"output_text","text":"[{ \"files\": [0] }]"}]}]}`)
	r := ExtractResponse(routing.ProtocolOpenAIResponses, body, 1<<20)
	if r.AssistantText != `[{ "files": [0] }]` {
		t.Fatalf("got %q", r.AssistantText)
	}
}

func TestExtractRequestOpenAI(t *testing.T) {
	body := []byte(`{"model":"m","messages":[{"role":"system","content":"sys"},{"role":"user","content":"Hello!"}]}`)
	p, err := ExtractRequest(routing.ProtocolOpenAIChat, body, 1<<20)
	if err != nil {
		t.Fatal(err)
	}
	if p.LastUserText != "Hello!" {
		t.Fatalf("last user: %q", p.LastUserText)
	}
	if !contains(p.MessagesJSON, "system") {
		t.Fatalf("messages json: %s", p.MessagesJSON)
	}
}

func TestExtractRequestMultipartUser(t *testing.T) {
	body := []byte(`{"model":"m","messages":[{"role":"user","content":[{"type":"text","text":"part1"},{"type":"text","text":"part2"}]}]}`)
	p, err := ExtractRequest(routing.ProtocolOpenAIChat, body, 1<<20)
	if err != nil {
		t.Fatal(err)
	}
	if p.LastUserText != "part1\npart2" {
		t.Fatalf("got %q", p.LastUserText)
	}
}

func TestExtractResponseOpenAI(t *testing.T) {
	body := []byte(`{"choices":[{"message":{"content":"hi there"}}]}`)
	r := ExtractResponse(routing.ProtocolOpenAIChat, body, 1<<20)
	if r.AssistantText != "hi there" {
		t.Fatalf("got %q", r.AssistantText)
	}
}

func TestExtractResponseOpenAIToolCalls(t *testing.T) {
	body := []byte(`{"choices":[{"message":{"tool_calls":[{"id":"c1","type":"function","function":{"name":"file_read","arguments":"{}"}}]}}]}`)
	r := ExtractResponse(routing.ProtocolOpenAIChat, body, 1<<20)
	if r.AssistantText != "" {
		t.Fatalf("unexpected text %q", r.AssistantText)
	}
	if !strings.Contains(r.ToolCallsJSON, "file_read") {
		t.Fatalf("tool_calls_json: %s", r.ToolCallsJSON)
	}
}

func TestStreamAccumulator(t *testing.T) {
	acc := NewStreamAccumulator(1<<20)
	acc.Feed(`data: {"choices":[{"delta":{"content":"hel"}}]}`)
	acc.Feed(`data: {"choices":[{"delta":{"content":"lo"}}]}`)
	r := acc.Result()
	if r.AssistantText != "hello" {
		t.Fatalf("got %q", r.AssistantText)
	}
}

func TestTruncateUTF8(t *testing.T) {
	s, trunc := TruncateUTF8("hello", 3)
	if !trunc || s != "hel" {
		t.Fatalf("got %q trunc=%v", s, trunc)
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 || indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
