import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Download, RefreshCw } from "lucide-react";
import { useState } from "react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import {
  exportMetadata,
  getStorageStats,
  runDatabaseMigrations,
  saveMetadataExport,
} from "@/features/proxy/api/proxyApi";
import { pingDatabase } from "@/features/settings/api/settingsApi";

export function DatabaseMaintenanceCard() {
  const qc = useQueryClient();
  const [exportFormat, setExportFormat] = useState<"jsonl" | "csv">("jsonl");
  const [since, setSince] = useState("");
  const [includeContent, setIncludeContent] = useState(false);

  const statsQ = useQuery({
    queryKey: ["storage-stats"],
    queryFn: getStorageStats,
    refetchInterval: 5000,
  });

  const pingMut = useMutation({
    mutationFn: pingDatabase,
    onSuccess: () => toast.success("Database connection OK"),
    onError: (e: Error) => toast.error(e.message),
  });

  const migrateMut = useMutation({
    mutationFn: runDatabaseMigrations,
    onSuccess: () => {
      toast.success("Migrations applied");
      qc.invalidateQueries({ queryKey: ["storage-stats"] });
    },
    onError: (e: Error) => toast.error(e.message),
  });

  const exportMut = useMutation({
    mutationFn: async () => {
      const payload = await exportMetadata(exportFormat, since.trim(), includeContent);
      const path = await saveMetadataExport(payload);
      return path;
    },
    onSuccess: (path) => {
      if (path) toast.success(`Saved to ${path}`);
      else toast.success("Export cancelled");
    },
    onError: (e: Error) => toast.error(e.message),
  });

  const stats = statsQ.data;

  return (
    <Card>
      <CardHeader>
        <CardTitle className="text-base">Database maintenance</CardTitle>
        <CardDescription>
          Connectivity check, schema migrations, storage stats, and metadata export (JSONL/CSV).
        </CardDescription>
      </CardHeader>
      <CardContent className="space-y-6">
        <div className="flex flex-wrap gap-2">
          <Button type="button" variant="outline" disabled={pingMut.isPending} onClick={() => pingMut.mutate()}>
            Test connection
          </Button>
          <Button type="button" variant="outline" disabled={migrateMut.isPending} onClick={() => migrateMut.mutate()}>
            Run migrations
          </Button>
          <Button
            type="button"
            variant="ghost"
            size="sm"
            disabled={statsQ.isFetching}
            onClick={() => statsQ.refetch()}
          >
            <RefreshCw className={`mr-2 h-4 w-4 ${statsQ.isFetching ? "animate-spin" : ""}`} />
            Refresh stats
          </Button>
        </div>

        {statsQ.error && (
          <p className="text-sm text-destructive">{(statsQ.error as Error).message}</p>
        )}
        {stats && (
          <dl className="grid grid-cols-2 gap-x-4 gap-y-1 text-sm sm:grid-cols-3">
            <div>
              <dt className="text-muted-foreground">ai_requests rows</dt>
              <dd className="font-mono">{stats.request_rows}</dd>
            </div>
            {stats.proxy_running && (
              <>
                <div>
                  <dt className="text-muted-foreground">Writer inserted</dt>
                  <dd className="font-mono">{stats.inserted}</dd>
                </div>
                <div>
                  <dt className="text-muted-foreground">Writer dropped</dt>
                  <dd className="font-mono">{stats.dropped}</dd>
                </div>
                <div>
                  <dt className="text-muted-foreground">Writer queue</dt>
                  <dd className="font-mono">{stats.queue_len}</dd>
                </div>
              </>
            )}
          </dl>
        )}

        <div className="space-y-3 border-t pt-4">
          <p className="text-sm font-medium">Export metadata</p>
          <div className="flex flex-wrap gap-4">
            <div className="space-y-2">
              <Label htmlFor="export-format">Format</Label>
              <select
                id="export-format"
                className="flex h-9 w-full rounded-md border border-input bg-transparent px-3 py-1 text-sm shadow-sm"
                value={exportFormat}
                onChange={(e) => setExportFormat(e.target.value as "jsonl" | "csv")}
              >
                <option value="jsonl">JSONL</option>
                <option value="csv">CSV</option>
              </select>
            </div>
            <div className="space-y-2">
              <Label htmlFor="export-since">Since (optional)</Label>
              <Input
                id="export-since"
                type="date"
                value={since}
                onChange={(e) => setSince(e.target.value)}
                className="w-[11rem]"
              />
            </div>
          </div>
          <div className="flex items-center gap-2">
            <Switch id="export-content" checked={includeContent} onCheckedChange={setIncludeContent} />
            <Label htmlFor="export-content">Include prompt/response content</Label>
          </div>
          <Button type="button" disabled={exportMut.isPending} onClick={() => exportMut.mutate()}>
            <Download className="mr-2 h-4 w-4" />
            Export to file…
          </Button>
        </div>
      </CardContent>
    </Card>
  );
}
