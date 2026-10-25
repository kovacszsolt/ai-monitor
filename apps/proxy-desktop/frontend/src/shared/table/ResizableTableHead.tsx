import { ArrowDown, ArrowUp, ArrowUpDown } from "lucide-react";
import { TableHead } from "@/components/ui/table";
import { HoverTooltip } from "@/shared/components/HoverTooltip";
import { cn } from "@/lib/utils";

type Props = React.ComponentProps<"th"> & {
  width: number;
  onResizeStart: (event: React.MouseEvent, width: number) => void;
  resizable?: boolean;
  sortable?: boolean;
  sortDirection?: "asc" | "desc" | false;
  onSortClick?: () => void;
};

export function ResizableTableHead({
  width,
  onResizeStart,
  resizable = true,
  sortable = false,
  sortDirection = false,
  onSortClick,
  className,
  children,
  ...props
}: Props) {
  const ariaSort =
    sortable && sortDirection ? (sortDirection === "asc" ? "ascending" : "descending") : undefined;
  const alignEnd = Boolean(className?.includes("text-right"));
  const sortTooltip =
    sortDirection === "asc"
      ? "Sorted ascending"
      : sortDirection === "desc"
        ? "Sorted descending"
        : "Click to sort";

  return (
    <TableHead
      className={cn("relative overflow-hidden px-2", className)}
      style={{ width, minWidth: width, maxWidth: width }}
      aria-sort={ariaSort}
      {...props}
    >
      {sortable ? (
        <button
          type="button"
          className={cn(
            "flex w-full min-w-0 items-center gap-0.5 truncate pr-1 font-medium hover:text-foreground",
            alignEnd ? "justify-end text-right" : "text-left",
            sortDirection && "text-foreground",
          )}
          onClick={(e) => {
            e.stopPropagation();
            onSortClick?.();
          }}
        >
          <span className="min-w-0 truncate">{children}</span>
          <HoverTooltip label={sortTooltip} side="bottom">
            <span className="inline-flex shrink-0 rounded-sm" aria-hidden>
              {sortDirection === "asc" ? (
                <ArrowUp className="h-3.5 w-3.5 opacity-80" />
              ) : sortDirection === "desc" ? (
                <ArrowDown className="h-3.5 w-3.5 opacity-80" />
              ) : (
                <ArrowUpDown className="h-3.5 w-3.5 opacity-40" />
              )}
            </span>
          </HoverTooltip>
        </button>
      ) : (
        <span className="block truncate pr-1">{children}</span>
      )}
      {resizable && (
        <HoverTooltip label="Resize column" side="bottom" className="absolute -right-0.5 top-0 z-10 h-full">
          <div
            role="separator"
            aria-orientation="vertical"
            aria-label="Resize column"
            className="h-full w-2 cursor-col-resize touch-none select-none hover:bg-primary/40"
            onMouseDown={(e) => onResizeStart(e, width)}
            onClick={(e) => e.stopPropagation()}
            onPointerDown={(e) => e.stopPropagation()}
          />
        </HoverTooltip>
      )}
    </TableHead>
  );
}
