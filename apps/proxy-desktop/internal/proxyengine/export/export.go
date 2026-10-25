package export

import (
	"context"
	"database/sql"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"strings"
	"time"
)

const SchemaVersion = "1.0"

type Filter struct {
	Since          *time.Time
	Until          *time.Time
	Project        string
	Provider       string
	Model          string
	State          string
	IncludeContent bool
}

type Record struct {
	SchemaVersion  string `json:"schema_version"`
	ID             string `json:"id"`
	ProjectID      string `json:"project_id"`
	ProxyAPIKeyID  string `json:"proxy_api_key_id"`
	ProviderID     string `json:"provider_id"`
	ModelRouteID   string `json:"model_route_id,omitempty"`
	RequestedModel string `json:"requested_model"`
	ActualModel    string `json:"actual_model,omitempty"`
	Endpoint       string `json:"endpoint"`
	Protocol       string `json:"protocol"`
	StartedAt      string `json:"started_at"`
	CompletedAt    string `json:"completed_at,omitempty"`
	DurationMS     *int64 `json:"duration_ms,omitempty"`
	StatusCode     *int   `json:"status_code,omitempty"`
	RequestState   string `json:"request_state"`
	Streaming      bool   `json:"streaming"`
	InputTokens    *int64 `json:"input_tokens,omitempty"`
	OutputTokens   *int64 `json:"output_tokens,omitempty"`
	UsageStatus    string `json:"usage_status"`
	RawUsageJSON   string `json:"raw_usage_json,omitempty"`
	RequestBytes   int64  `json:"request_bytes"`
	ResponseBytes  int64  `json:"response_bytes"`
	ErrorType      string `json:"error_type,omitempty"`
	Tag            string `json:"tag,omitempty"`
	LastUserText   string `json:"last_user_text,omitempty"`
	MessagesJSON   string `json:"messages_json,omitempty"`
	AssistantText      string `json:"assistant_text,omitempty"`
	ToolCallsJSON      string `json:"tool_calls_json,omitempty"`
	RawUpstreamBody    string `json:"raw_upstream_body,omitempty"`
}

func ExportJSONL(ctx context.Context, db *sql.DB, w io.Writer, f Filter) error {
	rows, err := query(ctx, db, f)
	if err != nil {
		return err
	}
	defer rows.Close()
	enc := json.NewEncoder(w)
	for rows.Next() {
		rec, err := scanExportRow(rows, f.IncludeContent)
		if err != nil {
			return err
		}
		if err := enc.Encode(rec); err != nil {
			return err
		}
	}
	return rows.Err()
}

func ExportCSV(ctx context.Context, db *sql.DB, w io.Writer, f Filter) error {
	rows, err := query(ctx, db, f)
	if err != nil {
		return err
	}
	defer rows.Close()
	cw := csv.NewWriter(w)
	_ = cw.Write([]string{
		"schema_version", "id", "project_id", "proxy_api_key_id", "provider_id", "model_route_id",
		"requested_model", "actual_model", "endpoint", "protocol", "started_at", "completed_at",
		"duration_ms", "status_code", "request_state", "streaming", "input_tokens", "output_tokens",
		"usage_status", "request_bytes", "response_bytes", "error_type", "tag",
	})
	for rows.Next() {
		rec, err := scanExportRow(rows, f.IncludeContent)
		if err != nil {
			return err
		}
		_ = cw.Write([]string{
			rec.SchemaVersion, rec.ID, rec.ProjectID, rec.ProxyAPIKeyID, rec.ProviderID, rec.ModelRouteID,
			rec.RequestedModel, rec.ActualModel, rec.Endpoint, rec.Protocol, rec.StartedAt, rec.CompletedAt,
			intStr(rec.DurationMS), intStrP(rec.StatusCode), rec.RequestState, fmt.Sprintf("%t", rec.Streaming),
			intStr(rec.InputTokens), intStr(rec.OutputTokens), rec.UsageStatus,
			fmt.Sprintf("%d", rec.RequestBytes), fmt.Sprintf("%d", rec.ResponseBytes), rec.ErrorType, rec.Tag,
		})
	}
	cw.Flush()
	return rows.Err()
}

