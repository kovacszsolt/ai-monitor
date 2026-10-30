package repository

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	_ "modernc.org/sqlite"
)

func TestClientErrorLogInsertAndTrim(t *testing.T) {
	db, err := sql.Open("sqlite", "file:client_errors_test?mode=memory&cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	ctx := context.Background()
	if _, err := db.ExecContext(ctx, `
		CREATE TABLE client_error_log (
			id TEXT PRIMARY KEY NOT NULL,
			level TEXT NOT NULL DEFAULT 'error',
			source TEXT NOT NULL DEFAULT '',
			message TEXT NOT NULL,
			detail TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL
		)`); err != nil {
		t.Fatal(err)
	}

	repo := NewClientErrorLog(db)
	for i := 0; i < clientErrorLogMaxRows+5; i++ {
		if _, err := repo.Insert(ctx, "error", "test", fmt.Sprintf("msg-%d", i), ""); err != nil {
			t.Fatalf("insert %d: %v", i, err)
		}
	}

	res, err := repo.List(ctx, 1, 200)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if res.Total != int64(clientErrorLogMaxRows) {
		t.Fatalf("expected total %d, got %d", clientErrorLogMaxRows, res.Total)
	}
	if len(res.Items) == 0 || res.Items[0].Message != fmt.Sprintf("msg-%d", clientErrorLogMaxRows+4) {
		t.Fatalf("expected newest message msg-%d, got %q", clientErrorLogMaxRows+4, res.Items[0].Message)
	}
}
