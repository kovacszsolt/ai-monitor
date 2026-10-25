import { Fragment, type MouseEvent } from "react";
import { MessageSquareText, Trash2 } from "lucide-react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Table, TableBody, TableCell, TableHeader, TableRow } from "@/components/ui/table";
import { RequestContentDialog } from "@/features/requests/components/RequestContentDialog";
import { useRequestsList } from "@/features/requests/hooks/useRequestsList";
import {
  formatCacheHitRate,
  formatDateTime,
  formatDuration,
  formatNumber,
  formatOptionalTokensPerSec,
  formatRequestTokenCount,
} from "@/lib/format";
import { HoverTooltip } from "@/shared/components/HoverTooltip";
import { EmptyState, ErrorState, LoadingState } from "@/shared/components/PageStates";
import { ResizableColGroup } from "@/shared/table/ResizableColGroup";
import { ResizableTableHead } from "@/shared/table/ResizableTableHead";
import { resizableCellStyle, useResizableColumns } from "@/shared/table/useResizableColumns";
import { service } from "../../../../wailsjs/go/models";

const STATES = ["", "started", "completed", "failed", "cancelled", "interrupted"] as const;

const REQUEST_TABLE_COLUMNS: Record<string, number> = {
  select: 40,
  started: 152,
  project: 108,
  model: 128,
  tag: 100,
  state: 104,
  error: 200,
  input: 80,
  cached: 80,
  output: 80,
  ttft: 72,
  duration: 88,
  gen: 80,
  prompt: 80,
  eff: 72,
  cache: 80,
  actions: 80,
};
const REQUEST_COLUMN_IDS = Object.keys(REQUEST_TABLE_COLUMNS);

const TIPS = {
  uncached: "Input tokens not served from cache (stored by ai-proxy).",
  gen: "Output tokens per second during decode. ~ prefix = estimated from stream chunk timing.",
  prompt: "Uncached prompt tokens per second during prefill (provider timing when available).",
  eff: "(Input + output) ÷ total duration. Cached tokens are not added again.",
  cache: "Cached input ÷ total input × 100.",
  ttft: "Time to first token (content or usage chunk), not prefill duration.",
} as const;

function stateBadge(state: string) {
  const variant =
    state === "completed" ? "success" : state === "failed" ? "destructive" : "secondary";
  return <Badge variant={variant}>{state}</Badge>;
}

