import { describe, expect, it } from "vitest";
import {
  buildRequestBody,
  chatEndpoint,
  extractAssistantText,
} from "@/features/chat/api/chatProtocol";
import type { ChatMessage } from "@/features/chat/types";

const msgs: ChatMessage[] = [
  { id: "1", role: "system", content: "Be helpful" },
  { id: "2", role: "user", content: "Hi" },
];

describe("chatEndpoint", () => {
  it("maps protocols to paths", () => {
    expect(chatEndpoint("openai_chat")).toBe("/v1/chat/completions");
    expect(chatEndpoint("anthropic_messages")).toBe("/v1/messages");
    expect(chatEndpoint("openai_responses")).toBe("/v1/responses");
  });
});

describe("buildRequestBody", () => {
  it("builds openai chat body", () => {
    const body = buildRequestBody("openai_chat", "gpt-test", msgs, false);
    expect(body.model).toBe("gpt-test");
    expect(body.stream).toBe(false);
    expect(body.messages).toHaveLength(2);
  });

  it("builds anthropic body with system split", () => {
    const body = buildRequestBody("anthropic_messages", "claude-test", msgs, true);
    expect(body.model).toBe("claude-test");
    expect(body.system).toBe("Be helpful");
    expect(body.messages).toEqual([{ role: "user", content: "Hi" }]);
    expect(body.stream).toBe(true);
  });

  it("builds responses body with instructions", () => {
    const body = buildRequestBody("openai_responses", "o-test", msgs, false);
    expect(body.instructions).toBe("Be helpful");
    expect(body.input).toBe("Hi");
  });
});

describe("extractAssistantText", () => {
  it("parses openai chat", () => {
    const text = extractAssistantText("openai_chat", {
      choices: [{ message: { content: "Hello" } }],
    });
    expect(text).toBe("Hello");
  });

  it("parses anthropic", () => {
    const text = extractAssistantText("anthropic_messages", {
      content: [{ type: "text", text: "Hi there" }],
    });
    expect(text).toBe("Hi there");
  });

  it("parses responses output_text", () => {
    const text = extractAssistantText("openai_responses", { output_text: "Done" });
    expect(text).toBe("Done");
  });
});
