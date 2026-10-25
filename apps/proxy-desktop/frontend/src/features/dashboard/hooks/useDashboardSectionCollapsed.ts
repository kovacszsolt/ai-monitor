import { useCallback, useState } from "react";
import { useDashboardExportForceExpanded } from "@/features/dashboard/export/DashboardExportContext";

const STORAGE_KEY = "proxy-desktop.dashboard-sections";

function readAll(): Record<string, boolean> {
  try {
    const raw = localStorage.getItem(STORAGE_KEY);
    if (!raw) return {};
    const parsed = JSON.parse(raw) as Record<string, boolean>;
    return parsed && typeof parsed === "object" ? parsed : {};
  } catch {
    return {};
  }
}

function writeSection(sectionId: string, collapsed: boolean) {
  const next = { ...readAll(), [sectionId]: collapsed };
  localStorage.setItem(STORAGE_KEY, JSON.stringify(next));
}

export function useDashboardSectionCollapsed(sectionId: string, defaultCollapsed = false) {
  const forceExpanded = useDashboardExportForceExpanded();
  const [collapsed, setCollapsed] = useState(() => {
    const stored = readAll()[sectionId];
    return stored ?? defaultCollapsed;
  });

  const toggle = useCallback(() => {
    setCollapsed((prev) => {
      const next = !prev;
      writeSection(sectionId, next);
      return next;
    });
  }, [sectionId]);

  const setCollapsedPersisted = useCallback(
    (value: boolean) => {
      setCollapsed(value);
      writeSection(sectionId, value);
    },
    [sectionId],
  );

  const effectiveCollapsed = forceExpanded ? false : collapsed;

  return { collapsed: effectiveCollapsed, toggle, setCollapsed: setCollapsedPersisted };
}
