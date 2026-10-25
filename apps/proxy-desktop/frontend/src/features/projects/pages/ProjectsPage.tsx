import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { Link } from "react-router-dom";
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
import { createProject, listProjects, setProjectEnabled, updateProject } from "@/features/projects/api/projectsApi";
import { formatDateTime, formatNumber } from "@/lib/format";
import { PageHeader } from "@/shared/components/PageHeader";
import { PageShell } from "@/shared/components/PageShell";
import { EmptyState, ErrorState, LoadingState } from "@/shared/components/PageStates";
import { service } from "../../../../wailsjs/go/models";

export function ProjectsPage() {
  const qc = useQueryClient();
  const { data, isLoading, error, refetch } = useQuery({ queryKey: ["projects"], queryFn: listProjects });
  const [createOpen, setCreateOpen] = useState(false);
  const [edit, setEdit] = useState<service.ProjectDTO | null>(null);
  const [name, setName] = useState("");
  const [description, setDescription] = useState("");

  const createMut = useMutation({
    mutationFn: () => createProject(name.trim(), description),
    onSuccess: () => {
      toast.success("Project created");
      setCreateOpen(false);
      setName("");
      setDescription("");
      qc.invalidateQueries({ queryKey: ["projects"] });
    },
    onError: (e: Error) => toast.error(e.message),
  });

  const updateMut = useMutation({
    mutationFn: () => updateProject(edit!.id, name.trim(), description),
    onSuccess: () => {
      toast.success("Project updated");
      setEdit(null);
      qc.invalidateQueries({ queryKey: ["projects"] });
    },
    onError: (e: Error) => toast.error(e.message),
  });

  const toggleMut = useMutation({
    mutationFn: ({ id, enabled }: { id: string; enabled: boolean }) => setProjectEnabled(id, enabled),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["projects"] }),
    onError: (e: Error) => toast.error(e.message),
  });

  function openEdit(p: service.ProjectDTO) {
    setEdit(p);
    setName(p.name);
    setDescription(p.description);
  }

  return (
    <PageShell>
      <PageHeader
        section="Configure"
        accent="Projects"
        title="Projects"
        description="Tenant boundaries for proxy API keys and usage."
        actions={
          <Button onClick={() => { setCreateOpen(true); setName(""); setDescription(""); }}>New project</Button>
        }
      />
    <Card>
      <CardContent className="pt-space-lg">
        {isLoading && <LoadingState />}
        {error && <ErrorState message={(error as Error).message} onRetry={() => refetch()} />}
        {!isLoading && !error && data?.length === 0 && (
          <EmptyState title="No projects yet." action={<Button onClick={() => setCreateOpen(true)}>Create one</Button>} />
        )}
        {data && data.length > 0 && (
          <Table>
            <TableHeader>
              <TableRow>
                <TableHead>Name</TableHead>
                <TableHead>Description</TableHead>
                <TableHead className="text-right">Requests</TableHead>
                <TableHead>Last request</TableHead>
                <TableHead>Status</TableHead>
                <TableHead className="text-right">Actions</TableHead>
              </TableRow>
            </TableHeader>
            <TableBody>
              {data.map((p) => (
                <TableRow key={p.id}>
                  <TableCell className="font-medium">{p.name}</TableCell>
                  <TableCell className="max-w-xs truncate text-muted-foreground">{p.description || "—"}</TableCell>
                  <TableCell className="text-right font-mono text-muted-foreground">
                    {formatNumber(p.request_count ?? 0)}
                  </TableCell>
                  <TableCell className="whitespace-nowrap text-muted-foreground">
                    {p.last_request_at ? formatDateTime(p.last_request_at) : "—"}
                  </TableCell>
                  <TableCell>
                    <div className="flex items-center gap-2">
                      <Switch
                        checked={p.enabled}
                        onCheckedChange={(v) => toggleMut.mutate({ id: p.id, enabled: v })}
                      />
                      <Badge variant={p.enabled ? "success" : "secondary"}>{p.enabled ? "Enabled" : "Disabled"}</Badge>
                    </div>
                  </TableCell>
                  <TableCell className="text-right">
                    <Button variant="ghost" size="sm" onClick={() => openEdit(p)}>Edit</Button>
                    <Button variant="ghost" size="sm" asChild>
                      <Link to={`/projects/${p.id}`}>API keys</Link>
                    </Button>
                  </TableCell>
                </TableRow>
              ))}
            </TableBody>
          </Table>
        )}
      </CardContent>

      <Dialog open={createOpen} onOpenChange={setCreateOpen}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>New project</DialogTitle>
            <DialogDescription>Unique name required.</DialogDescription>
          </DialogHeader>
          <div className="grid gap-4 py-2">
            <div className="grid gap-2">
              <Label htmlFor="p-name">Name</Label>
              <Input id="p-name" value={name} onChange={(e) => setName(e.target.value)} />
            </div>
            <div className="grid gap-2">
              <Label htmlFor="p-desc">Description</Label>
              <Input id="p-desc" value={description} onChange={(e) => setDescription(e.target.value)} />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setCreateOpen(false)}>Cancel</Button>
            <Button disabled={!name.trim() || createMut.isPending} onClick={() => createMut.mutate()}>Create</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>

      <Dialog open={!!edit} onOpenChange={(o) => !o && setEdit(null)}>
        <DialogContent>
          <DialogHeader>
            <DialogTitle>Edit project</DialogTitle>
          </DialogHeader>
          <div className="grid gap-4 py-2">
            <div className="grid gap-2">
              <Label>Name</Label>
              <Input value={name} onChange={(e) => setName(e.target.value)} />
            </div>
            <div className="grid gap-2">
              <Label>Description</Label>
              <Input value={description} onChange={(e) => setDescription(e.target.value)} />
            </div>
          </div>
          <DialogFooter>
            <Button variant="outline" onClick={() => setEdit(null)}>Cancel</Button>
            <Button disabled={!name.trim() || updateMut.isPending} onClick={() => updateMut.mutate()}>Save</Button>
          </DialogFooter>
        </DialogContent>
      </Dialog>
    </Card>
    </PageShell>
  );
}
