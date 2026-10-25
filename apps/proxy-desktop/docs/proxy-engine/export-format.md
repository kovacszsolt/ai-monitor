# Export format

JSONL (primary) and CSV via desktop **Settings → Database maintenance → Export**, or `pkg/export` from Go.

Each JSONL row includes `schema_version` (currently `1.0`), stable request `id`, project/key/provider identifiers, model names, timestamps, usage fields, `raw_usage_json` when present, transport metrics, and optional `tag` (from `X-AI-Proxy-Tag` when set).

CSV includes a `tag` column (empty when unset).

Secrets and message content are never exported.
