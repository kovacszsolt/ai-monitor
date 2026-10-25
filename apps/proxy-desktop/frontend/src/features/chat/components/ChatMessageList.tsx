import { Pin } from "lucide-react";
import { Button } from "@/components/ui/button";
import { cn } from "@/lib/utils";
import type { ChatMessage } from "@/features/chat/types";

type Props = {
  messages: ChatMessage[];
  streamingContent: string | null;
  sending: boolean;
  onPinUserMessage?: (content: string) => void;
};

export function ChatMessageList({ messages, streamingContent, sending, onPinUserMessage }: Props) {
  if (messages.length === 0 && !sending) {
    return (
      <p className="text-sm text-muted-foreground">Send a message to start the conversation.</p>
    );
  }

  const showStreamBubble = sending && streamingContent !== null;

  return (
    <ul className="flex flex-col gap-3">
      {messages.map((m) => (
        <li
          key={m.id}
          className={cn(
            "max-w-[85%] rounded-lg px-3 py-2 text-sm whitespace-pre-wrap",
            m.role === "user"
              ? "ml-auto bg-primary text-primary-foreground"
              : "mr-auto bg-muted text-foreground",
          )}
        >
          <div className="flex items-start justify-between gap-2">
            <span className="min-w-0 flex-1">{m.content}</span>
            {m.role === "user" && onPinUserMessage && (
              <Button
                type="button"
                size="icon"
                variant="ghost"
                className="h-7 w-7 shrink-0 text-primary-foreground/80 hover:bg-primary-foreground/10 hover:text-primary-foreground"
                aria-label="Pin question"
                onClick={() => onPinUserMessage(m.content)}
              >
                <Pin className="h-3.5 w-3.5" />
              </Button>
            )}
          </div>
        </li>
      ))}
      {showStreamBubble && (
        <li className="mr-auto max-w-[85%] rounded-lg bg-muted px-3 py-2 text-sm whitespace-pre-wrap text-foreground">
          {streamingContent || "…"}
        </li>
      )}
      {sending && streamingContent === null && (
        <li className="mr-auto max-w-[85%] rounded-lg bg-muted px-3 py-2 text-sm text-muted-foreground">
          …
        </li>
      )}
    </ul>
  );
}
