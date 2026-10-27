import type { DashboardLayoutUsagePrefs } from "@/features/dashboard/dashboardLayout";
import { hourlyChartBucketLabel } from "@/features/dashboard/components/hourlyChartBucket";
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
  const bucketLabel = hourlyChartBucketLabel(usagePrefs.bucket);
  return `Layout filters: project ${project} · tag ${tag} · ${bucketLabel} · metric ${metricLabel} · ${rangeLabel}`;
}
