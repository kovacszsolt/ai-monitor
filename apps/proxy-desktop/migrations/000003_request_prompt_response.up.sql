CREATE TABLE ai_request_prompts (
    request_id TEXT PRIMARY KEY NOT NULL REFERENCES ai_requests(id) ON DELETE CASCADE,
    last_user_text TEXT,
    messages_json TEXT,
    prompt_truncated INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL
);

CREATE TABLE ai_request_responses (
    request_id TEXT PRIMARY KEY NOT NULL REFERENCES ai_requests(id) ON DELETE CASCADE,
    assistant_text TEXT,
    response_truncated INTEGER NOT NULL DEFAULT 0,
    created_at TEXT NOT NULL
);

CREATE INDEX idx_ai_request_prompts_created_at ON ai_request_prompts(created_at);
CREATE INDEX idx_ai_request_responses_created_at ON ai_request_responses(created_at);
