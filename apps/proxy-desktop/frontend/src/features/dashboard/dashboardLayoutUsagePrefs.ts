import type { HourlyChartBucketUnit } from "@/features/dashboard/components/hourlyChartBucket";
import type { HourlyChartMetric } from "@/features/dashboard/components/hourlyChartMetric";
import {
  defaultCustomHourlyRange,
  type HourlyChartRangeMode,
} from "@/features/dashboard/components/hourlyChartRange";

export type DashboardLayoutUsagePrefs = {
  tag: string;
  projectID: string;
  apiKeyID: string;
  metric: HourlyChartMetric;
  bucket: HourlyChartBucketUnit;
  rangeMode: HourlyChartRangeMode;
  customRangeStart: string;
  customRangeEnd: string;
};

const METRICS: HourlyChartMetric[] = ["requests", "tokens", "cost", "generation_tps"];
const BUCKETS: HourlyChartBucketUnit[] = ["minute", "quarter", "hour", "day"];

function isRangeMode(value: unknown): value is HourlyChartRangeMode {
  return value === 1 || value === 7 || value === 30 || value === "custom";
}

export function defaultDashboardLayoutUsagePrefs(): DashboardLayoutUsagePrefs {
  const custom = defaultCustomHourlyRange();
  return {
    tag: "",
    projectID: "",
    apiKeyID: "",
    metric: "requests",
    bucket: "hour",
    rangeMode: 7,
    customRangeStart: custom.start,
    customRangeEnd: custom.end,
  };
}

export function normalizeDashboardLayoutUsagePrefs(raw: unknown): DashboardLayoutUsagePrefs {
  const defaults = defaultDashboardLayoutUsagePrefs();
  if (!raw || typeof raw !== "object") return defaults;
  const o = raw as Record<string, unknown>;
  const metric = METRICS.includes(o.metric as HourlyChartMetric)
    ? (o.metric as HourlyChartMetric)
    : defaults.metric;
  const bucket = BUCKETS.includes(o.bucket as HourlyChartBucketUnit)
    ? (o.bucket as HourlyChartBucketUnit)
    : defaults.bucket;
  const rangeMode = isRangeMode(o.rangeMode) ? o.rangeMode : defaults.rangeMode;
  const customRangeStart =
    typeof o.customRangeStart === "string" && /^\d{4}-\d{2}-\d{2}$/.test(o.customRangeStart)
      ? o.customRangeStart
      : defaults.customRangeStart;
  const customRangeEnd =
    typeof o.customRangeEnd === "string" && /^\d{4}-\d{2}-\d{2}$/.test(o.customRangeEnd)
      ? o.customRangeEnd
      : defaults.customRangeEnd;
  return {
    tag: typeof o.tag === "string" ? o.tag : defaults.tag,
    projectID: typeof o.projectID === "string" ? o.projectID : defaults.projectID,
    apiKeyID: typeof o.apiKeyID === "string" ? o.apiKeyID : defaults.apiKeyID,
    metric,
    bucket,
    rangeMode,
    customRangeStart,
    customRangeEnd,
  };
}

export function mergeDashboardLayoutUsagePrefs(
  current: DashboardLayoutUsagePrefs,
  patch: Partial<DashboardLayoutUsagePrefs>,
): DashboardLayoutUsagePrefs {
  return normalizeDashboardLayoutUsagePrefs({ ...current, ...patch });
}
