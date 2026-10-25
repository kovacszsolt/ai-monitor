CREATE TABLE chat_pinned_prompts (
    id TEXT PRIMARY KEY NOT NULL,
    title TEXT NOT NULL DEFAULT '',
    body TEXT NOT NULL,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE INDEX idx_chat_pinned_prompts_created_at ON chat_pinned_prompts (created_at DESC);
