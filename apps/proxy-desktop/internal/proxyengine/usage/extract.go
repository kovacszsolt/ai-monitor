package usage

import (
	"encoding/json"
	"strings"
)

type Usage struct {
	InputTokens       *int64
	OutputTokens      *int64
	CachedInputTokens *int64
	CacheWriteTokens  *int64
	ReasoningTokens   *int64
	RawJSON           string
	Status            string
}

func FromOpenAIJSON(body []byte) Usage {
	var v map[string]interface{}
	if err := json.Unmarshal(body, &v); err != nil {
		return Usage{Status: "parse_error"}
	}
	u, ok := v["usage"].(map[string]interface{})
	if !ok {
		return Usage{Status: "missing"}
	}
	raw, _ := json.Marshal(u)
	out := Usage{RawJSON: string(raw), Status: "ok"}
	// Chat Completions: prompt_tokens / completion_tokens
	out.InputTokens = intPtr(u["prompt_tokens"])
	out.OutputTokens = intPtr(u["completion_tokens"])
	if details, ok := u["prompt_tokens_details"].(map[string]interface{}); ok {
		out.CachedInputTokens = intPtr(details["cached_tokens"])
	}
	if details, ok := u["completion_tokens_details"].(map[string]interface{}); ok {
		out.ReasoningTokens = intPtr(details["reasoning_tokens"])
	}
	// Responses API (and some newer shapes): input_tokens / output_tokens
	if out.InputTokens == nil {
		out.InputTokens = intPtr(u["input_tokens"])
	}
	if out.OutputTokens == nil {
		out.OutputTokens = intPtr(u["output_tokens"])
	}
	if details, ok := u["input_tokens_details"].(map[string]interface{}); ok {
		if out.CachedInputTokens == nil {
			out.CachedInputTokens = intPtr(details["cached_tokens"])
		}
		if out.CacheWriteTokens == nil {
			out.CacheWriteTokens = intPtr(details["cache_write_tokens"])
		}
	}
	if details, ok := u["output_tokens_details"].(map[string]interface{}); ok && out.ReasoningTokens == nil {
		out.ReasoningTokens = intPtr(details["reasoning_tokens"])
	}
	return out
}

func FromAnthropicJSON(body []byte) Usage {
	var v map[string]interface{}
	if err := json.Unmarshal(body, &v); err != nil {
		return Usage{Status: "parse_error"}
	}
	u, ok := v["usage"].(map[string]interface{})
	if !ok {
		return Usage{Status: "missing"}
	}
	raw, _ := json.Marshal(u)
	out := Usage{RawJSON: string(raw), Status: "ok"}
	out.InputTokens = intPtr(u["input_tokens"])
	out.OutputTokens = intPtr(u["output_tokens"])
	if details, ok := u["cache_creation_input_tokens"]; ok {
		out.CacheWriteTokens = intPtr(details)
	}
	if details, ok := u["cache_read_input_tokens"]; ok {
		out.CachedInputTokens = intPtr(details)
	}
	return out
}

func FromOpenAIStreamLine(line string) *Usage {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "data:") {
		return nil
	}
	payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
	if payload == "[DONE]" {
		return nil
	}
	var v map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &v); err != nil {
		return nil
	}
	u, ok := v["usage"].(map[string]interface{})
	if !ok {
		return nil
	}
	raw, _ := json.Marshal(u)
	out := &Usage{RawJSON: string(raw), Status: "ok"}
	out.InputTokens = intPtr(u["prompt_tokens"])
	out.OutputTokens = intPtr(u["completion_tokens"])
	if details, ok := u["prompt_tokens_details"].(map[string]interface{}); ok {
		out.CachedInputTokens = intPtr(details["cached_tokens"])
	}
	if out.InputTokens == nil {
		out.InputTokens = intPtr(u["input_tokens"])
	}
	if out.OutputTokens == nil {
		out.OutputTokens = intPtr(u["output_tokens"])
	}
	if details, ok := u["input_tokens_details"].(map[string]interface{}); ok {
		if out.CachedInputTokens == nil {
			out.CachedInputTokens = intPtr(details["cached_tokens"])
		}
		if out.CacheWriteTokens == nil {
			out.CacheWriteTokens = intPtr(details["cache_write_tokens"])
		}
	}
	return out
}

func intPtr(v interface{}) *int64 {
	switch t := v.(type) {
	case float64:
		n := int64(t)
		return &n
	case json.Number:
		n, err := t.Int64()
		if err != nil {
			return nil
		}
		return &n
	default:
		return nil
	}
}
