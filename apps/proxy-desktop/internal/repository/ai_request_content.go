package repository

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/domain"
)

func (r *AIRequests) GetRequestContent(ctx context.Context, requestID string) (domain.AIRequestContent, error) {
	var requestState, proxyError string
	if err := r.db.QueryRowContext(ctx, `
		SELECT request_state, COALESCE(error_type, '')
		FROM ai_requests WHERE id = ?`, requestID).Scan(&requestState, &proxyError); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.AIRequestContent{}, ErrAIRequestNotFound
		}
		return domain.AIRequestContent{}, err
	}

	if !r.hasTable(ctx, "ai_request_prompts") {
		out := domain.AIRequestContent{
			RequestID:      requestID,
			RequestState:   requestState,
			ProxyError:     proxyError,
			CaptureEnabled: false,
		}
		out.CaptureErrors = mergeProxyErrorRows(proxyError, nil)
		return out, nil
	}

	out := domain.AIRequestContent{
		RequestID:      requestID,
		RequestState:   requestState,
		ProxyError:     proxyError,
		CaptureEnabled: true,
	}

	promptCols := `p.last_user_text, p.messages_json, p.prompt_truncated`
	responseCols := `res.assistant_text, res.response_truncated`
	if r.hasTableColumn(ctx, "ai_request_responses", "tool_calls_json") {
		responseCols += `, res.tool_calls_json, res.tool_calls_truncated`
	}
	if r.hasTableColumn(ctx, "ai_request_responses", "raw_body") {
		responseCols += `, res.raw_body, res.raw_body_truncated`
	}

	q := `SELECT ` + promptCols + `, ` + responseCols + `
		FROM ai_requests r
		LEFT JOIN ai_request_prompts p ON p.request_id = r.id
		LEFT JOIN ai_request_responses res ON res.request_id = r.id
		WHERE r.id = ?`

	row := r.db.QueryRowContext(ctx, q, requestID)

	var lastUser, messages sql.NullString
	var promptTrunc int
	var assistant sql.NullString
	var responseTrunc int
	var toolCalls sql.NullString
	var toolCallsTrunc int
	var rawBody sql.NullString
	var rawBodyTrunc int

	hasToolCols := r.hasTableColumn(ctx, "ai_request_responses", "tool_calls_json")
	hasRawBody := r.hasTableColumn(ctx, "ai_request_responses", "raw_body")

	scanDest := []interface{}{&lastUser, &messages, &promptTrunc, &assistant, &responseTrunc}
	if hasToolCols {
		scanDest = append(scanDest, &toolCalls, &toolCallsTrunc)
	}
	if hasRawBody {
		scanDest = append(scanDest, &rawBody, &rawBodyTrunc)
	}

	if err := row.Scan(scanDest...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return out, nil
		}
		return domain.AIRequestContent{}, err
	}

	out.LastUserText = nullString(lastUser)
	out.MessagesJSON = nullString(messages)
	out.PromptTruncated = promptTrunc != 0
	out.AssistantText = nullString(assistant)
	out.ResponseTruncated = responseTrunc != 0
	out.ToolCallsJSON = nullString(toolCalls)
	out.ToolCallsTruncated = toolCallsTrunc != 0
	out.RawBody = nullString(rawBody)
	out.RawBodyTruncated = rawBodyTrunc != 0
	out.HasContent = requestContentPresent(out)

	if r.hasTable(ctx, "ai_request_errors") {
		errs, err := r.listCaptureErrors(ctx, requestID)
		if err != nil {
			return domain.AIRequestContent{}, err
		}
		out.CaptureErrors = mergeProxyErrorRows(proxyError, errs)
	} else {
		out.CaptureErrors = mergeProxyErrorRows(proxyError, nil)
	}

	return out, nil
}

func mergeProxyErrorRows(proxyError string, rows []domain.AIRequestCaptureError) []domain.AIRequestCaptureError {
	if proxyError == "" {
		return rows
	}
	for _, e := range rows {
		if e.Stage == "proxy" && e.ErrorMessage == proxyError {
			return rows
		}
	}
	merged := []domain.AIRequestCaptureError{{
		Stage:        "proxy",
		ErrorCode:    "proxy",
		ErrorMessage: proxyError,
	}}
	return append(merged, rows...)
}

func (r *AIRequests) listCaptureErrors(ctx context.Context, requestID string) ([]domain.AIRequestCaptureError, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT stage, error_code, error_message
		FROM ai_request_errors
		WHERE request_id = ?
		ORDER BY created_at ASC`, requestID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.AIRequestCaptureError
	for rows.Next() {
		var e domain.AIRequestCaptureError
		if err := rows.Scan(&e.Stage, &e.ErrorCode, &e.ErrorMessage); err != nil {
			return nil, err
		}
		out = append(out, e)
	}
	return out, rows.Err()
}

func (r *AIRequests) hasTable(ctx context.Context, name string) bool {
	var n int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = ?`, name).Scan(&n)
	return err == nil && n > 0
}

func (r *AIRequests) hasTableColumn(ctx context.Context, table, name string) bool {
	var n int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM pragma_table_info(?) WHERE name = ?`, table, name).Scan(&n)
	return err == nil && n > 0
}

func nullString(s sql.NullString) string {
	if s.Valid {
		return s.String
	}
	return ""
}

func requestContentPresent(c domain.AIRequestContent) bool {
	return strings.TrimSpace(c.LastUserText) != "" ||
		strings.TrimSpace(c.MessagesJSON) != "" ||
		strings.TrimSpace(c.AssistantText) != "" ||
		strings.TrimSpace(c.ToolCallsJSON) != "" ||
		strings.TrimSpace(c.RawBody) != ""
}
