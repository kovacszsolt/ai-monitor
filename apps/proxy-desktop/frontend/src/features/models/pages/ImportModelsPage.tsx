import { Link } from "react-router-dom";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { useImportModelsPage } from "@/features/models/hooks/useImportModelsPage";
import { PageHeader } from "@/shared/components/PageHeader";
import { PageShell } from "@/shared/components/PageShell";
import { EmptyState, ErrorState, LoadingState } from "@/shared/components/PageStates";

const selectClassName =
  "h-9 min-w-[200px] rounded border border-outline-variant bg-surface-container px-space-md text-body-sm text-on-surface";

export function ImportModelsPage() {
  const {
    providersQ,
    providerId,
    onProviderChange,
    nameQuery,
    setNameQuery,
    publicPrefix,
    setPublicPrefix,
    models,
    selected,
    fetchError,
    fetchMut,
    importMut,
    filtered,
    importableFiltered,
    firstSelectedPreview,
    toggleId,
    selectAllImportable,
    clearSelection,
    canFetch,
    canImport,
  } = useImportModelsPage();

  const providers = providersQ.data ?? [];

  return (
    <PageShell>
      <PageHeader
        section="Configure"
        accent="Import"
        title="Import models from provider"
        description="Load upstream model ids, choose routes to add, and set a public name prefix."
        actions={
          <Button variant="outline" asChild>
            <Link to="/models">Back to models</Link>
          </Button>
        }
      />

      <Card>
        <CardHeader className="pb-3">
          <CardTitle className="text-base">Provider &amp; options</CardTitle>
        </CardHeader>
        <CardContent className="grid gap-4 lg:grid-cols-2">
          <div className="grid gap-2">
            <Label htmlFor="import-provider">Provider</Label>
            <select
              id="import-provider"
              className={selectClassName}
              value={providerId}
              disabled={providersQ.isLoading}
              onChange={(e) => onProviderChange(e.target.value)}
            >
              {providers.length === 0 && <option value="">No providers</option>}
              {providers.map((p) => (
                <option key={p.id} value={p.id}>
                  {p.name}
                </option>
              ))}
            </select>
          </div>
          <div className="grid gap-2">
            <Label htmlFor="import-prefix">Public name prefix</Label>
            <Input
              id="import-prefix"
              value={publicPrefix}
              onChange={(e) => setPublicPrefix(e.target.value)}
              placeholder="e.g. local/"
            />
            <p className="text-xs text-muted-foreground">
              Public name = prefix + upstream id. Defaults: openai_chat (anthropic →
              anthropic_messages), modality text-text.
            </p>
            {firstSelectedPreview && (
              <p className="text-xs text-muted-foreground">
                Example: <code className="text-code-inline">{firstSelectedPreview}</code>
              </p>
            )}
          </div>
          <div className="flex flex-wrap items-end gap-2 lg:col-span-2">
            <Button
              type="button"
              variant="secondary"
              disabled={!canFetch || fetchMut.isPending}
              onClick={() => fetchMut.mutate()}
            >
              {fetchMut.isPending ? "Loading…" : "Load models from provider"}
            </Button>
            <Button
              type="button"
              disabled={!canImport}
              onClick={() => importMut.mutate()}
            >
              {importMut.isPending ? "Importing…" : `Import selected (${selected.size})`}
            </Button>
          </div>
        </CardContent>
      </Card>

      <Card className="min-h-0 flex-1">
        <CardHeader className="flex flex-row flex-wrap items-center justify-between gap-3 space-y-0 pb-3">
          <CardTitle className="text-base">Upstream models</CardTitle>
          {models && importableFiltered.length > 0 && (
            <div className="flex flex-wrap gap-2">
              <Button type="button" variant="ghost" size="sm" onClick={selectAllImportable}>
                Select all visible
              </Button>
              <Button type="button" variant="ghost" size="sm" onClick={clearSelection}>
                Clear selection
              </Button>
            </div>
          )}
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="max-w-md">
            <Label htmlFor="import-search">Search by name</Label>
            <Input
              id="import-search"
              className="mt-1.5"
              value={nameQuery}
              onChange={(e) => setNameQuery(e.target.value)}
              placeholder="Filter upstream model id"
              disabled={!models}
            />
          </div>

          {providersQ.isLoading && <LoadingState />}
          {providersQ.error && (
            <ErrorState message={(providersQ.error as Error).message} onRetry={() => providersQ.refetch()} />
          )}

          {fetchMut.isPending && <LoadingState label="Fetching model list…" />}
          {fetchError && (
            <ErrorState message={fetchError} onRetry={() => fetchMut.mutate()} />
          )}

          {!models && !fetchMut.isPending && !fetchError && (
            <EmptyState title="Load models to see the upstream list." />
          )}

          {models && !fetchMut.isPending && !fetchError && filtered.length === 0 && (
            <EmptyState title="No models match your search." />
          )}

          {models && !fetchMut.isPending && !fetchError && filtered.length > 0 && (
            <div className="max-h-[min(60vh,640px)] overflow-auto rounded-md border border-border">
              <Table>
                <TableHeader>
                  <TableRow>
                    <TableHead className="w-10" />
                    <TableHead>Upstream model id</TableHead>
                    <TableHead className="w-[140px]">Status</TableHead>
                  </TableRow>
                </TableHeader>
                <TableBody>
                  {filtered.map((m) => (
                    <TableRow key={m.id} className={m.alreadyImported ? "opacity-70" : undefined}>
                      <TableCell>
                        <input
                          type="checkbox"
                          className="h-4 w-4"
                          disabled={m.alreadyImported}
                          checked={selected.has(m.id)}
                          onChange={() => toggleId(m.id, m.alreadyImported)}
                          aria-label={`Select ${m.id}`}
                        />
                      </TableCell>
                      <TableCell className="font-mono text-xs">{m.id}</TableCell>
                      <TableCell>
                        {m.alreadyImported ? (
                          <Badge variant="secondary">Already imported</Badge>
                        ) : (
                          <span className="text-xs text-muted-foreground">Available</span>
                        )}
                      </TableCell>
                    </TableRow>
                  ))}
                </TableBody>
              </Table>
            </div>
          )}
        </CardContent>
      </Card>
    </PageShell>
  );
}
