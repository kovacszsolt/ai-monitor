import * as DesktopService from "../../../../wailsjs/go/service/DesktopService";
import { unwrap, type ApiResult } from "@/lib/api";
import { service } from "../../../../wailsjs/go/models";

export async function listModelRoutes() {
  return unwrap((await DesktopService.ListModelRoutes()) as ApiResult<service.ModelRouteDTO[]>);
}

export async function createModelRoute(
  publicModelName: string,
  providerID: string,
  upstream: string,
  protocol: string,
  inputTokenUSDPer1M: string,
  cachedInputTokenUSDPer1M: string,
  outputTokenUSDPer1M: string,
) {
  return unwrap(
    (await DesktopService.CreateModelRoute(
      publicModelName,
      providerID,
      upstream,
      protocol,
      inputTokenUSDPer1M,
      cachedInputTokenUSDPer1M,
      outputTokenUSDPer1M,
    )) as ApiResult<service.ModelRouteDTO>,
  );
}

export async function updateModelRoute(
  id: string,
  upstream: string,
  protocol: string,
  inputTokenUSDPer1M: string,
  cachedInputTokenUSDPer1M: string,
  outputTokenUSDPer1M: string,
) {
  return unwrap(
    (await DesktopService.UpdateModelRoute(
      id,
      upstream,
      protocol,
      inputTokenUSDPer1M,
      cachedInputTokenUSDPer1M,
      outputTokenUSDPer1M,
    )) as ApiResult<service.ModelRouteDTO>,
  );
}

export async function setModelRouteEnabled(id: string, enabled: boolean) {
  return unwrap((await DesktopService.SetModelRouteEnabled(id, enabled)) as ApiResult<boolean>);
}
