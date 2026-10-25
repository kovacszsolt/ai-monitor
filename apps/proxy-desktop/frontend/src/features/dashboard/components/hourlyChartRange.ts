import { service } from "../../../../wailsjs/go/models";

export type HourlyChartPresetDays = 1 | 7 | 30;
export type HourlyChartRangeMode = HourlyChartPresetDays | "custom";

export const HOURLY_CHART_RANGE_PRESETS: { days: HourlyChartPresetDays; label: string }[] = [
  { days: 1, label: "Last 1d" },
  { days: 7, label: "Last 7d" },
  { days: 30, label: "Last 30d" },
];

export function utcDateInputValue(d: Date): string {
  return d.toISOString().slice(0, 10);
}

export function defaultCustomHourlyRange(): { start: string; end: string } {
  const end = new Date();
  const start = new Date();
  start.setUTCDate(start.getUTCDate() - 6);
  return { start: utcDateInputValue(start), end: utcDateInputValue(end) };
}

export function isCustomHourlyRangeValid(start: string, end: string): boolean {
  if (!/^\d{4}-\d{2}-\d{2}$/.test(start) || !/^\d{4}-\d{2}-\d{2}$/.test(end)) {
    return false;
  }
  return start <= end;
}

export function hourlyRangeQueryKeyPart(
  mode: HourlyChartRangeMode,
  customStart: string,
  customEnd: string,
): string {
  if (mode === "custom") {
    return `custom:${customStart}:${customEnd}`;
  }
  return `days:${mode}`;
}

export function hourlyRangeRequestParams(
  mode: HourlyChartRangeMode,
  customStart: string,
  customEnd: string,
): { days: number; rangeStart: string; rangeEnd: string } {
  if (mode === "custom") {
    return { days: 0, rangeStart: customStart, rangeEnd: customEnd };
  }
  return { days: mode, rangeStart: "", rangeEnd: "" };
}

export function formatHourlyChartRangeLabel(
  mode: HourlyChartRangeMode,
  customStart: string,
  customEnd: string,
  series?: service.HourlyModelRequestSeriesDTO,
): string {
  const fromApi =
    series?.range_start && series?.range_end
      ? `${series.range_start} – ${series.range_end} (UTC)`
      : null;
  if (fromApi) return fromApi;
  if (mode === "custom" && isCustomHourlyRangeValid(customStart, customEnd)) {
    return `${customStart} – ${customEnd} (UTC)`;
  }
  const days = mode === "custom" ? 7 : mode;
  return `last ${days} day${days === 1 ? "" : "s"}`;
}
