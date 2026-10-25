PRAGMA foreign_keys = ON;

CREATE TABLE projects (
    id TEXT PRIMARY KEY NOT NULL,
    name TEXT NOT NULL UNIQUE,
    description TEXT NOT NULL DEFAULT '',
    enabled INTEGER NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE TABLE proxy_api_keys (
    id TEXT PRIMARY KEY NOT NULL,
    project_id TEXT NOT NULL REFERENCES projects(id),
    name TEXT NOT NULL,
    key_hash TEXT NOT NULL UNIQUE,
    allow_all_models INTEGER NOT NULL DEFAULT 0,
    enabled INTEGER NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL,
    expires_at TEXT,
    last_used_at TEXT,
    UNIQUE(project_id, name)
);

CREATE TABLE providers (
    id TEXT PRIMARY KEY NOT NULL,
    name TEXT NOT NULL UNIQUE,
    provider_type TEXT NOT NULL,
    base_url TEXT NOT NULL,
    credential_ref TEXT NOT NULL DEFAULT '',
    enabled INTEGER NOT NULL DEFAULT 1,
    connection_timeout_ms INTEGER NOT NULL DEFAULT 10000,
    request_timeout_ms INTEGER NOT NULL DEFAULT 300000,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE TABLE model_routes (
    id TEXT PRIMARY KEY NOT NULL,
    public_model_name TEXT NOT NULL UNIQUE,
    provider_id TEXT NOT NULL REFERENCES providers(id),
    upstream_model_name TEXT NOT NULL,
    protocol TEXT NOT NULL,
    enabled INTEGER NOT NULL DEFAULT 1,
    created_at TEXT NOT NULL,
    updated_at TEXT NOT NULL
);

CREATE TABLE api_key_allowed_models (
    api_key_id TEXT NOT NULL REFERENCES proxy_api_keys(id) ON DELETE CASCADE,
    model_route_id TEXT NOT NULL REFERENCES model_routes(id),
    PRIMARY KEY (api_key_id, model_route_id)
);

CREATE TABLE ai_requests (
    id TEXT PRIMARY KEY NOT NULL,
    project_id TEXT NOT NULL,
    proxy_api_key_id TEXT NOT NULL,
    provider_id TEXT NOT NULL,
    model_route_id TEXT,
    requested_model TEXT NOT NULL,
    actual_model TEXT,
    endpoint TEXT NOT NULL,
    protocol TEXT NOT NULL,
    started_at TEXT NOT NULL,
    completed_at TEXT,
    duration_ms INTEGER,
    upstream_latency_ms INTEGER,
    ttfb_ms INTEGER,
    ttft_ms INTEGER,
    status_code INTEGER,
    request_state TEXT NOT NULL,
    streaming INTEGER NOT NULL DEFAULT 0,
    input_tokens INTEGER,
    output_tokens INTEGER,
    cached_input_tokens INTEGER,
    cache_write_tokens INTEGER,
    reasoning_tokens INTEGER,
    usage_status TEXT NOT NULL DEFAULT 'unknown',
    raw_usage_json TEXT,
    request_bytes INTEGER NOT NULL DEFAULT 0,
    response_bytes INTEGER NOT NULL DEFAULT 0,
    error_type TEXT,
    timeout_flag INTEGER NOT NULL DEFAULT 0,
    cancelled_flag INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX idx_ai_requests_started_at ON ai_requests(started_at);
CREATE INDEX idx_ai_requests_project_started ON ai_requests(project_id, started_at);
CREATE INDEX idx_ai_requests_provider_started ON ai_requests(provider_id, started_at);
CREATE INDEX idx_ai_requests_model_route_started ON ai_requests(model_route_id, started_at);
CREATE INDEX idx_ai_requests_api_key_started ON ai_requests(proxy_api_key_id, started_at);
CREATE INDEX idx_ai_requests_state ON ai_requests(request_state);
