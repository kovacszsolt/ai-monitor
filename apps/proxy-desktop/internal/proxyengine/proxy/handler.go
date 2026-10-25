package proxy

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/auth"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/capture"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/capture/content"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/config"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/metrics"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/providers"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/routing"
	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/usage"
)

type contextKey string

const authContextKey contextKey = "auth"

type Services struct {
	Keys      *auth.KeyStore
	Providers *providers.Store
	Routes    *routing.Store
	Writer    *capture.Writer
	MaxBody   int64
	Capture   config.CaptureConfig
}

type Handler struct {
	svc    Services
	client *http.Client
}

func NewHandler(svc Services) *Handler {
	return &Handler{svc: svc, client: NewHTTPClient(10000, 300000)}
}

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if strings.HasPrefix(r.URL.Path, "/health/") {
			next.ServeHTTP(w, r)
			return
		}
		hdr := r.Header.Get("Authorization")
		if !strings.HasPrefix(hdr, "Bearer ") {
			logProxyAuthFailure(r, "missing_bearer", nil)
			writeProxyJSONError(w, http.StatusUnauthorized, "missing bearer token")
			return
		}
		token := strings.TrimPrefix(hdr, "Bearer ")
		ak, err := h.svc.Keys.Authenticate(r.Context(), token)
		if err != nil {
			logProxyAuthFailure(r, authFailureReason(err), err)
			writeProxyJSONError(w, http.StatusUnauthorized, "invalid api key")
			return
		}
		slog.Debug("proxy client authenticated",
			"error_source", ErrorSourceProxy,
			"auth_stage", "proxy_api_key",
			"project_id", ak.ProjectID,
			"api_key_id", ak.ID,
			"api_key_name", ak.Name,
			"path", r.URL.Path,
			"request_id", middleware.GetReqID(r.Context()),
		)
		ctx := context.WithValue(r.Context(), authContextKey, ak)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func authFrom(ctx context.Context) auth.AuthenticatedKey {
	return ctx.Value(authContextKey).(auth.AuthenticatedKey)
}

func (h *Handler) Register(mux interface {
	Get(string, http.HandlerFunc)
	Post(string, http.HandlerFunc)
}) {
	mux.Get("/health/live", h.live)
	mux.Get("/health/ready", h.ready)
	mux.Get("/v1/models", h.models)
	mux.Post("/v1/chat/completions", h.chatCompletions)
	mux.Post("/v1/responses", h.responses)
	mux.Post("/v1/messages", h.messages)
}

func (h *Handler) live(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}

func (h *Handler) ready(w http.ResponseWriter, r *http.Request) {
	if h.svc.Writer == nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ready"))
}

func (h *Handler) models(w http.ResponseWriter, r *http.Request) {
	ak := authFrom(r.Context())
	routes, err := h.svc.Routes.List(r.Context())
	if err != nil {
		writeProxyJSONError(w, http.StatusInternalServerError, "internal error")
		return
	}
	allowed, err := h.allowedRouteSet(r.Context(), ak.APIKey)
	if err != nil {
		writeProxyJSONError(w, http.StatusInternalServerError, "internal error")
		return
	}
	type modelObj struct {
		ID      string `json:"id"`
		Object  string `json:"object"`
		OwnedBy string `json:"owned_by"`
	}
	var data []modelObj
	for _, rt := range routes {
		if !rt.Enabled {
			continue
		}
		if !allowed[rt.ID] {
			continue
		}
		data = append(data, modelObj{ID: rt.PublicModelName, Object: "model", OwnedBy: "proxy"})
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"object": "list",
		"data":   data,
	})
}

func (h *Handler) chatCompletions(w http.ResponseWriter, r *http.Request) {
	h.forward(w, r, routing.ProtocolOpenAIChat, "/v1/chat/completions")
}

func (h *Handler) responses(w http.ResponseWriter, r *http.Request) {
	h.forward(w, r, routing.ProtocolOpenAIResponses, "/v1/responses")
}

