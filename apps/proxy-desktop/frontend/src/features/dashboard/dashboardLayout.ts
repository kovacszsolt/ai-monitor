import { DASHBOARD_SECTION_IDS, LEGACY_DASHBOARD_SECTION_IDS } from "@/features/dashboard/dashboardSectionIds";
import {
  defaultDashboardLayoutUsagePrefs,
  mergeDashboardLayoutUsagePrefs,
  normalizeDashboardLayoutUsagePrefs,
  type DashboardLayoutUsagePrefs,
} from "@/features/dashboard/dashboardLayoutUsagePrefs";

export type { DashboardLayoutUsagePrefs } from "@/features/dashboard/dashboardLayoutUsagePrefs";

export type DashboardSectionId = (typeof DASHBOARD_SECTION_IDS)[keyof typeof DASHBOARD_SECTION_IDS];

export type DashboardSectionWidth = "full" | "half" | "third";

/** Tailwind col-span on the dashboard `lg:grid-cols-6` grid. */
export function dashboardSectionColSpanClass(width: DashboardSectionWidth): string {
  switch (width) {
    case "third":
      return "lg:col-span-2";
    case "half":
      return "lg:col-span-3";
    case "full":
      return "lg:col-span-6";
  }
}

export type DashboardLayoutEntry = {
  id: DashboardSectionId;
  width: DashboardSectionWidth;
};

export type DashboardLayoutPreset = {
  id: string;
  name: string;
  entries: DashboardLayoutEntry[];
  /** Usage-by-hour chart filters (also drive tag/project on related dashboard sections). */
  usage: DashboardLayoutUsagePrefs;
};

export type DashboardLayoutStore = {
  activePresetId: string;
  presets: DashboardLayoutPreset[];
};

const LEGACY_STORAGE_KEY = "proxy-desktop.dashboard-layout";
const STORAGE_KEY = "proxy-desktop.dashboard-layout-presets";

export const DEFAULT_PRESET_ID = "default";

const STAT_SECTION_IDS: DashboardSectionId[] = [
  DASHBOARD_SECTION_IDS.statTotalRequests,
  DASHBOARD_SECTION_IDS.statRequests24h,
  DASHBOARD_SECTION_IDS.statTokens,
  DASHBOARD_SECTION_IDS.statAvgDuration,
];

const BREAKDOWN_SECTION_IDS: DashboardSectionId[] = [
  DASHBOARD_SECTION_IDS.byState,
  DASHBOARD_SECTION_IDS.topModels,
  DASHBOARD_SECTION_IDS.topProjects,
];

export const DASHBOARD_SECTION_LABELS: Record<DashboardSectionId, string> = {
  [DASHBOARD_SECTION_IDS.usageByHour]: "Usage over time",
  [DASHBOARD_SECTION_IDS.usageByModel]: "Usage by model",
  [DASHBOARD_SECTION_IDS.statTotalRequests]: "Total requests",
  [DASHBOARD_SECTION_IDS.statRequests24h]: "Last 24 hours",
  [DASHBOARD_SECTION_IDS.statTokens]: "Tokens (in / cached / out)",
  [DASHBOARD_SECTION_IDS.statAvgDuration]: "Avg duration",
  [DASHBOARD_SECTION_IDS.usageByProxyTag]: "Usage by proxy tag",
  [DASHBOARD_SECTION_IDS.tokenMixChart]: "Token mix by model",
  [DASHBOARD_SECTION_IDS.byState]: "By state",
  [DASHBOARD_SECTION_IDS.topModels]: "Top models",
  [DASHBOARD_SECTION_IDS.topProjects]: "Top projects",
};

export const DEFAULT_DASHBOARD_LAYOUT: DashboardLayoutEntry[] = [
  { id: DASHBOARD_SECTION_IDS.usageByHour, width: "full" },
  { id: DASHBOARD_SECTION_IDS.usageByModel, width: "full" },
  { id: DASHBOARD_SECTION_IDS.statTotalRequests, width: "third" },
  { id: DASHBOARD_SECTION_IDS.statRequests24h, width: "third" },
  { id: DASHBOARD_SECTION_IDS.statTokens, width: "third" },
  { id: DASHBOARD_SECTION_IDS.statAvgDuration, width: "third" },
  { id: DASHBOARD_SECTION_IDS.usageByProxyTag, width: "third" },
  { id: DASHBOARD_SECTION_IDS.tokenMixChart, width: "third" },
  { id: DASHBOARD_SECTION_IDS.byState, width: "third" },
  { id: DASHBOARD_SECTION_IDS.topModels, width: "third" },
  { id: DASHBOARD_SECTION_IDS.topProjects, width: "third" },
];

function widthFromRaw(value: unknown): DashboardSectionWidth | null {
  if (value === "full" || value === "half" || value === "third") return value;
  return null;
}

