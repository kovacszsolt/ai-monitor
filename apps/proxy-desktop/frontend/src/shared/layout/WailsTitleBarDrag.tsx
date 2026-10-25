import type { MouseEvent, ReactNode } from "react";
import { useRef } from "react";
import { cn } from "@/lib/utils";
import { WindowToggleMaximise } from "../../../wailsjs/runtime/runtime";

const ZOOM_TOGGLE_DEBOUNCE_MS = 400;

type Props = {
  children?: ReactNode;
  className?: string;
  "aria-hidden"?: boolean;
};

/** Draggable title-bar region; double-click toggles window zoom / maximise (native title-bar behavior). */
export function WailsTitleBarDrag({ children, className, ...rest }: Props) {
  const lastZoomToggleAt = useRef(0);

  const toggleZoom = (e: MouseEvent<HTMLDivElement>) => {
    e.preventDefault();
    e.stopPropagation();
    window.getSelection()?.removeAllRanges();
    const now = Date.now();
    if (now - lastZoomToggleAt.current < ZOOM_TOGGLE_DEBOUNCE_MS) return;
    lastZoomToggleAt.current = now;
    WindowToggleMaximise();
  };

  // `-webkit-app-region: drag` often swallows `dblclick`; `mousedown` with detail === 2 still fires on macOS.
  const onMouseDown = (e: MouseEvent<HTMLDivElement>) => {
    if (e.button !== 0 || e.detail !== 2) return;
    toggleZoom(e);
  };

  const onDoubleClick = (e: MouseEvent<HTMLDivElement>) => {
    toggleZoom(e);
  };

  return (
    <div
      className={cn("wails-drag", className)}
      onMouseDown={onMouseDown}
      onDoubleClick={onDoubleClick}
      {...rest}
    >
      {children}
    </div>
  );
}
