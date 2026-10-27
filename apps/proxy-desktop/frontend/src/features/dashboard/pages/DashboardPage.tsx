import * as React from "react";
import { useQuery } from "@tanstack/react-query";
import { MaterialIcon } from "@/shared/components/MaterialIcon";
import { Badge } from "@/components/ui/badge";
import { Button } from "@/components/ui/button";
import { Table, TableBody, TableCell, TableHead, TableHeader, TableRow } from "@/components/ui/table";
import { DashboardCollapsibleSection } from "@/features/dashboard/components/DashboardCollapsibleSection";
import { DashboardRequestVolumeCard } from "@/features/dashboard/components/DashboardRequestVolumeCard";
import { DashboardExportControls } from "@/features/dashboard/components/DashboardExportControls";
import { DashboardLayoutHeaderControls } from "@/features/dashboard/components/DashboardLayoutHeaderControls";
import { DashboardExportProvider } from "@/features/dashboard/export/DashboardExportContext";
import { listProjects } from "@/features/projects/api/projectsApi";
import { DashboardLayoutPresetsBar } from "@/features/dashboard/components/DashboardLayoutPresetsBar";
import { DashboardSectionShell } from "@/features/dashboard/components/DashboardSectionShell";
import {
  ModelTokenUsageByTagChart,
  ModelUsageByTagSection,
  useModelTagLabel,
} from "@/features/dashboard/components/ModelUsageByTagSection";
import { UsageGroupTable } from "@/features/dashboard/components/UsageGroupTable";
import { DASHBOARD_SECTION_IDS } from "@/features/dashboard/dashboardSectionIds";
import type { DashboardSectionId } from "@/features/dashboard/dashboardLayout";
import { useDashboardLayout } from "@/features/dashboard/hooks/useDashboardLayout";
import { getDashboardStats } from "@/features/requests/api/requestsApi";
import { formatDuration, formatNumber } from "@/lib/format";
import { HoverTooltip } from "@/shared/components/HoverTooltip";
import { ErrorState, LoadingState } from "@/shared/components/PageStates";

function StatCardBody({
  value,
  hint,
  icon,
  iconTooltip,
}: {
  value: string;
  hint?: string;
  icon: string;
  iconTooltip: string;
}) {
  return (
    <div className="flex items-start justify-between gap-3">
      <div>
        <div className="font-headline-md text-headline-md tabular-nums tracking-tight text-on-surface">{value}</div>
        {hint && <p className="mt-1 text-body-sm text-on-surface-variant">{hint}</p>}
      </div>
      <HoverTooltip label={iconTooltip} side="left">
        <span className="mt-1 inline-flex shrink-0 rounded-sm">
          <MaterialIcon name={icon} size={18} className="text-on-surface-variant" />
        </span>
      </HoverTooltip>
    </div>
  );
}

function stateVariant(state: string): "success" | "destructive" | "secondary" | "default" {
  switch (state) {
    case "completed":
      return "success";
    case "failed":
      return "destructive";
    case "cancelled":
    case "interrupted":
      return "secondary";
    default:
      return "default";
  }
}

