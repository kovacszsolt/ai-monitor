# AI Proxy Manager (desktop)

Wails app for `proxy.db`: projects, providers, routes, API keys, request analytics, chat, and in-process HTTP proxy.

Go module: `github.com/kovacszso/ai-monitor/apps/proxy-desktop`

The HTTP gateway and migrations live under `internal/proxyengine/` and `migrations/` in this same module (no separate proxy package).

## Requirements

- Go 1.27+
- Node.js 20+
- [pnpm](https://pnpm.io/)
- [Wails v2](https://wails.io/docs/gettingstarted/installation)

## Configuration

| File | Purpose |
|------|---------|
| `config.yaml` | Desktop `storage.database_path` (empty → `~/.config/ai-proxy-manager/proxy.db`; override with `AI_PROXY_DESKTOP_DATABASE_PATH`) |
| `proxy.bootstrap.yaml` | Optional seed for server/logging/capture (`AI_PROXY_CONFIG`, `AI_PROXY_*`) |
| SQLite `settings` | Runtime server, logging, capture |

## Development

```bash
cd apps/proxy-desktop
make frontend-install
cp config.example.yaml config.yaml   # optional
cp proxy.bootstrap.example.yaml proxy.bootstrap.yaml   # optional
wails dev
```

## Tests

```bash
make test
go test ./internal/proxyengine/... ./proxyintegration/...
```