function expandRawLayoutEntry(id: string, width: DashboardSectionWidth): DashboardLayoutEntry[] {
  if (id === LEGACY_DASHBOARD_SECTION_IDS.summaryStats) {
    const w = width === "full" ? "third" : width;
    return STAT_SECTION_IDS.map((sectionId) => ({ id: sectionId, width: w }));
  }
  if (id === LEGACY_DASHBOARD_SECTION_IDS.breakdown) {
    const w = width === "full" ? "third" : width;
    return BREAKDOWN_SECTION_IDS.map((sectionId) => ({ id: sectionId, width: w }));
  }
  if (isSectionId(id)) {
    return [{ id, width }];
  }
  return [];
}

function parseLayoutEntriesRaw(entriesRaw: unknown[]): DashboardLayoutEntry[] {
  const out: DashboardLayoutEntry[] = [];
  for (const raw of entriesRaw) {
    if (!raw || typeof raw !== "object") continue;
    const id = (raw as { id?: string }).id;
    const width = widthFromRaw((raw as { width?: unknown }).width);
    if (!id || !width) continue;
    out.push(...expandRawLayoutEntry(id, width));
  }
  return out;
}

const ALL_SECTION_IDS = new Set<DashboardSectionId>(Object.values(DASHBOARD_SECTION_IDS));

function isSectionId(value: string): value is DashboardSectionId {
  return ALL_SECTION_IDS.has(value as DashboardSectionId);
}

/** Expand legacy composite section ids and drop unknown ids. */
function flattenLayoutEntries(saved: DashboardLayoutEntry[]): DashboardLayoutEntry[] {
  const out: DashboardLayoutEntry[] = [];
  for (const entry of saved) {
    const rawId = entry.id as string;
    out.push(...expandRawLayoutEntry(rawId, entry.width));
  }
  return out;
}

export function mergeDashboardLayout(saved: DashboardLayoutEntry[]): DashboardLayoutEntry[] {
  const seen = new Set<DashboardSectionId>();
  const merged: DashboardLayoutEntry[] = [];

  for (const entry of flattenLayoutEntries(saved)) {
    if (seen.has(entry.id)) continue;
    seen.add(entry.id);
    merged.push(entry);
  }

  for (const defaultEntry of DEFAULT_DASHBOARD_LAYOUT) {
    if (!seen.has(defaultEntry.id)) {
      merged.push(defaultEntry);
    }
  }

  return merged;
}

export function cloneDashboardLayout(entries: DashboardLayoutEntry[]): DashboardLayoutEntry[] {
  return mergeDashboardLayout(entries.map((e) => ({ ...e })));
}

function createDefaultStore(): DashboardLayoutStore {
  return {
    activePresetId: DEFAULT_PRESET_ID,
    presets: [
      {
        id: DEFAULT_PRESET_ID,
        name: "Default",
        entries: cloneDashboardLayout(DEFAULT_DASHBOARD_LAYOUT),
        usage: defaultDashboardLayoutUsagePrefs(),
      },
    ],
  };
}

function normalizePreset(raw: unknown): DashboardLayoutPreset | null {
  if (!raw || typeof raw !== "object") return null;
  const id = (raw as { id?: string }).id;
  const name = (raw as { name?: string }).name;
  const entriesRaw = (raw as { entries?: unknown }).entries;
  if (!id || typeof id !== "string" || !name || typeof name !== "string") return null;
  if (!Array.isArray(entriesRaw)) return null;
  const entries = parseLayoutEntriesRaw(entriesRaw);
  if (entries.length === 0) return null;
  const usageRaw = (raw as { usage?: unknown }).usage;
  return {
    id,
    name: name.trim() || "Layout",
    entries: mergeDashboardLayout(entries),
    usage: normalizeDashboardLayoutUsagePrefs(usageRaw),
  };
}

function migrateLegacyStore(): DashboardLayoutStore | null {
  try {
    const raw = localStorage.getItem(LEGACY_STORAGE_KEY);
    if (!raw) return null;
    const parsed = JSON.parse(raw) as unknown;
    if (!Array.isArray(parsed)) return null;
    const entries = parseLayoutEntriesRaw(parsed);
    if (entries.length === 0) return null;
    localStorage.removeItem(LEGACY_STORAGE_KEY);
    return {
      activePresetId: DEFAULT_PRESET_ID,
      presets: [
        {
          id: DEFAULT_PRESET_ID,
          name: "Default",
          entries: mergeDashboardLayout(entries),
          usage: defaultDashboardLayoutUsagePrefs(),
        },
      ],
    };
  } catch {
    return null;
  }
}

function normalizeStore(raw: unknown): DashboardLayoutStore | null {
  if (!raw || typeof raw !== "object") return null;
  const activePresetId = (raw as { activePresetId?: string }).activePresetId;
  const presetsRaw = (raw as { presets?: unknown }).presets;
  if (!activePresetId || !Array.isArray(presetsRaw)) return null;
  const presets = presetsRaw.map(normalizePreset).filter((p): p is DashboardLayoutPreset => p != null);
  if (presets.length === 0) return null;
  const activeExists = presets.some((p) => p.id === activePresetId);
  return {
    activePresetId: activeExists ? activePresetId : presets[0].id,
    presets,
  };
}

