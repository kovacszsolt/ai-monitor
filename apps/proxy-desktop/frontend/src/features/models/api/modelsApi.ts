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
  modality: string,
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
      modality,
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
  modality: string,
  inputTokenUSDPer1M: string,
  cachedInputTokenUSDPer1M: string,
  outputTokenUSDPer1M: string,
) {
  return unwrap(
    (await DesktopService.UpdateModelRoute(
      id,
      upstream,
      protocol,
      modality,
      inputTokenUSDPer1M,
      cachedInputTokenUSDPer1M,
      outputTokenUSDPer1M,
    )) as ApiResult<service.ModelRouteDTO>,
  );
}

export async function setModelRouteEnabled(id: string, enabled: boolean) {
  return unwrap((await DesktopService.SetModelRouteEnabled(id, enabled)) as ApiResult<boolean>);
}

export type UpstreamModelOption = {
  id: string;
  alreadyImported: boolean;
};

export async function listProviderUpstreamModels(providerID: string) {
  const rows = unwrap(
    (await DesktopService.ListProviderUpstreamModels(providerID)) as ApiResult<
      service.UpstreamModelDTO[]
    >,
  );
  return rows.map(
    (m): UpstreamModelOption => ({
      id: m.id,
      alreadyImported: m.already_imported,
    }),
  );
}

export type ImportModelRoutesResult = {
  created: number;
  skipped: string[];
};

export async function importModelRoutes(
  providerID: string,
  upstreamIDs: string[],
  publicNamePrefix: string,
): Promise<ImportModelRoutesResult> {
  const r = unwrap(
    (await DesktopService.ImportModelRoutes(
      providerID,
      upstreamIDs,
      publicNamePrefix,
    )) as ApiResult<service.ImportModelRoutesResultDTO>,
  );
  return {
    created: r.created,
    skipped: r.skipped ?? [],
  };
}
