package database

import (
	"context"
	"database/sql"
)

// MigrationVersion returns the golang-migrate version for this database, or 0 when unset.
func MigrationVersion(ctx context.Context, db *sql.DB) (uint, bool, error) {
	var version sql.NullInt64
	var dirty sql.NullInt64
	err := db.QueryRowContext(ctx, `SELECT version, dirty FROM schema_migrations LIMIT 1`).Scan(&version, &dirty)
	if err != nil {
		if err == sql.ErrNoRows {
			return 0, false, nil
		}
		// Table may not exist before first migration.
		return 0, false, nil
	}
	if !version.Valid {
		return 0, false, nil
	}
	return uint(version.Int64), dirty.Valid && dirty.Int64 == 1, nil
}
