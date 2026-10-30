package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/domain"
)

var (
	ErrAIRequestNotFound        = errors.New("request not found")
	ErrRequestTagColumnMissing  = errors.New("tag column not available; run migration 000008 (Settings → Database maintenance)")
	ErrNoRequestIDs             = errors.New("no request ids provided")
)

const maxRequestTagLen = 256

// NoTagLabel is the filter value for rows with an empty proxy tag (matches dashboard grouping).
const NoTagLabel = "(no tag)"

type AIRequestListFilter struct {
	ProjectID      string
	State          string
	Tag            string
	ModelRouteID   string
	ProxyAPIKeyID  string
	SortBy         string
	SortDesc       bool
}

// ModelRouteFilterOption is a distinct model route seen in request logs.
type ModelRouteFilterOption struct {
	ID   string
	Name string
}

const modelRouteJoinSQL = `LEFT JOIN model_routes mr ON mr.id = r.model_route_id`

const apiKeyJoinSQL = `LEFT JOIN proxy_api_keys ak ON ak.id = r.proxy_api_key_id`

const modelRouteDisplayNameSQL = `COALESCE(NULLIF(TRIM(mr.public_model_name), ''), '(unknown route)')`

const apiKeyDisplayNameSQL = `COALESCE(NULLIF(TRIM(ak.name), ''), '(deleted key)')`

type AIRequests struct {
	db *sql.DB
}

func NewAIRequests(db *sql.DB) *AIRequests {
	return &AIRequests{db: db}
}

func normalizeStoredRequestTag(tag string) string {
	tag = strings.TrimSpace(tag)
	if len(tag) > maxRequestTagLen {
		return tag[:maxRequestTagLen]
	}
	return tag
}

func (r *AIRequests) BulkSetRequestTags(ctx context.Context, ids []string, tag string) (int64, error) {
	if !r.hasAIRequestsColumn(ctx, "tag") {
		return 0, ErrRequestTagColumnMissing
	}
	cleanIDs := make([]string, 0, len(ids))
	for _, id := range ids {
		if id != "" {
			cleanIDs = append(cleanIDs, id)
		}
	}
	if len(cleanIDs) == 0 {
		return 0, ErrNoRequestIDs
	}
	tag = normalizeStoredRequestTag(tag)
	placeholders := strings.Repeat("?,", len(cleanIDs))
	placeholders = placeholders[:len(placeholders)-1]
	args := make([]interface{}, 0, len(cleanIDs)+1)
	args = append(args, tag)
	for _, id := range cleanIDs {
		args = append(args, id)
	}
	q := fmt.Sprintf(`UPDATE ai_requests SET tag = ? WHERE id IN (%s)`, placeholders)
	res, err := r.db.ExecContext(ctx, q, args...)
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return n, nil
}

func (r *AIRequests) Delete(ctx context.Context, id string) error {
	if id == "" {
		return ErrAIRequestNotFound
	}
	res, err := r.db.ExecContext(ctx, `DELETE FROM ai_requests WHERE id = ?`, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrAIRequestNotFound
	}
	return nil
}

func (r *AIRequests) BulkDelete(ctx context.Context, ids []string) (int64, error) {
	cleanIDs := make([]string, 0, len(ids))
	for _, id := range ids {
		if id != "" {
			cleanIDs = append(cleanIDs, id)
		}
	}
	if len(cleanIDs) == 0 {
		return 0, ErrNoRequestIDs
	}
	placeholders := strings.Repeat("?,", len(cleanIDs))
	placeholders = placeholders[:len(placeholders)-1]
	args := make([]interface{}, len(cleanIDs))
	for i, id := range cleanIDs {
		args[i] = id
	}
	q := fmt.Sprintf(`DELETE FROM ai_requests WHERE id IN (%s)`, placeholders)
	res, err := r.db.ExecContext(ctx, q, args...)
	if err != nil {
		return 0, err
	}
	n, _ := res.RowsAffected()
	return n, nil
}