func (h *Handler) messages(w http.ResponseWriter, r *http.Request) {
	h.forward(w, r, routing.ProtocolAnthropicMessages, "/v1/messages")
}

func (h *Handler) forward(w http.ResponseWriter, r *http.Request, expected routing.Protocol, endpoint string) {
	started := time.Now()
	reqID := middleware.GetReqID(r.Context())
	if reqID == "" {
		reqID = uuid.NewString()
	}
	ak := authFrom(r.Context())
	reqTag := requestTagFromHeader(r)

	body, err := readLimited(r.Body, h.svc.MaxBody)
	if err != nil {
		logProxyRequestError(ErrorSourceProxy, "read_body", "request_id", reqID, "err", err.Error())
		writeProxyJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	publicModel, stream, _, err := extractModel(body)
	if err != nil {
		logProxyRequestError(ErrorSourceProxy, "parse_model", "request_id", reqID, "err", err.Error())
		writeProxyJSONError(w, http.StatusBadRequest, err.Error())
		return
	}

	route, provider, err := h.resolveRoute(r.Context(), ak.APIKey, publicModel, expected)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, errUnauthorizedModel) {
			status = http.StatusForbidden
		}
		logProxyRequestError(ErrorSourceProxy, "routing",
			"request_id", reqID,
			"project_id", ak.ProjectID,
			"api_key_id", ak.ID,
			"public_model", publicModel,
			"protocol", string(expected),
			"http_status", status,
			"err", err.Error(),
		)
		writeProxyJSONError(w, status, err.Error())
		h.enqueueFailure(reqID, ak, route, provider, publicModel, endpoint, expected, started, body, err, reqTag)
		return
	}

	rewritten, err := rewriteModel(body, route.UpstreamModelName)
	if err != nil {
		logProxyRequestError(ErrorSourceProxy, "rewrite_model", "request_id", reqID, "err", err.Error())
		writeProxyJSONError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	upURL := joinUpstreamURL(provider.BaseURL, endpoint)
	logProxyRequestInfo(
		"request_id", reqID,
		"project_id", ak.ProjectID,
		"api_key_id", ak.ID,
		"public_model", publicModel,
		"upstream_model", route.UpstreamModelName,
		"provider_id", provider.ID,
		"provider_name", provider.Name,
		"provider_type", string(provider.ProviderType),
		"upstream_host", upstreamHost(upURL),
		"endpoint", endpoint,
		"streaming", stream,
	)

	upReq, err := http.NewRequestWithContext(r.Context(), http.MethodPost, upURL, bytes.NewReader(rewritten))
	if err != nil {
		logProxyRequestError(ErrorSourceProxy, "build_upstream_request", "request_id", reqID, "err", err.Error())
		writeProxyJSONError(w, http.StatusInternalServerError, "internal error")
		return
	}
	copySafeHeaders(r.Header, upReq.Header)
	upReq.Header.Set("Content-Type", "application/json")
	upReq.Header.Del("Authorization")
	if err := applyProviderAuth(provider, upReq); err != nil {
		logProxyRequestError(ErrorSourceProxy, "provider_credential",
			"request_id", reqID,
			"provider_id", provider.ID,
			"provider_name", provider.Name,
			"err", err.Error(),
		)
		writeProxyJSONError(w, http.StatusBadGateway, "provider auth failed")
		h.enqueueFailure(reqID, ak, route, provider, publicModel, endpoint, expected, started, body, err, reqTag)
		return
	}

	promptCap, captureErrs := h.extractPromptCapture(expected, body, reqID)
	h.enqueueRequestBegin(reqID, ak, route, provider, publicModel, endpoint, expected, started, stream, body, promptCap, captureErrs, reqTag)

	upStart := time.Now()
	resp, err := h.client.Do(upReq)
	upLatency := time.Since(upStart).Milliseconds()
	if err != nil {
		completed := time.Now()
		dur := completed.Sub(started).Milliseconds()
		if errors.Is(err, context.Canceled) {
			logProxyRequestError(ErrorSourceProxy, "client_cancelled",
				"request_id", reqID,
				"provider_name", provider.Name,
			)
			h.enqueueRequestFinishEarly(reqID, ak, route, provider, publicModel, endpoint, expected, started, &completed, &dur, &upLatency, body, capture.StateCancelled, true, "", reqTag)
			return
		}
		logProxyRequestError(ErrorSourceProvider, "upstream_connection",
			"request_id", reqID,
			"provider_id", provider.ID,
			"provider_name", provider.Name,
			"upstream_host", upstreamHost(upURL),
			"err", err.Error(),
		)
		writeProxyJSONError(w, http.StatusBadGateway, "upstream request failed")
		h.enqueueRequestFinishEarly(reqID, ak, route, provider, publicModel, endpoint, expected, started, &completed, &dur, &upLatency, body, capture.StateFailed, false, err.Error(), reqTag)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		logProxyRequestError(ErrorSourceProvider, classifyUpstreamError(resp.StatusCode),
			"request_id", reqID,
			"provider_id", provider.ID,
			"provider_name", provider.Name,
			"upstream_host", upstreamHost(upURL),
			"http_status", resp.StatusCode,
			"upstream_latency_ms", upLatency,
		)
	}

	var respBytes int64
	var ttfb, ttft *int64
	var streamM streamMetrics
	var u usage.Usage
	state := capture.StateCompleted
	errType := ""
	var respBody []byte
	var streamAcc *content.StreamAccumulator
	if h.svc.Capture.StoreContent && stream {
		streamAcc = content.NewStreamAccumulator(h.svc.Capture.MaxResponseBytes)
	}

	if !stream {
		respBody, err = io.ReadAll(io.LimitReader(resp.Body, h.svc.MaxBody))
		if err != nil {
			state = capture.StateFailed
			errType = "read_error"
			logProxyRequestError(ErrorSourceProvider, "read_upstream_body",
				"request_id", reqID,
				"provider_name", provider.Name,
				"err", err.Error(),
			)
		} else {
			decoded, decErr := maybeDecodeUpstreamBody(resp.Header, respBody)
			if decErr != nil {
				state = capture.StateFailed
				errType = "decode_error"
				logProxyRequestError(ErrorSourceProvider, "decode_upstream_body",
					"request_id", reqID,
					"provider_name", provider.Name,
					"err", decErr.Error(),
				)
			} else {
				respBody = decoded
			}
			respBytes = int64(len(respBody))
			ttfb = ptrInt64(0)
			if resp.StatusCode >= 400 {
				slog.Warn("upstream error body preview",
					"error_source", ErrorSourceProvider,
					"request_id", reqID,
					"provider_name", provider.Name,
					"http_status", resp.StatusCode,
					"message_preview", previewJSONError(respBody, 512),
				)
			}
			if state == capture.StateCompleted {
				if provider.ProviderType == providers.TypeAnthropic {
					u = usage.FromAnthropicJSON(respBody)
				} else {
					u = usage.FromOpenAIJSON(respBody)
				}
			}
		}
	}

	copyResponseHeaders(w, resp.Header)
	// Upstream error bodies are passed through; this header marks that the status is from the provider.
	if resp.StatusCode >= 400 {
		w.Header().Set("X-AI-Proxy-Error-Source", ErrorSourceProvider)
	}
	w.WriteHeader(resp.StatusCode)

	if stream {
		w.Header().Set("Content-Type", resp.Header.Get("Content-Type"))
		fl, _ := w.(http.Flusher)
		var copyErr error
		streamM, copyErr = copyStream(w, resp.Body, fl, func(line string) {
			if u2 := usageFromStreamLine(line); u2 != nil {
				u = *u2
			}
			if streamAcc != nil {
				streamAcc.Feed(line)
			}
		})
		respBytes = streamM.ResponseBytes
		ttfb = streamM.TTFBMS
		ttft = streamM.TTFTMS
		if streamM.LastUsage != nil {
			u = *streamM.LastUsage
		}
		if copyErr != nil {
			if errors.Is(copyErr, context.Canceled) {
				state = capture.StateCancelled
			} else {
				state = capture.StateInterrupted
				errType = "stream_error"
			}
		}
	} else if len(respBody) > 0 {
		_, _ = w.Write(respBody)
	}

	completed := time.Now()
	dur := completed.Sub(started).Milliseconds()
	rec := capture.Record{
		ID:                reqID,
		ProjectID:         ak.ProjectID,
		ProxyAPIKeyID:     ak.ID,
		ProviderID:        provider.ID,
		ModelRouteID:      route.ID,
		RequestedModel:    publicModel,
		ActualModel:       route.UpstreamModelName,
		Endpoint:          endpoint,
		Protocol:          string(expected),
		StartedAt:         started,
		CompletedAt:       &completed,
		DurationMS:        &dur,
		UpstreamLatencyMS: &upLatency,
		TTFBMS:            ttfb,
		TTFTMS:            ttft,
		StatusCode:        &resp.StatusCode,
		RequestState:      state,
		Streaming:         stream,
		InputTokens:       u.InputTokens,
		OutputTokens:      u.OutputTokens,
		CachedInputTokens: u.CachedInputTokens,
		CacheWriteTokens:  u.CacheWriteTokens,
		ReasoningTokens:   u.ReasoningTokens,
		UsageStatus:       u.Status,
		RawUsageJSON:      u.RawJSON,
		RequestBytes:      int64(len(body)),
		ResponseBytes:     respBytes,
		ErrorType:         errType,
		Tag:               reqTag,
	}
	metricsCompleted := state == capture.StateCompleted && resp.StatusCode < 400
	timing := metrics.TimingInput{
		TotalDurationMS: &dur,
		TTFTMS:          ttft,
	}
	genEstimated := false
	promptEstimated := false
	if stream {
		timing.DecodeDurationMS = metrics.DecodeDurationFromStream(streamM.FirstOutputTokenMS, streamM.LastOutputTokenMS)
		if timing.DecodeDurationMS != nil {
			genEstimated = true
		}
	} else if metricsCompleted && u.OutputTokens != nil && *u.OutputTokens > 0 && dur > 0 {
		// Non-stream: no per-token timing from provider; approximate output throughput as output / wall duration.
		timing.DecodeDurationMS = &dur
		genEstimated = true
	}
	capture.ApplyTokenMetrics(&rec, u, timing, metricsCompleted, genEstimated, promptEstimated)
	rec.Response = h.buildResponseCapture(expected, stream, streamAcc, respBody)
	h.attachRawUpstreamBody(&rec, resp.StatusCode, u.Status, stream, streamAcc, respBody)
	h.maybeDumpHTTPExchange(started, reqID, endpoint, publicModel, route.UpstreamModelName, provider.Name,
		stream, resp.StatusCode, u.Status, state, errType, body, respBody, streamAcc)
	h.enqueueRequestFinish(rec)
	if state == capture.StateCompleted && resp.StatusCode < 400 {
		slog.Info("proxy request completed",
			"request_id", reqID,
			"provider_name", provider.Name,
			"http_status", resp.StatusCode,
			"duration_ms", dur,
			"upstream_latency_ms", upLatency,
			"streaming", stream,
		)
	}
}

