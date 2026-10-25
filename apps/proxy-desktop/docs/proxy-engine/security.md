# Security

- Binds to `127.0.0.1` by default.
- Proxy keys: `crypto/rand` prefix `aip_…`, stored as plaintext in SQLite (`proxy_api_keys.api_key`). Protect the database file accordingly.
- Provider keys: env vars only; never logged or stored in SQLite.
- Metadata-only capture; no request/response bodies in the database.
- Upstream URLs come from the provider registry only (SSRF mitigation).
- Empty allowlist on a key (without `--allow-all-models`) denies all models.
