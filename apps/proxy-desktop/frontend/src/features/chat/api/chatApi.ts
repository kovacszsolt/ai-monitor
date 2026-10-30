import { unwrap, type ApiResult } from "@/lib/api";
import type { ChatImageAttachment } from "@/features/chat/chatImageAttachment";
import type { ChatMessage } from "@/features/chat/types";
import * as DesktopService from "../../../../wailsjs/go/service/DesktopService";
import { EventsOn } from "../../../../wailsjs/runtime/runtime";
import { service } from "../../../../wailsjs/go/models";

function toChatMessageDTO(m: ChatMessage): service.ChatMessageDTO {
  const dto = new service.ChatMessageDTO({
    role: m.role,
    content: m.content,
  });
  if (m.images?.length) {
    dto.images = m.images.map(
      (img) =>
        new service.ChatImageAttachmentDTO({
          media_type: img.mediaType,
          data: img.dataBase64,
        }),
    );
  }
  return dto;
}

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
  images?: ChatImageAttachment[];
};

export async function sendChatCompletion(params: SendChatParams): Promise<SendChatResult> {
  const dto = new service.ChatCompletionInputDTO({
    api_key_id: params.apiKeyId,
    model_route_id: params.modelRouteId,
    messages: params.messages.map(toChatMessageDTO),
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
      { context: "SendChatCompletion" },
    );
    const images = (result.images ?? [])
      .filter((img) => !!img?.data)
      .map((img) => ({
        mediaType: img.media_type || "image/png",
        dataBase64: img.data,
      }));
    return {
      content: result.content ?? "",
      images: images.length > 0 ? images : undefined,
    };
  } finally {
    unsubscribe?.();
  }
}
