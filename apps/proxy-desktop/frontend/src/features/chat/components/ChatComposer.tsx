import { useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import { Send } from "lucide-react";

type Props = {
  disabled: boolean;
  onSend: (text: string) => Promise<void>;
};

export function ChatComposer({ disabled, onSend }: Props) {
  const [draft, setDraft] = useState("");
  const [submitting, setSubmitting] = useState(false);

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    const text = draft.trim();
    if (!text || disabled || submitting) return;
    setSubmitting(true);
    try {
      await onSend(text);
      setDraft("");
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <form onSubmit={handleSubmit} className="flex gap-2">
      <Input
        value={draft}
        onChange={(e) => setDraft(e.target.value)}
        placeholder="Type a message…"
        disabled={disabled || submitting}
        autoComplete="off"
      />
      <Button type="submit" disabled={disabled || submitting || !draft.trim()}>
        <Send className="h-4 w-4" />
      </Button>
    </form>
  );
}
