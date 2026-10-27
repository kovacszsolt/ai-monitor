import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
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
  createProvider,
  listProviders,
  setProviderEnabled,
  updateProvider,
} from "@/features/providers/api/providersApi";
import { PageHeader } from "@/shared/components/PageHeader";
import { PageShell } from "@/shared/components/PageShell";
import { EmptyState, ErrorState, LoadingState } from "@/shared/components/PageStates";
import { service } from "../../../../wailsjs/go/models";

const TYPES = ["openai", "anthropic", "openai-compatible", "bifrost"] as const;

export function ProvidersPage() {
  const qc = useQueryClient();
  const { data, isLoading, error, refetch } = useQuery({ queryKey: ["providers"], queryFn: listProviders });
  const [createOpen, setCreateOpen] = useState(false);
  const [edit, setEdit] = useState<service.ProviderDTO | null>(null);
  const [name, setName] = useState("");
  const [providerType, setProviderType] = useState<string>(TYPES[0]);
  const [baseURL, setBaseURL] = useState("");
  const [credential, setCredential] = useState("");
  const [connMS, setConnMS] = useState(10000);
  const [reqMS, setReqMS] = useState(300000);

  const createMut = useMutation({
    mutationFn: () => createProvider(name.trim(), providerType, baseURL.trim(), credential.trim(), connMS, reqMS),
    onSuccess: () => {
      toast.success("Provider created");
      setCreateOpen(false);
      qc.invalidateQueries({ queryKey: ["providers"] });
    },
    onError: (e: Error) => toast.error(e.message),
  });

  const updateMut = useMutation({
    mutationFn: () => updateProvider(edit!.id, baseURL.trim(), credential.trim(), connMS, reqMS),
    onSuccess: () => {
      toast.success("Provider updated");
      setEdit(null);
      qc.invalidateQueries({ queryKey: ["providers"] });
    },
    onError: (e: Error) => toast.error(e.message),
  });

  const toggleMut = useMutation({
    mutationFn: ({ id, enabled }: { id: string; enabled: boolean }) => setProviderEnabled(id, enabled),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["providers"] }),
    onError: (e: Error) => toast.error(e.message),
  });

  function openEdit(p: service.ProviderDTO) {
    setEdit(p);
    setBaseURL(p.base_url);
    setCredential("");
    setConnMS(p.connection_timeout_ms);
    setReqMS(p.request_timeout_ms);
  }

  return (
    <PageShell>
      <PageHeader
        section="Configure"
        accent="Upstream"
        title="Providers"
        description="Upstream API endpoints. Credentials are stored in plain text in the database."
        actions={
          <Button onClick={() => { setCreateOpen(true); setName(""); setBaseURL(""); setCredential(""); }}>
            Add provider
          </Button>
        }
      />
    <Card>
      <CardContent className="pt-space-lg">
        {isLoading && <LoadingState />}
        {error && <ErrorState message={(error as Error).message} onRetry={() => refetch()} />}
        {!isLoading && !error && data?.length === 0 && <EmptyState title="No providers configured." />}
        {data && data.length > 0 && (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Name</TableHead>
                <TableHead>Type</TableHead>
                <TableHead>Base URL</TableHead>
                <TableHead>API key</TableHead>
                <TableHead>Status</TableHead>
                <TableHead className="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {data.map((p) => (
                <TableRow key={p.id}>
                  <TableCell className="font-medium">{p.name}</TableCell>
                  <TableCell><Badge variant="secondary">{p.provider_type}</Badge></TableCell>
                  <TableCell className="max-w-[200px] truncate font-mono text-xs">{p.base_url}</TableCell>
                  <TableCell className="text-sm text-muted-foreground">
                    {p.has_credential ? "Configured" : "—"}
                  </TableCell>
                  <TableCell>
                    <Switch checked={p.enabled} onCheckedChange={(v) => toggleMut.mutate({ id: p.id, enabled: v })} />
                  </TableCell>
                  <TableCell className="text-right">
                    <Button variant="ghost" size="sm" onClick={() => openEdit(p)}>Edit</Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        )}
      </CardContent>

      <Dialog open={createOpen} onOpenChange={setCreateOpen}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle>Add provider</DialogTitle>
          </DialogHeader>
          <FormFields
            name={name} setName={setName} showName
            providerType={providerType} setProviderType={setProviderType} showType
            baseURL={baseURL} setBaseURL={setBaseURL}
            credential={credential} setCredential={setCredential}
            connMS={connMS} setConnMS={setConnMS} reqMS={reqMS} setReqMS={setReqMS}
          />
          <DialogFooter>
            <Button variant="outline" onClick={() => setCreateOpen(false)}>Cancel</Button>
            <Button
              disabled={!name.trim() || !baseURL.trim() || !credential.trim() || createMut.isPending}
              onClick={() => createMut.mutate()}
            >
              Create
            </Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={!!edit} onOpenChange={(o) => !o && setEdit(null)}>
        <DialogContent className="max-w-md">
          <DialogHeader>
            <DialogTitle>Edit {edit?.name}</DialogTitle>
            <DialogDescription>Name and type cannot be changed here.</DialogDescription>
          </DialogHeader>
          <FormFields
            baseURL={baseURL} setBaseURL={setBaseURL}
            credential={credential} setCredential={setCredential}
            credentialOptional
            connMS={connMS} setConnMS={setConnMS} reqMS={reqMS} setReqMS={setReqMS}
          />
          <DialogFooter>
            <Button variant="outline" onClick={() => setEdit(null)}>Cancel</Button>
            <Button disabled={!baseURL.trim() || updateMut.isPending} onClick={() => updateMut.mutate()}>Save</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </Card>
    </PageShell>
  );
}