export function RequestsPage() {
  const {
    page,
    setPage,
    projectID,
    setProjectID,
    state,
    setState,
    tag,
    setTag,
    model,
    setModel,
    expandedId,
    setExpandedId,
    deletingId,
    selectedIds,
    setSelectedIds,
    moveTag,
    setMoveTag,
    contentRequest,
    setContentRequest,
    projectsQ,
    tagOptionsQ,
    tagBulkEnabled,
    modelOptionsQ,
    listQ,
    sortDirection,
    toggleSort,
    deleteMut,
    totalPages,
    allPageSelected,
    moveTagMut,
    bulkDeleteMut,
    toggleRowSelected,
    toggleSelectAllOnPage,
  } = useRequestsList();

  const tableColSpan = REQUEST_COLUMN_IDS.length;

  const { widths: colWidths, beginResize: beginColResize, totalWidth: tableTotalWidth } = useResizableColumns(
    "proxy-desktop.requests-table-columns",
    REQUEST_TABLE_COLUMNS,
  );

  function resizeCol(columnId: string) {
    return (event: MouseEvent, width: number) => beginColResize(event, columnId, width);
  }

  return (
    <div className="flex min-h-full flex-col">
      <div className="border-b border-outline-variant/30 bg-surface-container-lowest px-space-md py-space-sm shadow-sm">
        <div className="flex flex-wrap items-center justify-between gap-space-md">
          <div className="flex items-center gap-space-md">
            <span className="font-headline-md text-headline-md tracking-tight text-on-surface">Requests</span>
          </div>
        </div>
        <p className="mt-1 text-body-sm text-on-surface-variant">
          Request metadata from ai_requests. Use the message button to view prompt and response when ai-proxy has{" "}
          <code className="text-code-inline">capture.store_content: true</code>.
        </p>
      </div>
      <div className="flex flex-1 flex-col space-y-4 p-space-md">
        <div className="flex flex-wrap gap-3">
          <label className="flex flex-col gap-1 text-xs text-muted-foreground">
            Project
            <select
              className="h-9 min-w-[160px] rounded-md border border-border bg-card px-3 text-sm text-foreground"
              value={projectID}
              onChange={(e) => {
                setProjectID(e.target.value);
                setPage(1);
              }}
            >
              <option value="">All projects</option>
              {(projectsQ.data ?? []).map((p) => (
                <option key={p.id} value={p.id}>{p.name}</option>
              ))}
            </select>
          </label>
          <label className="flex flex-col gap-1 text-xs text-muted-foreground">
            State
            <div className="flex gap-2">
              <select
                className="h-9 min-w-[140px] rounded-md border border-border bg-card px-3 text-sm text-foreground"
                value={state}
                onChange={(e) => {
                  setState(e.target.value);
                  setPage(1);
                }}
              >
                {STATES.map((s) => (
                  <option key={s || "all"} value={s}>{s === "" ? "All states" : s}</option>
                ))}
              </select>
              <Button
                type="button"
                variant={state === "failed" ? "default" : "outline"}
                size="sm"
                className="h-9 shrink-0"
                onClick={() => {
                  setState((s) => (s === "failed" ? "" : "failed"));
                  setPage(1);
                }}
              >
                Errors only
              </Button>
            </div>
          </label>
          <label className="flex flex-col gap-1 text-xs text-muted-foreground">
            Model
            <select
              className="h-9 min-w-[180px] max-w-[280px] rounded-md border border-border bg-card px-3 text-sm text-foreground"
              value={model}
              onChange={(e) => {
                setModel(e.target.value);
                setPage(1);
              }}
            >
              <option value="">All models</option>
              {(modelOptionsQ.data ?? []).map((m) => (
                <option key={m} value={m}>{m}</option>
              ))}
            </select>
          </label>
          {tagOptionsQ.data?.enabled && (
            <label className="flex flex-col gap-1 text-xs text-muted-foreground">
              Tag
              <select
                className="h-9 min-w-[160px] rounded-md border border-border bg-card px-3 text-sm text-foreground"
                value={tag}
                onChange={(e) => {
                  setTag(e.target.value);
                  setPage(1);
                }}
              >
                <option value="">All tags</option>
                {tagOptionsQ.data.has_untagged && (
                  <option value={tagOptionsQ.data.no_tag_label}>{tagOptionsQ.data.no_tag_label}</option>
                )}
                {(tagOptionsQ.data.tags ?? []).map((t) => (
                  <option key={t} value={t}>{t}</option>
                ))}
              </select>
            </label>
          )}
        </div>

        {selectedIds.size > 0 && (
          <div className="flex flex-wrap items-end gap-3 rounded-md border border-border bg-muted/40 p-3">
            <p className="text-sm text-muted-foreground">
              {formatNumber(selectedIds.size)} selected
            </p>
            {tagBulkEnabled && (
              <>
                <label className="flex min-w-[200px] flex-1 flex-col gap-1 text-xs text-muted-foreground sm:max-w-xs">
                  Move to tag
                  <Input
                    value={moveTag}
                    onChange={(e) => setMoveTag(e.target.value)}
                    placeholder="e.g. staging"
                    maxLength={256}
                    className="h-9"
                  />
                </label>
                <Button
                  size="sm"
                  disabled={moveTagMut.isPending || bulkDeleteMut.isPending}
                  onClick={() => moveTagMut.mutate()}
                >
                  Apply tag
                </Button>
              </>
            )}
            <Button
              variant="destructive"
              size="sm"
              disabled={bulkDeleteMut.isPending || moveTagMut.isPending}
              onClick={() => bulkDeleteMut.mutate()}
            >
              <Trash2 className="mr-1.5 h-4 w-4" />
              Delete selected
            </Button>
            <Button variant="ghost" size="sm" onClick={() => setSelectedIds(new Set())}>
              Clear selection
            </Button>
          </div>
        )}

        {listQ.isLoading && <LoadingState />}
        {listQ.error && <ErrorState message={(listQ.error as Error).message} onRetry={() => listQ.refetch()} />}
        {!listQ.isLoading && !listQ.error && listQ.data?.items.length === 0 && (
          <EmptyState title="No requests match your filters." />
        )}

        {listQ.data && listQ.data.items.length > 0 && (
          <>
            <div className="overflow-x-auto">
              <Table className="table-fixed" style={{ width: tableTotalWidth, minWidth: "100%" }}>
                <ResizableColGroup columnIds={REQUEST_COLUMN_IDS} widths={colWidths} />
                <TableHeader>
                  <TableRow>
                    <ResizableTableHead
                      width={colWidths.select}
                      resizable={false}
                      onResizeStart={() => {}}
                      className="px-1"
                    >
                      <input
                        type="checkbox"
                        className="h-4 w-4 rounded border-border accent-primary"
                        aria-label="Select all on page"
                        checked={allPageSelected}
                        onChange={toggleSelectAllOnPage}
                      />
                    </ResizableTableHead>
                    <ResizableTableHead
                      width={colWidths.started}
                      onResizeStart={resizeCol("started")}
                      sortable
                      sortDirection={sortDirection("started")}
                      onSortClick={() => toggleSort("started")}
                    >
                      Started
                    </ResizableTableHead>
                    <ResizableTableHead
                      width={colWidths.project}
                      onResizeStart={resizeCol("project")}
                      sortable
                      sortDirection={sortDirection("project")}
                      onSortClick={() => toggleSort("project")}
                    >
                      Project
                    </ResizableTableHead>
                    <ResizableTableHead
                      width={colWidths.model}
                      onResizeStart={resizeCol("model")}
                      sortable
                      sortDirection={sortDirection("model")}
                      onSortClick={() => toggleSort("model")}
                    >
                      Model
                    </ResizableTableHead>
                    <ResizableTableHead
                      width={colWidths.tag}
                      onResizeStart={resizeCol("tag")}
                      title="X-AI-Proxy-Tag"
                      sortable={tagBulkEnabled}
                      sortDirection={tagBulkEnabled ? sortDirection("tag") : false}
                      onSortClick={tagBulkEnabled ? () => toggleSort("tag") : undefined}
                    >
                      Tag
                    </ResizableTableHead>
                    <ResizableTableHead
                      width={colWidths.state}
                      onResizeStart={resizeCol("state")}
                      sortable
                      sortDirection={sortDirection("state")}
                      onSortClick={() => toggleSort("state")}
                    >
                      State
                    </ResizableTableHead>
                    <ResizableTableHead
                      width={colWidths.error}
                      onResizeStart={resizeCol("error")}
                      title="Proxy / upstream error (ai_requests.error_type and ai_request_errors)"
                    >
                      Error
                    </ResizableTableHead>
                    <ResizableTableHead
                      width={colWidths.input}
                      onResizeStart={resizeCol("input")}
                      className="text-right whitespace-nowrap"
                      sortable
                      sortDirection={sortDirection("input")}
                      onSortClick={() => toggleSort("input")}
                    >
                      Input
                    </ResizableTableHead>
                    <ResizableTableHead
                      width={colWidths.cached}
                      onResizeStart={resizeCol("cached")}
                      className="text-right whitespace-nowrap"
                      sortable
                      sortDirection={sortDirection("cached")}
                      onSortClick={() => toggleSort("cached")}
                    >
                      Cached
                    </ResizableTableHead>
                    <ResizableTableHead
                      width={colWidths.output}
                      onResizeStart={resizeCol("output")}
                      className="text-right whitespace-nowrap"
                      sortable
                      sortDirection={sortDirection("output")}
                      onSortClick={() => toggleSort("output")}
                    >
                      Output
                    </ResizableTableHead>
                    <ResizableTableHead
                      width={colWidths.ttft}
                      onResizeStart={resizeCol("ttft")}
                      className="text-right"
                      title={TIPS.ttft}
                      sortable
                      sortDirection={sortDirection("ttft")}
                      onSortClick={() => toggleSort("ttft")}
                    >
                      TTFT
                    </ResizableTableHead>
                    <ResizableTableHead
                      width={colWidths.duration}
                      onResizeStart={resizeCol("duration")}
                      className="text-right"
                      sortable
                      sortDirection={sortDirection("duration")}
                      onSortClick={() => toggleSort("duration")}
                    >
                      Duration
                    </ResizableTableHead>
                    <ResizableTableHead
                      width={colWidths.gen}
                      onResizeStart={resizeCol("gen")}
                      className="text-right"
                      title={TIPS.gen}
                      sortable
                      sortDirection={sortDirection("gen")}
                      onSortClick={() => toggleSort("gen")}
                    >
                      Gen
                    </ResizableTableHead>
                    <ResizableTableHead
                      width={colWidths.prompt}
                      onResizeStart={resizeCol("prompt")}
                      className="text-right"
                      title={TIPS.prompt}
                      sortable
                      sortDirection={sortDirection("prompt")}
                      onSortClick={() => toggleSort("prompt")}
                    >
                      Prompt
                    </ResizableTableHead>
                    <ResizableTableHead
                      width={colWidths.eff}
                      onResizeStart={resizeCol("eff")}
                      className="text-right"
                      title={TIPS.eff}
                      sortable
                      sortDirection={sortDirection("eff")}
                      onSortClick={() => toggleSort("eff")}
                    >
                      Eff
                    </ResizableTableHead>
                    <ResizableTableHead
                      width={colWidths.cache}
                      onResizeStart={resizeCol("cache")}
                      className="text-right"
                      title={TIPS.cache}
                      sortable
                      sortDirection={sortDirection("cache")}
                      onSortClick={() => toggleSort("cache")}
                    >
                      Cache %
                    </ResizableTableHead>
                    <ResizableTableHead
                      width={colWidths.actions}
                      resizable={false}
                      onResizeStart={() => {}}
                      aria-label="Actions"
                    />
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {listQ.data.items.map((r) => {
                    const expanded = expandedId === r.id;
                    return (
                      <Fragment key={r.id}>
                        <TableRow
                          className="cursor-pointer"
                          onClick={() => setExpandedId(expanded ? null : r.id)}
                        >
                          <TableCell
                            className="px-1"
                            style={resizableCellStyle(colWidths.select)}
                            onClick={(e) => e.stopPropagation()}
                            onPointerDown={(e) => e.stopPropagation()}
                          >
                            <input
                              type="checkbox"
                              className="h-4 w-4 rounded border-border accent-primary"
                              aria-label={`Select request ${r.id}`}
                              checked={selectedIds.has(r.id)}
                              onChange={() => toggleRowSelected(r.id)}
                            />
                          </TableCell>
                          <TableCell
                            className="truncate whitespace-nowrap px-2 text-xs text-muted-foreground"
                            style={resizableCellStyle(colWidths.started)}
                          >
                            {formatDateTime(r.started_at)}
                          </TableCell>
                          <TableCell
                            className="truncate px-2 font-medium"
                            style={resizableCellStyle(colWidths.project)}
                          >
                            {r.project_name || r.project_id}
                          </TableCell>
                          <TableCell
                            className="truncate px-2 font-mono text-xs"
                            style={resizableCellStyle(colWidths.model)}
                          >
                            {r.requested_model}
                          </TableCell>
                          <TableCell
                            className="truncate px-2 font-mono text-xs text-muted-foreground"
                            style={resizableCellStyle(colWidths.tag)}
                          >
                            {r.tag?.trim() ? r.tag : "—"}
                          </TableCell>
                          <TableCell className="truncate px-2" style={resizableCellStyle(colWidths.state)}>
                            {stateBadge(r.request_state)}
                          </TableCell>
                          <TableCell
                            className="truncate px-2 font-mono text-xs text-destructive"
                            style={resizableCellStyle(colWidths.error)}
                            title={r.error_type || undefined}
                          >
                            {r.error_type || "—"}
                          </TableCell>
                          <TableCell
                            className="truncate px-2 text-right font-mono text-xs font-medium text-foreground"
                            style={resizableCellStyle(colWidths.input)}
                          >
                            {formatRequestTokenCount(r.has_input_tokens, r.input_tokens, r.usage_status)}
                          </TableCell>
                          <TableCell
                            className="truncate px-2 text-right font-mono text-xs font-medium text-foreground"
                            style={resizableCellStyle(colWidths.cached)}
                          >
                            {formatRequestTokenCount(r.has_cached_input_tokens, r.cached_input_tokens, r.usage_status)}
                          </TableCell>
                          <TableCell
                            className="truncate px-2 text-right font-mono text-xs font-medium text-foreground"
                            style={resizableCellStyle(colWidths.output)}
                          >
                            {formatRequestTokenCount(r.has_output_tokens, r.output_tokens, r.usage_status)}
                          </TableCell>
                          <TableCell
                            className="truncate px-2 text-right text-xs"
                            style={resizableCellStyle(colWidths.ttft)}
                          >
                            {r.has_ttft ? formatDuration(r.ttft_ms) : "—"}
                          </TableCell>
                          <TableCell
                            className="truncate px-2 text-right text-xs"
                            style={resizableCellStyle(colWidths.duration)}
                          >
                            {r.has_duration ? formatDuration(r.duration_ms) : "—"}
                          </TableCell>
                          <TableCell
                            className="truncate px-2 text-right font-mono text-xs"
                            style={resizableCellStyle(colWidths.gen)}
                          >
                            {formatOptionalTokensPerSec(r.has_generation_tps, r.generation_tps, r.generation_speed_estimated)}
                          </TableCell>
                          <TableCell
                            className="truncate px-2 text-right font-mono text-xs"
                            style={resizableCellStyle(colWidths.prompt)}
                          >
                            {formatOptionalTokensPerSec(r.has_prompt_tps, r.prompt_tps, r.prompt_speed_estimated)}
                          </TableCell>
                          <TableCell
                            className="truncate px-2 text-right font-mono text-xs"
                            style={resizableCellStyle(colWidths.eff)}
                          >
                            {formatOptionalTokensPerSec(r.has_effective_tps, r.effective_tps)}
                          </TableCell>
                          <TableCell
                            className="truncate px-2 text-right font-mono text-xs"
                            style={resizableCellStyle(colWidths.cache)}
                          >
                            {formatCacheHitRate(r.has_cache_hit_rate, r.cache_hit_rate)}
                          </TableCell>
                          <TableCell
                            className="px-1 text-right"
                            style={resizableCellStyle(colWidths.actions)}
                            onClick={(e) => e.stopPropagation()}
                            onPointerDown={(e) => e.stopPropagation()}
                          >
                            <div className="flex items-center justify-end gap-0.5">
                              <HoverTooltip label="Prompt & response">
                                <Button
                                  type="button"
                                  variant="ghost"
                                  size="sm"
                                  className="h-8 w-8 p-0 text-muted-foreground hover:text-foreground"
                                  aria-label="View prompt and response"
                                  onClick={() => setContentRequest(r)}
                                >
                                  <MessageSquareText className="h-4 w-4" />
                                </Button>
                              </HoverTooltip>
                              <HoverTooltip label="Delete request">
                                <Button
                                  type="button"
                                  variant="ghost"
                                  size="sm"
                                  className="h-8 w-8 p-0 text-muted-foreground hover:text-destructive"
                                  disabled={deletingId === r.id}
                                  aria-label="Delete request"
                                  onClick={() => deleteMut.mutate(r.id)}
                                >
                                  <Trash2 className="h-4 w-4" />
                                </Button>
                              </HoverTooltip>
                            </div>
                          </TableCell>
                        </TableRow>
                        {expanded && (
                          <TableRow className="bg-muted/30">
                            <TableCell colSpan={tableColSpan} className="text-xs text-muted-foreground">
                              <span className="font-medium text-foreground">Provider:</span> {r.provider_name || "—"}
                              {" · "}
                              <span className="font-medium text-foreground">Status:</span>{" "}
                              {r.has_status_code ? r.status_code : "—"}
                              {" · "}
                              <span className="font-medium text-foreground">Streaming:</span>{" "}
                              {r.streaming ? "yes" : "no"}
                              {" · "}
                              <span className="font-medium text-foreground">Usage:</span> {r.usage_status || "—"}
                              {" · "}
                              <span className="font-medium text-foreground">Uncached in:</span>{" "}
                              {formatRequestTokenCount(
                                r.has_uncached_input_tokens,
                                r.uncached_input_tokens,
                                r.usage_status,
                              )}
                              {r.error_type ? (
                                <>
                                  {" · "}
                                  <span className="font-medium text-destructive">Error:</span> {r.error_type}
                                </>
                              ) : null}
                            </TableCell>
                          </TableRow>
                        )}
                      </Fragment>
                    );
                  })}
                </TableBody>
              </Table>
            </div>

            <div className="flex flex-wrap items-center justify-between gap-2 text-sm text-muted-foreground">
              <span>
                {formatNumber(listQ.data.total)} total · page {listQ.data.page} of {totalPages}
              </span>
              <div className="flex gap-2">
                <Button variant="outline" size="sm" disabled={page <= 1} onClick={() => setPage((p) => p - 1)}>
                  Previous
                </Button>
                <Button variant="outline" size="sm" disabled={page >= totalPages} onClick={() => setPage((p) => p + 1)}>
                  Next
                </Button>
              </div>
            </div>
          </>
        )}
      </div>

      <RequestContentDialog
        request={contentRequest}
        open={contentRequest !== null}
        onOpenChange={(open) => {
          if (!open) {
            setContentRequest(null);
          }
        }}
      />
    </div>
  );
}
