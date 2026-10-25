-- Default local/OpenAI providers on first database bootstrap (empty providers table only).
INSERT INTO providers (
    id, name, provider_type, base_url, credential, enabled,
    connection_timeout_ms, request_timeout_ms, created_at, updated_at
)
SELECT id, name, provider_type, base_url, credential, enabled,
       connection_timeout_ms, request_timeout_ms, created_at, updated_at
FROM (
    SELECT
        'b1000001-0000-4000-8000-000000000001' AS id,
        'GS4' AS name,
        'openai-compatible' AS provider_type,
        'http://127.0.0.1:8003/v1' AS base_url,
        '' AS credential,
        1 AS enabled,
        10000 AS connection_timeout_ms,
        300000 AS request_timeout_ms,
        strftime('%Y-%m-%dT%H:%M:%fZ', 'now') AS created_at,
        strftime('%Y-%m-%dT%H:%M:%fZ', 'now') AS updated_at
    UNION ALL SELECT
        'b1000002-0000-4000-8000-000000000002',
        'MTPLX',
        'openai-compatible',
        'http://127.0.0.1:8001/v1',
        '', 1, 10000, 300000,
        strftime('%Y-%m-%dT%H:%M:%fZ', 'now'),
        strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
    UNION ALL SELECT
        'b1000003-0000-4000-8000-000000000003',
        'OMLX',
        'openai-compatible',
        'http://127.0.0.1:8002/v1',
        '', 1, 10000, 300000,
        strftime('%Y-%m-%dT%H:%M:%fZ', 'now'),
        strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
    UNION ALL SELECT
        'b1000004-0000-4000-8000-000000000004',
        'OPENAI',
        'openai',
        'https://api.openai.com/v1',
        '', 1, 10000, 300000,
        strftime('%Y-%m-%dT%H:%M:%fZ', 'now'),
        strftime('%Y-%m-%dT%H:%M:%fZ', 'now')
) AS defaults
WHERE (SELECT COUNT(*) FROM providers) = 0;
