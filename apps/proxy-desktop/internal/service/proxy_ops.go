package service

import (
	"bytes"
	"encoding/base64"
	"strings"
	"time"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/pkg/export"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/database"
)

func (s *DesktopService) ValidateProxyConfig() Result[ProxyConfigSummaryDTO] {
	cfg := s.deps.Bridge.ProxyMerged
	out := ProxyConfigSummaryDTO{
		ListenAddr:      cfg.ListenAddr(),
		DatabasePath:    cfg.Storage.DatabasePath,
		CaptureMode:     cfg.Capture.Mode,
		StoreContent:    cfg.Capture.StoreContent,
		ProxyConfigPath: s.deps.Bridge.ProxyPath,
		BridgeOK:        s.deps.Bridge.BridgeMessage == "",
		BridgeMessage:   s.deps.Bridge.BridgeMessage,
	}
	return ok(out)
}

func (s *DesktopService) PingDatabase() Result[bool] {
	if err := database.Ping(s.ctx(), s.deps.DB); err != nil {
		return fail[bool](err)
	}
	return ok(true)
}

func (s *DesktopService) RunDatabaseMigrations() Result[bool] {
	if err := database.RunMigrations(s.deps.DB); err != nil {
		return fail[bool](err)
	}
	return ok(true)
}

func (s *DesktopService) GetStorageStats() Result[StorageStatsDTO] {
	total, err := export.StorageStats(s.ctx(), s.deps.DB)
	if err != nil {
		return fail[StorageStatsDTO](err)
	}
	out := StorageStatsDTO{
		RequestRows: total,
		ProxyRunning: false,
	}
	if s.deps.Runtime != nil {
		ins, dropped, qlen, ok := s.deps.Runtime.WriterStats()
		if ok {
			out.ProxyRunning = true
			out.Inserted = ins
			out.Dropped = dropped
			out.QueueLen = qlen
		}
	}
	return ok(out)
}

func (s *DesktopService) ExportMetadata(format, since string, includeContent bool) Result[ExportPayloadDTO] {
	f := export.Filter{IncludeContent: includeContent}
	since = strings.TrimSpace(since)
	if since != "" {
		t, err := time.Parse("2006-01-02", since)
		if err != nil {
			return fail[ExportPayloadDTO](err)
		}
		f.Since = &t
	}
	var buf bytes.Buffer
	switch strings.ToLower(strings.TrimSpace(format)) {
	case "jsonl", "":
		if err := export.ExportJSONL(s.ctx(), s.deps.DB, &buf, f); err != nil {
			return fail[ExportPayloadDTO](err)
		}
		format = "jsonl"
	case "csv":
		if err := export.ExportCSV(s.ctx(), s.deps.DB, &buf, f); err != nil {
			return fail[ExportPayloadDTO](err)
		}
	default:
		return fail[ExportPayloadDTO](errInvalidExportFormat)
	}
	suggested := "proxy-metadata-export." + format
	return ok(ExportPayloadDTO{
		Filename:      suggested,
		ContentBase64: base64.StdEncoding.EncodeToString(buf.Bytes()),
		ByteCount:     buf.Len(),
	})
}