var errUnauthorizedModel = errors.New("model not allowed for this api key")

func protocolEndpointHint(protocol routing.Protocol) string {
	switch protocol {
	case routing.ProtocolOpenAIChat:
		return "POST /v1/chat/completions"
	case routing.ProtocolOpenAIResponses:
		return "POST /v1/responses"
	case routing.ProtocolAnthropicMessages:
		return "POST /v1/messages"
	default:
		return string(protocol)
	}
}

func (h *Handler) resolveRoute(ctx context.Context, key auth.APIKey, publicModel string, protocol routing.Protocol) (routing.Route, providers.Provider, error) {
	rt, err := h.svc.Routes.GetByPublicModel(ctx, publicModel)
	if err != nil {
		return routing.Route{}, providers.Provider{}, fmt.Errorf("unknown model")
	}
	if !rt.Enabled {
		return routing.Route{}, providers.Provider{}, fmt.Errorf("model disabled")
	}
	if rt.Protocol != protocol {
		return routing.Route{}, providers.Provider{}, fmt.Errorf(
			"protocol mismatch for model %q: route protocol is %s but this endpoint expects %s (%s)",
			publicModel, rt.Protocol, protocol, protocolEndpointHint(protocol),
		)
	}
	allowed, err := h.allowedRouteSet(ctx, key)
	if err != nil {
		return routing.Route{}, providers.Provider{}, err
	}
	if !allowed[rt.ID] {
		return routing.Route{}, providers.Provider{}, errUnauthorizedModel
	}
	p, err := h.svc.Providers.Get(ctx, rt.ProviderID)
	if err != nil {
		return routing.Route{}, providers.Provider{}, fmt.Errorf("provider not found")
	}
	if !p.Enabled {
		return routing.Route{}, providers.Provider{}, fmt.Errorf("provider disabled")
	}
	return rt, p, nil
}

