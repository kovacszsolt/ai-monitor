import { Link } from "react-router-dom";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Switch } from "@/components/ui/switch";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { formatUSDPer1M } from "@/features/models/components/modelRouteSort";
import { ModelRouteForm } from "@/features/models/components/ModelRouteForm";
import { useModelRoutesPage } from "@/features/models/hooks/useModelRoutesPage";
import { PageHeader } from "@/shared/components/PageHeader";
import { PageShell } from "@/shared/components/PageShell";
import { EmptyState, ErrorState, LoadingState } from "@/shared/components/PageStates";
import { SortableTableHead } from "@/shared/table/SortableTableHead";

const selectClassName =
  "h-9 min-w-[160px] rounded border border-outline-variant bg-surface-container px-space-md text-body-sm text-on-surface";

export function ModelsPage() {
  const {
    data,
    isLoading,
    error,
    refetch,
    providers,
    createOpen,
    setCreateOpen,
    edit,
    setEdit,
    publicName,
    setPublicName,
    providerID,
    setProviderID,
    upstream,
    setUpstream,
    protocol,
    setProtocol,
    modality,
    setModality,
    inputUSD,
    setInputUSD,
    cachedUSD,
    setCachedUSD,
    outputUSD,
    setOutputUSD,
    nameQuery,
    setNameQuery,
    providerFilter,
    setProviderFilter,
    displayRoutes,
    sortDirection,
    toggleSort,
    createMut,
    updateMut,
    toggleMut,
    openEdit,
    openCreateDialog,
  } = useModelRoutesPage();

  const importHref = providerFilter
    ? `/models/import?provider=${encodeURIComponent(providerFilter)}`
    : "/models/import";

  return (
    <PageShell>
      <PageHeader
        section="Configure"
        accent="Model routes"
        title="Models"
        description="Public model aliases routed to upstream provider models."
        actions={
          <div className="flex flex-wrap gap-2">
            <Button variant="outline" asChild>
              <Link to={importHref}>Import from provider</Link>
            </Button>
            <Button onClick={openCreateDialog}>Add model</Button>
          </div>
        }
      />
    <Card>
      <CardContent className="pt-space-lg">
        {isLoading && <LoadingState />}
        {error && <ErrorState message={(error as Error).message} onRetry={() => refetch()} />}
        {!isLoading && !error && data?.length === 0 && <EmptyState title="No model routes." />}
        {!isLoading && !error && data && data.length > 0 && (
          <div className="mb-4 flex flex-wrap items-end gap-3">
            <label className="flex min-w-[200px] flex-1 flex-col gap-1 text-xs text-muted-foreground sm:max-w-xs">
              Search name
              <Input
                value={nameQuery}
                onChange={(e) => setNameQuery(e.target.value)}
                placeholder="Public or upstream model"
                className="h-9"
              />
            </label>
            <label className="flex flex-col gap-1 text-xs text-muted-foreground">
              Provider
              <select
                className={`${selectClassName} min-w-[180px] max-w-[280px]`}
                value={providerFilter}
                onChange={(e) => setProviderFilter(e.target.value)}
              >
                <option value="">All providers</option>
                {(providers ?? []).map((p) => (
                  <option key={p.id} value={p.id}>{p.name}</option>
                ))}
              </select>
            </label>
          </div>
        )}
        {!isLoading && !error && data && data.length > 0 && displayRoutes.length === 0 && (
          <EmptyState title="No model routes match your filters." />
        )}
        {displayRoutes.length > 0 && (
          <Table>
            <TableHeader>
              <TableRow>
                <SortableTableHead
                  sortDirection={sortDirection("publicName")}
                  onSortClick={() => toggleSort("publicName")}
                >
                  Public name
                </SortableTableHead>
                <SortableTableHead
                  sortDirection={sortDirection("provider")}
                  onSortClick={() => toggleSort("provider")}
                >
                  Provider
                </SortableTableHead>
                <SortableTableHead
                  sortDirection={sortDirection("upstream")}
                  onSortClick={() => toggleSort("upstream")}
                >
                  Upstream
                </SortableTableHead>
                <SortableTableHead
                  sortDirection={sortDirection("protocol")}
                  onSortClick={() => toggleSort("protocol")}
                >
                  Protocol
                </SortableTableHead>
                <SortableTableHead
                  sortDirection={sortDirection("modality")}
                  onSortClick={() => toggleSort("modality")}
                >
                  Modality
                </SortableTableHead>
                <SortableTableHead
                  className="text-right"
                  sortDirection={sortDirection("inputUSD")}
                  onSortClick={() => toggleSort("inputUSD")}
                >
                  Input $/1M
                </SortableTableHead>
                <SortableTableHead
                  className="text-right"
                  sortDirection={sortDirection("cachedUSD")}
                  onSortClick={() => toggleSort("cachedUSD")}
                >
                  Cached $/1M
                </SortableTableHead>
                <SortableTableHead
                  className="text-right"
                  sortDirection={sortDirection("outputUSD")}
                  onSortClick={() => toggleSort("outputUSD")}
                >
                  Output $/1M
                </SortableTableHead>
                <SortableTableHead
                  sortDirection={sortDirection("enabled")}
                  onSortClick={() => toggleSort("enabled")}
                >
                  Enabled
                </SortableTableHead>
                <TableHead className="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {displayRoutes.map((m) => (
                <TableRow key={m.id}>
                  <TableCell className="font-medium">{m.public_model_name}</TableCell>
                  <TableCell>{m.provider_name || m.provider_id}</TableCell>
                  <TableCell className="font-mono text-xs">{m.upstream_model_name}</TableCell>
                  <TableCell><Badge variant="secondary">{m.protocol}</Badge></TableCell>
                  <TableCell><Badge variant="secondary">{m.modality || "text-text"}</Badge></TableCell>
                  <TableCell className="text-right font-mono text-xs">{formatUSDPer1M(m.input_token_usd_per_1m)}</TableCell>
                  <TableCell className="text-right font-mono text-xs">{formatUSDPer1M(m.cached_input_token_usd_per_1m)}</TableCell>
                  <TableCell className="text-right font-mono text-xs">{formatUSDPer1M(m.output_token_usd_per_1m)}</TableCell>
                  <TableCell>
                    <Switch checked={m.enabled} onCheckedChange={(v) => toggleMut.mutate({ id: m.id, enabled: v })} />
                  </TableCell>
                  <TableCell className="text-right">
                    <Button variant="ghost" size="sm" onClick={() => openEdit(m)}>Edit</Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        )}
      </CardContent>

      <Dialog open={createOpen} onOpenChange={setCreateOpen}>
        <DialogContent>
          <DialogHeader><DialogTitle>Add model route</DialogTitle></DialogHeader>
          <ModelRouteForm
            publicName={publicName} setPublicName={setPublicName} showPublic
            providerID={providerID} setProviderID={setProviderID} providers={providers ?? []}
            upstream={upstream} setUpstream={setUpstream}
            protocol={protocol} setProtocol={setProtocol}
            modality={modality} setModality={setModality}
            inputUSD={inputUSD} setInputUSD={setInputUSD}
            cachedUSD={cachedUSD} setCachedUSD={setCachedUSD}
            outputUSD={outputUSD} setOutputUSD={setOutputUSD}
          />
          <DialogFooter>
            <Button variant="outline" onClick={() => setCreateOpen(false)}>Cancel</Button>
            <Button
              disabled={!publicName.trim() || !providerID || !upstream.trim() || createMut.isPending}
              onClick={() => createMut.mutate()}
            >
              Create
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={!!edit} onOpenChange={(o) => !o && setEdit(null)}>
        <DialogContent>
          <DialogHeader><DialogTitle>Edit {edit?.public_model_name}</DialogTitle></DialogHeader>
          <ModelRouteForm
            upstream={upstream} setUpstream={setUpstream}
            protocol={protocol} setProtocol={setProtocol}
            modality={modality} setModality={setModality}
            inputUSD={inputUSD} setInputUSD={setInputUSD}
            cachedUSD={cachedUSD} setCachedUSD={setCachedUSD}
            outputUSD={outputUSD} setOutputUSD={setOutputUSD}
          />
          <DialogFooter>
            <Button variant="outline" onClick={() => setEdit(null)}>Cancel</Button>
            <Button disabled={!upstream.trim() || updateMut.isPending} onClick={() => updateMut.mutate()}>Save</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </Card>
    </PageShell>
  );
}
