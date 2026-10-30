package content

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/routing"
)

// PromptData is extracted request-side content for persistence.
type PromptData struct {
	LastUserText    string
	MessagesJSON    string
	PromptTruncated bool
}

// ExtractRequest parses the client request body (before model rewrite).
func ExtractRequest(protocol routing.Protocol, body []byte, maxFieldBytes int64) (PromptData, error) {
	var loose map[string]interface{}
	if err := json.Unmarshal(body, &loose); err != nil {
		return PromptData{}, fmt.Errorf("invalid json body")
	}
	if protocol == routing.ProtocolOpenAIImages {
		return extractOpenAIImagesRequest(loose, maxFieldBytes)
	}
	msgs := loose["messages"]
	if msgs == nil && protocol == routing.ProtocolOpenAIResponses {
		return extractOpenAIResponsesRequest(loose, maxFieldBytes)
	}
	if msgs == nil {
		return PromptData{}, fmt.Errorf("messages missing")
	}
	return promptFromMessages(protocol, msgs, maxFieldBytes)
}

func extractOpenAIImagesRequest(loose map[string]interface{}, maxFieldBytes int64) (PromptData, error) {
	prompt, _ := loose["prompt"].(string)
	if strings.TrimSpace(prompt) == "" {
		return PromptData{}, fmt.Errorf("prompt missing")
	}
	normalized := []map[string]string{{"role": "user", "content": prompt}}
	raw, err := json.Marshal(normalized)
	if err != nil {
		return PromptData{}, err
	}
	messagesJSON, truncJSON := TruncateUTF8(string(raw), maxFieldBytes)
	lastUser, truncUser := TruncateUTF8(prompt, maxFieldBytes)
	return PromptData{
		LastUserText:    lastUser,
		MessagesJSON:    messagesJSON,
		PromptTruncated: truncJSON || truncUser,
	}, nil
}

func extractOpenAIResponsesRequest(loose map[string]interface{}, maxFieldBytes int64) (PromptData, error) {
	input := loose["input"]
	if input == nil {
		return PromptData{}, fmt.Errorf("input missing")
	}
	normalized := normalizeResponsesInput(loose, input)
	raw, err := json.Marshal(normalized)
	if err != nil {
		return PromptData{}, err
	}
	messagesJSON, truncJSON := TruncateUTF8(string(raw), maxFieldBytes)
	lastUser, truncUser := TruncateUTF8(responsesLastUserText(input), maxFieldBytes)
	return PromptData{
		LastUserText:    lastUser,
		MessagesJSON:    messagesJSON,
		PromptTruncated: truncJSON || truncUser,
	}, nil
}

func normalizeResponsesInput(loose map[string]interface{}, input interface{}) []interface{} {
	out := make([]interface{}, 0, 4)
	if instr, ok := loose["instructions"].(string); ok && instr != "" {
		out = append(out, map[string]interface{}{
			"role":    "system",
			"content": instr,
		})
	}
	switch v := input.(type) {
	case string:
		out = append(out, map[string]interface{}{
			"role":    "user",
			"content": v,
		})
	case []interface{}:
		out = append(out, v...)
	}
	return out
}

func responsesLastUserText(input interface{}) string {
	if s, ok := input.(string); ok {
		return s
	}
	return lastOpenAIUserText(input)
}

func promptFromMessages(protocol routing.Protocol, msgs interface{}, maxFieldBytes int64) (PromptData, error) {
	raw, err := json.Marshal(msgs)
	if err != nil {
		return PromptData{}, err
	}
	messagesJSON, truncJSON := TruncateUTF8(string(raw), maxFieldBytes)

	lastUser := ""
	switch protocol {
	case routing.ProtocolAnthropicMessages:
		lastUser = lastAnthropicUserText(msgs)
	default:
		lastUser = lastOpenAIUserText(msgs)
	}
	lastUser, truncUser := TruncateUTF8(lastUser, maxFieldBytes)

	return PromptData{
		LastUserText:    lastUser,
		MessagesJSON:    messagesJSON,
		PromptTruncated: truncJSON || truncUser,
	}, nil
}

func lastOpenAIUserText(msgs interface{}) string {
	arr, ok := msgs.([]interface{})
	if !ok {
		return ""
	}
	for i := len(arr) - 1; i >= 0; i-- {
		m, ok := arr[i].(map[string]interface{})
		if !ok {
			continue
		}
		role, _ := m["role"].(string)
		if role != "user" {
			continue
		}
		return openAIContentString(m["content"])
	}
	return ""
}

func openAIContentString(content interface{}) string {
	if content == nil {
		return ""
	}
	if s, ok := content.(string); ok {
		return s
	}
	parts, ok := content.([]interface{})
	if !ok {
		return ""
	}
	var b []byte
	first := true
	for _, p := range parts {
		pm, ok := p.(map[string]interface{})
		if !ok {
			continue
		}
		t, _ := pm["text"].(string)
		if t == "" {
			continue
		}
		if !first {
			b = append(b, '\n')
		}
		b = append(b, t...)
		first = false
	}
	return string(b)
}

func lastAnthropicUserText(msgs interface{}) string {
	arr, ok := msgs.([]interface{})
	if !ok {
		return ""
	}
	for i := len(arr) - 1; i >= 0; i-- {
		m, ok := arr[i].(map[string]interface{})
		if !ok {
			continue
		}
		role, _ := m["role"].(string)
		if role != "user" {
			continue
		}
		return anthropicContentString(m["content"])
	}
	return ""
}

func anthropicContentString(content interface{}) string {
	if content == nil {
		return ""
	}
	if s, ok := content.(string); ok {
		return s
	}
	blocks, ok := content.([]interface{})
	if !ok {
		return ""
	}
	var b []byte
	first := true
	for _, bl := range blocks {
		bm, ok := bl.(map[string]interface{})
		if !ok {
			continue
		}
		t, _ := bm["text"].(string)
		if t == "" {
			continue
		}
		if !first {
			b = append(b, '\n')
		}
		b = append(b, t...)
		first = false
	}
	return string(b)
}
