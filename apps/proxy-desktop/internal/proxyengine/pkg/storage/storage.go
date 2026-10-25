package storage

import (
	"database/sql"

	istore "github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/storage"
)

func MigrateEmbedded(db *sql.DB) error { return istore.MigrateEmbedded(db) }
