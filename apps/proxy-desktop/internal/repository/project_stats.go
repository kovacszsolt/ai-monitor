package repository

import (
	"context"
	"database/sql"

	"github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/domain"
)

// ListProjectsWithStats returns projects with request counts from ai_requests (desktop-only query).
func ListProjectsWithStats(ctx context.Context, db *sql.DB) ([]domain.Project, error) {
	rows, err := db.QueryContext(ctx, `
		SELECT p.id, p.name, p.description, p.enabled, p.created_at, p.updated_at,
			COALESCE(s.request_count, 0),
			s.last_started_at
		FROM projects p
		LEFT JOIN (
			SELECT project_id, COUNT(*) AS request_count, MAX(started_at) AS last_started_at
			FROM ai_requests
			GROUP BY project_id
		) s ON s.project_id = p.id
		ORDER BY p.name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []domain.Project
	for rows.Next() {
		p, err := scanProjectWithStats(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func scanProjectWithStats(rows *sql.Rows) (domain.Project, error) {
	var p domain.Project
	var enabled int
	var created, updated string
	var lastStarted sql.NullString
	if err := rows.Scan(
		&p.ID, &p.Name, &p.Description, &enabled, &created, &updated,
		&p.RequestCount, &lastStarted,
	); err != nil {
		return domain.Project{}, err
	}
	p.Enabled = enabled == 1
	p.CreatedAt = parseTime(created)
	p.UpdatedAt = parseTime(updated)
	if lastStarted.Valid && lastStarted.String != "" {
		t := parseTime(lastStarted.String)
		p.LastRequestAt = &t
	}
	return p, nil
}
