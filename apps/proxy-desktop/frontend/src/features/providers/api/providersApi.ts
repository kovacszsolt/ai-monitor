import * as DesktopService from "../../../../wailsjs/go/service/DesktopService";
import { unwrap, type ApiResult } from "@/lib/api";
import { service } from "../../../../wailsjs/go/models";

export async function listProviders() {
  return unwrap((await DesktopService.ListProviders()) as ApiResult<service.ProviderDTO[]>);
}

export async function createProvider(
  name: string,
  providerType: string,
  baseURL: string,
  credential: string,
  connectionTimeoutMS: number,
  requestTimeoutMS: number,
) {
  return unwrap(
    (await DesktopService.CreateProvider(
      name,
      providerType,
      baseURL,
      credential,
      connectionTimeoutMS,
      requestTimeoutMS,
    )) as ApiResult<service.ProviderDTO>,
  );
}

export async function updateProvider(
  id: string,
  baseURL: string,
  credential: string,
  connectionTimeoutMS: number,
  requestTimeoutMS: number,
) {
  return unwrap(
    (await DesktopService.UpdateProvider(id, baseURL, credential, connectionTimeoutMS, requestTimeoutMS)) as ApiResult<
      service.ProviderDTO
    >,
  );
}

export async function setProviderEnabled(id: string, enabled: boolean) {
  return unwrap((await DesktopService.SetProviderEnabled(id, enabled)) as ApiResult<boolean>);
}
