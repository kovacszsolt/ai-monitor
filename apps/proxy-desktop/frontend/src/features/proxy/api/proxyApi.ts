import * as DesktopService from "../../../../wailsjs/go/service/DesktopService";
import * as ProxyRuntimeService from "../../../../wailsjs/go/service/ProxyRuntimeService";
import { SaveExportFile } from "../../../../wailsjs/go/main/App";
import { unwrap, type ApiResult } from "@/lib/api";
import { service } from "../../../../wailsjs/go/models";

export async function validateProxyConfig() {
  return unwrap(
    (await DesktopService.ValidateProxyConfig()) as ApiResult<service.ProxyConfigSummaryDTO>,
  );
}

export async function getStorageStats() {
  return unwrap((await DesktopService.GetStorageStats()) as ApiResult<service.StorageStatsDTO>);
}

export async function runDatabaseMigrations() {
  return unwrap((await DesktopService.RunDatabaseMigrations()) as ApiResult<boolean>);
}

export async function exportMetadata(format: string, since: string, includeContent: boolean) {
  return unwrap(
    (await DesktopService.ExportMetadata(format, since, includeContent)) as ApiResult<service.ExportPayloadDTO>,
  );
}

export async function saveMetadataExport(payload: service.ExportPayloadDTO) {
  const path = await SaveExportFile(payload.filename, payload.content_base64);
  return path;
}

export async function getProxyStatus() {
  return unwrap((await ProxyRuntimeService.GetProxyStatus()) as ApiResult<service.ProxyStatusDTO>);
}

export async function startProxy() {
  return unwrap(
    (await ProxyRuntimeService.StartProxy()) as ApiResult<service.StartProxyResultDTO>,
  );
}

export async function checkProviders() {
  return unwrap(
    (await ProxyRuntimeService.CheckProviders()) as ApiResult<service.ProviderHealthCheckDTO>,
  );
}

export async function stopProxy() {
  return unwrap((await ProxyRuntimeService.StopProxy()) as ApiResult<boolean>);
}
