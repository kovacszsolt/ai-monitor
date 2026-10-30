CREATE TABLE client_error_log (
    id TEXT PRIMARY KEY NOT NULL,
    level TEXT NOT NULL DEFAULT 'error',
    source TEXT NOT NULL DEFAULT '',
    message TEXT NOT NULL,
    detail TEXT NOT NULL DEFAULT '',
    created_at TEXT NOT NULL
);

CREATE INDEX idx_client_error_log_created_at ON client_error_log (created_at DESC);
