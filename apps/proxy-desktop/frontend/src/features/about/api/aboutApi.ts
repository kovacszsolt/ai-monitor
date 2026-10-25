import { GetAbout } from "../../../../wailsjs/go/service/DesktopService";
import { service } from "../../../../wailsjs/go/models";

export async function getAbout(): Promise<service.AboutDTO> {
  return GetAbout();
}
