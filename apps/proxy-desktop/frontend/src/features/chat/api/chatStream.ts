import { extractAssistantText } from "@/features/chat/api/chatProtocol";

export async function readStreamingAssistantText(
  protocol: string,
  response: Response,
  onDelta: (accumulated: string) => void,
): Promise<string> {
  const reader = response.body?.getReader();
  if (!reader) {
    throw new Error("Streaming not supported");
  }
  const decoder = new TextDecoder();
  let buffer = "";
  let accumulated = "";

  while (true) {
    const { done, value } = await reader.read();
    if (done) break;
    buffer += decoder.decode(value, { stream: true });
    const lines = buffer.split("\n");
    buffer = lines.pop() ?? "";
    for (const line of lines) {
      const delta = parseStreamLine(protocol, line);
      if (delta) {
        accumulated += delta;
        onDelta(accumulated);
      }
    }
  }
  if (buffer.trim()) {
    const delta = parseStreamLine(protocol, buffer);
    if (delta) {
      accumulated += delta;
      onDelta(accumulated);
    }
  }
  return accumulated;
}

function parseStreamLine(protocol: string, line: string): string {
  const trimmed = line.trim();
  if (!trimmed || trimmed.startsWith(":")) return "";

  if (protocol === "anthropic_messages") {
    return parseAnthropicStreamLine(trimmed);
  }
  return parseOpenAIStyleStreamLine(trimmed, protocol);
}

function parseOpenAIStyleStreamLine(line: string, protocol: string): string {
  if (!line.startsWith("data:")) return "";
  const payload = line.slice(5).trim();
  if (payload === "[DONE]") return "";
  try {
    const json = JSON.parse(payload) as unknown;
    return extractStreamDelta(protocol, json);
  } catch {
    return "";
  }
}

function parseAnthropicStreamLine(line: string): string {
  if (!line.startsWith("data:")) return "";
  const payload = line.slice(5).trim();
  try {
    const json = JSON.parse(payload) as Record<string, unknown>;
    if (json.type === "content_block_delta") {
      const delta = json.delta as Record<string, unknown> | undefined;
      if (delta?.type === "text_delta" && typeof delta.text === "string") {
        return delta.text;
      }
    }
  } catch {
    return "";
  }
  return "";
}

function extractStreamDelta(protocol: string, json: unknown): string {
  if (!json || typeof json !== "object") return "";
  const v = json as Record<string, unknown>;

  if (protocol === "openai_responses") {
    if (typeof v.delta === "string") return v.delta;
    const output = v.output;
    if (Array.isArray(output) && output.length > 0) {
      return extractAssistantText(protocol, { output });
    }
    const type = v.type;
    if (type === "response.output_text.delta" && typeof v.delta === "string") {
      return v.delta;
    }
    if (type === "response.output_text.done" && typeof v.text === "string") {
      return v.text;
    }
    return "";
  }

  const choices = v.choices;
  if (!Array.isArray(choices) || choices.length === 0) return "";
  const choice = choices[0] as Record<string, unknown>;
  const delta = choice.delta as Record<string, unknown> | undefined;
  if (delta && typeof delta.content === "string") return delta.content;
  return "";
}
