import type { ReactNode } from "react";
import { cn } from "@/lib/utils";

const tooltipBaseClassName =
  "pointer-events-none absolute z-[100] hidden whitespace-nowrap rounded-md border border-border bg-card px-2 py-1 text-xs font-medium text-foreground shadow-md group-hover:block group-focus-within:block";

const sideClassName = {
  top: "bottom-full left-1/2 mb-2 -translate-x-1/2",
  right: "left-full top-1/2 ml-2 -translate-y-1/2",
  bottom: "top-full left-1/2 mt-2 -translate-x-1/2",
  left: "right-full top-1/2 mr-2 -translate-y-1/2",
} as const;

export type HoverTooltipSide = keyof typeof sideClassName;

export function HoverTooltip({
  label,
  children,
  className,
  side = "top",
}: {
  label: string;
  children: ReactNode;
  className?: string;
  side?: HoverTooltipSide;
}) {
  return (
    <div className={cn("group relative inline-flex max-w-full", className)}>
      {children}
      <span role="tooltip" className={cn(tooltipBaseClassName, sideClassName[side])}>
        {label}
      </span>
    </div>
  );
}