function sanitizeStore(store: DashboardLayoutStore): DashboardLayoutStore {
  return {
    ...store,
    presets: store.presets.map((preset) => ({
      ...preset,
      entries: mergeDashboardLayout(preset.entries),
      usage: normalizeDashboardLayoutUsagePrefs(preset.usage),
    })),
  };
}

function persistSanitizedStore(store: DashboardLayoutStore): DashboardLayoutStore {
  const sanitized = sanitizeStore(store);
  try {
    const previous = localStorage.getItem(STORAGE_KEY);
    const next = JSON.stringify(sanitized);
    if (previous !== next) {
      localStorage.setItem(STORAGE_KEY, next);
    }
  } catch {
    // ignore quota / private mode
  }
  return sanitized;
}

export function readDashboardLayoutStore(): DashboardLayoutStore {
  const migrated = migrateLegacyStore();
  if (migrated) {
    return persistSanitizedStore(migrated);
  }
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) return createDefaultStore();
    const parsed = JSON.parse(raw) as unknown;
    const store = normalizeStore(parsed);
    return persistSanitizedStore(store ?? createDefaultStore());
  } catch {
    return createDefaultStore();
  }
}

export function writeDashboardLayoutStore(store: DashboardLayoutStore) {
  localStorage.setItem(STORAGE_KEY, JSON.stringify(store));
}

export function newLayoutPresetId(): string {
  return `layout-${Date.now().toString(36)}`;
}

export function uniquePresetName(name: string, presets: DashboardLayoutPreset[]): string {
  const base = name.trim() || "Layout";
  let candidate = base;
  let n = 2;
  while (presets.some((p) => p.name.toLowerCase() === candidate.toLowerCase())) {
    candidate = `${base} (${n})`;
    n += 1;
  }
  return candidate;
}

export function getActivePreset(store: DashboardLayoutStore): DashboardLayoutPreset {
  return store.presets.find((p) => p.id === store.activePresetId) ?? store.presets[0];
}

export function updateActivePresetEntries(
  store: DashboardLayoutStore,
  entries: DashboardLayoutEntry[],
): DashboardLayoutStore {
  const merged = mergeDashboardLayout(entries);
  return {
    ...store,
    presets: store.presets.map((p) =>
      p.id === store.activePresetId ? { ...p, entries: merged } : p,
    ),
  };
}

export function updateActivePresetUsage(
  store: DashboardLayoutStore,
  patch: Partial<DashboardLayoutUsagePrefs>,
): DashboardLayoutStore {
  return {
    ...store,
    presets: store.presets.map((p) =>
      p.id === store.activePresetId
        ? { ...p, usage: mergeDashboardLayoutUsagePrefs(p.usage, patch) }
        : p,
    ),
  };
}

export function setActivePreset(store: DashboardLayoutStore, presetId: string): DashboardLayoutStore {
  if (!store.presets.some((p) => p.id === presetId)) return store;
  return { ...store, activePresetId: presetId };
}

export function addPreset(
  store: DashboardLayoutStore,
  name: string,
  entries: DashboardLayoutEntry[],
  usage?: DashboardLayoutUsagePrefs,
): DashboardLayoutStore {
  const preset: DashboardLayoutPreset = {
    id: newLayoutPresetId(),
    name: uniquePresetName(name, store.presets),
    entries: cloneDashboardLayout(entries),
    usage: normalizeDashboardLayoutUsagePrefs(usage),
  };
  return {
    activePresetId: preset.id,
    presets: [...store.presets, preset],
  };
}

export function deletePreset(store: DashboardLayoutStore, presetId: string): DashboardLayoutStore {
  if (store.presets.length <= 1) return store;
  const presets = store.presets.filter((p) => p.id !== presetId);
  if (presets.length === 0) return store;
  const activePresetId =
    store.activePresetId === presetId ? presets[0].id : store.activePresetId;
  return { activePresetId, presets };
}

export function renamePreset(
  store: DashboardLayoutStore,
  presetId: string,
  name: string,
): DashboardLayoutStore {
  const trimmed = name.trim();
  if (!trimmed) return store;
  const others = store.presets.filter((p) => p.id !== presetId);
  const unique = uniquePresetName(trimmed, others);
  return {
    ...store,
    presets: store.presets.map((p) => (p.id === presetId ? { ...p, name: unique } : p)),
  };
}

export function moveDashboardSection(
  layout: DashboardLayoutEntry[],
  id: DashboardSectionId,
  direction: -1 | 1,
): DashboardLayoutEntry[] {
  const index = layout.findIndex((e) => e.id === id);
  if (index < 0) return layout;
  const target = index + direction;
  if (target < 0 || target >= layout.length) return layout;
  const next = [...layout];
  const [item] = next.splice(index, 1);
  next.splice(target, 0, item);
  return next;
}

export function setDashboardSectionWidth(
  layout: DashboardLayoutEntry[],
  id: DashboardSectionId,
  width: DashboardSectionWidth,
): DashboardLayoutEntry[] {
  return layout.map((e) => (e.id === id ? { ...e, width } : e));
}
