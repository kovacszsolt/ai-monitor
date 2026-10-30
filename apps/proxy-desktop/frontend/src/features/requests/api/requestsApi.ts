import * as DesktopService from "../../../../wailsjs/go/service/DesktopService";
import { unwrap, type ApiResult } from "@/lib/api";
import { service } from "../../../../wailsjs/go/models";

export async function getDashboardStats() {
  return unwrap((await DesktopService.GetDashboardStats()) as ApiResult<service.DashboardStatsDTO>);
}

export const REQUEST_SORT_COLUMNS = [
  "started",
  "project",
  "api_key",
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
  apiKeyID?: string;
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
  apiKeyID = "",
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
      apiKeyID,
      sortBy,
      sortDesc,
    )) as ApiResult<service.AIRequestListDTO>,
  );
}

export async function listAIRequestAPIKeyOptions(projectID = "") {
  return unwrap(
    (await DesktopService.ListAIRequestAPIKeyOptions(projectID)) as ApiResult<ModelRouteFilterOption[]>,
  );
}

export type ModelRouteFilterOption = {
  id: string;
  name: string;
};

export async function listAIRequestModelOptions() {
  return unwrap(
    (await DesktopService.ListAIRequestModelOptions()) as ApiResult<ModelRouteFilterOption[]>,
  );
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

export async function getModelUsageStats(tag = "", projectID = "", apiKeyID = "") {
  return unwrap(
    (await DesktopService.GetModelUsageStats(tag, projectID, apiKeyID)) as ApiResult<service.ModelUsageStatsDTO[]>,
  );
}

export const HOURLY_REQUEST_CHART_DAYS = 7;

export async function getHourlyRequestCountByModel(
  tag = "",
  projectID = "",
  apiKeyID = "",
  days = HOURLY_REQUEST_CHART_DAYS,
  rangeStart = "",
  rangeEnd = "",
  bucketUnit = "hour",
) {
  return unwrap(
    (await DesktopService.GetHourlyRequestCountByModel(
      tag,
      projectID,
      apiKeyID,
      days,
      rangeStart,
      rangeEnd,
      bucketUnit,
    )) as ApiResult<service.HourlyModelRequestSeriesDTO>,
  );
}

export async function getModelTokenUsageOverTime(tag: string, days = 14) {
  return unwrap(
    (await DesktopService.GetModelTokenUsageOverTime(tag, days)) as ApiResult<service.ModelTokenTimeSeriesDTO>,
  );
}
