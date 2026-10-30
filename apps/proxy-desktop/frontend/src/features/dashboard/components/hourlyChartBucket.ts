import type { HourlyChartRangeMode } from "@/features/dashboard/components/hourlyChartRange";

export type HourlyChartBucketUnit = "minute" | "quarter" | "hour" | "day";

export const HOURLY_CHART_BUCKET_OPTIONS: { id: HourlyChartBucketUnit; label: string }[] = [
  { id: "minute", label: "By minute" },
  { id: "quarter", label: "By 15 min" },
  { id: "hour", label: "By hour" },
  { id: "day", label: "By day" },
];

/** Max range (days) allowed per bucket on the backend — used for UI hints and auto-clamp. */
export const MAX_RANGE_DAYS_FOR_BUCKET: Record<HourlyChartBucketUnit, number> = {
  minute: 2,
  quarter: 14,
  hour: 90,
  day: 90,
};

export function hourlyChartBucketLabel(unit: HourlyChartBucketUnit): string {
  return HOURLY_CHART_BUCKET_OPTIONS.find((o) => o.id === unit)?.label ?? "By hour";
}

export function usageVolumeSectionTitle(unit: HourlyChartBucketUnit): string {
  switch (unit) {
    case "minute":
      return "Usage by minute";
    case "quarter":
      return "Usage by 15 minutes";
    case "day":
      return "Usage by day";
    default:
      return "Usage by hour";
  }
}

export function rangeSpanDays(mode: HourlyChartRangeMode, customStart: string, customEnd: string): number {
  if (mode === "custom") {
    if (!/^\d{4}-\d{2}-\d{2}$/.test(customStart) || !/^\d{4}-\d{2}-\d{2}$/.test(customEnd)) {
      return 7;
    }
    const start = Date.parse(`${customStart}T00:00:00Z`);
    const end = Date.parse(`${customEnd}T00:00:00Z`);
    if (!Number.isFinite(start) || !Number.isFinite(end) || end < start) {
      return 7;
    }
    return Math.floor((end - start) / 86400000) + 1;
  }
  return mode;
}

export function isBucketAllowedForRange(bucket: HourlyChartBucketUnit, spanDays: number): boolean {
  return spanDays <= MAX_RANGE_DAYS_FOR_BUCKET[bucket];
}

export function coerceBucketForRange(
  bucket: HourlyChartBucketUnit,
  spanDays: number,
): HourlyChartBucketUnit {
  if (isBucketAllowedForRange(bucket, spanDays)) {
    return bucket;
  }
  const order: HourlyChartBucketUnit[] = ["minute", "quarter", "hour", "day"];
  const start = order.indexOf(bucket);
  for (let i = start + 1; i < order.length; i++) {
    const candidate = order[i];
    if (isBucketAllowedForRange(candidate, spanDays)) {
      return candidate;
    }
  }
  return "day";
}
