import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { MaterialIcon } from "@/shared/components/MaterialIcon";
import { PageHeader } from "@/shared/components/PageHeader";
import { PageShell } from "@/shared/components/PageShell";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { DatabaseMaintenanceCard } from "@/features/settings/components/DatabaseMaintenanceCard";
import { getSettings, pickDatabaseFile, updateDatabasePath } from "@/features/settings/api/settingsApi";
import { ErrorState, LoadingState } from "@/shared/components/PageStates";

export function SettingsPage() {
  const qc = useQueryClient();
  const settingsQ = useQuery({ queryKey: ["settings"], queryFn: getSettings });

  const [autoResolve, setAutoResolve] = useState(true);
  const [databasePath, setDatabasePath] = useState("");

  useEffect(() => {
    if (!settingsQ.data) return;
    setAutoResolve(settingsQ.data.auto_resolve);
    setDatabasePath(settingsQ.data.database_path ?? "");
  }, [settingsQ.data]);

  const saveDbMut = useMutation({
    mutationFn: () => updateDatabasePath(autoResolve ? "" : databasePath.trim()),
    onSuccess: () => {
      toast.success("Database path saved");
      qc.invalidateQueries({ queryKey: ["settings", "connection", "proxy-config"] });
    },
    onError: (e: Error) => toast.error(e.message),
  });

  const browseMut = useMutation({
    mutationFn: pickDatabaseFile,
    onSuccess: (path) => {
      if (!path) return;
      setAutoResolve(false);
      setDatabasePath(path);
    },
    onError: (e: Error) => toast.error(e.message),
  });

  if (settingsQ.isLoading) return <LoadingState label="Loading settings…" />;
  if (settingsQ.error) {
    return <ErrorState message={(settingsQ.error as Error).message} onRetry={() => settingsQ.refetch()} />;
  }

  const s = settingsQ.data!;

  return (
    <PageShell>
      <PageHeader
        section="System"
        title="Settings"
        description={`Desktop configuration stored in ${s.config_file_path}`}
      />

      <Card>
        <CardHeader>
          <CardTitle className="text-base">Database</CardTitle>
          <CardDescription>
            Path to <code className="text-xs">proxy.db</code>. Saved to desktop config; proxy serve uses this path when
            started from the app.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="flex items-center gap-2">
            <Switch id="auto-db" checked={autoResolve} onCheckedChange={setAutoResolve} />
            <Label htmlFor="auto-db">Auto-detect monorepo data/proxy.db</Label>
          </div>

          {!autoResolve && (
            <div className="space-y-2">
              <Label htmlFor="db-path">Database path</Label>
              <div className="flex gap-2">
                <Input
                  id="db-path"
                  value={databasePath}
                  onChange={(e) => setDatabasePath(e.target.value)}
                  placeholder="/path/to/proxy.db"
                  className="font-mono text-xs"
                />
                <Button
                  type="button"
                  variant="outline"
                  size="icon"
                  aria-label="Browse"
                  disabled={browseMut.isPending}
                  onClick={() => browseMut.mutate()}
                >
                  <MaterialIcon name="folder_open" size={16} />
                </Button>
              </div>
            </div>
          )}

          <p className="text-xs text-muted-foreground">
            Active connection: <span className="font-mono">{s.effective_database_path}</span>
          </p>

          {!s.config_bridge_ok && s.bridge_message && (
            <p className="text-sm text-amber-500">{s.bridge_message}</p>
          )}

          <Button type="button" disabled={saveDbMut.isPending} onClick={() => saveDbMut.mutate()}>
            Save &amp; reconnect
          </Button>
        </CardContent>
      </Card>

      <DatabaseMaintenanceCard />

      <Card>
        <CardHeader>
          <CardTitle className="text-base">Proxy bootstrap config</CardTitle>
          <CardDescription>Optional YAML for first-time seed; runtime server/logging/capture live in the database.</CardDescription>
        </CardHeader>
        <CardContent className="space-y-2 text-sm">
          <p className="font-mono text-xs text-muted-foreground">{s.proxy_config_path || "(not found)"}</p>
          <div className="flex items-center gap-2">
            <span className="text-muted-foreground">Bridge</span>
            <Badge variant={s.config_bridge_ok ? "success" : "secondary"}>
              {s.config_bridge_ok ? "Aligned" : "Check paths"}
            </Badge>
          </div>
        </CardContent>
      </Card>
    </PageShell>
  );
}
