import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { Link, useParams } from "react-router-dom";
import { toast } from "sonner";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogFooter,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import {
  createAPIKey,
  deleteAPIKey,
  listAPIKeys,
  revokeAPIKey,
  setAPIKeyAllowAll,
} from "@/features/projects/api/keysApi";
import { getProject } from "@/features/projects/api/projectsApi";
import { HoverTooltip } from "@/shared/components/HoverTooltip";
import { EmptyState, ErrorState, LoadingState } from "@/shared/components/PageStates";
import { MaterialIcon } from "@/shared/components/MaterialIcon";
import { PageHeader } from "@/shared/components/PageHeader";
import { PageShell } from "@/shared/components/PageShell";
import { service } from "../../../../wailsjs/go/models";

type APIKeyRow = service.APIKeyDTO;

function isStoredProxyKey(value: string): boolean {
  return value.startsWith("aip_");
}

export function ProjectDetailPage() {
  const { id } = useParams<{ id: string }>();
  const projectID = id ?? "";
  const qc = useQueryClient();
  const projectQ = useQuery({ queryKey: ["project", projectID], queryFn: () => getProject(projectID), enabled: !!projectID });
  const keysQ = useQuery({ queryKey: ["keys", projectID], queryFn: () => listAPIKeys(projectID), enabled: !!projectID });

  const [createOpen, setCreateOpen] = useState(false);
  const [keyName, setKeyName] = useState("");
  const [allowAll, setAllowAll] = useState(true);
  const [secret, setSecret] = useState<string | null>(null);
  const [deleteTarget, setDeleteTarget] = useState<{ id: string; name: string; enabled: boolean } | null>(null);
  const [revealedKeyIds, setRevealedKeyIds] = useState<Set<string>>(() => new Set());

  const createMut = useMutation({
    mutationFn: () => createAPIKey(projectID, keyName.trim(), allowAll),
    onSuccess: (res) => {
      setSecret(res.plaintext);
      setCreateOpen(false);
      setKeyName("");
      qc.invalidateQueries({ queryKey: ["keys", projectID] });
      toast.success("API key created — listed in the table below");
    },
    onError: (e: Error) => toast.error(e.message),
  });

  const revokeMut = useMutation({
    mutationFn: (keyID: string) => revokeAPIKey(keyID),
    onSuccess: () => {
      toast.success("Key revoked");
      qc.invalidateQueries({ queryKey: ["keys", projectID] });
    },
    onError: (e: Error) => toast.error(e.message),
  });

  const allowMut = useMutation({
    mutationFn: ({ keyID, allow }: { keyID: string; allow: boolean }) => setAPIKeyAllowAll(keyID, allow),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["keys", projectID] }),
    onError: (e: Error) => toast.error(e.message),
  });

  const deleteMut = useMutation({
    mutationFn: (keyID: string) => deleteAPIKey(keyID),
    onMutate: async (keyID) => {
      await qc.cancelQueries({ queryKey: ["keys", projectID] });
      const previous = qc.getQueryData<APIKeyRow[]>(["keys", projectID]) ?? [];
      qc.setQueryData<APIKeyRow[]>(
        ["keys", projectID],
        previous.filter((k) => k.id !== keyID),
      );
      setRevealedKeyIds((prev) => {
        const next = new Set(prev);
        next.delete(keyID);
        return next;
      });
      return { previous };
    },
    onSuccess: () => {
      setDeleteTarget(null);
      toast.success("API key deleted");
    },
    onError: (e: Error, _keyID, context) => {
      if (context?.previous) {
        qc.setQueryData(["keys", projectID], context.previous);
      }
      toast.error(e.message);
    },
    onSettled: () => {
      void qc.invalidateQueries({ queryKey: ["keys", projectID] });
    },
  });

  async function copySecret() {
    if (!secret) return;
    await navigator.clipboard.writeText(secret);
    toast.success("Copied to clipboard");
  }

  if (!projectID) return <ErrorState message="Missing project id" />;

  return (
    <PageShell>
      <div className="flex items-center gap-space-sm">
        <Button variant="ghost" size="sm" asChild>
          <Link to="/projects">
            <MaterialIcon name="arrow_back" size={16} />
            Projects
          </Link>
        </Button>
      </div>

      {projectQ.isLoading && <LoadingState />}
      {projectQ.error && <ErrorState message={(projectQ.error as Error).message} onRetry={() => projectQ.refetch()} />}
      {projectQ.data && (
        <PageHeader
          section="Configure"
          accent="API Keys"
          title={projectQ.data.name}
          description={projectQ.data.description || "Project API keys and access control."}
        />
      )}

      <Card>
        <CardHeader className="flex flex-row items-start justify-between">
          <div>
            <CardTitle>API keys</CardTitle>
            <CardDescription>
              Proxy API keys are stored as plaintext in SQLite (`api_key` column). You can copy them anytime below.
            </CardDescription>
          </div>
          <Button onClick={() => setCreateOpen(true)}>Create key</Button>
        </CardHeader>
        <CardContent>
          {keysQ.isLoading && <LoadingState />}
          {keysQ.error && <ErrorState message={(keysQ.error as Error).message} onRetry={() => keysQ.refetch()} />}
          {!keysQ.isLoading && !keysQ.error && keysQ.data?.length === 0 && (
            <EmptyState title="No API keys for this project." />
          )}
          {keysQ.data && keysQ.data.length > 0 && (
            <Table>
              <TableHeader>
                <TableRow>
                  <TableHead>Name</TableHead>
                  <TableHead>API key</TableHead>
                  <TableHead>Allow all models</TableHead>
                  <TableHead>Status</TableHead>
                  <TableHead className="text-right">Actions</TableHead>
                </TableRow>
              </TableHeader>
              <TableBody>
                {keysQ.data.map((k) => (
                  <TableRow key={k.id}>
                    <TableCell className="font-medium">{k.name}</TableCell>
                    <TableCell className="min-w-[280px] max-w-md">
                      {!k.api_key ? (
                        <span className="text-xs text-muted-foreground">—</span>
                      ) : !isStoredProxyKey(k.api_key) ? (
                        <span className="text-xs text-amber-500/90">
                          Legacy hash — run migrations and create a new key
                        </span>
                      ) : (
                        <div className="flex items-start gap-1">
                          <code
                            className={`block flex-1 text-xs text-foreground ${revealedKeyIds.has(k.id) ? "break-all" : "truncate"}`}
                            title={k.api_key}
                          >
                            {revealedKeyIds.has(k.id) ? k.api_key : `${k.api_key.slice(0, 12)}…${k.api_key.slice(-6)}`}
                          </code>
                          <div className="flex shrink-0">
                            <Button
                              variant="ghost"
                              size="icon"
                              className="h-7 w-7"
                              onClick={() =>
                                setRevealedKeyIds((prev) => {
                                  const next = new Set(prev);
                                  if (next.has(k.id)) next.delete(k.id);
                                  else next.add(k.id);
                                  return next;
                                })
                              }
                            >
                              <MaterialIcon
                                name={revealedKeyIds.has(k.id) ? "visibility_off" : "visibility"}
                                size={14}
                              />
                            </Button>
                            <Button
                              variant="ghost"
                              size="icon"
                              className="h-7 w-7"
                              onClick={async () => {
                                await navigator.clipboard.writeText(k.api_key);
                                toast.success("Copied");
                              }}
                            >
                              <MaterialIcon name="content_copy" size={14} />
                            </Button>
                          </div>
                        </div>
                      )}
                    </TableCell>
                    <TableCell>
                      <Switch
                        checked={k.allow_all_models}
                        disabled={!k.enabled}
                        onCheckedChange={(v) => allowMut.mutate({ keyID: k.id, allow: v })}
                      />
                    </TableCell>
                    <TableCell>
                      <Badge variant={k.enabled ? "success" : "destructive"}>{k.enabled ? "Active" : "Revoked"}</Badge>
                    </TableCell>
                    <TableCell className="text-right">
                      <div className="flex flex-wrap justify-end gap-2">
                        {k.enabled && (
                          <Button
                            variant="outline"
                            size="sm"
                            disabled={revokeMut.isPending}
                            onClick={() => revokeMut.mutate(k.id)}
                          >
                            Revoke
                          </Button>
                        )}
                        <Button
                          variant="destructive"
                          size="sm"
                          disabled={deleteMut.isPending}
                          onClick={() => setDeleteTarget({ id: k.id, name: k.name, enabled: k.enabled })}
                        >
                          Delete
                        </Button>
                      </div>
                    </TableCell>
                  </TableRow>
                ))}
              </TableBody>
            </Table>
          )}
        </CardContent>
      </Card>

      <Dialog open={createOpen} onOpenChange={setCreateOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Create API key</DialogTitle>
          </DialogHeader>
          <div className="grid gap-3 py-2">
            <div className="grid gap-2">
              <Label>Name</Label>
              <Input value={keyName} onChange={(e) => setKeyName(e.target.value)} />
            </div>
            <div className="flex items-center gap-2">
              <Switch checked={allowAll} onCheckedChange={setAllowAll} />
              <Label>Allow all models</Label>
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setCreateOpen(false)}>Cancel</Button>
            <Button disabled={!keyName.trim() || createMut.isPending} onClick={() => createMut.mutate()}>Create</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={!!deleteTarget} onOpenChange={(o) => !o && setDeleteTarget(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Delete API key?</DialogTitle>
            <DialogDescription>
              Permanently remove <span className="font-medium text-foreground">{deleteTarget?.name}</span> from this
              project.{" "}
              {deleteTarget?.enabled
                ? "The key will stop working immediately and cannot be recovered."
                : "Past request history will still reference this key id."}
            </DialogDescription>
          </DialogHeader>
          <DialogFooter>
            <Button variant="outline" onClick={() => setDeleteTarget(null)}>Cancel</Button>
            <Button
              variant="destructive"
              disabled={!deleteTarget || deleteMut.isPending}
              onClick={() => deleteTarget && deleteMut.mutate(deleteTarget.id)}
            >
              Delete
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={!!secret} onOpenChange={(o) => !o && setSecret(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>API key created</DialogTitle>
            <DialogDescription>
              Copy it now if you need it elsewhere — it also stays visible in the table below.
            </DialogDescription>
          </DialogHeader>
          <div className="rounded-md border border-border bg-muted/50 p-3 font-mono text-sm break-all">{secret}</div>
          <DialogFooter>
            <Button onClick={copySecret}>
              <HoverTooltip label="Copy to clipboard" side="bottom" className="mr-2 inline-flex">
                <MaterialIcon name="content_copy" size={16} />
              </HoverTooltip>
              Copy
            </Button>
            <Button variant="outline" onClick={() => setSecret(null)}>Done</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </PageShell>
  );
}