func (h *Handler) allowedRouteSet(ctx context.Context, key auth.APIKey) (map[string]bool, error) {
	if key.AllowAllModels {
		routes, err := h.svc.Routes.List(ctx)
		if err != nil {
			return nil, err
		}
		m := make(map[string]bool, len(routes))
		for _, r := range routes {
			m[r.ID] = true
		}
		return m, nil
	}
	ids, err := h.svc.Keys.AllowedModelRouteIDs(ctx, key.ID)
	if err != nil {
		return nil, err
	}
	if len(ids) == 0 {
		return map[string]bool{}, nil
	}
	m := make(map[string]bool, len(ids))
	for _, id := range ids {
		m[id] = true
	}
	return m, nil
}

func applyProviderAuth(p providers.Provider, req *http.Request) error {
	key, err := p.APIKey()
	if err != nil {
		return err
	}
	switch p.ProviderType {
	case providers.TypeAnthropic:
		req.Header.Set("x-api-key", key)
		req.Header.Set("anthropic-version", "2023-06-01")
	default:
		if key != "" {
			req.Header.Set("Authorization", "Bearer "+key)
		}
	}
	return nil
}

func copySafeHeaders(src http.Header, dst http.Header) {
	for k, vals := range src {
		lk := strings.ToLower(k)
		if lk == "authorization" || lk == "host" || lk == "connection" || lk == "accept-encoding" || lk == "x-ai-proxy-tag" {
			continue
		}
		for _, v := range vals {
			dst.Add(k, v)
		}
	}
}

