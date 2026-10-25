ALTER TABLE ai_requests ADD COLUMN progress TEXT NOT NULL DEFAULT 'finished';

CREATE INDEX idx_ai_requests_progress ON ai_requests(progress);
