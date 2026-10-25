# AI Proxy Manager — design foundation (for redesign)

This document describes the **current** UI appearance, components, and patterns — not a new visual concept. Use it as a baseline for redesign, design-system consolidation, or UX audit.

Related: [PRODUCT_OVERVIEW.en.md](./PRODUCT_OVERVIEW.en.md)

---

## Implicit design goals (from code)

- **Dark, technical admin** UI — dense tables, monospace for technical values, minimal decoration.
- **Single window shell** — fixed header + sidebar; scrollable main only.
- **shadcn/Radix-inspired** — small UI primitives styled with Tailwind utilities.
- **Desktop (Wails)** — native macOS title bar integration.

---

## UI stack

| Layer | Choice |
|-------|--------|
| Styling | Tailwind CSS v4 (`@import "tailwindcss"`, `@theme`) |
| Components | Local `components/ui/*` (Button, Card, Dialog, Input, Label, Switch, Table, Badge) |
| Icons | Lucide React |
| Charts | Recharts (+ `.recharts-chart-readable` text overrides) |
| Routing | React Router 7 |
| Data | TanStack React Query (polling in several places) |
| Toasts | Sonner (`theme="dark"`, `richColors`, top-right) |
| Classes | `clsx` + `cn()`; CVA for Badge/Button variants |

---

## Color and tokens

Source: `frontend/src/index.css` `@theme` block.

| Token | HSL | Usage |
|-------|-----|--------|
| `background` | 240 10% 3.9% | App background |
| `foreground` | 0 0% 98% | Body text |
| `card` | 240 10% 6% | Cards, header |
| `muted` | 240 4% 16% | Hover, secondary surfaces |
| `muted-foreground` | 240 5% 64.9% | Secondary copy |
| `border` | 240 4% 16% | Dividers |
| `primary` | 217 91% 60% | CTAs, default badge |
| `destructive` | 0 62% 50% | Errors, delete |
| `--radius-lg` | 0.5rem | Border radius |

**No light theme** in code — single dark palette only.

### Badge variants

| Variant | Look |
|---------|------|
| `default` | Primary tint background |
| `secondary` | Muted |
| `success` | Emerald tint |
| `destructive` | Red tint |

Warnings often use `text-amber-500` (config bridge), not a dedicated badge variant.

### Charts

- Chart area uses app background.
- Axis/legend/tooltip text forced to `#fafafa` via utility class.

---

## Typography

| Element | Pattern |
|---------|---------|
| Body | System UI sans-serif |
| Page title | `text-2xl font-semibold tracking-tight` |
| Subcopy | `text-sm text-muted-foreground` |
| Card title | Often `text-base` inside `CardTitle` |
| Technical values | `font-mono text-xs` |
| Large stat | `text-2xl font-bold tracking-tight` |
| Header app name | `text-sm font-semibold` |

No separate type scale file — Tailwind utilities applied directly.

---

## Spacing and layout

### App shell

- Full viewport height; only `main` scrolls (`overflow-y-auto`, `p-6`).
- Header: `border-b`, `bg-card`, Wails drag regions; brand left padding 5.75rem on macOS.
- Sidebar: `border-r`, `bg-card/50`, collapsible with bottom toggle.

### Page widths

| Pattern | Class | Pages |
|---------|-------|-------|
| Full width | — | Dashboard, Requests, list pages |
| Narrow form | `mx-auto max-w-2xl` | Settings, Proxy variables |
| Medium | `mx-auto max-w-3xl` | Proxy, Chat empty states |
| Chat workspace | `max-w-6xl`, fixed height calc | Chat |

### Content patterns

- Default: **Card** stack with header + description + content.
- Dashboard: **6-column grid** with configurable section widths.
- Chat: **split column** layout (pins | messages) — exception to card-only lists.

---

## Component catalog (by usage)

- **Buttons:** primary / outline / ghost / destructive; sm and icon sizes.
- **Forms:** Input, Label, Switch; native `<select>` with duplicated class strings (no shared Select).
- **Tables:** base Table + ResizableTableHead (Requests) + SortableTableHead (Models).
- **Dialogs:** create/edit, confirm delete, secrets, request content viewer.
- **Shared:** HoverTooltip, PageStates, ConnectionBanner.

---

## Navigation and IA

**Nine flat sidebar items** — no grouping labels.

Implicit grouping:

| Group (not in UI) | Items |
|-------------------|-------|
| Observe | Dashboard, Requests |
| Configure | Projects, Providers, Models |
| Operate | Chat, Proxy, Proxy variables |
| App | Settings |

Proxy-related UI spans header Start, `/proxy`, `/proxy-variables`, and Settings maintenance/bridge.

---

## Page archetypes

| Archetype | Example | Traits |
|-----------|---------|--------|
| Analytics dashboard | Dashboard | Grid, collapsible sections, export, layout edit |
| Dense grid | Requests | Filters, pagination, bulk actions |
| CRUD list | Projects, Providers, Models | Card + table + modal |
| Detail | Project detail | Back nav, stacked cards |
| Control panel | Proxy | Status + few actions |
| Settings form | Settings, Proxy variables | Vertical form, max-w-2xl |
| Workspace | Chat | Split pane, composer |

---

## Interaction and feedback

- Toasts on mutations; polling intervals (connection 5s, proxy 3s, dashboard 30s).
- Optimistic delete for API keys; pending disables buttons.
- localStorage: sidebar, dashboard layouts, chat prefs, column widths.
- Lazy-loaded routes with Suspense fallback.

---

## Accessibility (current)

- `aria-label` on collapsed nav; `aria-expanded` on sidebar.
- Focus rings on inputs via `focus-visible:ring-1`.
- No documented a11y spec; contrast should be validated on redesign.
- Wide tables rely on horizontal scroll.

---

## Inconsistencies and redesign opportunities

1. Unify native `<select>` styling into one component.
2. Strong width contrast: full-bleed analytics vs narrow settings forms.
3. Consider nav grouping or nesting proxy settings.
4. Clarify two export flows (visual dashboard vs metadata export).
5. Mask provider credentials like API keys (reveal/copy).
6. Visually distinguish Chat workspace from admin pages.
7. Optional light theme or high-density mode for Requests.
8. Introduce shared `PageHeader` for consistent `h2` + actions.

---

## Suggested screenshot audit (manual)

Before redesign workshops, capture:

1. Shell — proxy off/on, config mismatch.
2. Dashboard — default and layout-edit mode.
3. Requests — filters, expanded row, content dialog.
4. Project detail — key reveal.
5. Chat — streaming session.
6. Settings — maintenance and export.

---

## Implementation references

| Area | Path |
|------|------|
| Theme | `frontend/src/index.css` |
| Shell | `frontend/src/shared/layout/AppLayout.tsx` |
| UI primitives | `frontend/src/components/ui/` |
| Dashboard | `frontend/src/features/dashboard/` |

---

## Suggested redesign sequence

1. Decide nav structure and whether to consolidate proxy settings.
2. Extend tokens (spacing, semantic colors, focus).
3. Add `PageHeader`, `Select`, optional shared data table.
4. Prototype Requests or Dashboard first (highest complexity).
