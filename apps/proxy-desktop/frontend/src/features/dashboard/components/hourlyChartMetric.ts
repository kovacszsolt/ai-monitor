import {
  formatChartAxisCount,
  formatChartAxisTokPerSec,
  formatChartAxisUSD,
  formatNumber,
  formatTokensPerSec,
  formatUSD,
} from "@/lib/format";
import { service } from "../../../../wailsjs/go/models";

export type HourlyChartMetric = "requests" | "tokens" | "cost" | "generation_tps";

export const HOURLY_CHART_METRICS: { id: HourlyChartMetric; label: string }[] = [
  { id: "requests", label: "Requests" },
  { id: "tokens", label: "Tokens" },
  { id: "cost", label: "Cost" },
  { id: "generation_tps", label: "Tok/s" },
];

export function hourlyPointValue(
  p: service.HourlyModelRequestCountDTO,
  metric: Exclude<HourlyChartMetric, "tokens">,
): number {
  switch (metric) {
    case "requests":
      return p.request_count;
    case "cost":
      return p.has_total_cost ? p.total_cost_usd : 0;
    case "generation_tps":
      return p.has_generation_tps ? p.generation_tps : 0;
  }
}

export function formatHourlyChartValue(metric: HourlyChartMetric, value: number): string {
  switch (metric) {
    case "requests":
    case "tokens":
      return formatNumber(value);
    case "cost":
      return formatUSD(true, value);
    case "generation_tps":
      return formatTokensPerSec(value);
  }
}

/** Compact Y-axis labels (tooltips keep full precision via formatHourlyChartValue). */
export function formatHourlyChartAxisValue(metric: HourlyChartMetric, value: number): string {
  switch (metric) {
    case "requests":
    case "tokens":
      return formatChartAxisCount(value);
    case "cost":
      return formatChartAxisUSD(value);
    case "generation_tps":
      return formatChartAxisTokPerSec(value);
  }
}

export function hourlyChartMetricDescription(metric: HourlyChartMetric): string {
  switch (metric) {
    case "requests":
      return "Request count per model";
    case "tokens":
      return "Per model: stacked input, cached input, and output (column height = total tokens)";
    case "cost":
      return "Estimated USD from model pricing";
    case "generation_tps":
      return "Output tok/s per model (sum output ÷ sum decode time in each hour)";
  }
}