var hopByHop = map[string]bool{
	"connection": true, "keep-alive": true, "proxy-authenticate": true,
	"proxy-authorization": true, "te": true, "trailers": true, "transfer-encoding": true, "upgrade": true,
}

func copyResponseHeaders(w http.ResponseWriter, src http.Header) {
	for k, vals := range src {
		if hopByHop[strings.ToLower(k)] {
			continue
		}
		for _, v := range vals {
			w.Header().Add(k, v)
		}
	}
}

func readLimited(r io.Reader, max int64) ([]byte, error) {
	return io.ReadAll(io.LimitReader(r, max))
}

func writeProxyJSONError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-AI-Proxy-Error-Source", ErrorSourceProxy)
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"error": map[string]string{
			"message": msg,
			"source":  ErrorSourceProxy,
		},
	})
}

func ptrInt64(v int64) *int64 { return &v }

func (h *Handler) extractPromptCapture(protocol routing.Protocol, body []byte, requestID string) (*capture.PromptCapture, []capture.CaptureError) {
	if !h.svc.Capture.StoreContent {
		return nil, nil
	}
	p, err := content.ExtractRequest(protocol, body, h.svc.Capture.MaxPromptBytes)
	if err != nil {
		code := content.PromptErrorCode(err)
		slog.Warn("prompt content not stored",
			"request_id", requestID,
			"protocol", string(protocol),
			"error_code", code,
			"err", err.Error(),
		)
		return nil, []capture.CaptureError{{
			Stage:   "prompt",
			Code:    code,
			Message: err.Error(),
		}}
	}
	return &capture.PromptCapture{
		LastUserText:    p.LastUserText,
		MessagesJSON:    p.MessagesJSON,
		PromptTruncated: p.PromptTruncated,
	}, nil
}

