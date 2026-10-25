import { Button } from "@/components/ui/button";
import {
  HOURLY_CHART_RANGE_PRESETS,
  type HourlyChartRangeMode,
  isCustomHourlyRangeValid,
} from "@/features/dashboard/components/hourlyChartRange";

const dateInputClass =
  "h-8 rounded-md border border-border bg-card px-2 text-xs text-foreground";

type Props = {
  mode: HourlyChartRangeMode;
  onModeChange: (mode: HourlyChartRangeMode) => void;
  customStart: string;
  customEnd: string;
  onCustomStartChange: (value: string) => void;
  onCustomEndChange: (value: string) => void;
};

export function HourlyChartRangePicker({
  mode,
  onModeChange,
  customStart,
  customEnd,
  onCustomStartChange,
  onCustomEndChange,
}: Props) {
  const customInvalid = mode === "custom" && !isCustomHourlyRangeValid(customStart, customEnd);

  return (
    <div className="flex flex-col items-start gap-2 sm:items-end">
      <div className="flex flex-wrap gap-1 rounded-md border border-border p-0.5">
        {HOURLY_CHART_RANGE_PRESETS.map((preset) => (
          <Button
            key={preset.days}
            type="button"
            size="sm"
            variant={mode === preset.days ? "default" : "ghost"}
            className="h-8 px-3 text-xs"
            onClick={() => onModeChange(preset.days)}
          >
            {preset.label}
          </Button>
        ))}
        <Button
          type="button"
          size="sm"
          variant={mode === "custom" ? "default" : "ghost"}
          className="h-8 px-3 text-xs"
          onClick={() => onModeChange("custom")}
        >
          Custom
        </Button>
      </div>
      {mode === "custom" && (
        <div className="flex flex-wrap items-center justify-end gap-2">
          <label className="flex items-center gap-1.5 text-xs text-muted-foreground">
            From
            <input
              type="date"
              className={dateInputClass}
              value={customStart}
              onChange={(e) => onCustomStartChange(e.target.value)}
            />
          </label>
          <label className="flex items-center gap-1.5 text-xs text-muted-foreground">
            To
            <input
              type="date"
              className={dateInputClass}
              value={customEnd}
              onChange={(e) => onCustomEndChange(e.target.value)}
            />
          </label>
          {customInvalid && (
            <span className="text-xs text-destructive">Pick a valid UTC date range (max 90 days).</span>
          )}
        </div>
      )}
    </div>
  );
}
