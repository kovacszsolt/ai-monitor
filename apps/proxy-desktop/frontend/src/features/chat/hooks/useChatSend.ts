import { useCallback, useState } from "react";
import { sendChatCompletion } from "@/features/chat/api/chatApi";
import type { ChatImageAttachment } from "@/features/chat/chatImageAttachment";
import type { ChatMessage } from "@/features/chat/types";

export type ChatSendConfig = {
  apiKeyId: string | undefined;
  modelRouteId: string | undefined;
  stream: boolean;
};

export function useChatSend(config: ChatSendConfig) {
  const [messages, setMessages] = useState<ChatMessage[]>([]);
  const [sending, setSending] = useState(false);
  const [streamingContent, setStreamingContent] = useState<string | null>(null);

  const clearChat = useCallback(() => {
    setMessages([]);
    setStreamingContent(null);
  }, []);

  const sendUserMessage = useCallback(
    async (text: string, images: ChatImageAttachment[] = []) => {
      const trimmed = text.trim();
      const hasImages = images.length > 0;
      if ((!trimmed && !hasImages) || sending) return;

      const { apiKeyId, modelRouteId, stream } = config;
      if (!apiKeyId || !modelRouteId) {
        throw new Error("Configure project, model, and API key before sending");
      }

      const userMsg: ChatMessage = {
        id: crypto.randomUUID(),
        role: "user",
        content: trimmed,
        images: hasImages ? images : undefined,
      };
      const history = [...messages, userMsg];
      setMessages(history);
      setSending(true);
      setStreamingContent(stream ? "" : null);

      const streamClientId = stream ? crypto.randomUUID() : undefined;
      try {
        const result = await sendChatCompletion({
          apiKeyId,
          modelRouteId,
          messages: history,
          stream,
          streamClientId,
          onStreamDelta: (accumulated) => setStreamingContent(accumulated),
        });
        const assistantMsg: ChatMessage = {
          id: crypto.randomUUID(),
          role: "assistant",
          content: result.content,
          images: result.images,
        };
        setMessages((prev) => [...prev, assistantMsg]);
      } finally {
        setStreamingContent(null);
        setSending(false);
      }
    },
    [config, messages, sending],
  );

  return {
    messages,
    sending,
    streamingContent,
    sendUserMessage,
    clearChat,
  };
}
