ALTER TABLE ai_requests ADD COLUMN uncached_input_tokens INTEGER;
ALTER TABLE ai_requests ADD COLUMN prefill_duration_ms INTEGER;
ALTER TABLE ai_requests ADD COLUMN decode_duration_ms INTEGER;
ALTER TABLE ai_requests ADD COLUMN prompt_tps REAL;
ALTER TABLE ai_requests ADD COLUMN generation_tps REAL;
ALTER TABLE ai_requests ADD COLUMN effective_tps REAL;
ALTER TABLE ai_requests ADD COLUMN cache_hit_rate REAL;
ALTER TABLE ai_requests ADD COLUMN generation_speed_estimated INTEGER NOT NULL DEFAULT 0;
ALTER TABLE ai_requests ADD COLUMN prompt_speed_estimated INTEGER NOT NULL DEFAULT 0;
