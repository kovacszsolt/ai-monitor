import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useEffect, useMemo } from "react";
import { Link } from "react-router-dom";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardHeader, CardTitle } from "@/components/ui/card";
import { Label } from "@/components/ui/label";
import { Switch } from "@/components/ui/switch";
import { createChatPinnedPrompt, deleteChatPinnedPrompt, listChatPinnedPrompts } from "@/features/chat/api/chatPinsApi";
import { ChatComposer } from "@/features/chat/components/ChatComposer";
import { ChatMessageList } from "@/features/chat/components/ChatMessageList";
import { ChatPinnedPromptsPanel } from "@/features/chat/components/ChatPinnedPromptsPanel";
import { useChatPreferences } from "@/features/chat/hooks/useChatPreferences";
import { useChatSessionPanelCollapsed } from "@/features/chat/hooks/useChatSessionPanelCollapsed";
import { useChatSend } from "@/features/chat/hooks/useChatSend";
import { getProxyVariables } from "@/features/settings/api/settingsApi";
import {
  maxChatImageFileBytesForRequestLimit,
  type ChatImageAttachment,
} from "@/features/chat/chatImageAttachment";
import { listModelRoutes } from "@/features/models/api/modelsApi";
import { getProxyStatus } from "@/features/proxy/api/proxyApi";
import { listAPIKeys } from "@/features/projects/api/keysApi";
import { listProjects } from "@/features/projects/api/projectsApi";
import { listProviders } from "@/features/providers/api/providersApi";
import { CollapsibleCardHeader } from "@/shared/components/CollapsibleCardHeader";
import { PageHeader } from "@/shared/components/PageHeader";
import { PageShell } from "@/shared/components/PageShell";
import { EmptyState, ErrorState, LoadingState } from "@/shared/components/PageStates";

const selectClass =
  "flex h-9 w-full min-w-0 rounded border border-outline-variant bg-surface-container px-space-md text-body-sm focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/40";

function pickValidId(saved: string, options: { id: string }[]): string {
  if (saved && options.some((o) => o.id === saved)) return saved;
  return options[0]?.id ?? "";
}

