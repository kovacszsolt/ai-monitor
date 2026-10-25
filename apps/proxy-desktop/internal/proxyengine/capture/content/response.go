package content

import (
	"encoding/json"
	"sort"
	"strings"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/routing"
)

// ResponseData is extracted assistant-side content for persistence.
type ResponseData struct {
	AssistantText      string
	ResponseTruncated  bool
	ToolCallsJSON      string
	ToolCallsTruncated bool
}

// ExtractResponse parses a non-streaming response body.
func ExtractResponse(protocol routing.Protocol, body []byte, maxBytes int64) ResponseData {
	if len(body) == 0 {
		return ResponseData{}
	}
	switch protocol {
	case routing.ProtocolAnthropicMessages:
		return extractAnthropicResponse(body, maxBytes)
	case routing.ProtocolOpenAIResponses:
		return extractOpenAIResponsesResponse(body, maxBytes)
	default:
		return extractOpenAIResponse(body, maxBytes)
	}
}

func extractOpenAIResponsesResponse(body []byte, maxBytes int64) ResponseData {
	var v map[string]interface{}
	if err := json.Unmarshal(body, &v); err != nil {
		return ResponseData{}
	}
	output, ok := v["output"].([]interface{})
	if !ok || len(output) == 0 {
		return ResponseData{}
	}
	var textParts []string
	for _, item := range output {
		m, ok := item.(map[string]interface{})
		if !ok {
			continue
		}
		typ, _ := m["type"].(string)
		if typ != "message" {
			continue
		}
		t := responsesOutputMessageText(m["content"])
		if t != "" {
			textParts = append(textParts, t)
		}
	}
	text := strings.Join(textParts, "\n")
	text, textTrunc := TruncateUTF8(text, maxBytes)
	return ResponseData{
		AssistantText:     text,
		ResponseTruncated: textTrunc,
	}
}

