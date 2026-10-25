import { useCallback, useEffect, useMemo, useState } from "react";
import {
  addPreset,
  cloneDashboardLayout,
  DEFAULT_DASHBOARD_LAYOUT,
  deletePreset,
  getActivePreset,
  type DashboardLayoutEntry,
  type DashboardLayoutPreset,
  type DashboardLayoutStore,
  type DashboardLayoutUsagePrefs,
  type DashboardSectionId,
  type DashboardSectionWidth,
  mergeDashboardLayout,
  moveDashboardSection,
  readDashboardLayoutStore,
  renamePreset,
  setActivePreset,
  setDashboardSectionWidth,
  updateActivePresetEntries,
  updateActivePresetUsage,
  writeDashboardLayoutStore,
} from "@/features/dashboard/dashboardLayout";
import { normalizeDashboardLayoutUsagePrefs } from "@/features/dashboard/dashboardLayoutUsagePrefs";

export function useDashboardLayout() {
  const [store, setStore] = useState<DashboardLayoutStore>(() => readDashboardLayoutStore());

  // Re-merge after hot reload when in-memory store still has legacy section ids.
  useEffect(() => {
    setStore((prev) => {
      const next: DashboardLayoutStore = {
        ...prev,
        presets: prev.presets.map((preset) => ({
          ...preset,
          entries: mergeDashboardLayout(preset.entries),
          usage: normalizeDashboardLayoutUsagePrefs(preset.usage),
        })),
      };
      const prevEntries = JSON.stringify(prev.presets.map((p) => p.entries));
      const nextEntries = JSON.stringify(next.presets.map((p) => p.entries));
      if (prevEntries === nextEntries) return prev;
      writeDashboardLayoutStore(next);
      return next;
    });
  }, []);

  const activePreset = useMemo(() => getActivePreset(store), [store]);
  const layout = useMemo(() => mergeDashboardLayout(activePreset.entries), [activePreset.entries]);
  const usagePrefs = useMemo(
    () => normalizeDashboardLayoutUsagePrefs(activePreset.usage),
    [activePreset.usage],
  );

  const persistStore = useCallback((updater: (prev: DashboardLayoutStore) => DashboardLayoutStore) => {
    setStore((prev) => {
      const next = updater(prev);
      writeDashboardLayoutStore(next);
      return next;
    });
  }, []);

  const persistEntries = useCallback(
    (entries: DashboardLayoutEntry[]) => {
      persistStore((prev) => updateActivePresetEntries(prev, entries));
    },
    [persistStore],
  );

  const switchPreset = useCallback(
    (presetId: string) => {
      persistStore((prev) => setActivePreset(prev, presetId));
    },
    [persistStore],
  );

  const moveSection = useCallback(
    (id: DashboardSectionId, direction: -1 | 1) => {
      persistEntries(moveDashboardSection(layout, id, direction));
    },
    [layout, persistEntries],
  );

  const setSectionWidth = useCallback(
    (id: DashboardSectionId, width: DashboardSectionWidth) => {
      persistEntries(setDashboardSectionWidth(layout, id, width));
    },
    [layout, persistEntries],
  );

  const resetLayout = useCallback(() => {
    persistEntries(cloneDashboardLayout(DEFAULT_DASHBOARD_LAYOUT));
  }, [persistEntries]);

  const createPreset = useCallback(
    (name: string) => {
      persistStore((prev) => addPreset(prev, name, cloneDashboardLayout(DEFAULT_DASHBOARD_LAYOUT)));
    },
    [persistStore],
  );

  const savePresetAs = useCallback(
    (name: string) => {
      persistStore((prev) => {
        const active = getActivePreset(prev);
        return addPreset(prev, name, active.entries, active.usage);
      });
    },
    [persistStore],
  );

  const copyLayout = useCallback(() => {
    persistStore((prev) => {
      const active = getActivePreset(prev);
      return addPreset(prev, `Copy of ${active.name}`, active.entries, active.usage);
    });
  }, [persistStore]);

  const updateUsagePrefs = useCallback(
    (patch: Partial<DashboardLayoutUsagePrefs>) => {
      persistStore((prev) => updateActivePresetUsage(prev, patch));
    },
    [persistStore],
  );

  const removePreset = useCallback(
    (presetId: string) => {
      persistStore((prev) => deletePreset(prev, presetId));
    },
    [persistStore],
  );

  const renameActiveOrPreset = useCallback(
    (presetId: string, name: string) => {
      persistStore((prev) => renamePreset(prev, presetId, name));
    },
    [persistStore],
  );

  const presetSummaries: Pick<DashboardLayoutPreset, "id" | "name">[] = useMemo(
    () => store.presets.map((p) => ({ id: p.id, name: p.name })),
    [store.presets],
  );

  return {
    layout,
    presets: presetSummaries,
    activePresetId: store.activePresetId,
    activePresetName: activePreset.name,
    switchPreset,
    createPreset,
    savePresetAs,
    copyLayout,
    removePreset,
    renamePreset: renameActiveOrPreset,
    moveSection,
    setSectionWidth,
    resetLayout,
    usagePrefs,
    updateUsagePrefs,
  };
}
