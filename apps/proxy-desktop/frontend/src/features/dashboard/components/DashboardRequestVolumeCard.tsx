import { useQuery } from "@tanstack/react-query";
import { useMemo } from "react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { CollapsibleCardHeader } from "@/shared/components/CollapsibleCardHeader";
import { DASHBOARD_SECTION_IDS } from "@/features/dashboard/dashboardSectionIds";
import { useDashboardSectionCollapsed } from "@/features/dashboard/hooks/useDashboardSectionCollapsed";
import { HourlyChartRangePicker } from "@/features/dashboard/components/HourlyChartRangePicker";
import {
  formatHourlyChartRangeLabel,
  hourlyRangeQueryKeyPart,
  hourlyRangeRequestParams,
  isCustomHourlyRangeValid,
  type HourlyChartRangeMode,
} from "@/features/dashboard/components/hourlyChartRange";
import { HourlyRequestsByModelChart } from "@/features/dashboard/components/HourlyRequestsByModelChart";
import {
  coerceBucketForRange,
  HOURLY_CHART_BUCKET_OPTIONS,
  hourlyChartBucketLabel,
  rangeSpanDays,
  usageVolumeSectionTitle,
  type HourlyChartBucketUnit,
} from "@/features/dashboard/components/hourlyChartBucket";
import { HOURLY_CHART_METRICS, type HourlyChartMetric } from "@/features/dashboard/components/hourlyChartMetric";
import { ModelUsageFilters } from "@/features/dashboard/components/ModelUsageByTagSection";
import { getHourlyRequestCountByModel, listAIRequestTagOptions } from "@/features/requests/api/requestsApi";
import { ErrorState, LoadingState } from "@/shared/components/PageStates";

type Props = {
  tag: string;
  onTagChange: (tag: string) => void;
  projectID: string;
  onProjectChange: (projectID: string) => void;
  metric: HourlyChartMetric;
  onMetricChange: (metric: HourlyChartMetric) => void;
  bucket: HourlyChartBucketUnit;
  onBucketChange: (bucket: HourlyChartBucketUnit) => void;
  rangeMode: HourlyChartRangeMode;
  onRangeModeChange: (mode: HourlyChartRangeMode) => void;
  customRangeStart: string;
  customRangeEnd: string;
  onCustomRangeStartChange: (start: string) => void;
  onCustomRangeEndChange: (end: string) => void;
};

export function DashboardRequestVolumeCard({
  tag,
  onTagChange,
  projectID,
  onProjectChange,
  metric,
  onMetricChange,
  bucket,
  onBucketChange,
  rangeMode,
  onRangeModeChange,
  customRangeStart,
  customRangeEnd,
  onCustomRangeStartChange,
  onCustomRangeEndChange,
}: Props) {
  const { collapsed, toggle } = useDashboardSectionCollapsed(DASHBOARD_SECTION_IDS.usageByHour);

  const tagOptionsQ = useQuery({ queryKey: ["request-tag-options"], queryFn: listAIRequestTagOptions });
  const tagEnabled = tagOptionsQ.data?.enabled ?? false;

  const rangeReady =
    rangeMode !== "custom" || isCustomHourlyRangeValid(customRangeStart, customRangeEnd);

  const requestParams = useMemo(
    () => hourlyRangeRequestParams(rangeMode, customRangeStart, customRangeEnd),
    [rangeMode, customRangeStart, customRangeEnd],
  );

  const spanDays = useMemo(
    () => rangeSpanDays(rangeMode, customRangeStart, customRangeEnd),
    [rangeMode, customRangeStart, customRangeEnd],
  );
  const effectiveBucket = useMemo(() => coerceBucketForRange(bucket, spanDays), [bucket, spanDays]);
  const bucketCoerced = effectiveBucket !== bucket;

  const hourlyQ = useQuery({
    queryKey: [
      "hourly-requests-by-model",
      tag,
      projectID,
      effectiveBucket,
      hourlyRangeQueryKeyPart(rangeMode, customRangeStart, customRangeEnd),
    ],
    queryFn: () =>
      getHourlyRequestCountByModel(
        tag,
        projectID,
        requestParams.days,
        requestParams.rangeStart,
        requestParams.rangeEnd,
        effectiveBucket,
      ),
    enabled: rangeReady && !collapsed,
    refetchInterval: 30000,
  });

  const rangeLabel = formatHourlyChartRangeLabel(
    rangeMode,
    customRangeStart,
    customRangeEnd,
    hourlyQ.data,
  );

  return (
    <Card className="w-full">
      <CardHeader className="pb-2">
        <CollapsibleCardHeader title={usageVolumeSectionTitle(bucket)} collapsed={collapsed} onToggle={toggle} />
      </CardHeader>
      {!collapsed && (
      <CardContent className="space-y-4">
        <div className="flex flex-wrap items-end gap-3">
          <ModelUsageFilters
            tag={tag}
            onTagChange={onTagChange}
            projectID={projectID}
            onProjectChange={onProjectChange}
            showTagFilter={tagEnabled}
          />
          <div className="flex flex-wrap gap-1 rounded-md border border-border p-0.5">
            {HOURLY_CHART_BUCKET_OPTIONS.map((opt) => (
              <Button
                key={opt.id}
                type="button"
                size="sm"
                variant={bucket === opt.id ? "toggle" : "ghost"}
                className="h-8 px-3 text-xs"
                onClick={() => onBucketChange(opt.id)}
              >
                {opt.label}
              </Button>
            ))}
          </div>
          <div className="flex flex-wrap gap-1 rounded-md border border-border p-0.5">
            {HOURLY_CHART_METRICS.map((opt) => (
              <Button
                key={opt.id}
                type="button"
                size="sm"
                variant={metric === opt.id ? "toggle" : "ghost"}
                className="h-8 px-3 text-xs"
                onClick={() => onMetricChange(opt.id)}
              >
                {opt.label}
              </Button>
            ))}
          </div>
          <div className="ms-auto flex flex-col items-end gap-2">
            <HourlyChartRangePicker
              mode={rangeMode}
              onModeChange={onRangeModeChange}
              customStart={customRangeStart}
              customEnd={customRangeEnd}
              onCustomStartChange={onCustomRangeStartChange}
              onCustomEndChange={onCustomRangeEndChange}
            />
          </div>
        </div>
        {!rangeReady && (
          <p className="text-sm text-muted-foreground">Choose a valid custom date range to load the chart.</p>
        )}
        {rangeReady && bucketCoerced && (
          <p className="text-sm text-muted-foreground">
            {hourlyChartBucketLabel(bucket)} is limited for this range; showing{" "}
            {hourlyChartBucketLabel(effectiveBucket).toLowerCase()} instead.
          </p>
        )}
        {rangeReady && hourlyQ.isLoading && <LoadingState label="Loading chart…" />}
        {rangeReady && hourlyQ.error && (
          <ErrorState message={(hourlyQ.error as Error).message} onRetry={() => hourlyQ.refetch()} />
        )}
        {rangeReady && !hourlyQ.isLoading && !hourlyQ.error && (
          <HourlyRequestsByModelChart series={hourlyQ.data} rangeLabel={rangeLabel} metric={metric} />
        )}
      </CardContent>
      )}
    </Card>
  );
}
