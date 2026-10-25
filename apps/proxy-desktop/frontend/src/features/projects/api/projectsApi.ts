import * as DesktopService from "../../../../wailsjs/go/service/DesktopService";
import { unwrap, type ApiResult } from "@/lib/api";
import { service } from "../../../../wailsjs/go/models";

export async function listProjects() {
  return unwrap((await DesktopService.ListProjects()) as ApiResult<service.ProjectDTO[]>);
}

export async function createProject(name: string, description: string) {
  return unwrap((await DesktopService.CreateProject(name, description)) as ApiResult<service.ProjectDTO>);
}

export async function updateProject(id: string, name: string, description: string) {
  return unwrap((await DesktopService.UpdateProject(id, name, description)) as ApiResult<service.ProjectDTO>);
}

export async function setProjectEnabled(id: string, enabled: boolean) {
  return unwrap((await DesktopService.SetProjectEnabled(id, enabled)) as ApiResult<boolean>);
}

export async function getProject(id: string) {
  return unwrap((await DesktopService.GetProject(id)) as ApiResult<service.ProjectDTO>);
}
