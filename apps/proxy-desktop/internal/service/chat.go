package service

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/pkg/routing"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

const chatRequestTag = "desktop-chat"
const chatStreamEvent = "chat:stream-delta"

type ChatMessageDTO struct {
	Role    string                    `json:"role"`
	Content string                    `json:"content"`
	Images  []ChatImageAttachmentDTO  `json:"images,omitempty"`
}

type ChatImageAttachmentDTO struct {
	MediaType string `json:"media_type"`
	Data      string `json:"data"`
}

type ChatCompletionInputDTO struct {
	APIKeyID       string           `json:"api_key_id"`
	ModelRouteID   string           `json:"model_route_id"`
	Messages       []ChatMessageDTO `json:"messages"`
	Stream         bool             `json:"stream"`
	StreamClientID string           `json:"stream_client_id"`
}

type ChatCompletionResultDTO struct {
	Content string                   `json:"content"`
	Images  []ChatImageAttachmentDTO `json:"images,omitempty"`
}

type ChatStreamDeltaDTO struct {
	ClientID string `json:"client_id"`
	Text     string `json:"text"`
}

func (s *DesktopService) SetWailsContext(ctx context.Context) {
	s.wailsCtx = ctx
}

func chatMessagesHaveImages(messages []ChatMessageDTO) bool {
	for _, m := range messages {
		if len(m.Images) > 0 {
			return true
		}
	}
	return false
}

