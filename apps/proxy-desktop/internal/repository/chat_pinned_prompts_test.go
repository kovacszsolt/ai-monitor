package repository

import (
	"context"
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func TestChatPinnedPrompts_CRUD(t *testing.T) {
	db, err := sql.Open("sqlite", "file:chat_pins_test?mode=memory&cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	ctx := context.Background()
	if _, err := db.ExecContext(ctx, `
		CREATE TABLE chat_pinned_prompts (
			id TEXT PRIMARY KEY NOT NULL,
			title TEXT NOT NULL DEFAULT '',
			body TEXT NOT NULL,
			created_at TEXT NOT NULL,
			updated_at TEXT NOT NULL
		)`); err != nil {
		t.Fatal(err)
	}

	r := NewChatPinnedPrompts(db)
	created, err := r.Create(ctx, "My pin", "What is 2+2?")
	if err != nil {
		t.Fatal(err)
	}
	list, err := r.List(ctx)
	if err != nil || len(list) != 1 {
		t.Fatalf("list: %v err %v", list, err)
	}
	if list[0].Body != "What is 2+2?" {
		t.Fatalf("unexpected body %q", list[0].Body)
	}
	if err := r.Delete(ctx, created.ID); err != nil {
		t.Fatal(err)
	}
	list, err = r.List(ctx)
	if err != nil || len(list) != 0 {
		t.Fatalf("expected empty list, got %v", list)
	}
	if err := r.Delete(ctx, created.ID); err != ErrChatPinnedPromptNotFound {
		t.Fatalf("expected not found, got %v", err)
	}
}
