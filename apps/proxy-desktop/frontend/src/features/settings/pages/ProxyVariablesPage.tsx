import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { getProxyVariables, updateProxyVariables } from "@/features/settings/api/settingsApi";
import { PageHeader } from "@/shared/components/PageHeader";
import { PageShell } from "@/shared/components/PageShell";
import { ErrorState, LoadingState } from "@/shared/components/PageStates";
import { service } from "../../../../wailsjs/go/models";

const selectClass =
  "flex h-9 w-full rounded border border-outline-variant bg-surface-container px-space-md py-1 text-body-sm transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/40";

export function ProxyVariablesPage() {
  const qc = useQueryClient();
  const proxyVarsQ = useQuery({ queryKey: ["settings", "proxy-variables"], queryFn: getProxyVariables });

  const [serverHost, setServerHost] = useState("");
  const [serverPort, setServerPort] = useState(8080);
  const [serverMaxRequestBytes, setServerMaxRequestBytes] = useState(10 * 1024 * 1024);
  const [loggingLevel, setLoggingLevel] = useState("info");
  const [loggingFormat, setLoggingFormat] = useState("console");
  const [captureStoreContent, setCaptureStoreContent] = useState(false);
  const [captureMaxPromptBytes, setCaptureMaxPromptBytes] = useState(262144);
  const [captureMaxResponseBytes, setCaptureMaxResponseBytes] = useState(1048576);

  useEffect(() => {
    const v = proxyVarsQ.data;
    if (!v) return;
    setServerHost(v.server_host);
    setServerPort(v.server_port);
    setServerMaxRequestBytes(v.server_max_request_bytes);
    setLoggingLevel(v.logging_level);
    setLoggingFormat(v.logging_format);
    setCaptureStoreContent(v.capture_store_content);
    setCaptureMaxPromptBytes(v.capture_max_prompt_bytes);
    setCaptureMaxResponseBytes(v.capture_max_response_bytes);
  }, [proxyVarsQ.data]);

  const saveMut = useMutation({
    mutationFn: () => {
      const dto = new service.UpdateProxyVariablesDTO({
        server_host: serverHost.trim(),
        server_port: serverPort,
        server_max_request_bytes: serverMaxRequestBytes,
        logging_level: loggingLevel,
        logging_format: loggingFormat,
        capture_store_content: captureStoreContent,
        capture_max_prompt_bytes: captureMaxPromptBytes,
        capture_max_response_bytes: captureMaxResponseBytes,
      });
      return updateProxyVariables(dto);
    },
    onSuccess: () => {
      toast.success("Proxy settings saved");
      qc.invalidateQueries({ queryKey: ["settings", "proxy-variables", "proxy-config"] });
    },
    onError: (e: Error) => toast.error(e.message),
  });

  if (proxyVarsQ.isLoading) return <LoadingState label="Loading proxy settings…" />;
  if (proxyVarsQ.error) {
    return (
      <ErrorState message={(proxyVarsQ.error as Error).message} onRetry={() => proxyVarsQ.refetch()} />
    );
  }

  const pv = proxyVarsQ.data!;

  return (
    <PageShell>
      <PageHeader
        section="System"
        accent="Runtime"
        title="Runtime Vars"
        description={
          <>
            Runtime server, logging, and capture — stored in SQLite{" "}
            <code className="text-code-inline">settings</code>.
          </>
        }
      />

      <Card>
        <CardHeader>
          <CardTitle className="text-base">Configuration</CardTitle>
          <CardDescription>Stop the proxy before saving. Restart proxy after changes take effect.</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="server-host">Server host</Label>
              <Input
                id="server-host"
                value={serverHost}
                onChange={(e) => setServerHost(e.target.value)}
                className="font-mono text-xs"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="server-port">Server port</Label>
              <Input
                id="server-port"
                type="number"
                value={serverPort}
                onChange={(e) => setServerPort(Number(e.target.value))}
              />
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="max-request">Max request bytes</Label>
            <Input
              id="max-request"
              type="number"
              value={serverMaxRequestBytes}
              onChange={(e) => setServerMaxRequestBytes(Number(e.target.value))}
            />
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="log-level">Logging level</Label>
              <select
                id="log-level"
                className={selectClass}
                value={loggingLevel}
                onChange={(e) => setLoggingLevel(e.target.value)}
              >
                <option value="debug">debug</option>
                <option value="info">info</option>
                <option value="warn">warn</option>
                <option value="error">error</option>
              </select>
            </div>
            <div className="space-y-2">
              <Label htmlFor="log-format">Logging format</Label>
              <select
                id="log-format"
                className={selectClass}
                value={loggingFormat}
                onChange={(e) => setLoggingFormat(e.target.value)}
              >
                <option value="console">console</option>
                <option value="json">json</option>
              </select>
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="capture-mode">Capture mode</Label>
            <Input id="capture-mode" value={pv.capture_mode} disabled className="font-mono text-xs" />
          </div>

          <div className="flex items-center gap-2">
            <Switch id="store-content" checked={captureStoreContent} onCheckedChange={setCaptureStoreContent} />
            <Label htmlFor="store-content">Store prompt/response content</Label>
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="max-prompt">Max prompt bytes</Label>
              <Input
                id="max-prompt"
                type="number"
                value={captureMaxPromptBytes}
                onChange={(e) => setCaptureMaxPromptBytes(Number(e.target.value))}
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="max-response">Max response bytes</Label>
              <Input
                id="max-response"
                type="number"
                value={captureMaxResponseBytes}
                onChange={(e) => setCaptureMaxResponseBytes(Number(e.target.value))}
              />
            </div>
          </div>

          <Button type="button" disabled={saveMut.isPending} onClick={() => saveMut.mutate()}>
            Save proxy settings
          </Button>
        </CardContent>
      </Card>
    </PageShell>
  );
}
