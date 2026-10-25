package capture

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/config"
)

// HTTPExchangeDump is written to disk when capture.file_dump is enabled.
type HTTPExchangeDump struct {
	RequestID      string `json:"request_id"`
	StartedAt      string `json:"started_at"`
	Endpoint       string `json:"endpoint"`
	RequestedModel string `json:"requested_model"`
	UpstreamModel  string `json:"upstream_model"`
	ProviderName   string `json:"provider_name"`
	Streaming      bool   `json:"streaming"`
	StatusCode     int    `json:"status_code"`
	UsageStatus    string `json:"usage_status"`
	RequestState   string `json:"request_state"`
	ErrorType      string `json:"error_type,omitempty"`
	Request        BodyPart `json:"request"`
	Response       BodyPart `json:"response"`
}

type BodyPart struct {
	Bytes     int    `json:"bytes"`
	Truncated bool   `json:"truncated,omitempty"`
	Body      string `json:"body"`
}

func ShouldFileDump(cfg config.FileDumpConfig, statusCode int, usageStatus string, requestState RequestState, errType string) bool {
	if !cfg.Enabled {
		return false
	}
	if !cfg.ErrorsOnly {
		return true
	}
	if statusCode >= 400 {
		return true
	}
	if usageStatus == "parse_error" || usageStatus == "missing" {
		return true
	}
	if requestState != StateCompleted {
		return true
	}
	if errType != "" {
		return true
	}
	return false
}

func WriteHTTPExchangeDump(cfg config.FileDumpConfig, dump HTTPExchangeDump) (string, error) {
	if !cfg.Enabled {
		return "", nil
	}
	dir := strings.TrimSpace(cfg.Directory)
	if dir == "" {
		return "", fmt.Errorf("capture.file_dump.directory is required when enabled")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return "", fmt.Errorf("create dump dir: %w", err)
	}
	max := cfg.MaxBodyBytes
	if max < 1 {
		max = 1024 * 1024
	}
	dump.Request = truncateBodyPart(dump.Request, max)
	dump.Response = truncateBodyPart(dump.Response, max)

	name := fmt.Sprintf("%s_%s.json",
		dump.StartedAt,
		sanitizeFileToken(dump.RequestID),
	)
	path := filepath.Join(dir, name)
	raw, err := json.MarshalIndent(dump, "", "  ")
	if err != nil {
		return "", err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, raw, 0o644); err != nil {
		return "", err
	}
	if err := os.Rename(tmp, path); err != nil {
		_ = os.Remove(tmp)
		return "", err
	}
	return path, nil
}

func truncateBodyPart(part BodyPart, max int64) BodyPart {
	if part.Body == "" {
		return part
	}
	if int64(len(part.Body)) <= max {
		return part
	}
	part.Body = truncateUTF8Bytes(part.Body, max)
	part.Truncated = true
	return part
}

func truncateUTF8Bytes(s string, max int64) string {
	if max <= 0 {
		return ""
	}
	if int64(len(s)) <= max {
		return s
	}
	end := int(max)
	for end > 0 && !utf8.ValidString(s[:end]) {
		end--
	}
	return s[:end]
}

func sanitizeFileToken(s string) string {
	s = strings.TrimSpace(s)
	if s == "" {
		return "unknown"
	}
	var b strings.Builder
	for _, r := range s {
		switch r {
		case '/', '\\', ':', '*', '?', '"', '<', '>', '|', ' ':
			b.WriteByte('_')
		default:
			b.WriteRune(r)
		}
	}
	out := b.String()
	if out == "" {
		return "unknown"
	}
	if len(out) > 120 {
		return out[:120]
	}
	return out
}

func FormatDumpTimestamp(t time.Time) string {
	return t.UTC().Format("20060102T150405.000000000")
}

func BodyPartFromBytes(body []byte, max int64) BodyPart {
	if len(body) == 0 {
		return BodyPart{}
	}
	s := string(body)
	truncated := false
	if max > 0 && int64(len(s)) > max {
		s = truncateUTF8Bytes(s, max)
		truncated = true
	}
	return BodyPart{
		Bytes:     len(body),
		Truncated: truncated,
		Body:      s,
	}
}
