export type ChatImageAttachment = {
  mediaType: string;
  dataBase64: string;
};

export const CHAT_IMAGE_MAX_BYTES = 40 * 1024 * 1024;
export const CHAT_IMAGE_MAX_COUNT = 4;

/** Raw file size limit so one image fits in proxy server.max_request_bytes (JSON base64 overhead). */
export function maxChatImageFileBytesForRequestLimit(serverMaxRequestBytes: number): number {
  if (!Number.isFinite(serverMaxRequestBytes) || serverMaxRequestBytes <= 0) {
    return CHAT_IMAGE_MAX_BYTES;
  }
  const jsonHeadroom = 512 * 1024;
  const budget = Math.max(0, serverMaxRequestBytes - jsonHeadroom);
  return Math.min(CHAT_IMAGE_MAX_BYTES, Math.floor((budget * 3) / 4));
}

const ALLOWED = new Set(["image/jpeg", "image/png", "image/webp", "image/gif"]);

export function isAllowedChatImageType(mediaType: string): boolean {
  return ALLOWED.has(mediaType);
}

export function chatImageDataUrl(attachment: ChatImageAttachment): string {
  return `data:${attachment.mediaType};base64,${attachment.dataBase64}`;
}

export async function readChatImageFile(
  file: File,
  maxFileBytes: number = CHAT_IMAGE_MAX_BYTES,
): Promise<ChatImageAttachment> {
  if (!isAllowedChatImageType(file.type)) {
    throw new Error(`Unsupported image type: ${file.type || "unknown"}`);
  }
  const cap = Math.min(CHAT_IMAGE_MAX_BYTES, maxFileBytes);
  if (file.size > cap) {
    const capMB = Math.round(cap / (1024 * 1024));
    throw new Error(
      cap < CHAT_IMAGE_MAX_BYTES
        ? `Image too large for proxy limit (max ~${capMB} MB — raise server.max_request_bytes in Runtime Vars)`
        : `Image too large (max ${capMB} MB)`,
    );
  }
  const dataBase64 = await readFileAsBase64(file);
  return { mediaType: file.type, dataBase64 };
}

function readFileAsBase64(file: File): Promise<string> {
  return new Promise((resolve, reject) => {
    const reader = new FileReader();
    reader.onload = () => {
      const result = reader.result;
      if (typeof result !== "string") {
        reject(new Error("Could not read image"));
        return;
      }
      const comma = result.indexOf(",");
      if (comma === -1) {
        reject(new Error("Invalid image data"));
        return;
      }
      resolve(result.slice(comma + 1));
    };
    reader.onerror = () => reject(new Error("Could not read image"));
    reader.readAsDataURL(file);
  });
}
