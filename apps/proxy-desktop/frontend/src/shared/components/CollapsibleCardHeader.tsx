import type { MouseEvent, ReactNode } from "react";
import { ChevronDown } from "lucide-react";
import { CardTitle } from "@/components/ui/card";
import { HoverTooltip } from "@/shared/components/HoverTooltip";
import { cn } from "@/lib/utils";

type Props = {
  title: string;
  collapsed: boolean;
  onToggle: () => void;
  /** Shown under the title; part of the same toggle control (full header hit target). */
  description?: ReactNode;
  className?: string;
};

function handleToggleClick(e: MouseEvent<HTMLButtonElement>, onToggle: () => void) {
  // A double-click fires two clicks; ignore the second so the section toggles once (native disclosure behavior).
  if (e.detail > 1) {
    e.preventDefault();
    return;
  }
  onToggle();
}

export function CollapsibleCardHeader({ title, collapsed, onToggle, description, className }: Props) {
  return (
    <button
      type="button"
      onClick={(e) => handleToggleClick(e, onToggle)}
      aria-expanded={!collapsed}
      className={cn(
        "flex w-full flex-col items-stretch gap-1 rounded-md py-0.5 text-left transition-colors hover:bg-muted/50",
        className,
      )}
    >
      <span className="flex w-full items-center justify-between gap-2">
        <CardTitle className="text-base font-semibold">{title}</CardTitle>
        <HoverTooltip label={collapsed ? "Expand section" : "Collapse section"} side="left">
          <span className="inline-flex shrink-0 rounded-sm p-0.5" aria-hidden>
            <ChevronDown
              className={cn(
                "h-4 w-4 text-muted-foreground transition-transform duration-200",
                collapsed && "-rotate-90",
              )}
            />
          </span>
        </HoverTooltip>
      </span>
      {description ? (
        <span className="text-sm font-normal text-muted-foreground">{description}</span>
      ) : null}
    </button>
  );
}
