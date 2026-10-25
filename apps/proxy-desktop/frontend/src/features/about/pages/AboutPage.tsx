import { useQuery } from "@tanstack/react-query";
import { Badge } from "@/components/ui/badge";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { getAbout } from "@/features/about/api/aboutApi";
import { formatBuildTimeLabel } from "@/lib/format";
import { PageHeader } from "@/shared/components/PageHeader";
import { PageShell } from "@/shared/components/PageShell";
import { ErrorState, LoadingState } from "@/shared/components/PageStates";

function AboutRow({ label, value }: { label: string; value: string }) {
  return (
    <div className="grid gap-1 border-b border-outline-variant/30 py-3 last:border-0 sm:grid-cols-[10rem_1fr] sm:gap-4">
      <dt className="font-label-sm text-label-sm text-on-surface-variant">{label}</dt>
      <dd className="text-body-sm text-on-surface">{value || "—"}</dd>
    </div>
  );
}

export function AboutPage() {
  const aboutQ = useQuery({ queryKey: ["about"], queryFn: getAbout });

  if (aboutQ.isLoading) return <LoadingState label="Loading about…" />;
  if (aboutQ.error) {
    return <ErrorState message={(aboutQ.error as Error).message} onRetry={() => aboutQ.refetch()} />;
  }

  const a = aboutQ.data!;

  return (
    <PageShell>
      <PageHeader
        section="System"
        title="About"
        description="Application version and build metadata."
      />

      <Card>
        <CardHeader>
          <CardTitle className="text-base">{a.app_name}</CardTitle>
          <CardDescription>Desktop client for the local AI proxy and request telemetry.</CardDescription>
        </CardHeader>
        <CardContent>
          <dl>
            <AboutRow label="Version" value={a.version} />
            <AboutRow label="Build" value={formatBuildTimeLabel(a.build_time)} />
            <AboutRow label="Author" value={a.author} />
            <div className="grid gap-1 border-b border-outline-variant/30 py-3 last:border-0 sm:grid-cols-[10rem_1fr] sm:gap-4">
              <dt className="font-label-sm text-label-sm text-on-surface-variant">License</dt>
              <dd className="space-y-1 text-body-sm text-on-surface">
                <a
                  href="https://opensource.org/license/mit"
                  target="_blank"
                  rel="noreferrer"
                  className="font-medium text-primary underline-offset-4 hover:underline"
                >
                  {a.license}
                </a>
                <p className="text-on-surface-variant">
                  Open source — use, modify, and distribute freely; include the copyright notice.
                </p>
              </dd>
            </div>
            <div className="grid gap-1 border-b border-outline-variant/30 py-3 last:border-0 sm:grid-cols-[10rem_1fr] sm:gap-4">
              <dt className="font-label-sm text-label-sm text-on-surface-variant">Database schema</dt>
              <dd className="flex flex-wrap items-center gap-2 text-body-sm">
                <span>Migration v{a.migration_version}</span>
                {a.schema_ready ? (
                  <Badge variant="secondary">Ready</Badge>
                ) : (
                  <Badge variant="destructive">Not ready</Badge>
                )}
                {a.migration_dirty ? <Badge variant="destructive">Dirty</Badge> : null}
              </dd>
            </div>
            <div className="grid gap-1 border-b border-outline-variant/30 py-3 last:border-0 sm:grid-cols-[10rem_1fr] sm:gap-4">
              <dt className="font-label-sm text-label-sm text-on-surface-variant">Proxy config bridge</dt>
              <dd className="flex flex-wrap items-center gap-2 text-body-sm">
                {a.config_bridge_ok ? (
                  <Badge variant="secondary">OK</Badge>
                ) : (
                  <Badge variant="destructive">Issue</Badge>
                )}
                {a.bridge_message ? <span className="text-on-surface-variant">{a.bridge_message}</span> : null}
              </dd>
            </div>
          </dl>
        </CardContent>
      </Card>
    </PageShell>
  );
}
