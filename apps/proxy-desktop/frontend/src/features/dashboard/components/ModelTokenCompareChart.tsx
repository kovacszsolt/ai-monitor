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
import { formatNumber } from "@/lib/format";
import { service } from "../../../../wailsjs/go/models";

const TOKEN_SERIES = [
  { key: "input_tokens", label: "Input", fill: "#60a5fa" },
  { key: "cached_input_tokens", label: "Cached", fill: "#c084fc" },
  { key: "output_tokens", label: "Output", fill: "#4ade80" },
] as const;

/** SVG does not resolve Tailwind theme vars — use fixed colors for dark UI. */
const CHART_TEXT = "#fafafa";
const CHART_AXIS = "#a1a1aa";
const CHART_GRID = "#3f3f46";

const AXIS_TICK = { fontSize: 11, fill: CHART_TEXT };
const AXIS_TICK_Y = { fontSize: 11, fill: CHART_TEXT, fontWeight: 500 };

const TOOLTIP_STYLE = {
  backgroundColor: "#18181b",
  border: "1px solid #3f3f46",
  borderRadius: "8px",
  color: CHART_TEXT,
};

type ChartRow = {
  model: string;
  modelShort: string;
  input_tokens: number;
  cached_input_tokens: number;
  output_tokens: number;
};

function shortModelName(name: string, max = 28): string {
  if (name.length <= max) return name;
  return `${name.slice(0, max - 1)}…`;
}

function toChartRows(rows: service.ModelUsageStatsDTO[], maxModels: number): ChartRow[] {
  const sorted = [...rows].sort(
    (a, b) =>
      b.input_tokens + b.cached_input_tokens + b.output_tokens -
      (a.input_tokens + a.cached_input_tokens + a.output_tokens),
  );
  return sorted.slice(0, maxModels).map((r) => ({
    model: r.name,
    modelShort: shortModelName(r.name),
    input_tokens: r.input_tokens,
    cached_input_tokens: r.cached_input_tokens,
    output_tokens: r.output_tokens,
  }));
}

export function ModelTokenCompareChart({
  tagLabel,
  rows,
  hideTitle = false,
}: {
  tagLabel: string;
  rows: service.ModelUsageStatsDTO[];
  hideTitle?: boolean;
}) {
  const data = useMemo(() => toChartRows(rows, 12), [rows]);
  const chartHeight = Math.max(280, data.length * 36 + 80);

  if (rows.length === 0) {
    return (
      <p className="text-sm text-muted-foreground">No token data for the selected filter.</p>
    );
  }

  return (
    <div className="space-y-4">
      {!hideTitle && (
        <div>
          <h3 className="text-base font-semibold">Token mix by model</h3>
          <p className="text-sm text-muted-foreground">
            Compare input, cached input, and output tokens per model for tag «{tagLabel}».
          </p>
        </div>
      )}
      {hideTitle && (
        <p className="text-sm text-muted-foreground">
          Input, cached input, and output per model for tag «{tagLabel}».
        </p>
      )}
      <div
        className="recharts-chart-readable h-[min(70vh,28rem)] w-full min-h-[17.5rem]"
        style={{ height: chartHeight }}
      >
        <ResponsiveContainer width="100%" height="100%">
          <BarChart
            data={data}
            layout="vertical"
            margin={{ top: 4, right: 12, left: 4, bottom: 4 }}
            barCategoryGap="18%"
            barGap={2}
          >
            <CartesianGrid stroke={CHART_GRID} strokeDasharray="3 3" strokeOpacity={0.8} horizontal={false} />
            <XAxis
              type="number"
              tick={AXIS_TICK}
              stroke={CHART_GRID}
              tickLine={{ stroke: CHART_AXIS }}
              axisLine={{ stroke: CHART_GRID }}
              tickFormatter={(v) => formatNumber(v)}
            />
            <YAxis
              type="category"
              dataKey="modelShort"
              tick={AXIS_TICK_Y}
              stroke={CHART_GRID}
              tickLine={{ stroke: CHART_AXIS }}
              axisLine={{ stroke: CHART_GRID }}
              width={148}
            />
            <Tooltip
              contentStyle={TOOLTIP_STYLE}
              labelStyle={{ color: CHART_TEXT, fontWeight: 600 }}
              itemStyle={{ color: CHART_TEXT }}
              cursor={{ fill: "rgba(250, 250, 250, 0.06)" }}
              labelFormatter={(_, payload) => {
                const row = payload?.[0]?.payload as ChartRow | undefined;
                return row?.model ?? "";
              }}
              formatter={(value, name) => [
                formatNumber(typeof value === "number" ? value : Number(value)),
                String(name),
              ]}
            />
            <Legend wrapperStyle={{ fontSize: 12, color: CHART_TEXT }} />
            {TOKEN_SERIES.map((s) => (
              <Bar key={s.key} dataKey={s.key} name={s.label} fill={s.fill} radius={[0, 3, 3, 0]} />
            ))}
          </BarChart>
        </ResponsiveContainer>
      </div>
    </div>
  );
}
