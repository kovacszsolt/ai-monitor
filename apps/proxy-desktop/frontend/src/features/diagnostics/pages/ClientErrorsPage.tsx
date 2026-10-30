import { useState } from "react";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { useClientErrorLog } from "@/features/diagnostics/hooks/useClientErrorLog";
import { formatDateTime } from "@/lib/format";
import type { ClientErrorLogEntry } from "@/lib/clientErrorCapture";
import { PageHeader } from "@/shared/components/PageHeader";
import { PageShell } from "@/shared/components/PageShell";
import { EmptyState, ErrorState, LoadingState } from "@/shared/components/PageStates";
import { toast } from "sonner";

const PAGE_SIZE = 50;

function levelBadge(level: string) {
  const variant =
    level === "error" ? "destructive" : level === "warn" ? "secondary" : "default";
  return <Badge variant={variant}>{level}</Badge>;
}

export function ClientErrorsPage() {
  const [page, setPage] = useState(1);
  const [selected, setSelected] = useState<ClientErrorLogEntry | null>(null);
  const { listQ, clearMut, totalPages } = useClientErrorLog(page, PAGE_SIZE);

  if (listQ.isLoading) {
    return <LoadingState label="Loading error log…" />;
  }
  if (listQ.error) {
    return (
      <ErrorState message={(listQ.error as Error).message} onRetry={() => listQ.refetch()} />
    );
  }

  const data = listQ.data!;
  const items = data.items;

  return (
    <PageShell>
      <PageHeader
        section="System"
        title="Error log"
        description="Client-side errors and console output saved for desktop debugging (newest first)."
        actions={
          <div className="flex flex-wrap gap-2">
            <Button variant="outline" size="sm" onClick={() => listQ.refetch()} disabled={listQ.isFetching}>
              Refresh
            </Button>
            <Button
              variant="destructive"
              size="sm"
              disabled={clearMut.isPending || data.total === 0}
              onClick={() => {
                clearMut.mutate(undefined, {
                  onSuccess: () => {
                    setPage(1);
                    toast.success("Error log cleared");
                  },
                  onError: (err) => toast.error((err as Error).message),
                });
              }}
            >
              Clear all
            </Button>
          </div>
        }
      />

      {items.length === 0 ? (
        <EmptyState title="No errors logged yet — console errors, API failures, and uncaught exceptions appear here." />
      ) : (
        <>
          <div className="rounded-lg border border-outline-variant/40 bg-surface-container-low">
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead className="w-[168px]">Time</TableHead>
                  <TableHead className="w-[72px]">Level</TableHead>
                  <TableHead className="w-[140px]">Source</TableHead>
                  <TableHead>Message</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {items.map((row) => (
                  <TableRow
                    key={row.id}
                    className="cursor-pointer"
                    onClick={() => setSelected(row)}
                  >
                    <TableCell className="whitespace-nowrap text-body-sm text-on-surface-variant">
                      {formatDateTime(row.created_at)}
                    </TableCell>
                    <TableCell>{levelBadge(row.level)}</TableCell>
                    <TableCell className="max-w-[140px] truncate font-mono text-body-sm">
                      {row.source || "—"}
                    </TableCell>
                    <TableCell className="max-w-md truncate text-body-sm">{row.message}</TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          </div>

          <div className="mt-4 flex items-center justify-between gap-2 text-body-sm text-on-surface-variant">
            <span>
              {data.total} entr{data.total === 1 ? "y" : "ies"} (keeping latest 500)
            </span>
            <div className="flex items-center gap-2">
              <Button
                variant="outline"
                size="sm"
                disabled={page <= 1}
                onClick={() => setPage((p) => Math.max(1, p - 1))}
              >
                Previous
              </Button>
              <span>
                Page {page} / {totalPages}
              </span>
              <Button
                variant="outline"
                size="sm"
                disabled={page >= totalPages}
                onClick={() => setPage((p) => p + 1)}
              >
                Next
              </Button>
            </div>
          </div>
        </>
      )}

      <Dialog open={selected != null} onOpenChange={(open) => !open && setSelected(null)}>
        <DialogContent className="max-h-[85vh] max-w-2xl overflow-y-auto">
          {selected ? (
            <>
              <DialogHeader>
                <DialogTitle className="flex flex-wrap items-center gap-2 text-base">
                  {levelBadge(selected.level)}
                  <span className="font-mono text-body-sm text-on-surface-variant">{selected.source}</span>
                </DialogTitle>
                <DialogDescription>{formatDateTime(selected.created_at)}</DialogDescription>
              </DialogHeader>
              <div className="space-y-4">
                <div>
                  <p className="font-label-sm text-label-sm text-on-surface-variant">Message</p>
                  <pre className="mt-1 whitespace-pre-wrap break-words rounded-md bg-surface-container-high p-3 text-body-sm">
                    {selected.message}
                  </pre>
                </div>
                {selected.detail ? (
                  <div>
                    <p className="font-label-sm text-label-sm text-on-surface-variant">Detail</p>
                    <pre className="mt-1 max-h-96 overflow-auto whitespace-pre-wrap break-words rounded-md bg-surface-container-high p-3 font-mono text-body-sm">
                      {selected.detail}
                    </pre>
                  </div>
                ) : null}
              </div>
            </>
          ) : null}
        </DialogContent>
      </Dialog>
    </PageShell>
  );
}