func (s *DesktopService) SendChatCompletion(in ChatCompletionInputDTO) Result[ChatCompletionResultDTO] {
	if s.deps.Runtime == nil || !s.deps.Runtime.Status().Running {
		return fail[ChatCompletionResultDTO](fmt.Errorf("proxy is not running"))
	}
	listen := strings.TrimSpace(s.deps.Runtime.Status().ListenAddr)
	if listen == "" {
		return fail[ChatCompletionResultDTO](fmt.Errorf("proxy listen address is not available"))
	}

	ctx := s.ctx()
	key, err := s.deps.Keys.Get(ctx, in.APIKeyID)
	if err != nil {
		return fail[ChatCompletionResultDTO](fmt.Errorf("api key not found"))
	}
	if !key.Enabled {
		return fail[ChatCompletionResultDTO](fmt.Errorf("api key is disabled"))
	}
	token := strings.TrimSpace(key.Secret)
	if token == "" {
		return fail[ChatCompletionResultDTO](fmt.Errorf("api key has no credential"))
	}

	route, err := s.deps.Routes.Get(ctx, in.ModelRouteID)
	if err != nil {
		return fail[ChatCompletionResultDTO](fmt.Errorf("model route not found"))
	}
	if !route.Enabled {
		return fail[ChatCompletionResultDTO](fmt.Errorf("model route is disabled"))
	}
	if chatMessagesHaveImages(in.Messages) && route.Modality != routing.ModalityTextImageText {
		return fail[ChatCompletionResultDTO](fmt.Errorf(
			"model modality %q does not accept image input (need text_image-text)",
			route.Modality,
		))
	}

	if route.Modality == routing.ModalityTextImage {
		return s.sendImageGeneration(ctx, listen, token, route.PublicModelName, in.Messages)
	}

	protocol := string(route.Protocol)
	body, err := buildChatRequestBody(protocol, route.PublicModelName, in.Messages, in.Stream)
	if err != nil {
		return fail[ChatCompletionResultDTO](err)
	}
	maxBody := s.deps.Runtime.MaxRequestBytes()
	if maxBody < 1 {
		maxBody = 10 * 1024 * 1024
	}
	if int64(len(body)) > maxBody {
		return fail[ChatCompletionResultDTO](fmt.Errorf(
			"chat request is %d bytes but server.max_request_bytes is %d — increase Runtime Vars or use a smaller image",
			len(body), maxBody,
		))
	}

	url := "http://" + listen + chatHTTPPath(protocol)
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fail[ChatCompletionResultDTO](err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-AI-Proxy-Tag", chatRequestTag)

	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return fail[ChatCompletionResultDTO](fmt.Errorf("proxy request failed: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errBody, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
		return fail[ChatCompletionResultDTO](parseProxyHTTPError(resp.StatusCode, errBody))
	}

	if in.Stream {
		text, err := readChatStream(ctx, s.wailsCtx, protocol, resp.Body, in.StreamClientID)
		if err != nil {
			return fail[ChatCompletionResultDTO](err)
		}
		return ok(ChatCompletionResultDTO{Content: text})
	}

	respBody, err := io.ReadAll(io.LimitReader(resp.Body, 8<<20))
	if err != nil {
		return fail[ChatCompletionResultDTO](fmt.Errorf("read response: %w", err))
	}
	text, err := extractChatAssistantText(protocol, respBody)
	if err != nil {
		return fail[ChatCompletionResultDTO](err)
	}
	return ok(ChatCompletionResultDTO{Content: text})
}

func (s *DesktopService) sendImageGeneration(
	ctx context.Context,
	listen, token, publicModel string,
	messages []ChatMessageDTO,
) Result[ChatCompletionResultDTO] {
	prompt := lastUserTextPrompt(messages)
	if prompt == "" {
		return fail[ChatCompletionResultDTO](fmt.Errorf("image generation requires a text prompt"))
	}
	body, err := buildImageGenerationBody(publicModel, prompt)
	if err != nil {
		return fail[ChatCompletionResultDTO](err)
	}

	url := "http://" + listen + "/v1/images/generations"
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewReader(body))
	if err != nil {
		return fail[ChatCompletionResultDTO](err)
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-AI-Proxy-Tag", chatRequestTag)

	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Do(req)
	if err != nil {
		return fail[ChatCompletionResultDTO](fmt.Errorf("proxy request failed: %w", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		errBody, _ := io.ReadAll(io.LimitReader(resp.Body, 64*1024))
		return fail[ChatCompletionResultDTO](parseProxyHTTPError(resp.StatusCode, errBody))
	}

	// Image payloads (base64) can be large; allow up to max_request_bytes or 32MiB.
	maxResp := s.deps.Runtime.MaxRequestBytes()
	if maxResp < 1 {
		maxResp = 32 << 20
	}
	if maxResp < 32<<20 {
		maxResp = 32 << 20
	}
	respBody, err := io.ReadAll(io.LimitReader(resp.Body, maxResp))
	if err != nil {
		return fail[ChatCompletionResultDTO](fmt.Errorf("read response: %w", err))
	}
	content, images, err := extractImageGenerationResult(respBody)
	if err != nil {
		return fail[ChatCompletionResultDTO](err)
	}
	return ok(ChatCompletionResultDTO{Content: content, Images: images})
}

func readChatStream(ctx, wailsCtx context.Context, protocol string, r io.Reader, clientID string) (string, error) {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 64*1024), 1024*1024)
	var accumulated string
	emit := func(text string) {
		if wailsCtx != nil && clientID != "" {
			runtime.EventsEmit(wailsCtx, chatStreamEvent, ChatStreamDeltaDTO{
				ClientID: clientID,
				Text:     text,
			})
		}
	}
	for scanner.Scan() {
		if ctx.Err() != nil {
			return "", ctx.Err()
		}
		line := scanner.Text()
		delta := parseStreamLine(protocol, line)
		if delta == "" {
			continue
		}
		accumulated += delta
		emit(accumulated)
	}
	if err := scanner.Err(); err != nil {
		return "", fmt.Errorf("read stream: %w", err)
	}
	if accumulated == "" {
		return "", fmt.Errorf("empty streaming response from model")
	}
	return accumulated, nil
}

func parseStreamLine(protocol, line string) string {
	line = strings.TrimSpace(line)
	if line == "" || strings.HasPrefix(line, ":") {
		return ""
	}
	if protocol == "anthropic_messages" {
		return parseAnthropicStreamLine(line)
	}
	return parseOpenAIStyleStreamLine(protocol, line)
}

func parseOpenAIStyleStreamLine(protocol, line string) string {
	if !strings.HasPrefix(line, "data:") {
		return ""
	}
	payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
	if payload == "[DONE]" {
		return ""
	}
	var v map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &v); err != nil {
		return ""
	}
	if protocol == "openai_responses" {
		if t, ok := v["delta"].(string); ok {
			return t
		}
		if v["type"] == "response.output_text.delta" {
			if t, ok := v["delta"].(string); ok {
				return t
			}
		}
		return ""
	}
	choices, ok := v["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return ""
	}
	first, ok := choices[0].(map[string]interface{})
	if !ok {
		return ""
	}
	delta, ok := first["delta"].(map[string]interface{})
	if !ok {
		return ""
	}
	if s, ok := delta["content"].(string); ok {
		return s
	}
	return ""
}

func parseAnthropicStreamLine(line string) string {
	if !strings.HasPrefix(line, "data:") {
		return ""
	}
	payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
	var v map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &v); err != nil {
		return ""
	}
	if v["type"] != "content_block_delta" {
		return ""
	}
	delta, ok := v["delta"].(map[string]interface{})
	if !ok || delta["type"] != "text_delta" {
		return ""
	}
	if t, ok := delta["text"].(string); ok {
		return t
	}
	return ""
}
