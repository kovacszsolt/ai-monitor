import * as DesktopService from "../../../../wailsjs/go/service/DesktopService";
import { unwrap, type ApiResult } from "@/lib/api";
import { service } from "../../../../wailsjs/go/models";

export async function getDashboardStats() {
  return unwrap((await DesktopService.GetDashboardStats()) as ApiResult<service.DashboardStatsDTO>);
}

export const REQUEST_SORT_COLUMNS = [
  "started",
  "project",
  "model",
  "tag",
  "state",
  "input",
  "cached",
  "output",
  "ttft",
  "duration",
  "gen",
  "prompt",
  "eff",
  "cache",
] as const;

export type RequestSortColumn = (typeof REQUEST_SORT_COLUMNS)[number];

export type ListAIRequestsParams = {
  page: number;
  pageSize: number;
  projectID?: string;
  state?: string;
  tag?: string;
  model?: string;
  sortBy?: RequestSortColumn;
  sortDesc?: boolean;
};

export async function listAIRequests({
  page,
  pageSize,
  projectID = "",
  state = "",
  tag = "",
  model = "",
  sortBy = "started",
  sortDesc = true,
}: ListAIRequestsParams) {
  return unwrap(
    (await DesktopService.ListAIRequests(
      page,
      pageSize,
      projectID,
      state,
      tag,
      model,
      sortBy,
      sortDesc,
    )) as ApiResult<service.AIRequestListDTO>,
  );
}

export async function listAIRequestModelOptions() {
  return unwrap((await DesktopService.ListAIRequestModelOptions()) as ApiResult<string[]>);
}

export async function deleteAIRequest(id: string) {
  return unwrap((await DesktopService.DeleteAIRequest(id)) as ApiResult<boolean>);
}

export async function getAIRequestContent(id: string) {
  return unwrap(
    (await DesktopService.GetAIRequestContent(id)) as ApiResult<service.AIRequestContentDTO>,
  );
}

export async function bulkSetAIRequestTag(ids: string[], tag: string) {
  return unwrap(
    (await DesktopService.BulkSetAIRequestTag(ids, tag)) as ApiResult<service.BulkSetAIRequestTagResultDTO>,
  );
}

export async function bulkDeleteAIRequests(ids: string[]) {
  return unwrap(
    (await DesktopService.BulkDeleteAIRequests(ids)) as ApiResult<service.BulkDeleteAIRequestResultDTO>,
  );
}

export async function listAIRequestTagOptions() {
  return unwrap(
    (await DesktopService.ListAIRequestTagOptions()) as ApiResult<service.AIRequestTagOptionsDTO>,
  );
}

export async function getModelUsageStats(tag = "", projectID = "") {
  return unwrap(
    (await DesktopService.GetModelUsageStats(tag, projectID)) as ApiResult<service.ModelUsageStatsDTO[]>,
  );
}

export const HOURLY_REQUEST_CHART_DAYS = 7;

export async function getHourlyRequestCountByModel(
  tag = "",
  projectID = "",
  days = HOURLY_REQUEST_CHART_DAYS,
  rangeStart = "",
  rangeEnd = "",
) {
  return unwrap(
    (await DesktopService.GetHourlyRequestCountByModel(
      tag,
      projectID,
      days,
      rangeStart,
      rangeEnd,
    )) as ApiResult<service.HourlyModelRequestSeriesDTO>,
  );
}

export async function getModelTokenUsageOverTime(tag: string, days = 14) {
  return unwrap(
    (await DesktopService.GetModelTokenUsageOverTime(tag, days)) as ApiResult<service.ModelTokenTimeSeriesDTO>,
  );
}
