package projects

import (
	"context"
	"database/sql"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
)

var ErrNotFound = errors.New("project not found")
var ErrDuplicate = errors.New("project name already exists")

type Project struct {
	ID          string
	Name        string
	Description string
	Enabled     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Store struct {
	db *sql.DB
}

func NewStore(db *sql.DB) *Store {
	return &Store{db: db}
}

func (s *Store) Create(ctx context.Context, name, description string) (Project, error) {
	now := time.Now().UTC()
	p := Project{
		ID:          uuid.NewString(),
		Name:        name,
		Description: description,
		Enabled:     true,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	_, err := s.db.ExecContext(ctx, `
		INSERT INTO projects (id, name, description, enabled, created_at, updated_at)
		VALUES (?, ?, ?, 1, ?, ?)`,
		p.ID, p.Name, p.Description, p.CreatedAt.Format(time.RFC3339Nano), p.UpdatedAt.Format(time.RFC3339Nano))
	if err != nil {
		if isUnique(err) {
			return Project{}, ErrDuplicate
		}
		return Project{}, err
	}
	return p, nil
}

func (s *Store) List(ctx context.Context) ([]Project, error) {
	rows, err := s.db.QueryContext(ctx, `
		SELECT id, name, description, enabled, created_at, updated_at FROM projects ORDER BY name`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var out []Project
	for rows.Next() {
		p, err := scanProject(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, p)
	}
	return out, rows.Err()
}

func (s *Store) GetByID(ctx context.Context, id string) (Project, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, name, description, enabled, created_at, updated_at FROM projects WHERE id = ?`, id)
	return scanProjectRow(row)
}

func (s *Store) GetByName(ctx context.Context, name string) (Project, error) {
	row := s.db.QueryRowContext(ctx, `
		SELECT id, name, description, enabled, created_at, updated_at FROM projects WHERE name = ?`, name)
	return scanProjectRow(row)
}

func (s *Store) Update(ctx context.Context, id, name, description string) (Project, error) {
	now := time.Now().UTC()
	res, err := s.db.ExecContext(ctx, `
		UPDATE projects SET name = ?, description = ?, updated_at = ? WHERE id = ?`,
		name, description, now.Format(time.RFC3339Nano), id)
	if err != nil {
		if isUnique(err) {
			return Project{}, ErrDuplicate
		}
		return Project{}, err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return Project{}, ErrNotFound
	}
	return s.GetByID(ctx, id)
}

func (s *Store) SetEnabled(ctx context.Context, id string, enabled bool) error {
	now := time.Now().UTC()
	res, err := s.db.ExecContext(ctx, `
		UPDATE projects SET enabled = ?, updated_at = ? WHERE id = ?`,
		boolInt(enabled), now.Format(time.RFC3339Nano), id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return ErrNotFound
	}
	return nil
}

func scanProject(rows *sql.Rows) (Project, error) {
	var p Project
	var enabled int
	var created, updated string
	if err := rows.Scan(&p.ID, &p.Name, &p.Description, &enabled, &created, &updated); err != nil {
		return Project{}, err
	}
	p.Enabled = enabled == 1
	p.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	p.UpdatedAt, _ = time.Parse(time.RFC3339Nano, updated)
	return p, nil
}

func scanProjectRow(row *sql.Row) (Project, error) {
	var p Project
	var enabled int
	var created, updated string
	if err := row.Scan(&p.ID, &p.Name, &p.Description, &enabled, &created, &updated); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return Project{}, ErrNotFound
		}
		return Project{}, err
	}
	p.Enabled = enabled == 1
	p.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	p.UpdatedAt, _ = time.Parse(time.RFC3339Nano, updated)
	return p, nil
}

func boolInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func isUnique(err error) bool {
	return err != nil && (strings.Contains(err.Error(), "UNIQUE") || strings.Contains(err.Error(), "unique"))
}
