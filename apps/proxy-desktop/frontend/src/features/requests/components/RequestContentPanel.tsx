import { useQuery } from "@tanstack/react-query";
import { Badge } from "@/components/ui/badge";
import { getAIRequestContent } from "@/features/requests/api/requestsApi";
import { ErrorState, LoadingState } from "@/shared/components/PageStates";

function formatJsonBlock(raw: string): string {
  const trimmed = raw.trim();
  if (!trimmed) {
    return "";
  }
  try {
    return JSON.stringify(JSON.parse(trimmed), null, 2);
  } catch {
    return raw;
  }
}

function TruncationBadge({ truncated, label }: { truncated: boolean; label: string }) {
  if (!truncated) {
    return null;
  }
  return (
    <Badge variant="secondary" className="ml-2 text-[10px] font-normal">
      {label} truncated
    </Badge>
  );
}

function ContentBlock({
  title,
  body,
  truncated,
  truncationLabel,
  blockClassName,
}: {
  title: string;
  body: string;
  truncated?: boolean;
  truncationLabel?: string;
  blockClassName?: string;
}) {
  if (!body.trim()) {
    return null;
  }
  return (
    <div className="space-y-1">
      <p className="font-medium text-foreground">
        {title}
        {truncated && truncationLabel ? (
          <TruncationBadge truncated={truncated} label={truncationLabel} />
        ) : null}
      </p>
      <pre
        className={`overflow-auto rounded-md border border-border bg-background p-3 font-mono text-[11px] leading-relaxed whitespace-pre-wrap break-words ${blockClassName ?? "max-h-64"}`}
      >
        {body}
      </pre>
    </div>
  );
}

export function RequestContentPanel({
  requestId,
  contentClassName,
}: {
  requestId: string;
  contentClassName?: string;
}) {
  const contentQ = useQuery({
    queryKey: ["request-content", requestId],
    queryFn: () => getAIRequestContent(requestId),
    staleTime: 60_000,
  });

  if (contentQ.isLoading) {
    return <LoadingState />;
  }
  if (contentQ.error) {
    return (
      <ErrorState
        message={(contentQ.error as Error).message}
        onRetry={() => contentQ.refetch()}
      />
    );
  }

  const data = contentQ.data;
  if (!data) {
    return null;
  }

  const captureErrors = data.capture_errors ?? [];
  const proxyError = data.proxy_error?.trim() ?? "";

  const promptBody = data.capture_enabled
    ? data.messages_json
      ? formatJsonBlock(data.messages_json)
      : (data.last_user_text ?? "")
    : "";

  const responseBody = data.capture_enabled ? (data.assistant_text ?? "") : "";
  const rawBody = data.capture_enabled && data.raw_body ? formatJsonBlock(data.raw_body) : "";

  return (
    <div className="space-y-4 pt-2">
      {(proxyError || captureErrors.length > 0) && (
        <div className="space-y-2 rounded-md border border-destructive/40 bg-destructive/10 p-3">
          <p className="text-sm font-medium text-destructive">Errors</p>
          {proxyError && captureErrors.length === 0 && (
            <p className="font-mono text-xs break-words text-foreground">{proxyError}</p>
          )}
          {captureErrors.length > 0 && (
            <ul className="space-y-2 text-xs text-foreground">
              {captureErrors.map((e, i) => (
                <li key={`${e.stage}-${e.error_code}-${i}`} className="font-mono break-words">
                  <span className="text-destructive">{e.error_code || "error"}</span>
                  {e.stage ? ` · ${e.stage}` : ""}
                  <div className="mt-0.5 font-sans text-foreground">{e.error_message}</div>
                </li>
              ))}
            </ul>
          )}
        </div>
      )}

      {!data.capture_enabled && (
        <p className="text-muted-foreground">
          Prompt and response text are not stored in this database. Enable{" "}
          <code className="rounded bg-muted px-1 py-0.5 text-[10px]">capture.store_content: true</code>{" "}
          in ai-proxy config and send new requests.
        </p>
      )}

      {data.capture_enabled && !data.has_content && captureErrors.length === 0 && !proxyError && (
        <p className="text-muted-foreground">No captured prompt or response for this request.</p>
      )}

      <ContentBlock
        title="Prompt"
        body={promptBody}
        truncated={data.prompt_truncated}
        truncationLabel="Prompt"
        blockClassName={contentClassName}
      />

      <ContentBlock
        title="Response"
        body={responseBody}
        truncated={data.response_truncated || data.tool_calls_truncated}
        truncationLabel="Response"
        blockClassName={contentClassName}
      />

      {data.tool_calls_json ? (
        <ContentBlock
          title="Tool calls"
          body={formatJsonBlock(data.tool_calls_json)}
          truncated={data.tool_calls_truncated}
          truncationLabel="Tool calls"
          blockClassName={contentClassName}
        />
      ) : null}

      {rawBody && (
        <ContentBlock
          title="Raw upstream body"
          body={rawBody}
          truncated={data.raw_body_truncated}
          truncationLabel="Body"
          blockClassName={contentClassName}
        />
      )}
    </div>
  );
}
