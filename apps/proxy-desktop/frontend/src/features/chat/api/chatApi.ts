import { unwrap, type ApiResult } from "@/lib/api";
import type { ChatMessage } from "@/features/chat/types";
import * as DesktopService from "../../../../wailsjs/go/service/DesktopService";
import { EventsOn } from "../../../../wailsjs/runtime/runtime";
import { service } from "../../../../wailsjs/go/models";

export type SendChatParams = {
  apiKeyId: string;
  modelRouteId: string;
  messages: ChatMessage[];
  stream: boolean;
  streamClientId?: string;
  onStreamDelta?: (text: string) => void;
};

export type SendChatResult = {
  content: string;
};

export async function sendChatCompletion(params: SendChatParams): Promise<SendChatResult> {
  const dto = new service.ChatCompletionInputDTO({
    api_key_id: params.apiKeyId,
    model_route_id: params.modelRouteId,
    messages: params.messages.map(
      (m) => new service.ChatMessageDTO({ role: m.role, content: m.content }),
    ),
    stream: params.stream,
    stream_client_id: params.streamClientId ?? "",
  });

  let unsubscribe: (() => void) | undefined;
  if (params.stream && params.streamClientId && params.onStreamDelta) {
    const clientId = params.streamClientId;
    unsubscribe = EventsOn("chat:stream-delta", (payload: { client_id?: string; text?: string }) => {
      if (payload?.client_id === clientId) {
        params.onStreamDelta?.(payload.text ?? "");
      }
    });
  }

  try {
    const result = unwrap(
      (await DesktopService.SendChatCompletion(dto)) as ApiResult<service.ChatCompletionResultDTO>,
    );
    return { content: result.content };
  } catch (err) {
    console.error("Chat completion failed", err);
    throw err;
  } finally {
    unsubscribe?.();
  }
}
