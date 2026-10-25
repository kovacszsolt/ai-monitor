import { useCallback, useEffect, useRef, useState, type CSSProperties, type MouseEvent as ReactMouseEvent } from "react";

const DEFAULT_MIN_WIDTH = 48;

function loadWidths(storageKey: string, defaults: Record<string, number>): Record<string, number> {
  try {
    const raw = localStorage.getItem(storageKey);
    if (!raw) {
      return { ...defaults };
    }
    const parsed = JSON.parse(raw) as Record<string, number>;
    return { ...defaults, ...parsed };
  } catch {
    return { ...defaults };
  }
}

function saveWidths(storageKey: string, widths: Record<string, number>) {
  try {
    localStorage.setItem(storageKey, JSON.stringify(widths));
  } catch {
    // ignore quota / private mode
  }
}

export function useResizableColumns(
  storageKey: string,
  defaults: Record<string, number>,
  minWidth = DEFAULT_MIN_WIDTH,
) {
  const [widths, setWidths] = useState<Record<string, number>>(() => loadWidths(storageKey, defaults));
  const dragRef = useRef<{ columnId: string; startX: number; startWidth: number } | null>(null);

  const setColumnWidth = useCallback(
    (columnId: string, width: number) => {
      const nextWidth = Math.max(minWidth, Math.round(width));
      setWidths((prev) => {
        const next = { ...prev, [columnId]: nextWidth };
        saveWidths(storageKey, next);
        return next;
      });
    },
    [minWidth, storageKey],
  );

  const beginResize = useCallback(
    (event: ReactMouseEvent, columnId: string, currentWidth: number) => {
      event.preventDefault();
      event.stopPropagation();
      dragRef.current = { columnId, startX: event.clientX, startWidth: currentWidth };
    },
    [],
  );

  useEffect(() => {
    const onMove = (event: MouseEvent) => {
      const drag = dragRef.current;
      if (!drag) {
        return;
      }
      const delta = event.clientX - drag.startX;
      setColumnWidth(drag.columnId, drag.startWidth + delta);
    };
    const onUp = () => {
      dragRef.current = null;
    };
    document.addEventListener("mousemove", onMove);
    document.addEventListener("mouseup", onUp);
    return () => {
      document.removeEventListener("mousemove", onMove);
      document.removeEventListener("mouseup", onUp);
    };
  }, [setColumnWidth]);

  const totalWidth = Object.values(widths).reduce((sum, w) => sum + w, 0);

  return { widths, beginResize, totalWidth };
}

export function resizableCellStyle(width: number): CSSProperties {
  return {
    width,
    minWidth: width,
    maxWidth: width,
  };
}
