import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useEffect, useState } from "react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { Input } from "@/components/ui/input";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { getProxyVariables, updateProxyVariables } from "@/features/settings/api/settingsApi";
import { ByteSizeInput } from "@/shared/components/ByteSizeInput";
import { PageHeader } from "@/shared/components/PageHeader";
import { PageShell } from "@/shared/components/PageShell";
import { ErrorState, LoadingState } from "@/shared/components/PageStates";
import { service } from "../../../../wailsjs/go/models";

const selectClass =
  "flex h-9 w-full rounded border border-outline-variant bg-surface-container px-space-md py-1 text-body-sm transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/40";

function SettingKey({ name }: { name: string }) {
  return (
    <p className="font-mono text-code-inline text-on-surface-variant">
      <code>{name}</code>
    </p>
  );
}

export function ProxyVariablesPage() {
  const qc = useQueryClient();
  const proxyVarsQ = useQuery({ queryKey: ["settings", "proxy-variables"], queryFn: getProxyVariables });

  const [serverHost, setServerHost] = useState("");
  const [serverPort, setServerPort] = useState(8080);
  const [serverMaxRequestBytes, setServerMaxRequestBytes] = useState(64 * 1024 * 1024);
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
      qc.invalidateQueries({ queryKey: ["settings", "proxy-variables"] });
      qc.invalidateQueries({ queryKey: ["proxy-variables"] });
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
            Proxy runtime settings in SQLite <code className="text-code-inline">settings</code>. Stop the proxy
            before saving; restart it to apply.
          </>
        }
      />

      <Card>
        <CardHeader>
          <CardTitle className="text-base">HTTP server</CardTitle>
          <CardDescription>Listen address and maximum accepted request body size (chat, API clients).</CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="server-host">Host</Label>
              <SettingKey name="server.host" />
              <Input
                id="server-host"
                value={serverHost}
                onChange={(e) => setServerHost(e.target.value)}
                className="font-mono text-xs"
              />
            </div>
            <div className="space-y-2">
              <Label htmlFor="server-port">Port</Label>
              <SettingKey name="server.port" />
              <Input
                id="server-port"
                type="number"
                value={serverPort}
                onChange={(e) => setServerPort(Number(e.target.value))}
              />
            </div>
          </div>

          <div className="space-y-2">
            <Label htmlFor="max-request">Max request body size</Label>
            <SettingKey name="server.max_request_bytes" />
            <ByteSizeInput
              id="max-request"
              value={serverMaxRequestBytes}
              onChange={setServerMaxRequestBytes}
            />
            <p className="text-body-sm text-on-surface-variant">
              Requests larger than this are rejected. Vision chat often needs 64 MiB or more.
            </p>
          </div>
        </CardContent>
      </Card>

      <Card className="mt-4">
        <CardHeader>
          <CardTitle className="text-base">Logging</CardTitle>
        </CardHeader>
        <CardContent className="grid gap-4 sm:grid-cols-2">
          <div className="space-y-2">
            <Label htmlFor="log-level">Level</Label>
            <SettingKey name="logging.level" />
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
            <Label htmlFor="log-format">Format</Label>
            <SettingKey name="logging.format" />
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
        </CardContent>
      </Card>

      <Card className="mt-4">
        <CardHeader>
          <CardTitle className="text-base">Request capture</CardTitle>
          <CardDescription>
            Metadata (timings, tokens, model) is always recorded when enabled. Content fields below apply only when
            storing prompt/response bodies.
          </CardDescription>
        </CardHeader>
        <CardContent className="space-y-4">
          <div className="space-y-2">
            <Label htmlFor="capture-mode">Capture mode</Label>
            <SettingKey name="capture.mode" />
            <Input id="capture-mode" value={pv.capture_mode} disabled className="font-mono text-xs" />
          </div>

          <div className="flex items-center gap-2">
            <Switch id="store-content" checked={captureStoreContent} onCheckedChange={setCaptureStoreContent} />
            <div>
              <Label htmlFor="store-content">Store prompt and response content</Label>
              <SettingKey name="capture.store_content" />
            </div>
          </div>

          <div className="grid gap-4 sm:grid-cols-2">
            <div className="space-y-2">
              <Label htmlFor="max-prompt">Do not store request body when larger than</Label>
              <SettingKey name="capture.max_prompt_bytes" />
              <ByteSizeInput
                id="max-prompt"
                disabled={!captureStoreContent}
                value={captureMaxPromptBytes}
                onChange={setCaptureMaxPromptBytes}
              />
              <p className="text-body-sm text-on-surface-variant">
                Above this size, prompt/messages are omitted from the DB (request row metadata remains).
              </p>
            </div>
            <div className="space-y-2">
              <Label htmlFor="max-response">Do not store response body when larger than</Label>
              <SettingKey name="capture.max_response_bytes" />
              <ByteSizeInput
                id="max-response"
                disabled={!captureStoreContent}
                value={captureMaxResponseBytes}
                onChange={setCaptureMaxResponseBytes}
              />
              <p className="text-body-sm text-on-surface-variant">
                Above this size, assistant text / tool output is not persisted.
              </p>
            </div>
          </div>

          <p className="text-body-sm text-on-surface-variant">
            Both capture limits must be ≤ <code className="text-code-inline">server.max_request_bytes</code>.
          </p>

          <Button type="button" disabled={saveMut.isPending} onClick={() => saveMut.mutate()}>
            Save proxy settings
          </Button>
        </CardContent>
      </Card>
    </PageShell>
  );
}
