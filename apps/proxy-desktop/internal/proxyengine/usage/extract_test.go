package usage

import "testing"

func TestFromOpenAIJSON(t *testing.T) {
	body := []byte(`{"usage":{"prompt_tokens":10,"completion_tokens":20}}`)
	u := FromOpenAIJSON(body)
	if u.InputTokens == nil || *u.InputTokens != 10 {
		t.Fatalf("input tokens")
	}
	if u.OutputTokens == nil || *u.OutputTokens != 20 {
		t.Fatalf("output tokens")
	}
}

func TestFromOpenAIJSONResponsesAPI(t *testing.T) {
	body := []byte(`{"usage":{"input_tokens":1338,"input_tokens_details":{"cache_write_tokens":1335,"cached_tokens":0},"output_tokens":393,"output_tokens_details":{"reasoning_tokens":236},"total_tokens":1731}}`)
	u := FromOpenAIJSON(body)
	if u.Status != "ok" {
		t.Fatalf("status %q", u.Status)
	}
	if u.InputTokens == nil || *u.InputTokens != 1338 {
		t.Fatalf("input tokens: %v", u.InputTokens)
	}
	if u.OutputTokens == nil || *u.OutputTokens != 393 {
		t.Fatalf("output tokens: %v", u.OutputTokens)
	}
	if u.CacheWriteTokens == nil || *u.CacheWriteTokens != 1335 {
		t.Fatalf("cache write: %v", u.CacheWriteTokens)
	}
	if u.CachedInputTokens == nil || *u.CachedInputTokens != 0 {
		t.Fatalf("cached read: %v", u.CachedInputTokens)
	}
	if u.ReasoningTokens == nil || *u.ReasoningTokens != 236 {
		t.Fatalf("reasoning: %v", u.ReasoningTokens)
	}
}

func TestFromOpenAIStreamLineCached(t *testing.T) {
	line := `data: {"usage":{"prompt_tokens":100,"completion_tokens":5,"prompt_tokens_details":{"cached_tokens":80}}}`
	u := FromOpenAIStreamLine(line)
	if u == nil || u.CachedInputTokens == nil || *u.CachedInputTokens != 80 {
		t.Fatalf("cached tokens in stream: %v", u)
	}
}