func query(ctx context.Context, db *sql.DB, f Filter) (*sql.Rows, error) {
	var clauses []string
	var args []interface{}
	if f.Since != nil {
		clauses = append(clauses, "started_at >= ?")
		args = append(args, f.Since.UTC().Format(time.RFC3339Nano))
	}
	if f.Until != nil {
		clauses = append(clauses, "started_at <= ?")
		args = append(args, f.Until.UTC().Format(time.RFC3339Nano))
	}
	if f.Project != "" {
		clauses = append(clauses, "project_id = ?")
		args = append(args, f.Project)
	}
	if f.Provider != "" {
		clauses = append(clauses, "provider_id = ?")
		args = append(args, f.Provider)
	}
	if f.Model != "" {
		clauses = append(clauses, "requested_model = ?")
		args = append(args, f.Model)
	}
	if f.State != "" {
		clauses = append(clauses, "request_state = ?")
		args = append(args, f.State)
	}
	var q string
	if f.IncludeContent {
		q = `SELECT r.id, r.project_id, r.proxy_api_key_id, r.provider_id, r.model_route_id, r.requested_model, r.actual_model,
		r.endpoint, r.protocol, r.started_at, r.completed_at, r.duration_ms, r.status_code, r.request_state, r.streaming,
		r.input_tokens, r.output_tokens, r.usage_status, r.raw_usage_json, r.request_bytes, r.response_bytes, r.error_type, r.tag,
		p.last_user_text, p.messages_json, res.assistant_text, res.tool_calls_json, res.raw_body
		FROM ai_requests r
		LEFT JOIN ai_request_prompts p ON p.request_id = r.id
		LEFT JOIN ai_request_responses res ON res.request_id = r.id`
		if len(clauses) > 0 {
			for i := range clauses {
				clauses[i] = strings.Replace(clauses[i], "started_at", "r.started_at", 1)
				clauses[i] = strings.Replace(clauses[i], "project_id", "r.project_id", 1)
				clauses[i] = strings.Replace(clauses[i], "provider_id", "r.provider_id", 1)
				clauses[i] = strings.Replace(clauses[i], "requested_model", "r.requested_model", 1)
				clauses[i] = strings.Replace(clauses[i], "request_state", "r.request_state", 1)
			}
		}
	} else {
		q = `SELECT id, project_id, proxy_api_key_id, provider_id, model_route_id, requested_model, actual_model,
		endpoint, protocol, started_at, completed_at, duration_ms, status_code, request_state, streaming,
		input_tokens, output_tokens, usage_status, raw_usage_json, request_bytes, response_bytes, error_type, tag
		FROM ai_requests`
	}
	if len(clauses) > 0 {
		q += " WHERE " + strings.Join(clauses, " AND ")
	}
	if f.IncludeContent {
		q += " ORDER BY r.started_at ASC"
	} else {
		q += " ORDER BY started_at ASC"
	}
	return db.QueryContext(ctx, q, args...)
}

func scanExport(rows *sql.Rows) (Record, error) {
	return scanExportRow(rows, false)
}

func scanExportRow(rows *sql.Rows, withContent bool) (Record, error) {
	var rec Record
	var modelRoute, actualModel, completed, rawUsage, errType sql.NullString
	var duration sql.NullInt64
	var status sql.NullInt64
	var streaming int
	var inTok, outTok sql.NullInt64
	var tag sql.NullString
	var lastUser, messagesJSON, assistant, toolCalls, rawBody sql.NullString
	if withContent {
		if err := rows.Scan(
			&rec.ID, &rec.ProjectID, &rec.ProxyAPIKeyID, &rec.ProviderID, &modelRoute,
			&rec.RequestedModel, &actualModel, &rec.Endpoint, &rec.Protocol, &rec.StartedAt, &completed,
			&duration, &status, &rec.RequestState, &streaming, &inTok, &outTok, &rec.UsageStatus, &rawUsage,
			&rec.RequestBytes, &rec.ResponseBytes, &errType, &tag,
			&lastUser, &messagesJSON, &assistant, &toolCalls, &rawBody,
		); err != nil {
			return Record{}, err
		}
	} else if err := rows.Scan(
		&rec.ID, &rec.ProjectID, &rec.ProxyAPIKeyID, &rec.ProviderID, &modelRoute,
		&rec.RequestedModel, &actualModel, &rec.Endpoint, &rec.Protocol, &rec.StartedAt, &completed,
		&duration, &status, &rec.RequestState, &streaming, &inTok, &outTok, &rec.UsageStatus, &rawUsage,
		&rec.RequestBytes, &rec.ResponseBytes, &errType, &tag,
	); err != nil {
		return Record{}, err
	}
	rec.SchemaVersion = SchemaVersion
	rec.ModelRouteID = modelRoute.String
	rec.ActualModel = actualModel.String
	rec.CompletedAt = completed.String
	if duration.Valid {
		rec.DurationMS = &duration.Int64
	}
	if status.Valid {
		v := int(status.Int64)
		rec.StatusCode = &v
	}
	rec.Streaming = streaming == 1
	if inTok.Valid {
		rec.InputTokens = &inTok.Int64
	}
	if outTok.Valid {
		rec.OutputTokens = &outTok.Int64
	}
	rec.RawUsageJSON = rawUsage.String
	rec.ErrorType = errType.String
	rec.Tag = tag.String
	if withContent {
		rec.LastUserText = lastUser.String
		rec.MessagesJSON = messagesJSON.String
		rec.AssistantText = assistant.String
		rec.ToolCallsJSON = toolCalls.String
		rec.RawUpstreamBody = rawBody.String
	}
	return rec, nil
}

func intStr(p *int64) string {
	if p == nil {
		return ""
	}
	return fmt.Sprintf("%d", *p)
}

func intStrP(p *int) string {
	if p == nil {
		return ""
	}
	return fmt.Sprintf("%d", *p)
}

func StorageStats(ctx context.Context, db *sql.DB) (total int64, err error) {
	row := db.QueryRowContext(ctx, `SELECT COUNT(*) FROM ai_requests`)
	if err := row.Scan(&total); err != nil {
		return 0, err
	}
	return total, nil
}
