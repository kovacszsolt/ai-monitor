# Database schema

SQLite with WAL, foreign keys, and migrations in `migrations/`.

Core tables: `projects`, `proxy_api_keys`, `providers`, `model_routes`, `api_key_allowed_models`, `ai_requests`.

**`proxy_api_keys.api_key`** (migration `000013_proxy_api_key_plaintext`): issued `aip_…` secret stored as plaintext (column renamed from `key_hash`). Recreate keys after upgrading from hashed storage.

**`model_routes` token pricing** (migration `000009_model_route_token_pricing`): optional USD per **1M tokens** — `input_token_usd_per_1m`, `cached_input_token_usd_per_1m`, `output_token_usd_per_1m` (nullable `REAL`; unset until configured).

`ai_requests` stores timing, usage, status, and byte counts only — not prompts or completions by default.

**`tag`** (migration `000008_request_tag`): optional client label from request header `X-AI-Proxy-Tag` (trimmed, max 256 chars). Empty string when omitted. Not forwarded upstream.

**`progress`** (`started` | `finished`): lifecycle while the proxy handles the request. In-flight rows have `progress = started` and `request_state = started` until the upstream response is done; then a finish **UPDATE** sets `progress = finished` and the final `request_state` (`completed`, `failed`, etc.). Legacy rows default to `finished`.

When `capture.store_content: true`, text is stored in separate tables (1:1 with `ai_requests.id`, `ON DELETE CASCADE`):

- **`ai_request_prompts`** — `last_user_text`, `messages_json`, `prompt_truncated`
- **`ai_request_responses`** — `assistant_text`, `response_truncated`, `tool_calls_json`, `tool_calls_truncated` (OpenAI `tool_calls` / Anthropic `tool_use` blocks), `raw_body` / `raw_body_truncated` (truncated upstream payload on HTTP errors, `usage_status=parse_error`, or when parsed assistant text is missing while `store_content` is enabled; controlled by `capture.store_upstream_body_on_error`)
- **`ai_request_errors`** — capture failures and **proxy errors** (`stage`, `error_code`, `error_message`). Proxy routing/upstream failures use `stage=proxy`, `error_code=proxy`, message from `ai_requests.error_type`.

Computed throughput fields (migration `000002_token_metrics`, nullable on older rows):

- `uncached_input_tokens`, `prefill_duration_ms`, `decode_duration_ms`
- `prompt_tps`, `generation_tps`, `effective_tps`, `cache_hit_rate`
- `generation_speed_estimated`, `prompt_speed_estimated` (0/1 flags)

Throughput semantics:

- **`effective_tps`** — `(input_tokens + output_tokens) / duration_ms` (end-to-end wall clock).
- **`generation_tps`** — streaming: `output_tokens / (last_output_token_ms - first_output_token_ms)` when SSE timing is available; non-stream: **estimate** `output_tokens / duration_ms` (`generation_speed_estimated = 1`). This is not isolated decode speed on non-stream (prefill is included in `duration_ms`).
- **`prompt_tps`** — only when `prefill_duration_ms` is set (not estimated for non-stream today).

Indexes cover `started_at`, project/provider/model/key + time, and `request_state`.