func (r *AIRequests) List(ctx context.Context, f AIRequestListFilter, limit, offset int) ([]domain.AIRequest, int64, error) {
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}
	if offset < 0 {
		offset = 0
	}
	if f.Tag != "" && !r.hasAIRequestsColumn(ctx, "tag") {
		f.Tag = ""
	}

	where, args := buildRequestFilters(f)
	countQ := `SELECT COUNT(*) FROM ai_requests r` + where
	var total int64
	if err := r.db.QueryRowContext(ctx, countQ, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	tagExpr := "''"
	if r.hasAIRequestsColumn(ctx, "tag") {
		tagExpr = "COALESCE(r.tag, '')"
	}
	q := `
		SELECT r.id, r.project_id, COALESCE(p.name, ''), r.proxy_api_key_id, COALESCE(ak.name, ''),
			r.provider_id, COALESCE(pr.name, ''),
			r.requested_model, COALESCE(r.actual_model, ''), r.endpoint, r.protocol,
			r.started_at, r.completed_at, r.duration_ms, r.ttft_ms, r.status_code, r.request_state,
			r.streaming, r.input_tokens, r.output_tokens, r.cached_input_tokens,
			r.uncached_input_tokens, r.prompt_tps, r.generation_tps, r.effective_tps, r.cache_hit_rate,
			r.generation_speed_estimated, r.prompt_speed_estimated,
			r.usage_status, r.request_bytes, r.response_bytes, COALESCE(r.error_type, ''), ` + tagExpr + `
		FROM ai_requests r
		LEFT JOIN projects p ON p.id = r.project_id
		` + apiKeyJoinSQL + `
		LEFT JOIN providers pr ON pr.id = r.provider_id
		` + modelRouteJoinSQL + where + `
		ORDER BY ` + requestListOrderBy(f, r.hasAIRequestsColumn(ctx, "tag")) + `
		LIMIT ? OFFSET ?`
	listArgs := append(args, limit, offset)
	rows, err := r.db.QueryContext(ctx, q, listArgs...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var out []domain.AIRequest
	for rows.Next() {
		req, err := scanAIRequest(rows)
		if err != nil {
			return nil, 0, err
		}
		out = append(out, req)
	}
	return out, total, rows.Err()
}

func requestListOrderBy(f AIRequestListFilter, hasTagColumn bool) string {
	cols := map[string]string{
		"started":  "r.started_at",
		"project":  "p.name",
		"api_key":  apiKeyDisplayNameSQL,
		"model":    modelRouteDisplayNameSQL,
		"state":    "r.request_state",
		"input":    "r.input_tokens",
		"cached":   "r.cached_input_tokens",
		"output":   "r.output_tokens",
		"ttft":     "r.ttft_ms",
		"duration": "r.duration_ms",
		"gen":      "r.generation_tps",
		"prompt":   "r.prompt_tps",
		"eff":      "r.effective_tps",
		"cache":    "r.cache_hit_rate",
	}
	sortBy := strings.TrimSpace(f.SortBy)
	expr := cols[sortBy]
	if sortBy == "tag" && hasTagColumn {
		expr = "r.tag"
	}
	if expr == "" {
		expr = "r.started_at"
	}
	dir := "ASC"
	if f.SortDesc {
		dir = "DESC"
	}
	// Stable tie-breaker for equal sort keys.
	return expr + " " + dir + ", r.id " + dir
}

func buildRequestFilters(f AIRequestListFilter) (string, []interface{}) {
	var clauses []string
	var args []interface{}
	if f.ProjectID != "" {
		clauses = append(clauses, "r.project_id = ?")
		args = append(args, f.ProjectID)
	}
	if f.State != "" {
		clauses = append(clauses, "r.request_state = ?")
		args = append(args, f.State)
	}
	if f.Tag != "" {
		if f.Tag == NoTagLabel {
			clauses = append(clauses, "TRIM(COALESCE(r.tag, '')) = ''")
		} else {
			clauses = append(clauses, "r.tag = ?")
			args = append(args, f.Tag)
		}
	}
	if f.ModelRouteID != "" {
		clauses = append(clauses, "r.model_route_id = ?")
		args = append(args, f.ModelRouteID)
	}
	if f.ProxyAPIKeyID != "" {
		clauses = append(clauses, "r.proxy_api_key_id = ?")
		args = append(args, f.ProxyAPIKeyID)
	}
	if len(clauses) == 0 {
		return "", args
	}
	return " WHERE " + strings.Join(clauses, " AND "), args
}

func (r *AIRequests) Dashboard(ctx context.Context) (domain.DashboardStats, error) {
	stats := domain.DashboardStats{}
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM ai_requests`).Scan(&stats.TotalRequests); err != nil {
		return domain.DashboardStats{}, err
	}

	since24h := time.Now().UTC().Add(-24 * time.Hour).Format(time.RFC3339Nano)
	if err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM ai_requests WHERE started_at >= ?`, since24h).
		Scan(&stats.RequestsLast24h); err != nil {
		return domain.DashboardStats{}, err
	}

	row := r.db.QueryRowContext(ctx, `
		SELECT
			COALESCE(SUM(input_tokens), 0),
			COALESCE(SUM(output_tokens), 0),
			COALESCE(SUM(cached_input_tokens), 0),
			COALESCE(AVG(duration_ms), 0),
			COALESCE(SUM(CASE WHEN streaming = 1 THEN 1 ELSE 0 END), 0)
		FROM ai_requests`)
	if err := row.Scan(&stats.SumInputTokens, &stats.SumOutputTokens, &stats.SumCachedInTokens, &stats.AvgDurationMS, &stats.StreamingCount); err != nil {
		return domain.DashboardStats{}, err
	}

	rows, err := r.db.QueryContext(ctx, `
		SELECT request_state, COUNT(*) FROM ai_requests GROUP BY request_state ORDER BY COUNT(*) DESC`)
	if err != nil {
		return domain.DashboardStats{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var st domain.RequestStateCount
		if err := rows.Scan(&st.State, &st.Count); err != nil {
			return domain.DashboardStats{}, err
		}
		stats.ByState = append(stats.ByState, st)
	}
	if err := rows.Err(); err != nil {
		return domain.DashboardStats{}, err
	}

	stats.TopModels, err = r.topNamed(ctx, `
		SELECT `+modelRouteDisplayNameSQL+`, COUNT(*)
		FROM ai_requests r
		`+modelRouteJoinSQL+`
		WHERE TRIM(COALESCE(r.model_route_id, '')) != ''
		GROUP BY r.model_route_id
		ORDER BY COUNT(*) DESC LIMIT 5`)
	if err != nil {
		return domain.DashboardStats{}, err
	}

	stats.TopProjects, err = r.topNamed(ctx, `
		SELECT COALESCE(p.name, r.project_id), COUNT(*)
		FROM ai_requests r
		LEFT JOIN projects p ON p.id = r.project_id
		GROUP BY r.project_id ORDER BY COUNT(*) DESC LIMIT 5`)
	if err != nil {
		return domain.DashboardStats{}, err
	}

	if r.hasAIRequestsColumn(ctx, "tag") {
		stats.UsageByTag, err = r.usageGroupStats(ctx, `
			SELECT COALESCE(NULLIF(TRIM(r.tag), ''), '(no tag)'),
				COUNT(*),
				COALESCE(SUM(COALESCE(r.input_tokens, 0)), 0),
				COALESCE(SUM(COALESCE(r.output_tokens, 0)), 0),
				COALESCE(SUM(COALESCE(r.cached_input_tokens, 0)), 0)
			FROM ai_requests r
			GROUP BY COALESCE(NULLIF(TRIM(r.tag), ''), '(no tag)')
			ORDER BY COUNT(*) DESC
			LIMIT 50`)
		if err != nil {
			return domain.DashboardStats{}, err
		}
	}

	return stats, nil
}

func (r *AIRequests) ModelUsageStats(ctx context.Context, tag, projectID, apiKeyID string) ([]domain.ModelUsageStats, error) {
	if tag != "" && !r.hasAIRequestsColumn(ctx, "tag") {
		tag = ""
	}
	filterSQL, filterArgs := analyticsFilterSQL(tag, projectID, apiKeyID)
	usage, err := r.modelRouteUsageGroupStatsWithArgs(ctx, `
		SELECT COALESCE(r.model_route_id, ''),
			`+modelRouteDisplayNameSQL+`,
			COUNT(*),
			COALESCE(SUM(COALESCE(r.input_tokens, 0)), 0),
			COALESCE(SUM(COALESCE(r.output_tokens, 0)), 0),
			COALESCE(SUM(COALESCE(r.cached_input_tokens, 0)), 0)
		FROM ai_requests r
		`+modelRouteJoinSQL+`
		WHERE TRIM(COALESCE(r.model_route_id, '')) != ''`+filterSQL+`
		GROUP BY r.model_route_id
		ORDER BY COUNT(*) DESC
		LIMIT 50`, filterArgs...)
	if err != nil {
		return nil, err
	}
	windows, err := r.modelRoutePromptWindowWithArgs(ctx, `
		SELECT COALESCE(r.model_route_id, ''),
			MIN(r.started_at),
			MAX(r.completed_at)
		FROM ai_requests r
		`+modelRouteJoinSQL+`
		WHERE TRIM(COALESCE(r.model_route_id, '')) != ''
			AND r.completed_at IS NOT NULL AND TRIM(r.completed_at) != ''`+filterSQL+`
		GROUP BY r.model_route_id
		LIMIT 50`, filterArgs...)
	if err != nil {
		return nil, err
	}
	throughput, err := r.modelRouteTokenThroughputWithArgs(ctx, `
		SELECT COALESCE(r.model_route_id, ''),
			`+modelRouteDisplayNameSQL+`,
			COALESCE(SUM(COALESCE(r.input_tokens, 0)), 0),
			COALESCE(SUM(COALESCE(r.output_tokens, 0)), 0),
			COALESCE(SUM(COALESCE(r.cached_input_tokens, 0)), 0),
			COALESCE(SUM(COALESCE(r.duration_ms, 0)), 0),
			COALESCE(SUM(COALESCE(r.decode_duration_ms, 0)), 0)
		FROM ai_requests r
		`+modelRouteJoinSQL+`
		WHERE TRIM(COALESCE(r.model_route_id, '')) != ''
			AND r.duration_ms IS NOT NULL AND r.duration_ms > 0`+filterSQL+`
		GROUP BY r.model_route_id
		HAVING SUM(r.duration_ms) > 0
		LIMIT 50`, filterArgs...)
	if err != nil {
		return nil, err
	}
	return mergeModelUsage(usage, throughput, windows), nil
}

func tagFilterAnd(tag string) (string, []interface{}) {
	if tag == "" {
		return "", nil
	}
	if tag == NoTagLabel {
		return " AND TRIM(COALESCE(r.tag, '')) = ''", nil
	}
	return " AND r.tag = ?", []interface{}{tag}
}

func analyticsFilterSQL(tag, projectID, apiKeyID string) (string, []interface{}) {
	var parts []string
	var args []interface{}
	if projectID != "" {
		parts = append(parts, "r.project_id = ?")
		args = append(args, projectID)
	}
	if apiKeyID != "" {
		parts = append(parts, "r.proxy_api_key_id = ?")
		args = append(args, apiKeyID)
	}
	tagSQL, tagArgs := tagFilterAnd(tag)
	if tagSQL != "" {
		parts = append(parts, strings.TrimPrefix(tagSQL, " AND "))
		args = append(args, tagArgs...)
	}
	if len(parts) == 0 {
		return "", nil
	}
	return " AND " + strings.Join(parts, " AND "), args
}

// HourlyRequestCountByModel returns time-bucketed aggregates per model for a rolling day window or custom UTC date range (YYYY-MM-DD).
func (r *AIRequests) HourlyRequestCountByModel(ctx context.Context, tag, projectID, apiKeyID string, days int, rangeStart, rangeEnd, bucketUnit string) ([]domain.HourlyModelRequestCount, error) {
	if tag != "" && !r.hasAIRequestsColumn(ctx, "tag") {
		tag = ""
	}
	bounds, err := resolveHourlyChartBounds(days, rangeStart, rangeEnd)
	if err != nil {
		return nil, err
	}
	unit := NormalizeChartBucketUnit(bucketUnit)
	if err := validateChartBucketForBounds(unit, bounds); err != nil {
		return nil, err
	}
	bucketSQL := chartBucketStartSQL(unit)
	filterSQL, filterArgs := analyticsFilterSQL(tag, projectID, apiKeyID)
	args := []interface{}{bounds.Since, bounds.Until}
	args = append(args, filterArgs...)

	q := `
		SELECT ` + bucketSQL + `,
			COALESCE(r.model_route_id, ''),
			` + modelRouteDisplayNameSQL + `,
			COUNT(*),
			COALESCE(SUM(COALESCE(r.input_tokens, 0)), 0),
			COALESCE(SUM(COALESCE(r.cached_input_tokens, 0)), 0),
			COALESCE(SUM(COALESCE(r.output_tokens, 0)), 0),
			COALESCE(SUM(COALESCE(r.decode_duration_ms, 0)), 0)
		FROM ai_requests r
		` + modelRouteJoinSQL + `
		WHERE TRIM(COALESCE(r.model_route_id, '')) != ''
			AND r.started_at >= ? AND r.started_at < ?` + filterSQL + `
		GROUP BY 1, r.model_route_id
		ORDER BY 1 ASC, 4 DESC`

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.HourlyModelRequestCount
	for rows.Next() {
		var p domain.HourlyModelRequestCount
		if err := rows.Scan(
			&p.BucketStart,
			&p.ModelRouteID,
			&p.Model,
			&p.RequestCount,
			&p.InputTokens,
			&p.CachedInputTokens,
			&p.OutputTokens,
			&p.DecodeDurationMS,
		); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

// ModelTokenUsageOverTime returns hourly token sums per model for the last days, filtered by tag.
func (r *AIRequests) ModelTokenUsageOverTime(ctx context.Context, tag string, days int) ([]domain.ModelTokenTimePoint, error) {
	if tag == "" {
		return nil, nil
	}
	if !r.hasAIRequestsColumn(ctx, "tag") {
		return nil, nil
	}
	if days <= 0 {
		days = 14
	}
	if days > 90 {
		days = 90
	}
	since := time.Now().UTC().AddDate(0, 0, -days).Format(time.RFC3339Nano)
	tagSQL, tagArgs := tagFilterAnd(tag)
	args := []interface{}{since}
	args = append(args, tagArgs...)

	q := `
		SELECT substr(r.started_at, 1, 13) || ':00:00Z',
			COALESCE(r.model_route_id, ''),
			` + modelRouteDisplayNameSQL + `,
			COALESCE(SUM(COALESCE(r.input_tokens, 0)), 0),
			COALESCE(SUM(COALESCE(r.output_tokens, 0)), 0)
		FROM ai_requests r
		` + modelRouteJoinSQL + `
		WHERE TRIM(COALESCE(r.model_route_id, '')) != ''
			AND r.started_at >= ?` + tagSQL + `
		GROUP BY substr(r.started_at, 1, 13), r.model_route_id
		ORDER BY 1 ASC, 3 ASC`

	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.ModelTokenTimePoint
	for rows.Next() {
		var p domain.ModelTokenTimePoint
		if err := rows.Scan(&p.BucketStart, &p.ModelRouteID, &p.Model, &p.InputTokens, &p.OutputTokens); err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (r *AIRequests) ListAPIKeyFilterOptions(ctx context.Context, projectID string) ([]ModelRouteFilterOption, error) {
	if projectID != "" {
		rows, err := r.db.QueryContext(ctx, `
			SELECT id, COALESCE(NULLIF(TRIM(name), ''), id)
			FROM proxy_api_keys
			WHERE project_id = ?
			ORDER BY 2 ASC`, projectID)
		if err != nil {
			return nil, err
		}
		defer rows.Close()
		var out []ModelRouteFilterOption
		for rows.Next() {
			var opt ModelRouteFilterOption
			if err := rows.Scan(&opt.ID, &opt.Name); err != nil {
				return nil, err
			}
			out = append(out, opt)
		}
		return out, rows.Err()
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT DISTINCT r.proxy_api_key_id, `+apiKeyDisplayNameSQL+`
		FROM ai_requests r
		`+apiKeyJoinSQL+`
		WHERE TRIM(COALESCE(r.proxy_api_key_id, '')) != ''
		ORDER BY 2 ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ModelRouteFilterOption
	for rows.Next() {
		var opt ModelRouteFilterOption
		if err := rows.Scan(&opt.ID, &opt.Name); err != nil {
			return nil, err
		}
		out = append(out, opt)
	}
	return out, rows.Err()
}

func (r *AIRequests) ListModelRouteFilterOptions(ctx context.Context) ([]ModelRouteFilterOption, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT DISTINCT r.model_route_id, `+modelRouteDisplayNameSQL+`
		FROM ai_requests r
		`+modelRouteJoinSQL+`
		WHERE TRIM(COALESCE(r.model_route_id, '')) != ''
		ORDER BY 2 ASC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []ModelRouteFilterOption
	for rows.Next() {
		var opt ModelRouteFilterOption
		if err := rows.Scan(&opt.ID, &opt.Name); err != nil {
			return nil, err
		}
		out = append(out, opt)
	}
	return out, rows.Err()
}

func (r *AIRequests) ListTagFilterOptions(ctx context.Context) (enabled bool, tags []string, hasUntagged bool, err error) {
	if !r.hasAIRequestsColumn(ctx, "tag") {
		return false, nil, false, nil
	}
	rows, err := r.db.QueryContext(ctx, `
		SELECT DISTINCT TRIM(tag) FROM ai_requests
		WHERE TRIM(COALESCE(tag, '')) != ''
		ORDER BY TRIM(tag) ASC`)
	if err != nil {
		return false, nil, false, err
	}
	defer rows.Close()
	for rows.Next() {
		var t string
		if err := rows.Scan(&t); err != nil {
			return false, nil, false, err
		}
		tags = append(tags, t)
	}
	if err := rows.Err(); err != nil {
		return false, nil, false, err
	}
	var n int64
	if err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM ai_requests WHERE TRIM(COALESCE(tag, '')) = ''`).Scan(&n); err != nil {
		return false, nil, false, err
	}
	return true, tags, n > 0, nil
}

func (r *AIRequests) hasAIRequestsColumn(ctx context.Context, name string) bool {
	var n int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM pragma_table_info('ai_requests') WHERE name = ?`, name).Scan(&n)
	return err == nil && n > 0
}

func (r *AIRequests) usageGroupStats(ctx context.Context, q string) ([]domain.UsageGroupStats, error) {
	return r.usageGroupStatsWithArgs(ctx, q)
}

func (r *AIRequests) usageGroupStatsWithArgs(ctx context.Context, q string, args ...interface{}) ([]domain.UsageGroupStats, error) {
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.UsageGroupStats
	for rows.Next() {
		var g domain.UsageGroupStats
		if err := rows.Scan(&g.Name, &g.RequestCount, &g.InputTokens, &g.OutputTokens, &g.CachedInputTokens); err != nil {
			return nil, err
		}
		out = append(out, g)
	}
	return out, rows.Err()
}

func (r *AIRequests) modelRouteUsageGroupStatsWithArgs(ctx context.Context, q string, args ...interface{}) ([]domain.UsageGroupStats, error) {
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.UsageGroupStats
	for rows.Next() {
		var g domain.UsageGroupStats
		if err := rows.Scan(
			&g.ModelRouteID,
			&g.Name,
			&g.RequestCount,
			&g.InputTokens,
			&g.OutputTokens,
			&g.CachedInputTokens,
		); err != nil {
			return nil, err
		}
		out = append(out, g)
	}
	return out, rows.Err()
}

func mergeModelUsage(usage []domain.UsageGroupStats, throughput []domain.TokenThroughput, windows []domain.ModelPromptWindow) []domain.ModelUsageStats {
	tpByRouteID := make(map[string]domain.TokenThroughput, len(throughput))
	for _, t := range throughput {
		tpByRouteID[t.ModelRouteID] = t
	}
	winByRouteID := make(map[string]domain.ModelPromptWindow, len(windows))
	for _, w := range windows {
		winByRouteID[w.ModelRouteID] = w
	}
	out := make([]domain.ModelUsageStats, 0, len(usage))
	for _, u := range usage {
		m := domain.ModelUsageStats{UsageGroupStats: u}
		if t, ok := tpByRouteID[u.ModelRouteID]; ok {
			m.EffectiveTPS = t.EffectiveTPS
			m.GenerationTPS = t.GenerationTPS
			m.HasGenerationTPS = t.HasGenerationTPS
		}
		if w, ok := winByRouteID[u.ModelRouteID]; ok {
			m.PromptWindowMS = w.PromptWindowMS
			m.HasPromptWindow = w.HasPromptWindow
		}
		out = append(out, m)
	}
	return out
}

func (r *AIRequests) modelRoutePromptWindowWithArgs(ctx context.Context, q string, args ...interface{}) ([]domain.ModelPromptWindow, error) {
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.ModelPromptWindow
	for rows.Next() {
		var routeID string
		var firstStarted, lastCompleted sql.NullString
		if err := rows.Scan(&routeID, &firstStarted, &lastCompleted); err != nil {
			return nil, err
		}
		ms, ok := promptWindowMS(firstStarted, lastCompleted)
		out = append(out, domain.ModelPromptWindow{
			ModelRouteID:    routeID,
			PromptWindowMS:  ms,
			HasPromptWindow: ok,
		})
	}
	return out, rows.Err()
}

func promptWindowMS(firstStarted, lastCompleted sql.NullString) (int64, bool) {
	if !firstStarted.Valid || !lastCompleted.Valid {
		return 0, false
	}
	start, err := parseStoredTime(firstStarted.String)
	if err != nil {
		return 0, false
	}
	end, err := parseStoredTime(lastCompleted.String)
	if err != nil {
		return 0, false
	}
	if !end.After(start) {
		return 0, false
	}
	return end.Sub(start).Milliseconds(), true
}

func parseStoredTime(s string) (time.Time, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}, fmt.Errorf("empty time")
	}
	t, err := time.Parse(time.RFC3339Nano, s)
	if err == nil {
		return t, nil
	}
	return time.Parse(time.RFC3339, s)
}

func (r *AIRequests) tokenThroughput(ctx context.Context, q string) ([]domain.TokenThroughput, error) {
	return r.tokenThroughputWithArgs(ctx, q)
}

func (r *AIRequests) modelRouteTokenThroughputWithArgs(ctx context.Context, q string, args ...interface{}) ([]domain.TokenThroughput, error) {
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.TokenThroughput
	for rows.Next() {
		var t domain.TokenThroughput
		if err := rows.Scan(
			&t.ModelRouteID,
			&t.Name,
			&t.InputTokens,
			&t.OutputTokens,
			&t.CachedInputTokens,
			&t.DurationMS,
			&t.DecodeDurationMS,
		); err != nil {
			return nil, err
		}
		if t.DurationMS > 0 {
			t.EffectiveTPS = float64(t.InputTokens+t.OutputTokens) / (float64(t.DurationMS) / 1000.0)
		}
		if t.DecodeDurationMS > 0 && t.OutputTokens > 0 {
			t.GenerationTPS = float64(t.OutputTokens) / (float64(t.DecodeDurationMS) / 1000.0)
			t.HasGenerationTPS = true
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *AIRequests) tokenThroughputWithArgs(ctx context.Context, q string, args ...interface{}) ([]domain.TokenThroughput, error) {
	rows, err := r.db.QueryContext(ctx, q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.TokenThroughput
	for rows.Next() {
		var t domain.TokenThroughput
		if err := rows.Scan(&t.Name, &t.InputTokens, &t.OutputTokens, &t.CachedInputTokens, &t.DurationMS, &t.DecodeDurationMS); err != nil {
			return nil, err
		}
		if t.DurationMS > 0 {
			t.EffectiveTPS = float64(t.InputTokens+t.OutputTokens) / (float64(t.DurationMS) / 1000.0)
		}
		if t.DecodeDurationMS > 0 && t.OutputTokens > 0 {
			t.GenerationTPS = float64(t.OutputTokens) / (float64(t.DecodeDurationMS) / 1000.0)
			t.HasGenerationTPS = true
		}
		out = append(out, t)
	}
	return out, rows.Err()
}

func (r *AIRequests) topNamed(ctx context.Context, q string) ([]domain.NamedCount, error) {
	rows, err := r.db.QueryContext(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.NamedCount
	for rows.Next() {
		var nc domain.NamedCount
		if err := rows.Scan(&nc.Name, &nc.Count); err != nil {
			return nil, err
		}
		out = append(out, nc)
	}
	return out, rows.Err()
}

func scanAIRequest(rows *sql.Rows) (domain.AIRequest, error) {
	var req domain.AIRequest
	var started, completed sql.NullString
	var duration sql.NullInt64
	var status sql.NullInt64
	var streaming int
	var inTok, outTok, cachedIn, uncached, ttft sql.NullInt64
	var promptTPS, genTPS, effTPS, cacheRate sql.NullFloat64
	var genEst, promptEst int
	if err := rows.Scan(
		&req.ID, &req.ProjectID, &req.ProjectName, &req.ProxyAPIKeyID, &req.ProxyAPIKeyName,
		&req.ProviderID, &req.ProviderName,
		&req.RequestedModel, &req.ActualModel, &req.Endpoint, &req.Protocol,
		&started, &completed, &duration, &ttft, &status, &req.RequestState,
		&streaming, &inTok, &outTok, &cachedIn, &uncached,
		&promptTPS, &genTPS, &effTPS, &cacheRate, &genEst, &promptEst,
		&req.UsageStatus, &req.RequestBytes, &req.ResponseBytes, &req.ErrorType, &req.Tag,
	); err != nil {
		return domain.AIRequest{}, err
	}
	if started.Valid {
		req.StartedAt = parseTime(started.String)
	}
	if completed.Valid && completed.String != "" {
		t := parseTime(completed.String)
		req.CompletedAt = &t
	}
	if duration.Valid {
		req.DurationMS = &duration.Int64
	}
	if status.Valid {
		v := int(status.Int64)
		req.StatusCode = &v
	}
	req.Streaming = streaming == 1
	if inTok.Valid {
		req.InputTokens = &inTok.Int64
	}
	if outTok.Valid {
		req.OutputTokens = &outTok.Int64
	}
	if cachedIn.Valid {
		req.CachedInputTokens = &cachedIn.Int64
	}
	if uncached.Valid {
		req.UncachedInputTokens = &uncached.Int64
	}
	if ttft.Valid {
		req.TTFTMS = &ttft.Int64
	}
	if promptTPS.Valid {
		req.PromptTPS = &promptTPS.Float64
	}
	if genTPS.Valid {
		req.GenerationTPS = &genTPS.Float64
	}
	if effTPS.Valid {
		req.EffectiveTPS = &effTPS.Float64
	}
	if cacheRate.Valid {
		req.CacheHitRate = &cacheRate.Float64
	}
	req.GenerationSpeedEstimated = genEst == 1
	req.PromptSpeedEstimated = promptEst == 1
	ensureTokenPointersForUsage(&req)
	return req, nil
}

func ensureTokenPointersForUsage(req *domain.AIRequest) {
	if req.UsageStatus != "ok" {
		return
	}
	if req.InputTokens == nil {
		z := int64(0)
		req.InputTokens = &z
	}
	if req.OutputTokens == nil {
		z := int64(0)
		req.OutputTokens = &z
	}
	if req.CachedInputTokens == nil {
		z := int64(0)
		req.CachedInputTokens = &z
	}
}

// TableExists reports whether ai_requests is present (older DBs may lack usage data only).
func (r *AIRequests) TableExists(ctx context.Context) bool {
	var n int
	err := r.db.QueryRowContext(ctx, `
		SELECT COUNT(*) FROM sqlite_master WHERE type = 'table' AND name = 'ai_requests'`).Scan(&n)
	return err == nil && n == 1
}

func (r *AIRequests) PingTable(ctx context.Context) error {
	if !r.TableExists(ctx) {
		return fmt.Errorf("table ai_requests not found")
	}
	return nil
}
