import type { DashboardLayoutUsagePrefs } from "@/features/dashboard/dashboardLayout";
import { HOURLY_CHART_METRICS } from "@/features/dashboard/components/hourlyChartMetric";
import { formatHourlyChartRangeLabel } from "@/features/dashboard/components/hourlyChartRange";

export function buildExportFilterSummary(
  usagePrefs: DashboardLayoutUsagePrefs,
  projectLabel: string,
  tagLabel: string,
): string {
  const metricLabel =
    HOURLY_CHART_METRICS.find((m) => m.id === usagePrefs.metric)?.label ?? usagePrefs.metric;
  const rangeLabel = formatHourlyChartRangeLabel(
    usagePrefs.rangeMode,
    usagePrefs.customRangeStart,
    usagePrefs.customRangeEnd,
  );
  const project = projectLabel || "All projects";
  const tag = tagLabel || "All tags";
  return `Layout filters: project ${project} · tag ${tag} · metric ${metricLabel} · ${rangeLabel}`;
}
