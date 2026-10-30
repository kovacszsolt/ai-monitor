package service

import (
	"encoding/json"
	"testing"
)

func TestBuildImageGenerationBody(t *testing.T) {
	body, err := buildImageGenerationBody("dall-e-3", "a red cube")
	if err != nil {
		t.Fatal(err)
	}
	var v map[string]interface{}
	if err := json.Unmarshal(body, &v); err != nil {
		t.Fatal(err)
	}
	if v["model"] != "dall-e-3" || v["prompt"] != "a red cube" || v["response_format"] != "b64_json" {
		t.Fatalf("unexpected body: %#v", v)
	}
}

func TestExtractImageGenerationResult(t *testing.T) {
	raw := []byte(`{"data":[{"b64_json":"abc","revised_prompt":"a cube"}]}`)
	content, images, err := extractImageGenerationResult(raw)
	if err != nil {
		t.Fatal(err)
	}
	if content != "a cube" {
		t.Fatalf("content=%q", content)
	}
	if len(images) != 1 || images[0].Data != "abc" || images[0].MediaType != "image/png" {
		t.Fatalf("images=%#v", images)
	}
}

func TestLastUserTextPrompt(t *testing.T) {
	msgs := []ChatMessageDTO{
		{Role: "user", Content: "first"},
		{Role: "assistant", Content: "ok"},
		{Role: "user", Content: "  second  "},
	}
	if got := lastUserTextPrompt(msgs); got != "second" {
		t.Fatalf("got %q", got)
	}
}

func TestBuildChatRequestBody_OpenAIWithImage(t *testing.T) {
	msgs := []ChatMessageDTO{
		{
			Role:    "user",
			Content: "What is this?",
			Images: []ChatImageAttachmentDTO{
				{MediaType: "image/png", Data: "abc123"},
			},
		},
	}
	body, err := buildChatRequestBody("openai_chat", "gpt-4o", msgs, false)
	if err != nil {
		t.Fatal(err)
	}
	var v map[string]interface{}
	if err := json.Unmarshal(body, &v); err != nil {
		t.Fatal(err)
	}
	rawMsgs, ok := v["messages"].([]interface{})
	if !ok || len(rawMsgs) != 1 {
		t.Fatalf("messages: %#v", v["messages"])
	}
	first, ok := rawMsgs[0].(map[string]interface{})
	if !ok {
		t.Fatal("first message not object")
	}
	parts, ok := first["content"].([]interface{})
	if !ok || len(parts) != 2 {
		t.Fatalf("content parts: %#v", first["content"])
	}
}

func TestBuildChatRequestBody_AnthropicWithImage(t *testing.T) {
	msgs := []ChatMessageDTO{
		{
			Role:    "user",
			Content: "Describe",
			Images: []ChatImageAttachmentDTO{
				{MediaType: "image/jpeg", Data: "xyz"},
			},
		},
	}
	body, err := buildChatRequestBody("anthropic_messages", "claude-3", msgs, false)
	if err != nil {
		t.Fatal(err)
	}
	var v map[string]interface{}
	if err := json.Unmarshal(body, &v); err != nil {
		t.Fatal(err)
	}
	rawMsgs, ok := v["messages"].([]interface{})
	if !ok || len(rawMsgs) != 1 {
		t.Fatalf("messages: %#v", v["messages"])
	}
	first, ok := rawMsgs[0].(map[string]interface{})
	if !ok {
		t.Fatal("first message not object")
	}
	blocks, ok := first["content"].([]interface{})
	if !ok || len(blocks) != 2 {
		t.Fatalf("content blocks: %#v", first["content"])
	}
}
