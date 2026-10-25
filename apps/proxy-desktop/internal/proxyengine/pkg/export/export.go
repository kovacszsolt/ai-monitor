package export

import (
	"context"
	"database/sql"
	"io"

	iexport "github.com/kovacszso/ai-monitor/apps/proxy-desktop/internal/proxyengine/export"
)

type Filter = iexport.Filter

func StorageStats(ctx context.Context, db *sql.DB) (int64, error) {
	return iexport.StorageStats(ctx, db)
}

func ExportJSONL(ctx context.Context, db *sql.DB, w io.Writer, f Filter) error {
	return iexport.ExportJSONL(ctx, db, w, f)
}

func ExportCSV(ctx context.Context, db *sql.DB, w io.Writer, f Filter) error {
	return iexport.ExportCSV(ctx, db, w, f)
}