func responsesOutputMessageText(content interface{}) string {
	blocks, ok := content.([]interface{})
	if !ok {
		return openAIContentString(content)
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

func extractOpenAIResponse(body []byte, maxBytes int64) ResponseData {
	var v map[string]interface{}
	if err := json.Unmarshal(body, &v); err != nil {
		return ResponseData{}
	}
	choices, ok := v["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		return ResponseData{}
	}
	ch, ok := choices[0].(map[string]interface{})
	if !ok {
		return ResponseData{}
	}
	msg, ok := ch["message"].(map[string]interface{})
	if !ok {
		return ResponseData{}
	}
	text := openAIContentString(msg["content"])
	text, textTrunc := TruncateUTF8(text, maxBytes)
	toolJSON, toolTrunc := marshalToolCallsField(msg["tool_calls"], maxBytes)
	return ResponseData{
		AssistantText:      text,
		ResponseTruncated:  textTrunc,
		ToolCallsJSON:      toolJSON,
		ToolCallsTruncated: toolTrunc,
	}
}

func extractAnthropicResponse(body []byte, maxBytes int64) ResponseData {
	var v map[string]interface{}
	if err := json.Unmarshal(body, &v); err != nil {
		return ResponseData{}
	}
	text := anthropicContentString(v["content"])
	text, textTrunc := TruncateUTF8(text, maxBytes)
	tools := anthropicToolUseBlocks(v["content"])
	toolJSON, toolTrunc := marshalToolCallsValue(tools, maxBytes)
	return ResponseData{
		AssistantText:      text,
		ResponseTruncated:  textTrunc,
		ToolCallsJSON:      toolJSON,
		ToolCallsTruncated: toolTrunc,
	}
}

func anthropicToolUseBlocks(content interface{}) []interface{} {
	blocks, ok := content.([]interface{})
	if !ok {
		return nil
	}
	var out []interface{}
	for _, bl := range blocks {
		m, ok := bl.(map[string]interface{})
		if !ok {
			continue
		}
		if typ, _ := m["type"].(string); typ == "tool_use" {
			out = append(out, m)
		}
	}
	return out
}

func marshalToolCallsField(v interface{}, maxBytes int64) (string, bool) {
	if v == nil {
		return "", false
	}
	arr, ok := v.([]interface{})
	if !ok || len(arr) == 0 {
		return "", false
	}
	return marshalToolCallsValue(arr, maxBytes)
}

func marshalToolCallsValue(v interface{}, maxBytes int64) (string, bool) {
	if v == nil {
		return "", false
	}
	raw, err := json.Marshal(v)
	if err != nil {
		return "", false
	}
	s, trunc := TruncateUTF8(string(raw), maxBytes)
	if s == "" {
		return "", false
	}
	return s, trunc
}

func anthropicResponseText(body []byte) string {
	var v map[string]interface{}
	if err := json.Unmarshal(body, &v); err != nil {
		return ""
	}
	return anthropicContentString(v["content"])
}

// StreamAccumulator collects assistant text and tool call deltas from SSE lines.
type StreamAccumulator struct {
	maxBytes     int64
	buf          strings.Builder
	truncated    bool
	rawBuf       strings.Builder
	rawTruncated bool
	tools        *streamToolCalls
}

func NewStreamAccumulator(maxBytes int64) *StreamAccumulator {
	return &StreamAccumulator{
		maxBytes: maxBytes,
		tools:    newStreamToolCalls(maxBytes),
	}
}

func (a *StreamAccumulator) Feed(line string) {
	a.appendRawLine(line)
	if a.truncated && a.maxBytes <= 0 {
		return
	}
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "data:") {
		return
	}
	payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
	if payload == "" || payload == "[DONE]" {
		return
	}
	var v map[string]interface{}
	if err := json.Unmarshal([]byte(payload), &v); err != nil {
		return
	}
	if choices, ok := v["choices"].([]interface{}); ok && len(choices) > 0 {
		ch, _ := choices[0].(map[string]interface{})
		if ch != nil {
			delta, _ := ch["delta"].(map[string]interface{})
			if delta != nil {
				if tcs, ok := delta["tool_calls"].([]interface{}); ok {
					a.tools.feed(tcs)
				}
				if s, ok := delta["content"].(string); ok && s != "" {
					a.append(s)
				}
			}
		}
	}
	if typ, _ := v["type"].(string); typ == "content_block_delta" {
		d, _ := v["delta"].(map[string]interface{})
		if d != nil {
			if s, ok := d["text"].(string); ok && s != "" {
				a.append(s)
			}
		}
	}
	if typ, _ := v["type"].(string); typ == "content_block_start" {
		block, _ := v["content_block"].(map[string]interface{})
		if block != nil {
			if bt, _ := block["type"].(string); bt == "tool_use" {
				a.tools.feedAnthropicBlock(block)
			}
		}
	}
}

func (a *StreamAccumulator) append(s string) {
	if s == "" {
		return
	}
	if a.maxBytes <= 0 {
		a.truncated = true
		return
	}
	remaining := a.maxBytes - int64(a.buf.Len())
	if remaining <= 0 {
		a.truncated = true
		return
	}
	part, trunc := TruncateUTF8(s, remaining)
	a.buf.WriteString(part)
	if trunc || int64(len(part)) < int64(len(s)) {
		a.truncated = true
	}
}

func (a *StreamAccumulator) Result() ResponseData {
	toolJSON, toolTrunc := a.tools.json()
	return ResponseData{
		AssistantText:      a.buf.String(),
		ResponseTruncated:  a.truncated,
		ToolCallsJSON:      toolJSON,
		ToolCallsTruncated: toolTrunc,
	}
}

// RawBody returns a truncated copy of streamed upstream lines for debugging.
func (a *StreamAccumulator) RawBody() (string, bool) {
	return a.rawBuf.String(), a.rawTruncated
}

func (a *StreamAccumulator) appendRawLine(line string) {
	if a.maxBytes <= 0 {
		a.rawTruncated = true
		return
	}
	if line == "" {
		return
	}
	chunk := line
	if !strings.HasSuffix(chunk, "\n") {
		chunk += "\n"
	}
	remaining := a.maxBytes - int64(a.rawBuf.Len())
	if remaining <= 0 {
		a.rawTruncated = true
		return
	}
	part, trunc := TruncateUTF8(chunk, remaining)
	a.rawBuf.WriteString(part)
	if trunc || int64(len(part)) < int64(len(chunk)) {
		a.rawTruncated = true
	}
}

type streamToolCall struct {
	ID        string
	Type      string
	Index     int
	Name      string
	Arguments string
}

type streamToolCalls struct {
	maxBytes   int64
	byIndex    map[int]*streamToolCall
	anthropic  []map[string]interface{}
	truncated  bool
}

func newStreamToolCalls(maxBytes int64) *streamToolCalls {
	return &streamToolCalls{maxBytes: maxBytes, byIndex: make(map[int]*streamToolCall)}
}

func (s *streamToolCalls) feed(chunks []interface{}) {
	for _, c := range chunks {
		m, ok := c.(map[string]interface{})
		if !ok {
			continue
		}
		idx := 0
		if v, ok := m["index"].(float64); ok {
			idx = int(v)
		}
		tc := s.byIndex[idx]
		if tc == nil {
			tc = &streamToolCall{Index: idx}
			s.byIndex[idx] = tc
		}
		if id, ok := m["id"].(string); ok && id != "" {
			tc.ID = id
		}
		if typ, ok := m["type"].(string); ok && typ != "" {
			tc.Type = typ
		}
		if fn, ok := m["function"].(map[string]interface{}); ok {
			if name, ok := fn["name"].(string); ok && name != "" {
				tc.Name = name
			}
			if args, ok := fn["arguments"].(string); ok {
				tc.Arguments += args
			}
		}
	}
}

func (s *streamToolCalls) feedAnthropicBlock(block map[string]interface{}) {
	s.anthropic = append(s.anthropic, block)
}

func (s *streamToolCalls) json() (string, bool) {
	if len(s.byIndex) > 0 {
		idxs := make([]int, 0, len(s.byIndex))
		for i := range s.byIndex {
			idxs = append(idxs, i)
		}
		sort.Ints(idxs)
		arr := make([]map[string]interface{}, 0, len(idxs))
		for _, i := range idxs {
			tc := s.byIndex[i]
			item := map[string]interface{}{
				"id":    tc.ID,
				"type":  tc.Type,
				"index": tc.Index,
				"function": map[string]interface{}{
					"name":      tc.Name,
					"arguments": tc.Arguments,
				},
			}
			arr = append(arr, item)
		}
		return marshalToolCallsValue(arr, s.maxBytes)
	}
	if len(s.anthropic) > 0 {
		return marshalToolCallsValue(s.anthropic, s.maxBytes)
	}
	return "", false
}

// openAIResponseText kept for tests / legacy callers.
func openAIResponseText(body []byte) string {
	return extractOpenAIResponse(body, 1<<20).AssistantText
}
