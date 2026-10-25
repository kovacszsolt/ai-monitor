package proxy

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"time"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/usage"
)

type streamMetrics struct {
	TTFBMS              *int64
	TTFTMS              *int64
	FirstOutputTokenMS  *int64
	LastOutputTokenMS   *int64
	ResponseBytes       int64
	LastUsage           *usage.Usage
}

func copyStream(w io.Writer, r io.Reader, flusher interface{ Flush() }, onLine func(string)) (streamMetrics, error) {
	m := streamMetrics{}
	start := time.Now()
	br := bufio.NewReader(r)
	firstByte := true
	var buf bytes.Buffer
	for {
		line, err := br.ReadBytes('\n')
		if len(line) > 0 {
			m.ResponseBytes += int64(len(line))
			if firstByte {
				ms := time.Since(start).Milliseconds()
				m.TTFBMS = &ms
				firstByte = false
			}
			buf.Write(line)
			lineStr := string(line)
			if onLine != nil {
				onLine(lineStr)
			}
			if _, werr := w.Write(line); werr != nil {
				return m, werr
			}
			if flusher != nil {
				flusher.Flush()
			}
			nowMS := time.Since(start).Milliseconds()
			if u := usageFromStreamLine(lineStr); u != nil {
				if u.OutputTokens != nil || u.InputTokens != nil {
					if m.TTFTMS == nil {
						m.TTFTMS = ptrMS(time.Since(start))
					}
					m.LastUsage = u
				}
			} else if isContentDelta(lineStr) {
				if m.TTFTMS == nil {
					m.TTFTMS = ptrMS(time.Since(start))
				}
				if m.FirstOutputTokenMS == nil {
					m.FirstOutputTokenMS = &nowMS
				}
				m.LastOutputTokenMS = &nowMS
			}
		}
		if err != nil {
			if err == io.EOF {
				return m, nil
			}
			return m, err
		}
	}
}

func usageFromStreamLine(line string) *usage.Usage {
	line = strings.TrimSpace(line)
	if strings.HasPrefix(line, "data:") {
		if u := usage.FromOpenAIStreamLine(line); u != nil {
			return u
		}
	}
	return nil
}

func isContentDelta(line string) bool {
	line = strings.TrimSpace(line)
	if !strings.HasPrefix(line, "data:") {
		return false
	}
	payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
	if payload == "[DONE]" {
		return false
	}
	if strings.Contains(payload, `"usage"`) && !strings.Contains(payload, `"delta"`) && !strings.Contains(payload, `"content_block"`) {
		return false
	}
	return strings.Contains(payload, `"delta"`) || strings.Contains(payload, `"content_block"`)
}

func ptrMS(d time.Duration) *int64 {
	ms := d.Milliseconds()
	return &ms
}
