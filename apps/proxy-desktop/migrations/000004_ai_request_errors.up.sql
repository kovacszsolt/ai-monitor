CREATE TABLE ai_request_errors (
    id TEXT PRIMARY KEY NOT NULL,
    request_id TEXT NOT NULL REFERENCES ai_requests(id) ON DELETE CASCADE,
    stage TEXT NOT NULL,
    error_code TEXT NOT NULL,
    error_message TEXT NOT NULL,
    created_at TEXT NOT NULL
);

CREATE INDEX idx_ai_request_errors_request_id ON ai_request_errors(request_id);
CREATE INDEX idx_ai_request_errors_created_at ON ai_request_errors(created_at);
