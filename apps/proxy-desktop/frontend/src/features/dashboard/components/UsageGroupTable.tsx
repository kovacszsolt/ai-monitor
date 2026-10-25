import { Card, CardContent, CardDescription, CardHeader } from "@/components/ui/card";
import { CollapsibleCardHeader } from "@/shared/components/CollapsibleCardHeader";
import { DASHBOARD_SECTION_IDS } from "@/features/dashboard/dashboardSectionIds";
import { useDashboardSectionCollapsed } from "@/features/dashboard/hooks/useDashboardSectionCollapsed";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { formatNumber } from "@/lib/format";
import { service } from "../../../../wailsjs/go/models";

type Row = service.UsageGroupStatsDTO;

export function UsageGroupTable({
  title,
  description,
  rows,
  emptyLabel,
  selectedTag,
  onSelectTag,
  collapsible = false,
}: {
  title: string;
  description: string;
  rows: Row[];
  emptyLabel?: string;
  selectedTag?: string;
  onSelectTag?: (tag: string) => void;
  collapsible?: boolean;
}) {
  const { collapsed, toggle } = useDashboardSectionCollapsed(DASHBOARD_SECTION_IDS.usageByProxyTag);
  const isCollapsed = collapsible && collapsed;

  return (
    <Card>
      <CardHeader className="pb-2">
        {collapsible ? (
          <CollapsibleCardHeader title={title} collapsed={isCollapsed} onToggle={toggle} />
        ) : (
          <h3 className="text-base font-semibold">{title}</h3>
        )}
        {!isCollapsed && <CardDescription className="pt-1">{description}</CardDescription>}
      </CardHeader>
      {!isCollapsed && (
      <CardContent>
        {rows.length === 0 ? (
          <p className="text-sm text-muted-foreground">
            {emptyLabel ?? "No requests recorded yet."}
          </p>
        ) : (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Name</TableHead>
                <TableHead className="text-right">Requests</TableHead>
                <TableHead className="text-right">Input</TableHead>
                <TableHead className="text-right">Cached</TableHead>
                <TableHead className="text-right">Output</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {rows.map((r) => (
                <TableRow
                  key={r.name}
                  className={
                    onSelectTag
                      ? "cursor-pointer hover:bg-muted/50 " +
                        (selectedTag === r.name ? "bg-muted/40" : "")
                      : undefined
                  }
                  onClick={onSelectTag ? () => onSelectTag(r.name) : undefined}
                >
                  <TableCell className="max-w-[200px] truncate font-medium">{r.name}</TableCell>
                  <TableCell className="text-right font-mono text-sm font-semibold text-primary">
                    {formatNumber(r.request_count)}
                  </TableCell>
                  <TableCell className="text-right font-mono text-xs text-muted-foreground">
                    {formatNumber(r.input_tokens)}
                  </TableCell>
                  <TableCell className="text-right font-mono text-xs text-muted-foreground">
                    {formatNumber(r.cached_input_tokens)}
                  </TableCell>
                  <TableCell className="text-right font-mono text-xs text-muted-foreground">
                    {formatNumber(r.output_tokens)}
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        )}
      </CardContent>
      )}
    </Card>
  );
}
