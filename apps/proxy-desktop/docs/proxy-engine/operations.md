# Operations

## Backup

Copy `storage.database_path` (default `./data/proxy.db`) while the process is stopped or rely on SQLite WAL-aware backup tools.

## Graceful shutdown

`serve` drains the HTTP server and attempts to flush the metadata writer queue on SIGINT/SIGTERM.

## Queue saturation

If the persistence queue is full, new metadata returns `ErrQueueFull` and increments a dropped counter (`storage stats`).

## Durability

SQLite `synchronous` defaults to `NORMAL` (see config). Abrupt kill may lose records still in the in-memory queue.

## Recovery

On restart, rows with `progress = started` (or `request_state = started`) from a previous crash are not auto-reconciled (known limitation).