func shouldStoreRawUpstreamBody(cfg config.CaptureConfig, statusCode int, usageStatus string, parsedContent bool) bool {
	if !cfg.StoreUpstreamBodyOnError {
		return false
	}
	if statusCode >= 400 || usageStatus == "parse_error" {
		return true
	}
	if !cfg.StoreContent {
		return false
	}
	return !parsedContent
}

func (h *Handler) attachRawUpstreamBody(rec *capture.Record, statusCode int, usageStatus string, stream bool, acc *content.StreamAccumulator, respBody []byte) {
	parsed := rec.Response != nil
	if !shouldStoreRawUpstreamBody(h.svc.Capture, statusCode, usageStatus, parsed) {
		return
	}
	var raw string
	var truncated bool
	switch {
	case stream && acc != nil:
		raw, truncated = acc.RawBody()
	case len(respBody) > 0:
		raw, truncated = content.TruncateUTF8(string(respBody), h.svc.Capture.MaxResponseBytes)
	default:
		return
	}
	if raw == "" {
		return
	}
	if rec.Response == nil {
		rec.Response = &capture.ResponseCapture{}
	}
	rec.Response.RawBody = raw
	rec.Response.RawBodyTruncated = truncated
}

func (h *Handler) buildResponseCapture(protocol routing.Protocol, stream bool, acc *content.StreamAccumulator, respBody []byte) *capture.ResponseCapture {
	if !h.svc.Capture.StoreContent {
		return nil
	}
	var data content.ResponseData
	switch {
	case stream && acc != nil:
		data = acc.Result()
	case len(respBody) > 0:
		data = content.ExtractResponse(protocol, respBody, h.svc.Capture.MaxResponseBytes)
	default:
		return nil
	}
	if data.AssistantText == "" && data.ToolCallsJSON == "" {
		return nil
	}
	return &capture.ResponseCapture{
		AssistantText:      data.AssistantText,
		ResponseTruncated:  data.ResponseTruncated,
		ToolCallsJSON:      data.ToolCallsJSON,
		ToolCallsTruncated: data.ToolCallsTruncated,
	}
}

func (h *Handler) enqueueFailure(reqID string, ak auth.AuthenticatedKey, route routing.Route, provider providers.Provider,
	publicModel, endpoint string, protocol routing.Protocol, started time.Time, body []byte, err error, tag string) {
	completed := time.Now()
	dur := completed.Sub(started).Milliseconds()
	code := http.StatusBadGateway
	rec := capture.Record{
		ID:             reqID,
		ProjectID:      ak.ProjectID,
		ProxyAPIKeyID:  ak.ID,
		ProviderID:     provider.ID,
		ModelRouteID:   route.ID,
		RequestedModel: publicModel,
		ActualModel:    route.UpstreamModelName,
		Endpoint:       endpoint,
		Protocol:       string(protocol),
		StartedAt:      started,
		CompletedAt:    &completed,
		DurationMS:     &dur,
		StatusCode:     &code,
		RequestState:   capture.StateFailed,
		UsageStatus:    "n/a",
		RequestBytes:   int64(len(body)),
		ErrorType:      err.Error(),
		Progress:       capture.ProgressFinished,
		Phase:          capture.PhaseOneShot,
		Tag:            tag,
	}
	h.attachPromptCapture(&rec, protocol, body, reqID)
	if err := h.svc.Writer.Enqueue(rec); err != nil {
		slog.Warn("metadata queue full", "request_id", reqID)
	}
}

