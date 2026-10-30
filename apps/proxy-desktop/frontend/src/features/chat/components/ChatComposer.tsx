import { useEffect, useRef, useState } from "react";
import { ImagePlus, Send, X } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";
import {
  CHAT_IMAGE_MAX_COUNT,
  chatImageDataUrl,
  readChatImageFile,
  type ChatImageAttachment,
} from "@/features/chat/chatImageAttachment";

type Props = {
  disabled: boolean;
  /** When false, image attach is hidden (text-only / image-gen models). */
  allowImageAttach?: boolean;
  maxImageFileBytes?: number;
  onSend: (text: string, images: ChatImageAttachment[]) => Promise<void>;
};

export function ChatComposer({
  disabled,
  allowImageAttach = true,
  maxImageFileBytes,
  onSend,
}: Props) {
  const [draft, setDraft] = useState("");
  const [attachments, setAttachments] = useState<ChatImageAttachment[]>([]);
  const [submitting, setSubmitting] = useState(false);
  const [attachError, setAttachError] = useState<string | null>(null);
  const fileInputRef = useRef<HTMLInputElement>(null);

  useEffect(() => {
    if (!allowImageAttach && attachments.length > 0) {
      setAttachments([]);
      setAttachError(null);
    }
  }, [allowImageAttach, attachments.length]);

  const canSend =
    !disabled && !submitting && (draft.trim().length > 0 || attachments.length > 0);

  async function handleFiles(files: FileList | null) {
    if (!allowImageAttach || !files?.length) return;
    setAttachError(null);
    const next = [...attachments];
    for (const file of Array.from(files)) {
      if (next.length >= CHAT_IMAGE_MAX_COUNT) {
        setAttachError(`At most ${CHAT_IMAGE_MAX_COUNT} images per message`);
        break;
      }
      try {
        next.push(await readChatImageFile(file, maxImageFileBytes));
      } catch (e) {
        setAttachError(e instanceof Error ? e.message : "Could not attach image");
      }
    }
    setAttachments(next);
    if (fileInputRef.current) {
      fileInputRef.current.value = "";
    }
  }

  async function handleSubmit(e: React.FormEvent) {
    e.preventDefault();
    if (!canSend) return;
    setSubmitting(true);
    try {
      await onSend(draft, attachments);
      setDraft("");
      setAttachments([]);
      setAttachError(null);
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <div className="space-y-2">
      {attachments.length > 0 && (
        <div className="flex flex-wrap gap-2">
          {attachments.map((img, index) => (
            <div key={`${img.mediaType}-${index}`} className="relative">
              <img
                src={chatImageDataUrl(img)}
                alt=""
                className="h-16 w-16 rounded-md border border-border object-cover"
              />
              <Button
                type="button"
                size="icon"
                variant="secondary"
                className="absolute -right-1 -top-1 h-6 w-6 rounded-full"
                aria-label="Remove image"
                disabled={submitting}
                onClick={() => setAttachments((prev) => prev.filter((_, i) => i !== index))}
              >
                <X className="h-3 w-3" />
              </Button>
            </div>
          ))}
        </div>
      )}
      {attachError && <p className="text-xs text-destructive">{attachError}</p>}
      <form onSubmit={handleSubmit} className="flex gap-2">
        {allowImageAttach && (
          <>
            <input
              ref={fileInputRef}
              type="file"
              accept="image/jpeg,image/png,image/webp,image/gif"
              multiple
              className="hidden"
              onChange={(e) => void handleFiles(e.target.files)}
            />
            <Button
              type="button"
              variant="outline"
              size="icon"
              disabled={disabled || submitting || attachments.length >= CHAT_IMAGE_MAX_COUNT}
              aria-label="Attach image"
              onClick={() => fileInputRef.current?.click()}
            >
              <ImagePlus className="h-4 w-4" />
            </Button>
          </>
        )}
        <Input
          value={draft}
          onChange={(e) => setDraft(e.target.value)}
          placeholder="Type a message…"
          disabled={disabled || submitting}
          autoComplete="off"
          className="min-w-0 flex-1"
        />
        <Button type="submit" disabled={!canSend}>
          <Send className="h-4 w-4" />
        </Button>
      </form>
    </div>
  );
}
