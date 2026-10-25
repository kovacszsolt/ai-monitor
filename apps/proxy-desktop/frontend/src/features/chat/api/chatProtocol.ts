import type { ChatMessage } from "@/features/chat/types";

export type ChatProtocol = "openai_chat" | "openai_responses" | "anthropic_messages";

export function chatEndpoint(protocol: string): string {
  switch (protocol) {
    case "anthropic_messages":
      return "/v1/messages";
    case "openai_responses":
      return "/v1/responses";
    default:
      return "/v1/chat/completions";
  }
}

export function buildRequestBody(
  protocol: string,
  publicModelName: string,
  messages: ChatMessage[],
  stream: boolean,
): Record<string, unknown> {
  switch (protocol) {
    case "anthropic_messages":
      return buildAnthropicBody(publicModelName, messages, stream);
    case "openai_responses":
      return buildResponsesBody(publicModelName, messages, stream);
    default:
      return {
        model: publicModelName,
        messages: messages.map((m) => ({ role: m.role, content: m.content })),
        stream,
      };
  }
}

function buildAnthropicBody(
  publicModelName: string,
  messages: ChatMessage[],
  stream: boolean,
): Record<string, unknown> {
  const systemParts: string[] = [];
  const convo: { role: string; content: string }[] = [];
  for (const m of messages) {
    if (m.role === "system") {
      systemParts.push(m.content);
      continue;
    }
    if (m.role === "user" || m.role === "assistant") {
      convo.push({ role: m.role, content: m.content });
    }
  }
  const body: Record<string, unknown> = {
    model: publicModelName,
    max_tokens: 4096,
    messages: convo,
    stream,
  };
  if (systemParts.length > 0) {
    body.system = systemParts.join("\n\n");
  }
  return body;
}

function buildResponsesBody(
  publicModelName: string,
  messages: ChatMessage[],
  stream: boolean,
): Record<string, unknown> {
  let instructions = "";
  const inputItems: { role: string; content: string }[] = [];
  for (const m of messages) {
    if (m.role === "system") {
      instructions = instructions ? `${instructions}\n\n${m.content}` : m.content;
      continue;
    }
    if (m.role === "user" || m.role === "assistant") {
      inputItems.push({ role: m.role, content: m.content });
    }
  }
  const body: Record<string, unknown> = {
    model: publicModelName,
    stream,
  };
  if (instructions) {
    body.instructions = instructions;
  }
  if (inputItems.length === 1 && inputItems[0].role === "user") {
    body.input = inputItems[0].content;
  } else if (inputItems.length > 0) {
    body.input = inputItems;
  } else {
    body.input = "";
  }
  return body;
}

export function extractAssistantText(protocol: string, body: unknown): string {
  if (!body || typeof body !== "object") return "";
  const v = body as Record<string, unknown>;

  switch (protocol) {
    case "anthropic_messages":
      return extractAnthropicText(v);
    case "openai_responses":
      return extractResponsesText(v);
    default:
      return extractOpenAIChatText(v);
  }
}

function extractOpenAIChatText(v: Record<string, unknown>): string {
  const choices = v.choices;
  if (!Array.isArray(choices) || choices.length === 0) return "";
  const first = choices[0] as Record<string, unknown>;
  const message = first.message as Record<string, unknown> | undefined;
  if (message && typeof message.content === "string") return message.content;
  if (typeof first.text === "string") return first.text;
  return "";
}

function extractAnthropicText(v: Record<string, unknown>): string {
  const content = v.content;
  if (!Array.isArray(content)) return "";
  const parts: string[] = [];
  for (const block of content) {
    if (!block || typeof block !== "object") continue;
    const b = block as Record<string, unknown>;
    if (b.type === "text" && typeof b.text === "string") {
      parts.push(b.text);
    }
  }
  return parts.join("");
}

function extractResponsesText(v: Record<string, unknown>): string {
  if (typeof v.output_text === "string" && v.output_text) {
    return v.output_text;
  }
  const output = v.output;
  if (!Array.isArray(output)) return "";
  const parts: string[] = [];
  for (const item of output) {
    if (!item || typeof item !== "object") continue;
    const m = item as Record<string, unknown>;
    if (m.type !== "message") continue;
    const text = responsesOutputMessageText(m.content);
    if (text) parts.push(text);
  }
  return parts.join("\n");
}

function responsesOutputMessageText(content: unknown): string {
  if (typeof content === "string") return content;
  if (!Array.isArray(content)) return "";
  const parts: string[] = [];
  for (const bl of content) {
    if (!bl || typeof bl !== "object") continue;
    const bm = bl as Record<string, unknown>;
    if (typeof bm.text === "string" && bm.text) parts.push(bm.text);
  }
  return parts.join("");
}

export function parseProxyError(body: unknown, status: number): string {
  if (body && typeof body === "object") {
    const err = (body as Record<string, unknown>).error;
    if (err && typeof err === "object") {
      const msg = (err as Record<string, unknown>).message;
      if (typeof msg === "string" && msg) return msg;
    }
    const message = (body as Record<string, unknown>).message;
    if (typeof message === "string" && message) return message;
  }
  if (status === 401) return "Invalid API key";
  if (status === 403) return "API key cannot use this model";
  return `Request failed (HTTP ${status})`;
}