func (h *Handler) enqueueRequestBegin(reqID string, ak auth.AuthenticatedKey, route routing.Route, provider providers.Provider,
	publicModel, endpoint string, protocol routing.Protocol, started time.Time, stream bool, body []byte,
	prompt *capture.PromptCapture, captureErrs []capture.CaptureError, tag string) {
	rec := capture.Record{
		ID:             reqID,
		ProjectID:      ak.ProjectID,
		ProxyAPIKeyID:  ak.ID,
		ProviderID:     provider.ID,
		ModelRouteID:   route.ID,
		RequestedModel: publicModel,
		ActualModel:    route.UpstreamModelName,
		Endpoint:       endpoint,
		Protocol:       string(protocol),
		StartedAt:      started,
		RequestState:   capture.StateStarted,
		Progress:       capture.ProgressStarted,
		Phase:          capture.PhaseBegin,
		Streaming:      stream,
		UsageStatus:    "pending",
		RequestBytes:   int64(len(body)),
		Prompt:         prompt,
		Errors:         captureErrs,
		Tag:            tag,
	}
	if err := h.svc.Writer.Enqueue(rec); err != nil {
		slog.Warn("metadata queue full on begin", "request_id", reqID)
	}
}

func (h *Handler) enqueueRequestFinish(rec capture.Record) {
	rec.Phase = capture.PhaseFinish
	rec.Progress = capture.ProgressFinished
	rec.Prompt = nil
	rec.Errors = nil
	if err := h.svc.Writer.Enqueue(rec); err != nil {
		slog.Warn("metadata queue full", "request_id", rec.ID)
	}
}

func (h *Handler) enqueueRequestFinishEarly(reqID string, ak auth.AuthenticatedKey, route routing.Route, provider providers.Provider,
	publicModel, endpoint string, protocol routing.Protocol, started time.Time, completed *time.Time, dur, upLatency *int64,
	body []byte, state capture.RequestState, cancelled bool, errType string, tag string) {
	code := http.StatusBadGateway
	rec := capture.Record{
		ID:                reqID,
		ProjectID:         ak.ProjectID,
		ProxyAPIKeyID:     ak.ID,
		ProviderID:        provider.ID,
		ModelRouteID:      route.ID,
		RequestedModel:    publicModel,
		ActualModel:       route.UpstreamModelName,
		Endpoint:          endpoint,
		Protocol:          string(protocol),
		StartedAt:         started,
		CompletedAt:       completed,
		DurationMS:        dur,
		UpstreamLatencyMS: upLatency,
		StatusCode:        &code,
		RequestState:      state,
		Streaming:         false,
		UsageStatus:       "n/a",
		RequestBytes:      int64(len(body)),
		ErrorType:         errType,
		CancelledFlag:     cancelled,
		Tag:               tag,
	}
	h.enqueueRequestFinish(rec)
}

func (h *Handler) attachPromptCapture(rec *capture.Record, protocol routing.Protocol, body []byte, requestID string) {
	prompt, errs := h.extractPromptCapture(protocol, body, requestID)
	rec.Prompt = prompt
	rec.Errors = errs
}

func (h *Handler) enqueueCancelled(reqID string, ak auth.AuthenticatedKey, route routing.Route, provider providers.Provider,
	publicModel, endpoint string, protocol routing.Protocol, started time.Time, body []byte) {
	completed := time.Now()
	dur := completed.Sub(started).Milliseconds()
	rec := capture.Record{
		ID:             reqID,
		ProjectID:      ak.ProjectID,
		ProxyAPIKeyID:  ak.ID,
		ProviderID:     provider.ID,
		ModelRouteID:   route.ID,
		RequestedModel: publicModel,
		Endpoint:       endpoint,
		Protocol:       string(protocol),
		StartedAt:      started,
		CompletedAt:    &completed,
		DurationMS:     &dur,
		RequestState:   capture.StateCancelled,
		CancelledFlag:  true,
		UsageStatus:    "n/a",
		RequestBytes:   int64(len(body)),
		Progress:       capture.ProgressFinished,
		Phase:          capture.PhaseOneShot,
	}
	h.attachPromptCapture(&rec, protocol, body, reqID)
	if err := h.svc.Writer.Enqueue(rec); err != nil {
		slog.Warn("metadata queue full", "request_id", reqID)
	}
}
