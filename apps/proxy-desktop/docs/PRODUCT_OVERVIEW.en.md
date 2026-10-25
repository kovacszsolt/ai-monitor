# AI Proxy Manager — product and user overview

Desktop app (Wails) for managing the local **ai-proxy**: SQLite database (`proxy.db`), admin UI, in-process proxy server, and built-in chat test.

Related: [DESIGN_FOUNDATION.en.md](./DESIGN_FOUNDATION.en.md)

---

## What is this app?

| | |
|---|---|
| **Display name** | AI Proxy Manager |
| **Platform** | macOS desktop (Wails v2 + React) |
| **Data** | Single shared SQLite file — projects, providers, model routes, API keys, `ai_requests` log |
| **Proxy** | Same logic as `proxy-cli` — the desktop runs the HTTP proxy **in-process** |

**Core value:** configure upstream AI providers, issue per-project API keys, run the proxy, inspect request metrics, and smoke-test models in the built-in chat — all in one place.

---

## Who is the user?

- Single-user **local admin** — no sign-in, no roles in the UI.
- Typical user: developer or power user running the proxy on their machine and pointing Cursor, CLI, or other clients at it.

---

## Application shell

Every page shares the same frame:

### Header (title bar)

- App name and **build time** (UTC tooltip).
- macOS: draggable region for native window chrome.
- **Start proxy** button when the proxy is stopped (quick start).
- **Connection strip:** Proxy on/off, Connected / Database issue, optional “Config mismatch”, active database path, error text.

### Sidebar

| Menu | Route | Purpose |
|------|-------|---------|
| Dashboard | `/` | Analytics overview |
| Requests | `/requests` | Request log |
| Projects | `/projects` | Tenants |
| Providers | `/providers` | Upstream APIs |
| Models | `/models` | Model routes |
| Chat | `/chat` | Test chat via proxy |
| Proxy | `/proxy` | Start/stop proxy |
| Proxy variables | `/proxy-variables` | Runtime server/log/capture |
| Settings | `/settings` | Desktop DB path, maintenance |

- Sidebar is **collapsible**; state persists in localStorage.
- Collapsed: icon + tooltip only.

### Main content

- Scrollable area with ~24px padding.
- Toasts: top-right, dark theme (Sonner).

### Shared states

- **Loading** — “Loading…” label.
- **Error** — message + retry.
- **Empty** — short copy, often a primary action (e.g. “Create one”).

---

## Pages

### Dashboard (`/`)

**Purpose:** summary of proxied traffic.

**User can:**

- View aggregate stats (total requests, last 24h, tokens, avg duration, by state, top models/projects).
- Filter hourly volume and model usage charts by **project** and **proxy tag** (`X-AI-Proxy-Tag` header).
- Click table rows to set tag filters in usage sections.
- **Customize layout:** section order, width (⅓ / ½ / full row on large screens), collapse sections.
- **Layout presets:** save, switch, copy, rename, delete.
- **Export** dashboard visuals as PNG / HTML / PDF (with filter summary).

**Sections include:** request volume by hour, usage by model, stat cards, usage by proxy tag table, token mix chart (hidden when no tag selected), by state, top models, top projects.

**UI:** responsive grid, cards, Recharts. Data refreshes about every 30 seconds.

---

### Requests (`/requests`)

**Purpose:** browse and manage `ai_requests` rows.

**User can:**

- **Filter** by project, state, tag, model.
- **Paginate** (50 per page).
- **Sort** by column; **resize** columns.
- **Expand** rows; open **content** dialog (prompt/response when capture allows).
- **Delete** one request.
- **Multi-select:** bulk tag move, bulk delete (when DB supports tags).

**UI:** dense resizable table; metric tooltips (TTFT, cache rate, tok/s, etc.); state badges.

---

### Projects (`/projects`)

**Purpose:** tenant boundaries for API keys.

**User can:** create/edit projects (name, description), enable/disable, view request counts and last request time, open **API keys** detail.

---

### Project detail (`/projects/:id`)

**Purpose:** API keys for one project.

**User can:** create keys (`aip_…`), reveal/copy, allow-all-models toggle, revoke, delete (with confirmation). Legacy hashed keys show a migration warning.

**UI:** notes that keys are stored **plaintext** in SQLite.

---

### Providers (`/providers`)

**Purpose:** upstream LLM endpoints.

**User can:** add/edit providers (`openai`, `anthropic`, `openai-compatible`), base URL, credential, timeouts, enable/disable.

---

### Models (`/models`)

**Purpose:** public model name → upstream model + provider.

**User can:** CRUD routes with protocol (`openai_chat`, `openai_responses`, `anthropic_messages`), optional USD pricing fields, filter and sort table, enable/disable.

---

### Chat (`/chat`)

**Purpose:** quick local proxy test.

**Requires:** proxy **running**, enabled project, API key, and model route.

**User can:** pick session (project, key, provider, model — persisted locally), toggle streaming, send messages, pin/unpin prompts, run pinned prompts, clear chat. Requests logged with tag `desktop-chat`.

**UI:** two-column layout (pinned prompts | messages + composer).

---

### Proxy (`/proxy`)

**User can:** start/stop, copy base URL, re-check provider health while running.

---

### Proxy variables (`/proxy-variables`)

**Purpose:** runtime settings in SQLite `settings` (not desktop YAML).

**User can:** server host/port, max request bytes, logging level/format, capture toggles and byte limits. **Stop proxy before save;** restart to apply.

---

### Settings (`/settings`)

**User can:** database auto-detect vs manual path, file browse, save & reconnect; database maintenance (ping, migrations, stats, JSONL/CSV metadata export); view proxy-cli YAML bridge status.

---

## Typical flows

**First-time setup:** Settings (DB path, migrations) → Providers → Models → Projects & keys → Proxy variables → Start proxy → external client or Chat.

**Daily ops:** Dashboard for trends → Requests for failures and content.

**Config change:** Stop proxy → edit variables or routes → start again.

---

## Limits and notes

| Topic | User impact |
|-------|-------------|
| No auth | Anyone with machine/DB access sees credentials and keys. |
| Plaintext secrets | Provider credentials and `aip_` keys are readable in DB and often in UI. |
| Single DB | Desktop and proxy-cli should point at the same file (bridge warning). |
| Two config layers | Desktop YAML = DB **path**; proxy runtime = DB **settings** (+ optional YAML seed). |
| Chat | Local smoke test only; tag `desktop-chat`. |
| `/upstream` | Redirects to `/providers`. |

---

## Outside the UI

Clients call the running proxy HTTP API with project API keys and configured public model names. The desktop configures and logs that traffic; see `proxy-cli` for API details.

---

## Related docs

- [README.md](../README.md) — development and build.
- [DESIGN_FOUNDATION.en.md](./DESIGN_FOUNDATION.en.md) — visual system and redesign baseline.
