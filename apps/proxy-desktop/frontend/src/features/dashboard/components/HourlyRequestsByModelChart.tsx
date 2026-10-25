import { useMemo } from "react";
import {
  Bar,
  BarChart,
  CartesianGrid,
  Legend,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from "recharts";
import {
  formatHourlyChartValue,
  hourlyPointValue,
  type HourlyChartMetric,
} from "@/features/dashboard/components/hourlyChartMetric";
import {
  parseTokenStackDataKey,
  TOKEN_STACK_SERIES,
  tokenStackDataKey,
  totalTokenCountFromPoint,
  type TokenStackKey,
} from "@/features/dashboard/components/tokenStackSeries";
import { formatNumber } from "@/lib/format";
import { service } from "../../../../wailsjs/go/models";

const CHART_TEXT = "#fafafa";
const CHART_AXIS = "#a1a1aa";
const CHART_GRID = "#3f3f46";
const OTHER_MODEL_KEY = "__other__";
const MAX_MODEL_SERIES = 8;

const MODEL_COLORS = [
  "#60a5fa",
  "#4ade80",
  "#c084fc",
  "#fbbf24",
  "#f87171",
  "#2dd4bf",
  "#fb923c",
  "#a78bfa",
  "#94a3b8",
] as const;

const AXIS_TICK = { fontSize: 11, fill: CHART_TEXT };

type ChartRow = Record<string, number | string> & {
  bucket: string;
  hourLabel: string;
};

type TokenModelStack = {
  stackKey: string;
  label: string;
};

function shortModelName(name: string, max = 24): string {
  if (name.length <= max) return name;
  return `${name.slice(0, max - 1)}…`;
}

function formatHourLabel(bucketISO: string): string {
  const d = new Date(bucketISO);
  if (Number.isNaN(d.getTime())) return bucketISO;
  return d.toLocaleString(undefined, {
    month: "short",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

function pivotHourlyTokensByModel(
  points: service.HourlyModelRequestCountDTO[],
  maxModels: number,
): { rows: ChartRow[]; models: TokenModelStack[] } {
  const totals = new Map<string, number>();
  for (const p of points) {
    totals.set(p.model, (totals.get(p.model) ?? 0) + totalTokenCountFromPoint(p));
  }
  const ranked = [...totals.entries()].sort((a, b) => b[1] - a[1]);
  const topModels = ranked.slice(0, maxModels).map(([name]) => name);
  const topSet = new Set(topModels);
  const hasOther = ranked.length > topModels.length;

  const models: TokenModelStack[] = topModels.map((name) => ({
    stackKey: name,
    label: shortModelName(name),
  }));
  if (hasOther) {
    models.push({ stackKey: OTHER_MODEL_KEY, label: "Other" });
  }

  const bucketOrder: string[] = [];
  const bucketRows = new Map<string, ChartRow>();

  const addTokens = (row: ChartRow, modelKey: string, segment: TokenStackKey, amount: number) => {
    if (amount <= 0) return;
    const field = tokenStackDataKey(modelKey, segment);
    row[field] = (Number(row[field] ?? 0) + amount) as number;
  };

  for (const p of points) {
    const tokens = totalTokenCountFromPoint(p);
    if (tokens <= 0) continue;

    if (!bucketRows.has(p.bucket_start)) {
      bucketOrder.push(p.bucket_start);
      bucketRows.set(p.bucket_start, {
        bucket: p.bucket_start,
        hourLabel: formatHourLabel(p.bucket_start),
      });
    }
    const row = bucketRows.get(p.bucket_start)!;
    const modelKey = topSet.has(p.model) ? p.model : OTHER_MODEL_KEY;
    addTokens(row, modelKey, "input_tokens", p.input_tokens ?? 0);
    addTokens(row, modelKey, "cached_input_tokens", p.cached_input_tokens ?? 0);
    addTokens(row, modelKey, "output_tokens", p.output_tokens ?? 0);
  }

  const rows = bucketOrder
    .map((b) => bucketRows.get(b)!)
    .filter((row) =>
      models.some((m) =>
        TOKEN_STACK_SERIES.some((seg) => Number(row[tokenStackDataKey(m.stackKey, seg.key)] ?? 0) > 0),
      ),
    );

  return { rows, models };
}

function pivotHourlyByMetric(
  points: service.HourlyModelRequestCountDTO[],
  maxModels: number,
  metric: Exclude<HourlyChartMetric, "tokens">,
): { rows: ChartRow[]; series: { key: string; label: string; color: string }[] } {
  const totals = new Map<string, number>();
  for (const p of points) {
    totals.set(p.model, (totals.get(p.model) ?? 0) + hourlyPointValue(p, metric));
  }
  const ranked = [...totals.entries()].sort((a, b) => b[1] - a[1]);
  const topModels = ranked.slice(0, maxModels).map(([name]) => name);
  const topSet = new Set(topModels);
  const hasOther = ranked.length > topModels.length;

  const bucketOrder: string[] = [];
  const bucketRows = new Map<string, ChartRow>();

  for (const p of points) {
    const value = hourlyPointValue(p, metric);
    if (!Number.isFinite(value) || value <= 0) {
      continue;
    }
    if (!bucketRows.has(p.bucket_start)) {
      bucketOrder.push(p.bucket_start);
      bucketRows.set(p.bucket_start, {
        bucket: p.bucket_start,
        hourLabel: formatHourLabel(p.bucket_start),
      });
    }
    const row = bucketRows.get(p.bucket_start)!;
    const seriesKey = topSet.has(p.model) ? p.model : OTHER_MODEL_KEY;
    row[seriesKey] = (Number(row[seriesKey] ?? 0) + value) as number;
  }

  const series: { key: string; label: string; color: string }[] = topModels.map((model, i) => ({
    key: model,
    label: shortModelName(model),
    color: MODEL_COLORS[i % MODEL_COLORS.length],
  }));
  if (hasOther) {
    series.push({
      key: OTHER_MODEL_KEY,
      label: "Other",
      color: MODEL_COLORS[MODEL_COLORS.length - 1],
    });
  }

  const rows = bucketOrder.map((b) => ({ ...bucketRows.get(b)! }));

  return { rows, series };
}

type HourlyTooltipPayload = {
  dataKey?: string | number;
  name?: string;
  value?: number | string;
  color?: string;
};

function HourlyChartTooltip({
  active,
  payload,
  label,
  metric,
}: {
  active?: boolean;
  payload?: readonly HourlyTooltipPayload[];
  label?: string;
  metric: Exclude<HourlyChartMetric, "tokens">;
}) {
  if (!active || !payload?.length) {
    return null;
  }
  const items = payload.filter((entry) => {
    const value = Number(entry.value);
    return Number.isFinite(value) && value > 0;
  });
  if (items.length === 0) {
    return null;
  }

  return (
    <div className="rounded-lg border border-border bg-card px-3 py-2 text-xs shadow-md">
      <p className="mb-1.5 font-semibold text-foreground">{label}</p>
      <ul className="space-y-1">
        {items.map((entry) => (
          <li key={String(entry.dataKey)} className="flex items-center justify-between gap-4">
            <span className="flex min-w-0 items-center gap-1.5">
              <span
                className="inline-block h-2 w-2 shrink-0 rounded-sm"
                style={{ backgroundColor: entry.color }}
              />
              <span className="truncate text-foreground">{entry.name}</span>
            </span>
            <span className="shrink-0 font-mono text-foreground">
              {formatHourlyChartValue(
                metric,
                typeof entry.value === "number" ? entry.value : Number(entry.value),
              )}
            </span>
          </li>
        ))}
      </ul>
    </div>
  );
}

function HourlyTokenStackByModelTooltip({
  active,
  payload,
  label,
  modelLabels,
}: {
  active?: boolean;
  payload?: readonly HourlyTooltipPayload[];
  label?: string;
  modelLabels: Map<string, string>;
}) {
  if (!active || !payload?.length) {
    return null;
  }

  const byModel = new Map<string, Partial<Record<TokenStackKey, number>>>();
  for (const entry of payload) {
    const value = Number(entry.value);
    if (!Number.isFinite(value) || value <= 0) continue;
    const parsed = parseTokenStackDataKey(String(entry.dataKey));
    if (!parsed) continue;
    const bucket = byModel.get(parsed.modelKey) ?? {};
    bucket[parsed.segment] = value;
    byModel.set(parsed.modelKey, bucket);
  }
  if (byModel.size === 0) {
    return null;
  }

  const modelOrder = [...byModel.keys()].sort((a, b) => {
    const totalA = TOKEN_STACK_SERIES.reduce((s, seg) => s + (byModel.get(a)?.[seg.key] ?? 0), 0);
    const totalB = TOKEN_STACK_SERIES.reduce((s, seg) => s + (byModel.get(b)?.[seg.key] ?? 0), 0);
    return totalB - totalA;
  });

  return (
    <div className="max-h-[min(50vh,16rem)] overflow-y-auto rounded-lg border border-border bg-card px-3 py-2 text-xs shadow-md">
      <p className="mb-1.5 font-semibold text-foreground">{label}</p>
      <ul className="space-y-2">
        {modelOrder.map((modelKey) => {
          const parts = byModel.get(modelKey)!;
          const modelTotal = TOKEN_STACK_SERIES.reduce((s, seg) => s + (parts[seg.key] ?? 0), 0);
          const display =
            modelKey === OTHER_MODEL_KEY ? "Other" : (modelLabels.get(modelKey) ?? modelKey);
          return (
            <li key={modelKey}>
              <p className="mb-0.5 font-medium text-foreground">{display}</p>
              <ul className="space-y-0.5 pl-0.5">
                {TOKEN_STACK_SERIES.map((seg) => {
                  const v = parts[seg.key];
                  if (v == null || v <= 0) return null;
                  return (
                    <li key={seg.key} className="flex items-center justify-between gap-4">
                      <span className="flex min-w-0 items-center gap-1.5 text-muted-foreground">
                        <span
                          className="inline-block h-2 w-2 shrink-0 rounded-sm"
                          style={{ backgroundColor: seg.fill }}
                        />
                        {seg.label}
                      </span>
                      <span className="shrink-0 font-mono text-foreground">{formatNumber(v)}</span>
                    </li>
                  );
                })}
              </ul>
              <p className="mt-0.5 text-muted-foreground">
                Total <span className="font-mono text-foreground">{formatNumber(modelTotal)}</span>
              </p>
            </li>
          );
        })}
      </ul>
    </div>
  );
}

function TokenStackLegend() {
  return (
    <ul className="flex flex-wrap justify-center gap-x-4 gap-y-1 pt-2 text-xs text-foreground">
      {TOKEN_STACK_SERIES.map((seg) => (
        <li key={seg.key} className="flex items-center gap-1.5">
          <span className="inline-block h-2 w-2 rounded-sm" style={{ backgroundColor: seg.fill }} />
          {seg.label}
        </li>
      ))}
    </ul>
  );
}

export function HourlyRequestsByModelChart({
  series,
  rangeLabel,
  metric,
}: {
  series: service.HourlyModelRequestSeriesDTO | undefined;
  rangeLabel: string;
  metric: HourlyChartMetric;
}) {
  const tokenChart = useMemo(
    () =>
      metric === "tokens"
        ? pivotHourlyTokensByModel(series?.points ?? [], MAX_MODEL_SERIES)
        : null,
    [series?.points, metric],
  );

  const tokenModelLabels = useMemo(() => {
    const map = new Map<string, string>();
    if (!tokenChart) return map;
    for (const m of tokenChart.models) {
      map.set(m.stackKey, m.stackKey === OTHER_MODEL_KEY ? "Other" : m.stackKey);
    }
    return map;
  }, [tokenChart]);

  const modelChart = useMemo(
    () =>
      metric === "tokens"
        ? null
        : pivotHourlyByMetric(series?.points ?? [], MAX_MODEL_SERIES, metric),
    [series?.points, metric],
  );

  if (!series || series.points.length === 0) {
    return (
      <p className="text-sm text-muted-foreground">No requests for {rangeLabel} with this filter.</p>
    );
  }

  if (metric === "tokens") {
    const tokenRows = tokenChart?.rows ?? [];
    const tokenModels = tokenChart?.models ?? [];

    if (tokenRows.length === 0) {
      return (
        <p className="text-sm text-muted-foreground">
          No token usage for {rangeLabel} with this filter.
        </p>
      );
    }

    return (
      <div className="recharts-chart-readable h-[min(50vh,22rem)] w-full min-h-[16rem]">
        <ResponsiveContainer width="100%" height="100%">
          <BarChart
            data={tokenRows}
            margin={{ top: 8, right: 8, left: 0, bottom: 4 }}
            barCategoryGap="12%"
            barGap={2}
          >
            <CartesianGrid stroke={CHART_GRID} strokeDasharray="3 3" strokeOpacity={0.8} vertical={false} />
            <XAxis
              dataKey="hourLabel"
              tick={AXIS_TICK}
              stroke={CHART_GRID}
              tickLine={{ stroke: CHART_AXIS }}
              axisLine={{ stroke: CHART_GRID }}
              interval="preserveStartEnd"
              minTickGap={24}
            />
            <YAxis
              allowDecimals={false}
              tick={AXIS_TICK}
              stroke={CHART_GRID}
              tickLine={{ stroke: CHART_AXIS }}
              axisLine={{ stroke: CHART_GRID }}
              width={48}
              tickFormatter={(v) => formatNumber(Number(v))}
            />
            <Tooltip
              content={(props) => (
                <HourlyTokenStackByModelTooltip
                  active={props.active}
                  label={typeof props.label === "string" ? props.label : undefined}
                  payload={props.payload as unknown as readonly HourlyTooltipPayload[] | undefined}
                  modelLabels={tokenModelLabels}
                />
              )}
              cursor={{ fill: "rgba(250, 250, 250, 0.06)" }}
            />
            <Legend content={<TokenStackLegend />} />
            {tokenModels.flatMap((model) =>
              TOKEN_STACK_SERIES.map((seg) => (
                <Bar
                  key={`${model.stackKey}-${seg.key}`}
                  dataKey={tokenStackDataKey(model.stackKey, seg.key)}
                  name={seg.label}
                  stackId={model.stackKey}
                  fill={seg.fill}
                  legendType="none"
                  radius={[0, 0, 0, 0]}
                />
              )),
            )}
          </BarChart>
        </ResponsiveContainer>
      </div>
    );
  }

  const rows = modelChart!.rows;
  const modelSeries = modelChart!.series;
  const hasAnyValue = rows.length > 0;

  if (!hasAnyValue) {
    const emptyHint =
      metric === "cost"
        ? "No priced usage"
        : metric === "generation_tps"
          ? "No decode timing for throughput"
          : "No usage";
    return (
      <p className="text-sm text-muted-foreground">
        {emptyHint} for {rangeLabel} with this filter.
      </p>
    );
  }

  const stackBars = metric !== "generation_tps";
  const yAxisWidth = metric === "cost" ? 52 : 44;

  return (
    <div className="recharts-chart-readable h-[min(50vh,22rem)] w-full min-h-[16rem]">
      <ResponsiveContainer width="100%" height="100%">
        <BarChart data={rows} margin={{ top: 8, right: 8, left: 0, bottom: 4 }}>
          <CartesianGrid stroke={CHART_GRID} strokeDasharray="3 3" strokeOpacity={0.8} vertical={false} />
          <XAxis
            dataKey="hourLabel"
            tick={AXIS_TICK}
            stroke={CHART_GRID}
            tickLine={{ stroke: CHART_AXIS }}
            axisLine={{ stroke: CHART_GRID }}
            interval="preserveStartEnd"
            minTickGap={24}
          />
          <YAxis
            allowDecimals={metric !== "requests"}
            tick={AXIS_TICK}
            stroke={CHART_GRID}
            tickLine={{ stroke: CHART_AXIS }}
            axisLine={{ stroke: CHART_GRID }}
            width={yAxisWidth}
            tickFormatter={(v) => formatHourlyChartValue(metric, Number(v))}
          />
          <Tooltip
            content={(props) => (
              <HourlyChartTooltip
                active={props.active}
                label={typeof props.label === "string" ? props.label : undefined}
                payload={props.payload as unknown as readonly HourlyTooltipPayload[] | undefined}
                metric={metric}
              />
            )}
            cursor={{ fill: "rgba(250, 250, 250, 0.06)" }}
          />
          <Legend wrapperStyle={{ fontSize: 12, color: CHART_TEXT }} />
          {modelSeries.map((s) => (
            <Bar
              key={s.key}
              dataKey={s.key}
              name={s.key === OTHER_MODEL_KEY ? "Other" : s.key}
              stackId={stackBars ? "volume" : undefined}
              fill={s.color}
              radius={[0, 0, 0, 0]}
            />
          ))}
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
}
