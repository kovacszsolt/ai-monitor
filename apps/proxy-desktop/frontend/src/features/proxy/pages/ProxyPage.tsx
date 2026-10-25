import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { Link } from "react-router-dom";
import { useState } from "react";
import { toast } from "sonner";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Card, CardContent } from "@/components/ui/card";
import {
  checkProviders,
  getProxyStatus,
  startProxy,
  stopProxy,
  validateProxyConfig,
} from "@/features/proxy/api/proxyApi";
import { MaterialIcon } from "@/shared/components/MaterialIcon";
import { PageHeader } from "@/shared/components/PageHeader";
import { PageShell } from "@/shared/components/PageShell";
import { StatusPulse } from "@/shared/components/StatusPulse";
import { SurfaceCard } from "@/shared/components/SurfaceCard";
import { ErrorState, LoadingState } from "@/shared/components/PageStates";
import { service } from "../../../../wailsjs/go/models";

function applyHealthCheckResult(
  failed: service.UnhealthyProviderDTO[],
  setUnhealthy: (v: service.UnhealthyProviderDTO[]) => void,
  setChecked: (v: boolean) => void,
  toastOnSuccess = true,
) {
  setUnhealthy(failed);
  setChecked(true);
  if (!toastOnSuccess) return;
  if (failed.length > 0) {
    toast.warning(`${failed.length} provider(s) failed the health check`);
  } else {
    toast.success("All enabled providers are reachable");
  }
}