function FormFields(props: {
  name?: string; setName?: (v: string) => void; showName?: boolean;
  providerType?: string; setProviderType?: (v: string) => void; showType?: boolean;
  baseURL: string; setBaseURL: (v: string) => void;
  credential: string; setCredential: (v: string) => void;
  /** Edit: leave blank to keep the stored upstream API key. */
  credentialOptional?: boolean;
  connMS: number; setConnMS: (v: number) => void;
  reqMS: number; setReqMS: (v: number) => void;
}) {
  return (
    <div className="grid gap-3 py-2">
      {props.showName && (
        <div className="grid gap-2">
          <Label>Name</Label>
          <Input value={props.name ?? ""} onChange={(e) => props.setName?.(e.target.value)} />
        </div>
      )}
      {props.showType && (
        <div className="grid gap-2">
          <Label>Type</Label>
          <select
            className="h-9 rounded-md border border-border bg-card px-3 text-sm"
            value={props.providerType}
            onChange={(e) => props.setProviderType?.(e.target.value)}
          >
            {TYPES.map((t) => <option key={t} value={t}>{t}</option>)}
          </select>
        </div>
      )}
      <div className="grid gap-2">
        <Label>Base URL</Label>
        <Input value={props.baseURL} onChange={(e) => props.setBaseURL(e.target.value)} />
      </div>
      <div className="grid gap-2">
        <Label>API key</Label>
        <Input
          type="password"
          value={props.credential}
          onChange={(e) => props.setCredential(e.target.value)}
          placeholder={
            props.credentialOptional
              ? "Leave blank to keep the current key"
              : props.providerType === "bifrost"
                ? "Bifrost virtual key (x-bf-vk)"
                : "sk-… or upstream API key"
          }
          className="font-mono text-xs"
          autoComplete="new-password"
        />
        {props.credentialOptional ? (
          <p className="text-xs text-muted-foreground">
            The stored key is never shown. Enter a new value only to replace it.
          </p>
        ) : null}
      </div>
      <div className="grid grid-cols-2 gap-3">
        <div className="grid gap-2">
          <Label>Conn timeout ms</Label>
          <Input type="number" value={props.connMS} onChange={(e) => props.setConnMS(Number(e.target.value))} />
        </div>
        <div className="grid gap-2">
          <Label>Request timeout ms</Label>
          <Input type="number" value={props.reqMS} onChange={(e) => props.setReqMS(Number(e.target.value))} />
        </div>
      </div>
    </div>
  );
}
