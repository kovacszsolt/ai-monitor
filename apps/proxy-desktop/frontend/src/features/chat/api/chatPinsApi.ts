import * as DesktopService from "../../../../wailsjs/go/service/DesktopService";
import { unwrap, type ApiResult } from "@/lib/api";
import { service } from "../../../../wailsjs/go/models";

export async function listChatPinnedPrompts() {
  const result = (await DesktopService.ListChatPinnedPrompts()) as ApiResult<
    service.ChatPinnedPromptDTO[]
  >;
  return unwrap(result);
}

export async function createChatPinnedPrompt(body: string, title = "") {
  return unwrap(
    (await DesktopService.CreateChatPinnedPrompt(body, title)) as ApiResult<service.ChatPinnedPromptDTO>,
  );
}

export async function deleteChatPinnedPrompt(id: string) {
  return unwrap((await DesktopService.DeleteChatPinnedPrompt(id)) as ApiResult<boolean>);
}
