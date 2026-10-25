import * as DesktopService from "../../../../wailsjs/go/service/DesktopService";
import { PickDatabaseFile } from "../../../../wailsjs/go/main/App";
import { unwrap, type ApiResult } from "@/lib/api";
import { service } from "../../../../wailsjs/go/models";

export async function getSettings() {
  return unwrap((await DesktopService.GetSettings()) as ApiResult<service.SettingsDTO>);
}

export async function updateDatabasePath(databasePath: string) {
  return unwrap(
    (await DesktopService.UpdateDatabasePath(databasePath)) as ApiResult<service.ConnectionInfoDTO>,
  );
}

export async function getProxyVariables() {
  return unwrap(
    (await DesktopService.GetProxyVariables()) as ApiResult<service.ProxyVariablesDTO>,
  );
}

export async function updateProxyVariables(payload: service.UpdateProxyVariablesDTO) {
  return unwrap(
    (await DesktopService.UpdateProxyVariables(payload)) as ApiResult<service.ProxyVariablesDTO>,
  );
}

export async function pickDatabaseFile() {
  return PickDatabaseFile();
}

export async function pingDatabase() {
  return unwrap((await DesktopService.PingDatabase()) as ApiResult<boolean>);
}
