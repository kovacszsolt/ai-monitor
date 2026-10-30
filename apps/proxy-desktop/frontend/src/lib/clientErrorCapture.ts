import * as DesktopService from "../../wailsjs/go/service/DesktopService";
import type { ApiResult } from "@/lib/api";

export type ClientErrorPayload = {
  level: "error" | "warn" | "info";
  source: string;
  message: string;
  detail?: string;
};

let installDone = false;
let forwarding = false;

const pending: ClientErrorPayload[] = [];
let flushTimer: ReturnType<typeof setTimeout> | null = null;

function formatConsoleArg(value: unknown): string {
  if (value instanceof Error) {
    const stack = value.stack?.trim();
    if (stack) {
      return stack;
    }
    return value.message;
  }
  if (typeof value === "string") {
    return value;
  }
  try {
    return JSON.stringify(value);
  } catch {
    return String(value);
  }
}

export function formatConsoleArgs(args: unknown[]): { message: string; detail: string } {
  if (args.length === 0) {
    return { message: "(console.error with no arguments)", detail: "" };
  }
  if (args.length === 1 && args[0] instanceof Error) {
    const err = args[0];
    return {
      message: err.message || "Error",
      detail: err.stack?.trim() ?? "",
    };
  }
  const first = args[0];
  const rest = args.slice(1);
  if (typeof first === "string" && rest.length === 1 && rest[0] instanceof Error) {
    const err = rest[0];
    const message = err.message ? `${first}: ${err.message}` : first;
    return { message, detail: err.stack?.trim() ?? "" };
  }
  const parts = args.map(formatConsoleArg);
  const message = parts[0] ?? "(empty)";
  const detail = parts.length > 1 ? parts.slice(1).join("\n") : "";
  return { message, detail };
}

function scheduleFlush() {
  if (flushTimer != null) {
    return;
  }
  flushTimer = setTimeout(() => {
    flushTimer = null;
    void flushPending();
  }, 200);
}

async function flushPending() {
  if (forwarding || pending.length === 0) {
    return;
  }
  const batch = pending.splice(0, 8);
  forwarding = true;
  try {
    for (const item of batch) {
      await DesktopService.LogClientError(
        item.level,
        item.source,
        item.message,
        item.detail ?? "",
      );
    }
  } catch {
    // Avoid infinite loops; original console output already happened.
  } finally {
    forwarding = false;
    if (pending.length > 0) {
      scheduleFlush();
    }
  }
}

function wailsReady(): boolean {
  return typeof window !== "undefined" && "go" in window;
}

/** Persist a client-side error for the diagnostics table (fire-and-forget). */
export function enqueueClientError(payload: ClientErrorPayload) {
  if (!wailsReady()) {
    return;
  }
  pending.push({
    level: payload.level,
    source: payload.source.slice(0, 256),
    message: payload.message.slice(0, 8192),
    detail: (payload.detail ?? "").slice(0, 65536),
  });
  scheduleFlush();
}

export function installClientErrorCapture() {
  if (installDone || typeof window === "undefined") {
    return;
  }
  installDone = true;

  const origError = console.error.bind(console);
  const origWarn = console.warn.bind(console);

  console.error = (...args: unknown[]) => {
    origError(...args);
    if (forwarding) {
      return;
    }
    const { message, detail } = formatConsoleArgs(args);
    enqueueClientError({ level: "error", source: "console.error", message, detail });
  };

  console.warn = (...args: unknown[]) => {
    origWarn(...args);
    if (forwarding) {
      return;
    }
    const { message, detail } = formatConsoleArgs(args);
    enqueueClientError({ level: "warn", source: "console.warn", message, detail });
  };

  window.addEventListener("error", (event) => {
    const message = event.message || "Uncaught error";
    const detailParts = [
      event.filename ? `file: ${event.filename}` : "",
      event.lineno ? `line: ${event.lineno}` : "",
      event.colno ? `col: ${event.colno}` : "",
      event.error instanceof Error && event.error.stack ? event.error.stack : "",
    ].filter(Boolean);
    enqueueClientError({
      level: "error",
      source: "window.error",
      message,
      detail: detailParts.join("\n"),
    });
  });

  window.addEventListener("unhandledrejection", (event) => {
    const reason = event.reason;
    let message = "Unhandled promise rejection";
    let detail = "";
    if (reason instanceof Error) {
      message = reason.message || message;
      detail = reason.stack ?? "";
    } else if (reason != null) {
      detail = formatConsoleArg(reason);
    }
    enqueueClientError({
      level: "error",
      source: "unhandledrejection",
      message,
      detail,
    });
  });
}

export type ClientErrorLogEntry = {
  id: string;
  level: string;
  source: string;
  message: string;
  detail: string;
  created_at: string;
};

export type ClientErrorLogList = {
  items: ClientErrorLogEntry[];
  total: number;
  page: number;
  page_size: number;
};

export async function listClientErrors(page: number, pageSize: number) {
  const result = (await DesktopService.ListClientErrors(page, pageSize)) as ApiResult<ClientErrorLogList>;
  if (!result.ok || !result.data) {
    throw new Error(result.error?.message ?? "Failed to load error log");
  }
  return result.data;
}

export async function clearClientErrors() {
  const result = (await DesktopService.ClearClientErrors()) as ApiResult<boolean>;
  if (!result.ok) {
    throw new Error(result.error?.message ?? "Failed to clear error log");
  }
  return result.data ?? true;
}
