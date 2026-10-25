ALTER TABLE ai_request_responses ADD COLUMN raw_body TEXT;
ALTER TABLE ai_request_responses ADD COLUMN raw_body_truncated INTEGER NOT NULL DEFAULT 0;
