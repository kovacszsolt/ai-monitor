import { type MouseEvent } from "react";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Table, TableBody, TableCell, TableHeader, TableRow } from "@/components/ui/table";
import { formatNumber, formatOptionalTokensPerSec, formatUSD } from "@/lib/format";
import { ResizableColGroup } from "@/shared/table/ResizableColGroup";
import { ResizableTableHead } from "@/shared/table/ResizableTableHead";
import { resizableCellStyle, useResizableColumns } from "@/shared/table/useResizableColumns";
import { service } from "../../../../wailsjs/go/models";

type Row = service.ModelUsageStatsDTO;

const GEN_TIP =
  "Generation throughput: sum(output tokens) ÷ sum(decode duration). Empty when decode timing was not captured.";
const COST_TIP =
  "Estimated total USD from model route pricing (uncached input, cached input, output) × token sums. Set prices on the Models page.";

const MODEL_USAGE_COLUMNS: Record<string, number> = {
  model: 200,
  requests: 96,
  input: 88,
  cached: 88,
  output: 88,
  cost: 104,
  gen: 96,
};
const MODEL_USAGE_COLUMN_IDS = Object.keys(MODEL_USAGE_COLUMNS);

export function ModelUsageTable({
  title,
  description,
  rows,
  embedded = false,
}: {
  title: string;
  description: string;
  rows: Row[];
  embedded?: boolean;
}) {
  const { widths: colWidths, beginResize: beginColResize, totalWidth: tableTotalWidth } = useResizableColumns(
    "proxy-desktop.model-usage-table-columns",
    MODEL_USAGE_COLUMNS,
  );

  function resizeCol(columnId: string) {
    return (event: MouseEvent, width: number) => beginColResize(event, columnId, width);
  }

  const tableBody = rows.length === 0 ? (
    <p className="text-sm text-muted-foreground">No requests recorded yet.</p>
  ) : (
    <div className="overflow-x-auto">
      <Table className="table-fixed" style={{ width: tableTotalWidth, minWidth: "100%" }}>
        <ResizableColGroup columnIds={MODEL_USAGE_COLUMN_IDS} widths={colWidths} />
        <TableHeader>
          <TableRow>
            <ResizableTableHead width={colWidths.model} onResizeStart={resizeCol("model")}>
              Model
            </ResizableTableHead>
            <ResizableTableHead
              width={colWidths.requests}
              onResizeStart={resizeCol("requests")}
              className="text-right"
            >
              Requests
            </ResizableTableHead>
            <ResizableTableHead width={colWidths.input} onResizeStart={resizeCol("input")} className="text-right">
              Input
            </ResizableTableHead>
            <ResizableTableHead width={colWidths.cached} onResizeStart={resizeCol("cached")} className="text-right">
              Cached
            </ResizableTableHead>
            <ResizableTableHead width={colWidths.output} onResizeStart={resizeCol("output")} className="text-right">
              Output
            </ResizableTableHead>
            <ResizableTableHead
              width={colWidths.cost}
              onResizeStart={resizeCol("cost")}
              className="text-right"
              title={COST_TIP}
            >
              Cost (USD)
            </ResizableTableHead>
            <ResizableTableHead
              width={colWidths.gen}
              onResizeStart={resizeCol("gen")}
              className="text-right"
              title={GEN_TIP}
            >
              Gen tok/s
            </ResizableTableHead>
          </TableRow>
        </TableHeader>
        <TableBody>
          {rows.map((r) => (
            <TableRow key={r.name}>
              <TableCell className="truncate px-2 font-medium" style={resizableCellStyle(colWidths.model)}>
                {r.name}
              </TableCell>
              <TableCell
                className="truncate px-2 text-right font-mono text-sm font-semibold text-primary"
                style={resizableCellStyle(colWidths.requests)}
              >
                {formatNumber(r.request_count)}
              </TableCell>
              <TableCell
                className="truncate px-2 text-right font-mono text-xs text-muted-foreground"
                style={resizableCellStyle(colWidths.input)}
              >
                {formatNumber(r.input_tokens)}
              </TableCell>
              <TableCell
                className="truncate px-2 text-right font-mono text-xs text-muted-foreground"
                style={resizableCellStyle(colWidths.cached)}
              >
                {formatNumber(r.cached_input_tokens)}
              </TableCell>
              <TableCell
                className="truncate px-2 text-right font-mono text-xs text-muted-foreground"
                style={resizableCellStyle(colWidths.output)}
              >
                {formatNumber(r.output_tokens)}
              </TableCell>
              <TableCell
                className="truncate px-2 text-right font-mono text-sm font-semibold text-primary"
                style={resizableCellStyle(colWidths.cost)}
              >
                {formatUSD(r.has_total_cost, r.total_cost_usd)}
              </TableCell>
              <TableCell
                className="truncate px-2 text-right font-mono text-sm font-semibold text-primary"
                style={resizableCellStyle(colWidths.gen)}
              >
                {formatOptionalTokensPerSec(r.has_generation_tps, r.generation_tps)}
              </TableCell>
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </div>
  );

  if (embedded) {
    return (
      <div className="space-y-2 border-t border-border pt-4">
        <div>
          <h3 className="text-base font-semibold">{title}</h3>
          <p className="text-sm text-muted-foreground">{description}</p>
        </div>
        {tableBody}
      </div>
    );
  }

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-base">{title}</CardTitle>
        <CardDescription>{description}</CardDescription>
      </CardHeader>
      <CardContent>
        {tableBody}
      </CardContent>
    </Card>
  );
}
