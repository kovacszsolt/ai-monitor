# Architecture

The proxy engine lives under [`apps/proxy-desktop/internal/proxyengine`](../../internal/proxyengine) and is embedded in the Wails desktop app. Paths below are relative to `internal/proxyengine/`.

## Proxy engine (in-process HTTP)

- **HTTP server** (`server`) — chi router, auth middleware, health endpoints.
- **Proxy** (`proxy`) — routing, upstream forwarding, streaming (SSE), metadata enqueue.
- **Auth** (`auth`) — proxy API keys in SQLite (`proxy_api_keys`); desktop stores plaintext keys for copy/reuse.
- **Projects / providers / routing** — SQLite-backed registries (`projects`, `providers`, `model_routes`).
- **Capture** (`capture`) — bounded channel + batch writer; optional prompt/response body persistence when `capture.store_content` is enabled.
- **Settings** (`settings`) — DB-backed overlay for server, logging, and capture options (merged with bootstrap YAML).
- **Export** (`export` / `pkg/export`) — streaming JSONL/CSV reads from `ai_requests`.

Request path: Client → auth → model resolve → provider adapter → upstream → client (+ async SQLite writer).

## Desktop shell

- **Wails UI** (`frontend/`) — React feature modules call `DesktopService` via wailsjs bindings.
- **Desktop service** (`internal/service`) — façade over domain packages:
  - `catalog` — projects, providers, model routes, API keys
  - `analytics` — dashboard and chart aggregates
  - `airequests` — paginated request list, tags, content
  - `settings`, `chat`, `proxy_ops` — app configuration and in-app chat against the local proxy
- **Config bridge** (`internal/config/bridge.go`) — merges desktop YAML, proxy bootstrap YAML, and DB settings; aligns `database_path`.
- **Runtime** (`internal/proxyruntime`) — starts `proxyengine` in-process with a **second** SQLite connection to the same file as the UI (WAL). Avoid long write transactions on both sides at once.

## Data

- Single SQLite database (default under app config dir); schema via embedded migrations in [`migrations/`](../../migrations).
- Desktop-only queries (e.g. project list with request stats, bulk tag edits) live in `internal/repository` without duplicating proxyengine store CRUD.
