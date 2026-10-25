ALTER TABLE ai_request_responses ADD COLUMN tool_calls_json TEXT;
ALTER TABLE ai_request_responses ADD COLUMN tool_calls_truncated INTEGER NOT NULL DEFAULT 0;
