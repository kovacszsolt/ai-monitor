import type { ReactNode } from "react";
import { ArrowDown, ArrowUp, Columns2, Columns3, Square } from "lucide-react";
import { Button } from "@/components/ui/button";
import {
  DASHBOARD_SECTION_LABELS,
  dashboardSectionColSpanClass,
  type DashboardSectionId,
  type DashboardSectionWidth,
} from "@/features/dashboard/dashboardLayout";
import { HoverTooltip } from "@/shared/components/HoverTooltip";
import { cn } from "@/lib/utils";

type Props = {
  sectionId: DashboardSectionId;
  width: DashboardSectionWidth;
  editMode: boolean;
  canMoveUp: boolean;
  canMoveDown: boolean;
  onMoveUp: () => void;
  onMoveDown: () => void;
  onSetWidth: (width: DashboardSectionWidth) => void;
  children: ReactNode;
};

export function DashboardSectionShell({
  sectionId,
  width,
  editMode,
  canMoveUp,
  canMoveDown,
  onMoveUp,
  onMoveDown,
  onSetWidth,
  children,
}: Props) {
  return (
    <div
      className={cn(
        "min-w-0",
        dashboardSectionColSpanClass(width),
        editMode && "rounded-lg ring-1 ring-dashed ring-border/80",
      )}
    >
      {editMode && (
        <div className="mb-2 flex flex-wrap items-center gap-1 rounded-t-lg border border-b-0 border-dashed border-border bg-muted/40 px-2 py-1.5">
          <span className="me-auto text-xs font-medium text-muted-foreground">
            {DASHBOARD_SECTION_LABELS[sectionId]}
          </span>
          <HoverTooltip label="Move section up">
            <Button
              type="button"
              size="sm"
              variant="ghost"
              className="h-7 w-7 p-0"
              disabled={!canMoveUp}
              onClick={onMoveUp}
              aria-label="Move section up"
            >
              <ArrowUp className="h-3.5 w-3.5" />
            </Button>
          </HoverTooltip>
          <HoverTooltip label="Move section down">
            <Button
              type="button"
              size="sm"
              variant="ghost"
              className="h-7 w-7 p-0"
              disabled={!canMoveDown}
              onClick={onMoveDown}
              aria-label="Move section down"
            >
              <ArrowDown className="h-3.5 w-3.5" />
            </Button>
          </HoverTooltip>
          <span className="mx-0.5 h-4 w-px bg-border" aria-hidden />
          <HoverTooltip label="One-third width">
            <Button
              type="button"
              size="sm"
              variant={width === "third" ? "default" : "ghost"}
              className="h-7 gap-1 px-2 text-xs"
              onClick={() => onSetWidth("third")}
              aria-label="One-third width"
            >
              <Columns3 className="h-3.5 w-3.5" />
              ⅓
            </Button>
          </HoverTooltip>
          <HoverTooltip label="Half width">
            <Button
              type="button"
              size="sm"
              variant={width === "half" ? "default" : "ghost"}
              className="h-7 gap-1 px-2 text-xs"
              onClick={() => onSetWidth("half")}
              aria-label="Half width"
            >
              <Columns2 className="h-3.5 w-3.5" />
              ½
            </Button>
          </HoverTooltip>
          <HoverTooltip label="Full width">
            <Button
              type="button"
              size="sm"
              variant={width === "full" ? "default" : "ghost"}
              className="h-7 gap-1 px-2 text-xs"
              onClick={() => onSetWidth("full")}
              aria-label="Full width"
            >
              <Square className="h-3.5 w-3.5" />
              Full
            </Button>
          </HoverTooltip>
        </div>
      )}
      <div className={cn(editMode && "px-0.5 pb-0.5")}>{children}</div>
    </div>
  );
}
