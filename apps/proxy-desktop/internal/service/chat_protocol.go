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
	case "openai_images":
		return "/v1/images/generations"
	default:
		return "/v1/chat/completions"
	}
}

func lastUserTextPrompt(messages []ChatMessageDTO) string {
	for i := len(messages) - 1; i >= 0; i-- {
		if messages[i].Role != "user" {
			continue
		}
		if t := strings.TrimSpace(messages[i].Content); t != "" {
			return t
		}
	}
	return ""
}

func buildImageGenerationBody(publicModel, prompt string) ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"model":           publicModel,
		"prompt":          prompt,
		"n":               1,
		"response_format": "b64_json",
	})
}

func extractImageGenerationResult(body []byte) (string, []ChatImageAttachmentDTO, error) {
	var v map[string]interface{}
	if err := json.Unmarshal(body, &v); err != nil {
		return "", nil, fmt.Errorf("invalid response json: %w", err)
	}
	data, ok := v["data"].([]interface{})
	if !ok || len(data) == 0 {
		return "", nil, fmt.Errorf("empty image generation response")
	}
	images := make([]ChatImageAttachmentDTO, 0, len(data))
	var revisedParts []string
	for _, item := range data {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		if revised, ok := m["revised_prompt"].(string); ok && strings.TrimSpace(revised) != "" {
			revisedParts = append(revisedParts, strings.TrimSpace(revised))
		}
		if b64, ok := m["b64_json"].(string); ok && b64 != "" {
			images = append(images, ChatImageAttachmentDTO{
				MediaType: "image/png",
				Data:      b64,
			})
			continue
		}
		if url, ok := m["url"].(string); ok && strings.TrimSpace(url) != "" {
			return "", nil, fmt.Errorf("provider returned image URL only; request b64_json (response_format)")
		}
	}
	if len(images) == 0 {
		return "", nil, fmt.Errorf("no image data in generation response")
	}
	return strings.Join(revisedParts, "\n"), images, nil
}

func buildChatRequestBody(protocol, publicModel string, messages []ChatMessageDTO, stream bool) ([]byte, error) {
	switch protocol {
	case "anthropic_messages":
		return json.Marshal(buildAnthropicChatBody(publicModel, messages, stream))
	case "openai_responses":
		return json.Marshal(buildResponsesChatBody(publicModel, messages, stream))
	default:
		msgs := make([]map[string]interface{}, 0, len(messages))
		for _, m := range messages {
			msgs = append(msgs, map[string]interface{}{
				"role":    m.Role,
				"content": openAIChatMessageContent(m),
			})
		}
		return json.Marshal(map[string]interface{}{
			"model":    publicModel,
			"messages": msgs,
			"stream":   stream,
		})
	}
}

func openAIChatMessageContent(m ChatMessageDTO) interface{} {
	if m.Role != "user" || len(m.Images) == 0 {
		return m.Content
	}
	parts := make([]interface{}, 0, 1+len(m.Images))
	if strings.TrimSpace(m.Content) != "" {
		parts = append(parts, map[string]interface{}{
			"type": "text",
			"text": m.Content,
		})
	}
	for _, img := range m.Images {
		mt := normalizeImageMediaType(img.MediaType)
		parts = append(parts, map[string]interface{}{
			"type": "image_url",
			"image_url": map[string]string{
				"url": fmt.Sprintf("data:%s;base64,%s", mt, img.Data),
			},
		})
	}
	if len(parts) == 0 {
		return m.Content
	}
	return parts
}

func buildAnthropicChatBody(publicModel string, messages []ChatMessageDTO, stream bool) map[string]interface{} {
	var systemParts []string
	convo := make([]map[string]interface{}, 0, len(messages))
	for _, m := range messages {
		switch m.Role {
		case "system":
			systemParts = append(systemParts, m.Content)
		case "user", "assistant":
			convo = append(convo, map[string]interface{}{
				"role":    m.Role,
				"content": anthropicMessageContent(m),
			})
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

func anthropicMessageContent(m ChatMessageDTO) interface{} {
	if m.Role != "user" || len(m.Images) == 0 {
		return m.Content
	}
	blocks := make([]interface{}, 0, 1+len(m.Images))
	if strings.TrimSpace(m.Content) != "" {
		blocks = append(blocks, map[string]interface{}{
			"type": "text",
			"text": m.Content,
		})
	}
	for _, img := range m.Images {
		blocks = append(blocks, map[string]interface{}{
			"type": "image",
			"source": map[string]string{
				"type":       "base64",
				"media_type": normalizeImageMediaType(img.MediaType),
				"data":       img.Data,
			},
		})
	}
	if len(blocks) == 0 {
		return m.Content
	}
	return blocks
}

func buildResponsesChatBody(publicModel string, messages []ChatMessageDTO, stream bool) map[string]interface{} {
	var instructions string
	inputItems := make([]map[string]interface{}, 0, len(messages))
	for _, m := range messages {
		switch m.Role {
		case "system":
			if instructions == "" {
				instructions = m.Content
			} else {
				instructions += "\n\n" + m.Content
			}
		case "user", "assistant":
			inputItems = append(inputItems, map[string]interface{}{
				"role":    m.Role,
				"content": responsesMessageContent(m),
			})
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

func responsesMessageContent(m ChatMessageDTO) interface{} {
	if m.Role != "user" || len(m.Images) == 0 {
		return m.Content
	}
	parts := make([]interface{}, 0, 1+len(m.Images))
	if strings.TrimSpace(m.Content) != "" {
		parts = append(parts, map[string]interface{}{
			"type": "input_text",
			"text": m.Content,
		})
	}
	for _, img := range m.Images {
		mt := normalizeImageMediaType(img.MediaType)
		parts = append(parts, map[string]interface{}{
			"type":      "input_image",
			"image_url": fmt.Sprintf("data:%s;base64,%s", mt, img.Data),
		})
	}
	if len(parts) == 0 {
		return m.Content
	}
	return parts
}

func normalizeImageMediaType(mt string) string {
	mt = strings.TrimSpace(mt)
	if mt == "" {
		return "image/jpeg"
	}
	return mt
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