export function ProxyPage() {
  const qc = useQueryClient();
  const [unhealthyProviders, setUnhealthyProviders] = useState<service.UnhealthyProviderDTO[]>([]);
  const [providersChecked, setProvidersChecked] = useState(false);

  const configQ = useQuery({ queryKey: ["proxy-config"], queryFn: validateProxyConfig });
  const statusQ = useQuery({
    queryKey: ["proxy-status"],
    queryFn: getProxyStatus,
    refetchInterval: 3000,
  });

  const recheckMut = useMutation({
    mutationFn: checkProviders,
    onSuccess: (result) => {
      applyHealthCheckResult(result.unhealthy_providers ?? [], setUnhealthyProviders, setProvidersChecked);
    },
    onError: (e: Error) => toast.error(e.message),
  });

  const startMut = useMutation({
    mutationFn: startProxy,
    onSuccess: (result) => {
      const failed = result.unhealthy_providers ?? [];
      applyHealthCheckResult(failed, setUnhealthyProviders, setProvidersChecked, false);
      if (failed.length > 0) {
        toast.warning(`Proxy started; ${failed.length} provider(s) failed the health check`);
      } else {
        toast.success("Proxy started");
      }
      qc.invalidateQueries({ queryKey: ["proxy-status"] });
    },
    onError: (e: Error) => toast.error(e.message),
  });

  const stopMut = useMutation({
    mutationFn: stopProxy,
    onSuccess: () => {
      setUnhealthyProviders([]);
      setProvidersChecked(false);
      toast.success("Proxy stopped");
      qc.invalidateQueries({ queryKey: ["proxy-status"] });
    },
    onError: (e: Error) => toast.error(e.message),
  });

  if (configQ.isLoading) return <LoadingState label="Loading proxy config…" />;
  if (configQ.error) return <ErrorState message={(configQ.error as Error).message} onRetry={() => configQ.refetch()} />;

  const cfg = configQ.data;
  const status = statusQ.data;
  const showProviderHealth = Boolean(status?.running && providersChecked);
  const running = Boolean(status?.running);
  const listenUrl = status?.listen_addr ? `http://${status.listen_addr}` : cfg?.listen_addr ? `http://${cfg.listen_addr}` : null;

  return (
    <PageShell>
      <PageHeader
        section="Operate"
        accent="HTTP Gateway"
        title="Proxy Server"
        description="Control, orchestrate, and inspect the local in-process AI reverse proxy runtime."
        actions={
          <Button type="button" variant="secondary" asChild>
            <Link to="/proxy-variables">
              <MaterialIcon name="tune" size={16} className="text-secondary" />
              Proxy Variables
            </Link>
          </Button>
        }
      />

      <Card className="relative overflow-hidden p-space-lg">
        <div className="pointer-events-none absolute -right-20 -top-20 h-72 w-72 rounded-full bg-tertiary/5 blur-3xl" />
        <CardContent className="relative z-10 space-y-space-lg p-0">
          <div className="flex flex-col justify-between gap-space-lg lg:flex-row lg:items-center">
            <div className="flex flex-col gap-space-xs">
              <div className="flex flex-wrap items-center gap-space-sm">
                <span className="font-headline-md text-headline-md text-on-surface">HTTP Server</span>
                <div className="inline-flex items-center gap-1.5 rounded-full bg-surface-container-lowest px-2.5 py-1 shadow-sm">
                  <StatusPulse active={running} variant={running ? "tertiary" : "error"} />
                  <span className="font-label-sm text-label-sm font-semibold tracking-wide text-tertiary">
                    {running ? "RUNNING" : "STOPPED"}
                  </span>
                </div>
                {listenUrl && (
                  <div className="flex items-center gap-1.5 rounded-lg bg-surface-container-low px-2.5 py-1">
                    <MaterialIcon name="lan" size={14} className="text-outline" />
                    <code className="text-code-inline text-secondary-fixed">{listenUrl}</code>
                    <button
                      type="button"
                      className="ml-1 rounded p-0.5 text-on-surface-variant transition-colors hover:bg-surface-container hover:text-on-surface"
                      title="Copy endpoint"
                      onClick={() => {
                        void navigator.clipboard.writeText(listenUrl);
                        toast.success("Copied base URL");
                      }}
                    >
                      <MaterialIcon name="content_copy" size={14} />
                    </button>
                  </div>
                )}
              </div>
              <p className="max-w-2xl text-body-sm text-on-surface-variant">
                Listen interface, capture, and logging are loaded from{" "}
                <code className="text-code-inline text-on-surface">proxy.db</code>.
              </p>
              {status?.error && <p className="text-body-sm text-error">{status.error}</p>}
            </div>
            <div className="flex flex-wrap items-center gap-space-sm">
              {running ? (
                <Button
                  type="button"
                  variant="destructive"
                  disabled={stopMut.isPending}
                  onClick={() => stopMut.mutate()}
                >
                  <MaterialIcon name="stop_circle" size={18} />
                  Stop Server
                </Button>
              ) : (
                <Button type="button" disabled={startMut.isPending} onClick={() => startMut.mutate()}>
                  <MaterialIcon name="play_circle" size={18} />
                  Start Server
                </Button>
              )}
              <Button
                type="button"
                variant="secondary"
                disabled={!running || recheckMut.isPending}
                onClick={() => recheckMut.mutate()}
              >
                <MaterialIcon
                  name="network_check"
                  size={16}
                  className={recheckMut.isPending ? "animate-spin" : "text-tertiary"}
                />
                Probe Upstreams
              </Button>
            </div>
          </div>
        </CardContent>
      </Card>

      {showProviderHealth && (
        <SurfaceCard>
          <div className="flex flex-wrap items-center justify-between gap-2">
            <p
              className={
                unhealthyProviders.length > 0
                  ? "text-body-sm font-medium text-error"
                  : "text-body-sm font-medium text-on-surface"
              }
            >
              {unhealthyProviders.length > 0
                ? "Providers not reachable"
                : "All enabled providers are reachable"}
            </p>
            <Badge variant={unhealthyProviders.length > 0 ? "destructive" : "success"}>
              {unhealthyProviders.length > 0 ? `${unhealthyProviders.length} failed` : "Healthy"}
            </Badge>
          </div>
          {unhealthyProviders.length > 0 && (
            <ul className="mt-space-md space-y-2 text-body-sm">
              {unhealthyProviders.map((p) => (
                <li key={p.id}>
                  <span className="font-medium text-on-surface">{p.name}</span>
                  <span className="text-on-surface-variant"> — {p.message}</span>
                </li>
              ))}
            </ul>
          )}
        </SurfaceCard>
      )}
    </PageShell>
  );
}
