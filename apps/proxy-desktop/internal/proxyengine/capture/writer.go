package capture

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/config"
)

var ErrQueueFull = errors.New("persistence queue full")

type Writer struct {
	db       *sql.DB
	cfg      config.PersistenceConfig
	ch       chan Record
	done     chan struct{}
	wg       sync.WaitGroup
	inserted int64
	dropped  int64
}

func NewWriter(db *sql.DB, cfg config.PersistenceConfig) *Writer {
	return &Writer{
		db:   db,
		cfg:  cfg,
		ch:   make(chan Record, cfg.QueueCapacity),
		done: make(chan struct{}),
	}
}

func (w *Writer) Start() {
	w.wg.Add(1)
	go w.loop()
}

func (w *Writer) Enqueue(r Record) error {
	select {
	case w.ch <- r:
		return nil
	default:
		w.dropped++
		return ErrQueueFull
	}
}

func (w *Writer) Shutdown(ctx context.Context) error {
	close(w.ch)
	done := make(chan struct{})
	go func() {
		w.wg.Wait()
		close(done)
	}()
	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (w *Writer) Stats() (inserted, dropped int64, queueLen int) {
	return w.inserted, w.dropped, len(w.ch)
}

func (w *Writer) loop() {
	defer w.wg.Done()
	batch := make([]Record, 0, w.cfg.BatchSize)
	timer := time.NewTimer(time.Duration(w.cfg.BatchDelayMS) * time.Millisecond)
	defer timer.Stop()

	flush := func() {
		if len(batch) == 0 {
			return
		}
		if err := w.persistBatch(batch); err != nil {
			slog.Error("capture batch persist failed", "err", err, "count", len(batch))
		} else {
			w.inserted += int64(len(batch))
		}
		batch = batch[:0]
	}

	for {
		select {
		case r, ok := <-w.ch:
			if !ok {
				flush()
				return
			}
			batch = append(batch, r)
			if len(batch) >= w.cfg.BatchSize {
				flush()
				if !timer.Stop() {
					select {
					case <-timer.C:
					default:
					}
				}
				timer.Reset(time.Duration(w.cfg.BatchDelayMS) * time.Millisecond)
			}
		case <-timer.C:
			flush()
			timer.Reset(time.Duration(w.cfg.BatchDelayMS) * time.Millisecond)
		}
	}
}

type persistStmts struct {
	insertFull   *sql.Stmt
	insertBegin  *sql.Stmt
	updateFinish *sql.Stmt
	prompt       *sql.Stmt
	response     *sql.Stmt
	errors       *sql.Stmt
}

func (w *Writer) persistBatch(batch []Record) error {
	tx, err := w.db.Begin()
	if err != nil {
		return err
	}
	stmts, err := preparePersistStmts(tx)
	if err != nil {
		tx.Rollback()
		return err
	}
	defer stmts.close()

	for _, r := range batch {
		if err := w.persistOne(tx, stmts, r); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}

func preparePersistStmts(tx *sql.Tx) (persistStmts, error) {
	var s persistStmts
	var err error
	s.insertFull, err = tx.Prepare(`
		INSERT INTO ai_requests (
			id, project_id, proxy_api_key_id, provider_id, model_route_id,
			requested_model, actual_model, endpoint, protocol,
			started_at, completed_at, duration_ms, upstream_latency_ms, ttfb_ms, ttft_ms,
			status_code, request_state, progress, streaming,
			input_tokens, output_tokens, cached_input_tokens, cache_write_tokens, reasoning_tokens,
			uncached_input_tokens, prefill_duration_ms, decode_duration_ms,
			prompt_tps, generation_tps, effective_tps, cache_hit_rate,
			generation_speed_estimated, prompt_speed_estimated,
			usage_status, raw_usage_json, request_bytes, response_bytes,
			error_type, timeout_flag, cancelled_flag, tag
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return s, err
	}
	s.insertBegin, err = tx.Prepare(`
		INSERT INTO ai_requests (
			id, project_id, proxy_api_key_id, provider_id, model_route_id,
			requested_model, actual_model, endpoint, protocol,
			started_at, request_state, progress, streaming,
			usage_status, request_bytes, response_bytes,
			timeout_flag, cancelled_flag, tag
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return s, err
	}
	s.updateFinish, err = tx.Prepare(`
		UPDATE ai_requests SET
			completed_at = ?, duration_ms = ?, upstream_latency_ms = ?, ttfb_ms = ?, ttft_ms = ?,
			status_code = ?, request_state = ?, progress = ?, streaming = ?,
			input_tokens = ?, output_tokens = ?, cached_input_tokens = ?, cache_write_tokens = ?, reasoning_tokens = ?,
			uncached_input_tokens = ?, prefill_duration_ms = ?, decode_duration_ms = ?,
			prompt_tps = ?, generation_tps = ?, effective_tps = ?, cache_hit_rate = ?,
			generation_speed_estimated = ?, prompt_speed_estimated = ?,
			usage_status = ?, raw_usage_json = ?, request_bytes = ?, response_bytes = ?,
			error_type = ?, timeout_flag = ?, cancelled_flag = ?
		WHERE id = ?`)
	if err != nil {
		return s, err
	}
	s.prompt, err = tx.Prepare(`
		INSERT INTO ai_request_prompts (request_id, last_user_text, messages_json, prompt_truncated, created_at)
		VALUES (?, ?, ?, ?, ?)`)
	if err != nil {
		return s, err
	}
	s.response, err = tx.Prepare(`
		INSERT INTO ai_request_responses (request_id, assistant_text, response_truncated, tool_calls_json, tool_calls_truncated, raw_body, raw_body_truncated, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return s, err
	}
	s.errors, err = tx.Prepare(`
		INSERT INTO ai_request_errors (id, request_id, stage, error_code, error_message, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`)
	return s, err
}

func (s persistStmts) close() {
	for _, stmt := range []*sql.Stmt{s.insertFull, s.insertBegin, s.updateFinish, s.prompt, s.response, s.errors} {
		if stmt != nil {
			stmt.Close()
		}
	}
}

func (w *Writer) persistOne(tx *sql.Tx, stmts persistStmts, r Record) error {
	switch r.Phase {
	case PhaseBegin:
		return w.persistBegin(stmts, r)
	case PhaseFinish:
		return w.persistFinish(stmts, r)
	default:
		return w.persistOneShot(stmts, r)
	}
}

func (w *Writer) persistBegin(stmts persistStmts, r Record) error {
	createdAt := r.StartedAt.UTC().Format(time.RFC3339Nano)
	progress := string(r.Progress)
	if progress == "" {
		progress = string(ProgressStarted)
	}
	state := string(r.RequestState)
	if state == "" {
		state = string(StateStarted)
	}
	usage := r.UsageStatus
	if usage == "" {
		usage = "pending"
	}
	_, err := stmts.insertBegin.Exec(
		r.ID, r.ProjectID, r.ProxyAPIKeyID, r.ProviderID, nullStr(r.ModelRouteID),
		r.RequestedModel, r.ActualModel, r.Endpoint, r.Protocol,
		createdAt, state, progress, boolInt(r.Streaming),
		usage, r.RequestBytes, 0,
		boolInt(r.TimeoutFlag), boolInt(r.CancelledFlag), r.Tag,
	)
	if err != nil {
		return err
	}
	return w.persistSidecars(stmts, r, createdAt, true)
}

func (w *Writer) persistFinish(stmts persistStmts, r Record) error {
	var completed *string
	if r.CompletedAt != nil {
		v := r.CompletedAt.UTC().Format(time.RFC3339Nano)
		completed = &v
	}
	progress := string(r.Progress)
	if progress == "" {
		progress = string(ProgressFinished)
	}
	res, err := stmts.updateFinish.Exec(
		completed,
		r.DurationMS, r.UpstreamLatencyMS, r.TTFBMS, r.TTFTMS,
		r.StatusCode, string(r.RequestState), progress, boolInt(r.Streaming),
		r.InputTokens, r.OutputTokens, r.CachedInputTokens, r.CacheWriteTokens, r.ReasoningTokens,
		r.UncachedInputTokens, r.PrefillDurationMS, r.DecodeDurationMS,
		r.PromptTPS, r.GenerationTPS, r.EffectiveTPS, r.CacheHitRate,
		boolInt(r.GenerationSpeedEstimated), boolInt(r.PromptSpeedEstimated),
		r.UsageStatus, nullStr(r.RawUsageJSON), r.RequestBytes, r.ResponseBytes,
		nullStr(r.ErrorType), boolInt(r.TimeoutFlag), boolInt(r.CancelledFlag),
		r.ID,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		slog.Warn("capture finish update missed row, inserting one-shot fallback", "request_id", r.ID)
		r.Phase = PhaseOneShot
		if r.Progress == "" {
			r.Progress = ProgressFinished
		}
		return w.persistOneShot(stmts, r)
	}
	createdAt := r.StartedAt.UTC().Format(time.RFC3339Nano)
	return w.persistSidecars(stmts, r, createdAt, false)
}

func (w *Writer) persistOneShot(stmts persistStmts, r Record) error {
	var completed *string
	if r.CompletedAt != nil {
		v := r.CompletedAt.UTC().Format(time.RFC3339Nano)
		completed = &v
	}
	createdAt := r.StartedAt.UTC().Format(time.RFC3339Nano)
	progress := string(r.Progress)
	if progress == "" {
		progress = string(ProgressFinished)
	}
	usage := r.UsageStatus
	if usage == "" {
		usage = "unknown"
	}
	_, err := stmts.insertFull.Exec(
		r.ID, r.ProjectID, r.ProxyAPIKeyID, r.ProviderID, nullStr(r.ModelRouteID),
		r.RequestedModel, r.ActualModel, r.Endpoint, r.Protocol,
		createdAt, completed,
		r.DurationMS, r.UpstreamLatencyMS, r.TTFBMS, r.TTFTMS,
		r.StatusCode, string(r.RequestState), progress, boolInt(r.Streaming),
		r.InputTokens, r.OutputTokens, r.CachedInputTokens, r.CacheWriteTokens, r.ReasoningTokens,
		r.UncachedInputTokens, r.PrefillDurationMS, r.DecodeDurationMS,
		r.PromptTPS, r.GenerationTPS, r.EffectiveTPS, r.CacheHitRate,
		boolInt(r.GenerationSpeedEstimated), boolInt(r.PromptSpeedEstimated),
		usage, nullStr(r.RawUsageJSON), r.RequestBytes, r.ResponseBytes,
		nullStr(r.ErrorType), boolInt(r.TimeoutFlag), boolInt(r.CancelledFlag), r.Tag,
	)
	if err != nil {
		return err
	}
	return w.persistSidecars(stmts, r, createdAt, true)
}

// includePrompt: true for begin/one-shot; finish skips prompt (already inserted).
func (w *Writer) persistSidecars(stmts persistStmts, r Record, createdAt string, includePrompt bool) error {
	if includePrompt && r.Prompt != nil {
		_, err := stmts.prompt.Exec(
			r.ID, nullStr(r.Prompt.LastUserText), nullStr(r.Prompt.MessagesJSON),
			boolInt(r.Prompt.PromptTruncated), createdAt,
		)
		if err != nil {
			return err
		}
	}
	if r.Response != nil {
		_, err := stmts.response.Exec(
			r.ID, nullStr(r.Response.AssistantText), boolInt(r.Response.ResponseTruncated),
			nullStr(r.Response.ToolCallsJSON), boolInt(r.Response.ToolCallsTruncated),
			nullStr(r.Response.RawBody), boolInt(r.Response.RawBodyTruncated), createdAt,
		)
		if err != nil {
			return err
		}
	}
	for _, ce := range proxyAndCaptureErrors(r) {
		_, err := stmts.errors.Exec(
			uuid.NewString(), r.ID, ce.Stage, ce.Code, ce.Message, createdAt,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

// proxyAndCaptureErrors persists ai_request_errors rows from capture failures and proxy error_type.
func proxyAndCaptureErrors(r Record) []CaptureError {
	out := append([]CaptureError{}, r.Errors...)
	if r.ErrorType == "" {
		return out
	}
	for _, e := range out {
		if e.Stage == "proxy" && e.Message == r.ErrorType {
			return out
		}
	}
	out = append(out, CaptureError{Stage: "proxy", Code: "proxy", Message: r.ErrorType})
	return out
}

func nullStr(s string) interface{} {
	if s == "" {
		return nil
	}
	return s
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
