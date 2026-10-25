import * as DesktopService from "../../../../wailsjs/go/service/DesktopService";
import { unwrap, type ApiResult } from "@/lib/api";
import { service } from "../../../../wailsjs/go/models";

export async function listAPIKeys(projectID: string) {
  const result = (await DesktopService.ListAPIKeys(projectID)) as ApiResult<service.APIKeyDTO[]>;
  if (!result.ok) {
    throw new Error(result.error?.message ?? "Failed to load API keys");
  }
  return result.data ?? [];
}

export async function createAPIKey(projectID: string, name: string, allowAllModels: boolean) {
  return unwrap(
    (await DesktopService.CreateAPIKey(projectID, name, allowAllModels)) as ApiResult<service.CreateAPIKeyResultDTO>,
  );
}

export async function revokeAPIKey(id: string) {
  return unwrap((await DesktopService.RevokeAPIKey(id)) as ApiResult<boolean>);
}

export async function deleteAPIKey(id: string) {
  const result = (await DesktopService.DeleteAPIKey(id)) as ApiResult<boolean>;
  if (!result.ok) {
    throw new Error(result.error?.message ?? "Delete failed");
  }
  return true;
}

export async function setAPIKeyAllowAll(id: string, allowAll: boolean) {
  return unwrap((await DesktopService.SetAPIKeyAllowAllModels(id, allowAll)) as ApiResult<boolean>);
}
