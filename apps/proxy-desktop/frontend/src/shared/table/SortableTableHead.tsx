import { ArrowDown, ArrowUp, ArrowUpDown } from "lucide-react";
import { TableHead } from "@/components/ui/table";
import { HoverTooltip } from "@/shared/components/HoverTooltip";
import { cn } from "@/lib/utils";

type Props = React.ComponentProps<"th"> & {
  sortDirection?: "asc" | "desc" | false;
  onSortClick?: () => void;
};

export function SortableTableHead({
  sortDirection = false,
  onSortClick,
  className,
  children,
  ...props
}: Props) {
  const ariaSort = sortDirection ? (sortDirection === "asc" ? "ascending" : "descending") : undefined;
  const alignEnd = Boolean(className?.includes("text-right"));
  const sortTooltip =
    sortDirection === "asc"
      ? "Sorted ascending"
      : sortDirection === "desc"
        ? "Sorted descending"
        : "Click to sort";

  return (
    <TableHead className={cn("px-2", className)} aria-sort={ariaSort} {...props}>
      <button
        type="button"
        className={cn(
          "flex w-full min-w-0 items-center gap-0.5 truncate font-medium hover:text-foreground",
          alignEnd ? "justify-end text-right" : "text-left",
          sortDirection && "text-foreground",
        )}
        onClick={onSortClick}
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
    </TableHead>
  );
}
