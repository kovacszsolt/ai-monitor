package proxy

import (
	"log/slog"
	"time"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/capture"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/capture/content"
)

func (h *Handler) maybeDumpHTTPExchange(
	started time.Time,
	reqID, endpoint, publicModel, upstreamModel, providerName string,
	stream bool,
	statusCode int,
	usageStatus string,
	state capture.RequestState,
	errType string,
	requestBody []byte,
	respBody []byte,
	streamAcc *content.StreamAccumulator,
) {
	cfg := h.svc.Capture.FileDump
	if !capture.ShouldFileDump(cfg, statusCode, usageStatus, state, errType) {
		return
	}
	responsePayload := respBody
	if stream && streamAcc != nil {
		raw, _ := streamAcc.RawBody()
		if raw != "" {
			responsePayload = []byte(raw)
		}
	}
	max := cfg.MaxBodyBytes
	dump := capture.HTTPExchangeDump{
		RequestID:      reqID,
		StartedAt:      capture.FormatDumpTimestamp(started),
		Endpoint:       endpoint,
		RequestedModel: publicModel,
		UpstreamModel:  upstreamModel,
		ProviderName:   providerName,
		Streaming:      stream,
		StatusCode:     statusCode,
		UsageStatus:    usageStatus,
		RequestState:   string(state),
		ErrorType:      errType,
		Request:        capture.BodyPartFromBytes(requestBody, max),
		Response:       capture.BodyPartFromBytes(responsePayload, max),
	}
	go func(d capture.HTTPExchangeDump) {
		path, err := capture.WriteHTTPExchangeDump(cfg, d)
		if err != nil {
			slog.Warn("http exchange file dump failed", "request_id", d.RequestID, "err", err.Error())
			return
		}
		slog.Info("http exchange file dump written", "request_id", d.RequestID, "path", path)
	}(dump)
}
