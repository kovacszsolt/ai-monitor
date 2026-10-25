# ai-monitor

Open-source monorepo for the **AI Proxy Manager** macOS desktop app.

Run a local AI HTTP proxy, configure upstream providers, issue per-project API keys, and inspect request metrics in one place — without relying on a hosted admin panel or a separate telemetry SaaS.

## What is it?

**AI Proxy Manager** (Wails + Go + React) is a **single-user, local admin** tool:

- **One SQLite file** (`proxy.db`) holds projects, providers, model routes, API keys, and `ai_requests` rows logged by the proxy.
- The **HTTP proxy runs inside the desktop process** — the same engine you could run via a standalone CLI; the UI configures it, starts/stops it, and surfaces the data.
- **No sign-in** in the app: anyone with access to your machine and database can see credentials and keys (designed for local / developer use).

Screen-by-screen details: [PRODUCT_OVERVIEW.en.md](apps/proxy-desktop/docs/PRODUCT_OVERVIEW.en.md).

## What can you do?

| Area | What you get |
|------|----------------|
| **Proxy** | OpenAI-compatible local endpoint: clients (Cursor, CLI, your app) call with an `aip_…` key and a public model name; the proxy forwards to the right upstream. |
| **Projects & keys** | Tenant-style boundaries: API keys per project, model allow lists, enable/disable. |
| **Providers** | Multiple upstreams (OpenAI, Anthropic, OpenAI-compatible / custom base URL), credentials, timeouts, enable/disable. |
| **Model routes** | Public name → provider + upstream model + protocol (`openai_chat`, `openai_responses`, `anthropic_messages`), optional token pricing fields. |
| **Dashboard** | Aggregated traffic, tokens, states, hourly charts; filter by **project** and **`X-AI-Proxy-Tag`**; layout presets; export PNG / HTML / PDF. |
| **Requests** | Paginated, filterable request log; metrics (TTFT, tok/s, cache); optional prompt/response view when capture is enabled; bulk tag / delete. |
| **Chat** | Built-in smoke test through the running proxy (streaming, pinned prompts); requests are logged with the `desktop-chat` tag. |
| **Proxy variables** | Runtime settings in SQLite `settings`: host/port, log level, capture limits — save after stopping the proxy. |
| **Settings** | Database path, migrations, maintenance, JSONL/CSV metadata export; warning when desktop and bootstrap YAML disagree on the DB path. |

## Typical setups

**One machine, multiple “users” (projects)**  
Separate project + API key per team, experiment, or app; filter the dashboard and request list by project.

**Cost and experiment tracking with tags**  
Send `X-AI-Proxy-Tag` from the client (e.g. `cursor`, `eval-v2`); dashboards and tables break down usage by tag.

**Several upstreams behind one proxy**  
Local MLX/Ollama-compatible endpoint and cloud OpenAI/Anthropic side by side: providers + model routes; clients always point at the same proxy URL; routing picks the upstream.

**OpenAI-compatible toolchain**  
Cursor, `curl`, LangChain, etc. use the proxy base URL and an issued `aip_…` key; the public model name comes from the route table.

**Telemetry in your database, not a third-party SaaS**  
Metrics and (optional) capture stay in your `proxy.db`; export or back up the file as needed.

**First-time setup (recommended order)**  
Settings (DB, migrations) → Providers → Models → Projects & keys → Proxy variables → Start proxy → external client or Chat.

**Day to day**  
Dashboard for trends → Requests for failures and content; stop the proxy before changing runtime config.

## Quick start (developers)

Requirements: Go 1.27+, Node.js 20+, [pnpm](https://pnpm.io/), [Wails v2](https://wails.io/docs/gettingstarted/installation).

```bash
cd apps/proxy-desktop
make frontend-install
cp config.example.yaml config.yaml   # optional
wails dev
```

Build, tests, and proxy engine architecture: [apps/proxy-desktop/README.md](apps/proxy-desktop/README.md) and [docs/proxy-engine/](apps/proxy-desktop/docs/proxy-engine/).

## Repository layout

| Path | Description |
|------|-------------|
| [`apps/proxy-desktop/`](apps/proxy-desktop/) | Desktop app and in-process proxy |
| [`data/`](data/) | Suggested location for local `proxy.db` (gitignored) |

## Important limits

- **Local security:** provider credentials and API keys are stored in plaintext in SQLite; do not share the DB file or run on untrusted machines with sensitive keys.
- **Single database:** desktop and optional bootstrap YAML should point at the same `proxy.db` (the app warns on mismatch).
- **Chat:** quick test only — not a replacement for your production client.

## License

MIT — see [LICENSE](LICENSE).
