# Configuration

Primary file: `config.yaml` (see `config.example.yaml`).

Environment overrides:

| Variable | Effect |
|----------|--------|
| `AI_PROXY_CONFIG` | Config file path |
| `AI_PROXY_DATABASE` | SQLite path |
| `AI_PROXY_HOST` / `AI_PROXY_PORT` | Listen address |
| `AI_PROXY_MIGRATIONS` | Migrations directory |
| `AI_PROXY_LOG_FORMAT` | `console` (colored, [tint](https://github.com/lmittmann/tint)) or `json` |
| `AI_PROXY_HTTP_DUMP_DIR` | Enables `capture.file_dump` and sets `directory` |

### HTTP exchange file dump

Config is **YAML** (`config.yaml`), not TOML. Optional mirror of each proxied request/response to JSON files (independent of SQLite):

```yaml
capture:
  file_dump:
    enabled: true
    directory: ./data/http-dumps
    errors_only: false   # true: only 4xx/5xx, parse_error, failed/interrupted
    max_body_bytes: 1048576
```

Files: `{directory}/{timestamp}_{request_id}.json` with `request` and `response` bodies (truncated). Log line: `http exchange file dump written`.

`logging.format` in YAML: `console` for local dev, `json` for Docker/production log collectors.

Provider upstream API keys are stored in SQLite (`providers.credential`). Set them in the desktop **Providers** UI.
