# API compatibility

Authenticated routes (Bearer `aip_...`):

| Method | Path | Protocol |
|--------|------|----------|
| POST | `/v1/chat/completions` | `openai_chat` |
| POST | `/v1/responses` | `openai_responses` |
| POST | `/v1/messages` | `anthropic_messages` |
| GET | `/v1/models` | lists routes allowed for the key |

Unauthenticated:

| Method | Path |
|--------|------|
| GET | `/health/live` |
| GET | `/health/ready` |

Model routes declare a protocol; mismatched endpoint/model pairs return a clear error. Client `Authorization` is never forwarded upstream.

Optional request metadata (not forwarded upstream):

| Header | Purpose |
|--------|---------|
| `X-AI-Proxy-Tag` | Stored on `ai_requests.tag` for grouping/filtering in export (trimmed, max 256 chars). Omitted or whitespace-only → empty tag. |
