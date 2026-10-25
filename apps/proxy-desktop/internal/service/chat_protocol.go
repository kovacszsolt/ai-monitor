package service

import (
	"encoding/json"
	"fmt"
	"strings"
)

func chatHTTPPath(protocol string) string {
	switch protocol {
	case "anthropic_messages":
		return "/v1/messages"
	case "openai_responses":
		return "/v1/responses"
	default:
		return "/v1/chat/completions"
	}
}

func buildChatRequestBody(protocol, publicModel string, messages []ChatMessageDTO, stream bool) ([]byte, error) {
	switch protocol {
	case "anthropic_messages":
		return json.Marshal(buildAnthropicChatBody(publicModel, messages, stream))
	case "openai_responses":
		return json.Marshal(buildResponsesChatBody(publicModel, messages, stream))
	default:
		msgs := make([]map[string]string, 0, len(messages))
		for _, m := range messages {
			msgs = append(msgs, map[string]string{"role": m.Role, "content": m.Content})
		}
		return json.Marshal(map[string]interface{}{
			"model":    publicModel,
			"messages": msgs,
			"stream":   stream,
		})
	}
}

func buildAnthropicChatBody(publicModel string, messages []ChatMessageDTO, stream bool) map[string]interface{} {
	var systemParts []string
	convo := make([]map[string]string, 0, len(messages))
	for _, m := range messages {
		switch m.Role {
		case "system":
			systemParts = append(systemParts, m.Content)
		case "user", "assistant":
			convo = append(convo, map[string]string{"role": m.Role, "content": m.Content})
		}
	}
	body := map[string]interface{}{
		"model":      publicModel,
		"max_tokens": 4096,
		"messages":   convo,
		"stream":     stream,
	}
	if len(systemParts) > 0 {
		body["system"] = strings.Join(systemParts, "\n\n")
	}
	return body
}

func buildResponsesChatBody(publicModel string, messages []ChatMessageDTO, stream bool) map[string]interface{} {
	var instructions string
	inputItems := make([]map[string]string, 0, len(messages))
	for _, m := range messages {
		switch m.Role {
		case "system":
			if instructions == "" {
				instructions = m.Content
			} else {
				instructions += "\n\n" + m.Content
			}
		case "user", "assistant":
			inputItems = append(inputItems, map[string]string{"role": m.Role, "content": m.Content})
		}
	}
	body := map[string]interface{}{
		"model":  publicModel,
		"stream": stream,
	}
	if instructions != "" {
		body["instructions"] = instructions
	}
	switch len(inputItems) {
	case 0:
		body["input"] = ""
	case 1:
		if inputItems[0]["role"] == "user" {
			body["input"] = inputItems[0]["content"]
		} else {
			body["input"] = inputItems
		}
	default:
		body["input"] = inputItems
	}
	return body
}

func extractChatAssistantText(protocol string, body []byte) (string, error) {
	var v map[string]interface{}
	if err := json.Unmarshal(body, &v); err != nil {
		return "", fmt.Errorf("invalid response json: %w", err)
	}
	text := extractAssistantFromValue(protocol, v)
	if text == "" {
		return "", fmt.Errorf("empty response from model")
	}
	return text, nil
}

func extractAssistantFromValue(protocol string, v map[string]interface{}) string {
	switch protocol {
	case "anthropic_messages":
		return extractAnthropicAssistant(v)
	case "openai_responses":
		return extractResponsesAssistant(v)
	default:
		return extractOpenAIChatAssistant(v)
	}
}

func extractOpenAIChatAssistant(v map[string]interface{}) string {
	choices, ok := v["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return ""
	}
	first, ok := choices[0].(map[string]interface{})
	if !ok {
		return ""
	}
	if msg, ok := first["message"].(map[string]interface{}); ok {
		if s, ok := msg["content"].(string); ok {
			return s
		}
	}
	if s, ok := first["text"].(string); ok {
		return s
	}
	return ""
}

func extractAnthropicAssistant(v map[string]interface{}) string {
	content, ok := v["content"].([]interface{})
	if !ok {
		return ""
	}
	var parts []string
	for _, block := range content {
		bm, ok := block.(map[string]interface{})
		if !ok {
			continue
		}
		if bm["type"] == "text" {
			if t, ok := bm["text"].(string); ok {
				parts = append(parts, t)
			}
		}
	}
	return strings.Join(parts, "")
}

func extractResponsesAssistant(v map[string]interface{}) string {
	if s, ok := v["output_text"].(string); ok && s != "" {
		return s
	}
	output, ok := v["output"].([]interface{})
	if !ok {
		return ""
	}
	var parts []string
	for _, item := range output {
		m, ok := item.(map[string]interface{})
		if !ok || m["type"] != "message" {
			continue
		}
		if t := responsesOutputText(m["content"]); t != "" {
			parts = append(parts, t)
		}
	}
	return strings.Join(parts, "\n")
}

func responsesOutputText(content interface{}) string {
	if s, ok := content.(string); ok {
		return s
	}
	blocks, ok := content.([]interface{})
	if !ok {
		return ""
	}
	var parts []string
	for _, bl := range blocks {
		bm, ok := bl.(map[string]interface{})
		if !ok {
			continue
		}
		if t, ok := bm["text"].(string); ok && t != "" {
			parts = append(parts, t)
		}
	}
	return strings.Join(parts, "")
}

func parseProxyHTTPError(status int, body []byte) error {
	if len(body) > 0 {
		var v map[string]interface{}
		if err := json.Unmarshal(body, &v); err == nil {
			if errObj, ok := v["error"].(map[string]interface{}); ok {
				if msg, ok := errObj["message"].(string); ok && msg != "" {
					return fmt.Errorf("%s", msg)
				}
			}
		}
	}
	switch status {
	case 401:
		return fmt.Errorf("invalid API key")
	case 403:
		return fmt.Errorf("API key cannot use this model")
	default:
		return fmt.Errorf("request failed (HTTP %d)", status)
	}
}