export function ChatPage() {
  const qc = useQueryClient();
  const { preferences, setPreference, setPreferences } = useChatPreferences();
  const { collapsed: sessionCollapsed, toggle: toggleSessionPanel } = useChatSessionPanelCollapsed();

  const statusQ = useQuery({
    queryKey: ["proxy-status"],
    queryFn: getProxyStatus,
    refetchInterval: 3000,
  });
  const proxyVarsQ = useQuery({ queryKey: ["proxy-variables"], queryFn: getProxyVariables });
  const maxImageFileBytes = useMemo(
    () => maxChatImageFileBytesForRequestLimit(proxyVarsQ.data?.server_max_request_bytes ?? 0),
    [proxyVarsQ.data?.server_max_request_bytes],
  );
  const projectsQ = useQuery({ queryKey: ["projects"], queryFn: listProjects });
  const providersQ = useQuery({ queryKey: ["providers"], queryFn: listProviders });
  const modelsQ = useQuery({ queryKey: ["models"], queryFn: listModelRoutes });
  const keysQ = useQuery({
    queryKey: ["keys", preferences.projectId],
    queryFn: () => listAPIKeys(preferences.projectId),
    enabled: Boolean(preferences.projectId),
  });

  const enabledProjects = useMemo(
    () => (projectsQ.data ?? []).filter((p) => p.enabled),
    [projectsQ.data],
  );
  const enabledProviders = useMemo(
    () => (providersQ.data ?? []).filter((p) => p.enabled),
    [providersQ.data],
  );
  const modelsForProvider = useMemo(() => {
    const all = (modelsQ.data ?? []).filter((m) => m.enabled);
    if (!preferences.providerId) return all;
    return all.filter((m) => m.provider_id === preferences.providerId);
  }, [modelsQ.data, preferences.providerId]);

  const enabledKeys = useMemo(
    () => (keysQ.data ?? []).filter((k) => k.enabled),
    [keysQ.data],
  );

  useEffect(() => {
    if (!projectsQ.data || !providersQ.data || !modelsQ.data) return;

    let projectId = preferences.projectId;
    let providerId = preferences.providerId;
    let modelRouteId = preferences.modelRouteId;
    let apiKeyId = preferences.apiKeyId;

    const nextProject = pickValidId(projectId, enabledProjects);
    if (nextProject !== projectId) {
      projectId = nextProject;
      apiKeyId = "";
      if (projectId && projectId !== preferences.projectId) {
        toast.message("Chat project selection updated");
      }
    }

    const nextProvider = pickValidId(providerId, enabledProviders);
    if (nextProvider !== providerId) providerId = nextProvider;

    const routes =
      providerId
        ? (modelsQ.data ?? []).filter((m) => m.enabled && m.provider_id === providerId)
        : (modelsQ.data ?? []).filter((m) => m.enabled);
    const nextModel = pickValidId(modelRouteId, routes);
    if (nextModel !== modelRouteId) modelRouteId = nextModel;

    const patch: Partial<typeof preferences> = {};
    if (projectId !== preferences.projectId) patch.projectId = projectId;
    if (providerId !== preferences.providerId) patch.providerId = providerId;
    if (modelRouteId !== preferences.modelRouteId) patch.modelRouteId = modelRouteId;
    if (Object.keys(patch).length > 0) setPreferences(patch);
  }, [
    projectsQ.data,
    providersQ.data,
    modelsQ.data,
    enabledProjects,
    enabledProviders,
    preferences.projectId,
    preferences.providerId,
    preferences.modelRouteId,
    setPreferences,
  ]);

  useEffect(() => {
    if (!keysQ.data) return;
    const nextKey = pickValidId(preferences.apiKeyId, enabledKeys);
    if (nextKey !== preferences.apiKeyId) {
      setPreference("apiKeyId", nextKey);
    }
  }, [keysQ.data, enabledKeys, preferences.apiKeyId, setPreference]);

  const selectedModel = modelsForProvider.find((m) => m.id === preferences.modelRouteId);
  const selectedKey = enabledKeys.find((k) => k.id === preferences.apiKeyId);
  const selectedProject = enabledProjects.find((p) => p.id === preferences.projectId);
  const selectedProvider = enabledProviders.find((p) => p.id === preferences.providerId);

  const sessionSummary = useMemo(() => {
    const parts = [
      selectedProject?.name,
      selectedKey?.name,
      selectedProvider?.name,
      selectedModel?.public_model_name,
      preferences.stream ? "streaming" : "non-streaming",
    ].filter(Boolean);
    return parts.join(" · ");
  }, [
    selectedProject?.name,
    selectedKey?.name,
    selectedProvider?.name,
    selectedModel?.public_model_name,
    preferences.stream,
  ]);

  const proxyRunning = statusQ.data?.running ?? false;

  const isImageGenModel = (selectedModel?.modality || "text-text") === "text-image";
  const chatSend = useChatSend({
    apiKeyId: selectedKey?.id,
    modelRouteId: selectedModel?.id,
    stream: isImageGenModel ? false : preferences.stream,
  });

  const pinsQ = useQuery({ queryKey: ["chat-pins"], queryFn: listChatPinnedPrompts });

  const sendMut = useMutation({
    mutationFn: ({ text, images }: { text: string; images: ChatImageAttachment[] }) =>
      chatSend.sendUserMessage(text, images),
    onError: (e: Error) => {
      toast.error(e.message || "Chat request failed");
    },
  });

  const pinMut = useMutation({
    mutationFn: (body: string) => createChatPinnedPrompt(body, ""),
    onSuccess: () => {
      toast.success("Pinned");
      qc.invalidateQueries({ queryKey: ["chat-pins"] });
    },
    onError: (e: Error) => {
      toast.error(e.message);
    },
  });

  const unpinMut = useMutation({
    mutationFn: (id: string) => deleteChatPinnedPrompt(id),
    onSuccess: () => {
      toast.success("Unpinned");
      qc.invalidateQueries({ queryKey: ["chat-pins"] });
    },
    onError: (e: Error) => {
      toast.error(e.message);
    },
  });

  const loading =
    projectsQ.isLoading || providersQ.isLoading || modelsQ.isLoading || statusQ.isLoading;
  const loadError =
    projectsQ.error ?? providersQ.error ?? modelsQ.error ?? statusQ.error;

  if (loading) return <LoadingState label="Loading chat…" />;
  if (loadError) {
    return (
      <ErrorState
        message={(loadError as Error).message}
        onRetry={() => {
          void projectsQ.refetch();
          void providersQ.refetch();
          void modelsQ.refetch();
          void statusQ.refetch();
        }}
      />
    );
  }

  if (!proxyRunning) {
    return (
      <PageShell className="space-y-4">
        <h2 className="text-2xl font-semibold tracking-tight">Chat</h2>
        <ErrorState
          message="Start the local proxy before chatting."
          onRetry={() => statusQ.refetch()}
        />
        <Button asChild variant="outline">
          <Link to="/proxy">Open Proxy</Link>
        </Button>
      </PageShell>
    );
  }

  if (enabledProjects.length === 0) {
    return (
      <PageShell>
        <EmptyState title="No enabled projects." />
        <Button asChild className="mt-4" variant="outline">
          <Link to="/projects">Projects</Link>
        </Button>
      </PageShell>
    );
  }

  const composerDisabled =
    !selectedModel || !selectedKey?.id || sendMut.isPending || chatSend.sending;

  return (
    <PageShell className="flex h-[calc(100vh-3.5rem)] w-full gap-space-md">
      <PageHeader
        section="Operate"
        accent="Smoke Test"
        title="Chat Testbed"
        description={
          <>
            Messages go through the local proxy and are recorded under the selected project (tag{" "}
            <code className="text-code-inline">desktop-chat</code>).{" "}
            <code className="text-code-inline">text_image-text</code> allows image attach;{" "}
            <code className="text-code-inline">text-image</code> generates images from the prompt.
          </>
        }
      />

      <Card>
        <CardHeader className="pb-3">
          <CollapsibleCardHeader
            title="Session"
            collapsed={sessionCollapsed}
            onToggle={toggleSessionPanel}
            description={
              sessionCollapsed
                ? sessionSummary
                : "Selections are remembered on this device."
            }
          />
        </CardHeader>
        {!sessionCollapsed && (
        <CardContent className="grid gap-3 sm:grid-cols-2">
          <div className="space-y-1.5">
            <Label htmlFor="chat-project">Project</Label>
            <select
              id="chat-project"
              className={selectClass}
              value={preferences.projectId}
              onChange={(e) => {
                setPreferences({ projectId: e.target.value, apiKeyId: "" });
              }}
            >
              {enabledProjects.map((p) => (
                <option key={p.id} value={p.id}>
                  {p.name}
                </option>
              ))}
            </select>
          </div>
          <div className="space-y-1.5">
            <Label htmlFor="chat-api-key">API key</Label>
            {enabledKeys.length === 0 ? (
              <p className="text-sm text-muted-foreground">
                No keys —{" "}
                <Link className="underline" to={`/projects/${preferences.projectId}`}>
                  create one
                </Link>
              </p>
            ) : (
              <select
                id="chat-api-key"
                className={selectClass}
                value={preferences.apiKeyId}
                onChange={(e) => setPreference("apiKeyId", e.target.value)}
              >
                {enabledKeys.map((k) => (
                  <option key={k.id} value={k.id}>
                    {k.name}
                  </option>
                ))}
              </select>
            )}
          </div>
          <div className="space-y-1.5">
            <Label htmlFor="chat-provider">Provider</Label>
            <select
              id="chat-provider"
              className={selectClass}
              value={preferences.providerId}
              onChange={(e) => {
                const providerId = e.target.value;
                const routes = (modelsQ.data ?? []).filter(
                  (m) => m.enabled && m.provider_id === providerId,
                );
                const modelRouteId = pickValidId("", routes);
                setPreferences({ providerId, modelRouteId });
              }}
            >
              {enabledProviders.map((p) => (
                <option key={p.id} value={p.id}>
                  {p.name}
                </option>
              ))}
            </select>
          </div>
          <div className="space-y-1.5">
            <Label htmlFor="chat-model">Model</Label>
            <select
              id="chat-model"
              className={selectClass}
              value={preferences.modelRouteId}
              onChange={(e) => setPreference("modelRouteId", e.target.value)}
            >
              {modelsForProvider.map((m) => (
                <option key={m.id} value={m.id}>
                  {m.public_model_name} ({m.modality || "text-text"} · {m.protocol})
                </option>
              ))}
            </select>
          </div>
          <div className="flex items-center gap-2 sm:col-span-2">
            <Switch
              id="chat-stream"
              checked={isImageGenModel ? false : preferences.stream}
              disabled={isImageGenModel}
              onCheckedChange={(v) => setPreference("stream", v)}
            />
            <Label htmlFor="chat-stream">
              Stream responses{isImageGenModel ? " (not available for text-image)" : ""}
            </Label>
          </div>
        </CardContent>
        )}
      </Card>

      <div className="grid min-h-0 flex-1 gap-4 lg:grid-cols-[minmax(240px,320px)_1fr]">
        <ChatPinnedPromptsPanel
          pins={pinsQ.data}
          isLoading={pinsQ.isLoading}
          runDisabled={composerDisabled}
          onRun={(body) => sendMut.mutate({ text: body, images: [] })}
          onUnpin={(id) => unpinMut.mutate(id)}
        />
        <Card className="flex min-h-0 flex-col">
          <CardHeader className="flex flex-row items-center justify-between space-y-0 pb-2">
            <CardTitle className="text-base">Messages</CardTitle>
            <Button
              type="button"
              variant="ghost"
              size="sm"
              disabled={chatSend.messages.length === 0 && !chatSend.sending}
              onClick={() => chatSend.clearChat()}
            >
              Clear chat
            </Button>
          </CardHeader>
          <CardContent className="flex min-h-0 flex-1 flex-col gap-3">
            <div className="min-h-0 flex-1 overflow-y-auto pr-1">
              <ChatMessageList
                messages={chatSend.messages}
                streamingContent={chatSend.streamingContent}
                sending={chatSend.sending}
                onPinUserMessage={(content) => pinMut.mutate(content)}
              />
            </div>
            <ChatComposer
              disabled={composerDisabled}
              allowImageAttach={(selectedModel?.modality || "text-text") === "text_image-text"}
              maxImageFileBytes={maxImageFileBytes}
              onSend={(text, images) => sendMut.mutateAsync({ text, images })}
            />
          </CardContent>
        </Card>
      </div>
    </PageShell>
  );
}
