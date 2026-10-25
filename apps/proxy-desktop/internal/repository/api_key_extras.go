package repository

import (
	"context"
	"database/sql"
	"time"
)

// DeleteAPIKey removes a key row (desktop-only; CLI exposes revoke only).
func DeleteAPIKey(ctx context.Context, db *sql.DB, id string) error {
	res, err := db.ExecContext(ctx, `DELETE FROM proxy_api_keys WHERE id = ?`, id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrAPIKeyNotFound
	}
	return nil
}

// SetAPIKeyAllowAllModels updates allow_all_models (desktop-only).
func SetAPIKeyAllowAllModels(ctx context.Context, db *sql.DB, id string, allowAll bool) error {
	now := time.Now().UTC()
	res, err := db.ExecContext(ctx, `
		UPDATE proxy_api_keys SET allow_all_models = ?, updated_at = ? WHERE id = ?`,
		boolInt(allowAll), formatTime(now), id)
	if err != nil {
		return err
	}
	if n, _ := res.RowsAffected(); n == 0 {
		return ErrAPIKeyNotFound
	}
	return nil
}
