import { enqueueClientError } from "@/lib/clientErrorCapture";

export type ApiError = { code: string; message: string };

export type ApiResult<T> = {
  ok: boolean;
  data?: T;
  error?: ApiError;
};

export type UnwrapOptions = {
  /** Shown in error log detail, e.g. SendChatCompletion */
  context?: string;
  /** Default true — set false when the caller logs explicitly */
  log?: boolean;
};

function buildApiErrorDetail(
  apiErr: ApiError | undefined,
  context: string | undefined,
  stack: string | undefined,
): string {
  const lines: string[] = [];
  if (apiErr?.code) {
    lines.push(`code: ${apiErr.code}`);
  }
  if (context) {
    lines.push(`context: ${context}`);
  }
  if (stack) {
    lines.push(stack);
  }
  return lines.join("\n");
}

export function unwrap<T>(result: ApiResult<T>, options?: UnwrapOptions): T {
  if (!result.ok) {
    const message = result.error?.message ?? "Request failed";
    const err = new Error(message);
    if (options?.log !== false) {
      enqueueClientError({
        level: "error",
        source: "api",
        message,
        detail: buildApiErrorDetail(result.error, options?.context, err.stack),
      });
    }
    throw err;
  }
  if (result.data === undefined) {
    const message = result.error?.message ?? "Request failed";
    const err = new Error(message);
    if (options?.log !== false) {
      enqueueClientError({
        level: "error",
        source: "api",
        message,
        detail: buildApiErrorDetail(result.error, options?.context, err.stack),
      });
    }
    throw err;
  }
  return result.data;
}