export function DashboardPage() {
  const [editLayout, setEditLayout] = React.useState(false);
  const [exporting, setExporting] = React.useState(false);
  const exportRootRef = React.useRef<HTMLDivElement>(null);
  const {
    layout,
    presets,
    activePresetId,
    activePresetName,
    moveSection,
    setSectionWidth,
    resetLayout,
    switchPreset,
    copyLayout,
    removePreset,
    renamePreset,
    usagePrefs,
    updateUsagePrefs,
  } = useDashboardLayout();
  const selectedTag = usagePrefs.tag;
  const selectedProjectID = usagePrefs.projectID;
  const tagLabel = useModelTagLabel(selectedTag);
  const projectsQ = useQuery({ queryKey: ["projects"], queryFn: listProjects });
  const projectLabel =
    selectedProjectID === ""
      ? ""
      : (projectsQ.data?.find((p) => p.id === selectedProjectID)?.name ?? selectedProjectID);
  const { data, isLoading, error, refetch } = useQuery({
    queryKey: ["dashboard"],
    queryFn: getDashboardStats,
    refetchInterval: 30000,
  });

  const renderSection = React.useCallback(
    (id: DashboardSectionId): React.ReactNode => {
      if (!data) return null;
      switch (id) {
        case DASHBOARD_SECTION_IDS.usageByHour:
          return (
            <DashboardRequestVolumeCard
              tag={selectedTag}
              onTagChange={(tag) => updateUsagePrefs({ tag })}
              projectID={selectedProjectID}
              onProjectChange={(projectID) => updateUsagePrefs({ projectID })}
              metric={usagePrefs.metric}
              onMetricChange={(metric) => updateUsagePrefs({ metric })}
              bucket={usagePrefs.bucket}
              onBucketChange={(bucket) => updateUsagePrefs({ bucket })}
              rangeMode={usagePrefs.rangeMode}
              onRangeModeChange={(rangeMode) => updateUsagePrefs({ rangeMode })}
              customRangeStart={usagePrefs.customRangeStart}
              customRangeEnd={usagePrefs.customRangeEnd}
              onCustomRangeStartChange={(customRangeStart) => updateUsagePrefs({ customRangeStart })}
              onCustomRangeEndChange={(customRangeEnd) => updateUsagePrefs({ customRangeEnd })}
            />
          );
        case DASHBOARD_SECTION_IDS.usageByModel:
          return (
            <ModelUsageByTagSection
              tag={selectedTag}
              onTagChange={(tag) => updateUsagePrefs({ tag })}
              projectID={selectedProjectID}
              onProjectChange={(projectID) => updateUsagePrefs({ projectID })}
              hideFilters
            />
          );
        case DASHBOARD_SECTION_IDS.statTotalRequests:
          return (
            <DashboardCollapsibleSection sectionId={id} title="Total requests">
              <StatCardBody
                value={formatNumber(data.total_requests)}
                icon="insights"
                iconTooltip="All-time request count"
              />
            </DashboardCollapsibleSection>
          );
        case DASHBOARD_SECTION_IDS.statRequests24h:
          return (
            <DashboardCollapsibleSection sectionId={id} title="Last 24 hours">
              <StatCardBody
                value={formatNumber(data.requests_last_24h)}
                hint="Rolling window"
                icon="schedule"
                iconTooltip="Requests in the last 24 hours"
              />
            </DashboardCollapsibleSection>
          );
        case DASHBOARD_SECTION_IDS.statTokens:
          return (
            <DashboardCollapsibleSection sectionId={id} title="Tokens">
              <StatCardBody
                value={`${formatNumber(data.sum_input_tokens)} / ${formatNumber(data.sum_cached_input_tokens ?? 0)} / ${formatNumber(data.sum_output_tokens)}`}
                icon="token"
                iconTooltip="Input / cached / output tokens"
              />
            </DashboardCollapsibleSection>
          );
        case DASHBOARD_SECTION_IDS.statAvgDuration:
          return (
            <DashboardCollapsibleSection sectionId={id} title="Avg duration">
              <StatCardBody
                value={data.avg_duration_ms > 0 ? formatDuration(Math.round(data.avg_duration_ms)) : "—"}
                hint={`${formatNumber(data.streaming_count)} streaming`}
                icon="speed"
                iconTooltip="Average request duration"
              />
            </DashboardCollapsibleSection>
          );
        case DASHBOARD_SECTION_IDS.usageByProxyTag:
          return (
            <UsageGroupTable
              collapsible
              title="Usage by proxy tag"
              description='From X-AI-Proxy-Tag on proxied requests (empty tag shown as “(no tag)”). Click a row to filter.'
              rows={data.usage_by_tag ?? []}
              emptyLabel="No tag column in database yet — run migrations (Settings → Database maintenance)."
              selectedTag={selectedTag}
              onSelectTag={(tag) => updateUsagePrefs({ tag })}
            />
          );
        case DASHBOARD_SECTION_IDS.tokenMixChart:
          return (
            <ModelTokenUsageByTagChart
              tag={selectedTag}
              tagLabel={tagLabel}
              projectID={selectedProjectID}
            />
          );
        case DASHBOARD_SECTION_IDS.byState:
          return (
            <DashboardCollapsibleSection sectionId={id} title="By state">
              {data.by_state.length === 0 ? (
                <p className="text-sm text-muted-foreground">No requests yet.</p>
              ) : (
                <ul className="space-y-2">
                  {data.by_state.map((s) => (
                    <li key={s.state} className="flex items-center justify-between text-sm">
                      <Badge variant={stateVariant(s.state)}>{s.state}</Badge>
                      <span className="font-mono text-muted-foreground">{formatNumber(s.count)}</span>
                    </li>
                  ))}
                </ul>
              )}
            </DashboardCollapsibleSection>
          );
        case DASHBOARD_SECTION_IDS.topModels:
          return (
            <DashboardCollapsibleSection sectionId={id} title="Top models">
              <p className="mb-3 text-sm text-muted-foreground">By request count</p>
              <MiniCountTable rows={data.top_models} />
            </DashboardCollapsibleSection>
          );
        case DASHBOARD_SECTION_IDS.topProjects:
          return (
            <DashboardCollapsibleSection sectionId={id} title="Top projects">
              <p className="mb-3 text-sm text-muted-foreground">By request count</p>
              <MiniCountTable rows={data.top_projects} />
            </DashboardCollapsibleSection>
          );
        default:
          return null;
      }
    },
    [data, selectedTag, selectedProjectID, tagLabel, usagePrefs, updateUsagePrefs],
  );

  if (isLoading) return <LoadingState label="Loading dashboard…" />;
  if (error) return <ErrorState message={(error as Error).message} onRetry={() => refetch()} />;
  if (!data) return null;

  const visibleLayout = layout.filter((entry) => {
    if (entry.id === DASHBOARD_SECTION_IDS.tokenMixChart && selectedTag === "") {
      return false;
    }
    return true;
  });

  return (
    <DashboardExportProvider forceExpanded={exporting}>
    <div className="flex flex-col">
      <div className="border-b border-outline-variant/30 bg-surface-container-lowest/80 px-space-xl py-space-lg shadow-sm backdrop-blur-md">
        <div className="flex flex-col gap-space-md lg:flex-row lg:items-center lg:justify-between">
          <div>
            <h1 className="font-headline-lg text-headline-lg tracking-tight text-on-surface">Analytics &amp; Telemetry</h1>
            <p className="text-body-sm text-on-surface-variant">
              Real-time proxied traffic, model distribution, and latency metrics across local workloads.
            </p>
          </div>
        <div className="flex flex-wrap items-center gap-2">
          <DashboardExportControls
            exportRootRef={exportRootRef}
            layoutName={activePresetName}
            usagePrefs={usagePrefs}
            projectLabel={projectLabel}
            tagLabel={tagLabel}
            onExportStart={() => setExporting(true)}
            onExportEnd={() => setExporting(false)}
          />
          <DashboardLayoutHeaderControls
            presets={presets}
            activePresetId={activePresetId}
            onSwitchPreset={switchPreset}
            onCopyLayout={copyLayout}
          />
          {editLayout && (
            <Button type="button" size="sm" variant="ghost" onClick={resetLayout}>
              Reset layout
            </Button>
          )}
          <Button
            type="button"
            size="sm"
            variant={editLayout ? "default" : "secondary"}
            onClick={() => setEditLayout((v) => !v)}
          >
            <MaterialIcon name="dashboard_customize" size={16} />
            {editLayout ? "Done" : "Customize Grid"}
          </Button>
        </div>
        </div>
      </div>

      <div className="space-y-space-xl p-space-xl">
      {editLayout && (
        <>
          <DashboardLayoutPresetsBar
            presets={presets}
            activePresetId={activePresetId}
            onDeletePreset={removePreset}
            onRenamePreset={renamePreset}
          />
          <p className="text-sm text-muted-foreground">
            Section order and width save automatically (⅓, ½, or full row on large screens). Edit the layout name
            and press Save to rename.
          </p>
        </>
      )}

      <div
        ref={exportRootRef}
        className="dashboard-export-root space-y-4 rounded-xl bg-background"
        data-exporting={exporting ? "true" : undefined}
      >
        <div className="space-y-1">
          <p className="text-lg font-semibold tracking-tight">Dashboard — {activePresetName}</p>
          <p className="text-sm text-muted-foreground">
            {projectLabel || "All projects"} · {tagLabel || "All tags"}
          </p>
        </div>
        <div className="grid grid-cols-1 gap-6 lg:grid-cols-6">
          {visibleLayout.map((entry, index) => (
            <DashboardSectionShell
              key={entry.id}
              sectionId={entry.id}
              width={entry.width}
              editMode={editLayout && !exporting}
              canMoveUp={index > 0}
              canMoveDown={index < visibleLayout.length - 1}
              onMoveUp={() => moveSection(entry.id, -1)}
              onMoveDown={() => moveSection(entry.id, 1)}
              onSetWidth={(width) => setSectionWidth(entry.id, width)}
            >
              {renderSection(entry.id)}
            </DashboardSectionShell>
          ))}
        </div>
      </div>
      </div>
    </div>
    </DashboardExportProvider>
  );
}

function MiniCountTable({ rows }: { rows: { name: string; count: number }[] }) {
  if (rows.length === 0) {
    return <p className="text-sm text-muted-foreground">No data.</p>;
  }
  return (
    <Table>
      <TableHeader>
        <TableRow>
          <TableHead>Name</TableHead>
          <TableHead className="text-right">Count</TableHead>
        </TableRow>
      </TableHeader>
      <TableBody>
        {rows.map((r) => (
          <TableRow key={r.name}>
            <TableCell className="max-w-[180px] truncate font-medium">{r.name}</TableCell>
            <TableCell className="text-right font-mono text-muted-foreground">{formatNumber(r.count)}</TableCell>
          </TableRow>
        ))}
      </TableBody>
    </Table>
  );
}
