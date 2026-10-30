import { useQuery } from "@tanstack/react-query";
import { useEffect, useMemo } from "react";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { CollapsibleCardHeader } from "@/shared/components/CollapsibleCardHeader";
import { ModelTokenCompareChart } from "@/features/dashboard/components/ModelTokenCompareChart";
import { DASHBOARD_SECTION_IDS } from "@/features/dashboard/dashboardSectionIds";
import { useDashboardSectionCollapsed } from "@/features/dashboard/hooks/useDashboardSectionCollapsed";
import { ModelUsageTable } from "@/features/dashboard/components/ModelUsageTable";
import { listProjects } from "@/features/projects/api/projectsApi";
import { getModelUsageStats, listAIRequestTagOptions } from "@/features/requests/api/requestsApi";
import { useApiKeyFilterOptions } from "@/features/requests/hooks/useApiKeyFilterOptions";
import { ErrorState, LoadingState } from "@/shared/components/PageStates";

const selectClassName =
  "h-9 min-w-[160px] rounded-md border border-border bg-card px-3 text-sm text-foreground sm:max-w-xs";

type Props = {
  tag: string;
  onTagChange: (tag: string) => void;
  projectID: string;
  onProjectChange: (projectID: string) => void;
  apiKeyID: string;
  onApiKeyChange: (apiKeyID: string) => void;
  hideFilters?: boolean;
};

function useModelUsageFilters(tag: string, projectID: string, apiKeyID: string) {
  const tagOptionsQ = useQuery({ queryKey: ["request-tag-options"], queryFn: listAIRequestTagOptions });
  const projectsQ = useQuery({ queryKey: ["projects"], queryFn: listProjects });
  const { options: apiKeyOptions } = useApiKeyFilterOptions(projectID);

  const tagLabel =
    tag === ""
      ? "all tags"
      : tag === tagOptionsQ.data?.no_tag_label
        ? tagOptionsQ.data.no_tag_label
        : tag;

  const projectLabel =
    projectID === ""
      ? "all projects"
      : (projectsQ.data?.find((p) => p.id === projectID)?.name ?? projectID);

  const apiKeyLabel =
    apiKeyID === ""
      ? "all API keys"
      : (apiKeyOptions.find((k) => k.id === apiKeyID)?.name ?? apiKeyID);

  const filterParts = [
    projectID !== "" ? projectLabel : null,
    apiKeyID !== "" ? `API key ${apiKeyLabel}` : null,
    tag !== "" ? `tag ${tagLabel}` : null,
  ].filter(Boolean);

  const filterSummary = filterParts.length === 0 ? "all requests" : filterParts.join(", ");

  return { tagOptionsQ, projectsQ, apiKeyOptions, tagLabel, projectLabel, apiKeyLabel, filterSummary };
}

export function ModelUsageFilters({
  tag,
  onTagChange,
  projectID,
  onProjectChange,
  apiKeyID,
  onApiKeyChange,
  showTagFilter,
}: {
  tag: string;
  onTagChange: (tag: string) => void;
  projectID: string;
  onProjectChange: (projectID: string) => void;
  apiKeyID: string;
  onApiKeyChange: (apiKeyID: string) => void;
  showTagFilter: boolean;
}) {
  const { tagOptionsQ, projectsQ, apiKeyOptions } = useModelUsageFilters(tag, projectID, apiKeyID);
  const activeProjects = useMemo(
    () => (projectsQ.data ?? []).filter((p) => p.enabled),
    [projectsQ.data],
  );

  useEffect(() => {
    if (projectID && !activeProjects.some((p) => p.id === projectID)) {
      onProjectChange("");
    }
  }, [projectID, activeProjects, onProjectChange]);

  useEffect(() => {
    onApiKeyChange("");
  }, [projectID]);

  useEffect(() => {
    if (apiKeyID && !apiKeyOptions.some((k) => k.id === apiKeyID)) {
      onApiKeyChange("");
    }
  }, [apiKeyID, apiKeyOptions, onApiKeyChange]);

  return (
    <div className="flex flex-wrap items-end gap-3">
      <label className="flex flex-col gap-1 text-xs text-muted-foreground">
        Project
        <select
          className={selectClassName}
          value={projectID}
          onChange={(e) => onProjectChange(e.target.value)}
        >
          <option value="">All projects</option>
          {activeProjects.map((p) => (
            <option key={p.id} value={p.id}>{p.name}</option>
          ))}
        </select>
      </label>
      <label className="flex flex-col gap-1 text-xs text-muted-foreground">
        API key
        <select
          className={selectClassName}
          value={apiKeyID}
          onChange={(e) => onApiKeyChange(e.target.value)}
          disabled={projectID === ""}
          title={projectID === "" ? "Select a project to filter by API key" : undefined}
        >
          <option value="">{projectID === "" ? "Select project first" : "All API keys"}</option>
          {apiKeyOptions.map((k) => (
            <option key={k.id} value={k.id}>{k.name}</option>
          ))}
        </select>
      </label>
      {showTagFilter && (
        <label className="flex flex-col gap-1 text-xs text-muted-foreground">
          Proxy tag
          <select
            className={selectClassName}
            value={tag}
            onChange={(e) => onTagChange(e.target.value)}
          >
            <option value="">All tags</option>
            {tagOptionsQ.data?.has_untagged && (
              <option value={tagOptionsQ.data.no_tag_label}>{tagOptionsQ.data.no_tag_label}</option>
            )}
            {(tagOptionsQ.data?.tags ?? []).map((t) => (
              <option key={t} value={t}>{t}</option>
            ))}
          </select>
        </label>
      )}
    </div>
  );
}

export function ModelUsageByTagSection({
  tag,
  onTagChange,
  projectID,
  onProjectChange,
  apiKeyID,
  onApiKeyChange,
  hideFilters = false,
}: Props) {
  const { collapsed, toggle } = useDashboardSectionCollapsed(DASHBOARD_SECTION_IDS.usageByModel);
  const { tagOptionsQ, filterSummary } = useModelUsageFilters(tag, projectID, apiKeyID);
  const tagEnabled = tagOptionsQ.data?.enabled ?? false;
  const modelUsageQ = useQuery({
    queryKey: ["model-usage-stats", tag, projectID, apiKeyID],
    queryFn: () => getModelUsageStats(tag, projectID, apiKeyID),
    enabled: !collapsed,
    refetchInterval: 30000,
  });

  if (!tagEnabled) {
    return (
      <Card className="w-full">
        <CardHeader className="pb-2">
          <CollapsibleCardHeader title="Usage by model" collapsed={collapsed} onToggle={toggle} />
        </CardHeader>
        {!collapsed && (
        <CardContent className="space-y-4">
          {!hideFilters && (
            <ModelUsageFilters
              tag={tag}
              onTagChange={onTagChange}
              projectID={projectID}
              onProjectChange={onProjectChange}
              apiKeyID={apiKeyID}
              onApiKeyChange={onApiKeyChange}
              showTagFilter={false}
            />
          )}
          {modelUsageQ.isLoading && <LoadingState label="Loading model usage…" />}
          {!modelUsageQ.isLoading && (
            <ModelUsageTable
              embedded
              title="Usage by model"
              description={`Aggregated for ${filterSummary}; tag filter unavailable (run migration 000008)`}
              rows={modelUsageQ.data ?? []}
            />
          )}
        </CardContent>
        )}
      </Card>
    );
  }

  return (
    <Card className="w-full">
      <CardHeader className="pb-2">
        <CollapsibleCardHeader title="Usage by model" collapsed={collapsed} onToggle={toggle} />
      </CardHeader>
      {!collapsed && (
      <CardContent className="space-y-4">
        {!hideFilters && (
          <ModelUsageFilters
            tag={tag}
            onTagChange={onTagChange}
            projectID={projectID}
            onProjectChange={onProjectChange}
            apiKeyID={apiKeyID}
            onApiKeyChange={onApiKeyChange}
            showTagFilter
          />
        )}

        {modelUsageQ.isLoading && <LoadingState label="Loading model usage…" />}
        {!modelUsageQ.isLoading && (
          <ModelUsageTable
            embedded
            title="Usage by model"
            description={`Aggregated for ${filterSummary}; cost from Models pricing; tok/s from requests with duration`}
            rows={modelUsageQ.data ?? []}
          />
        )}
      </CardContent>
      )}
    </Card>
  );
}

export function ModelTokenUsageByTagChart({
  tag,
  tagLabel,
  projectID,
  apiKeyID,
}: {
  tag: string;
  tagLabel: string;
  projectID: string;
  apiKeyID: string;
}) {
  const { collapsed, toggle } = useDashboardSectionCollapsed(DASHBOARD_SECTION_IDS.tokenMixChart);
  const modelUsageQ = useQuery({
    queryKey: ["model-usage-stats", tag, projectID, apiKeyID],
    queryFn: () => getModelUsageStats(tag, projectID, apiKeyID),
    enabled: tag !== "" && !collapsed,
    refetchInterval: 30000,
  });

  if (tag === "") return null;

  return (
    <Card>
      <CardHeader className="pb-2">
        <CollapsibleCardHeader title="Token mix by model" collapsed={collapsed} onToggle={toggle} />
      </CardHeader>
      {!collapsed && (
        <CardContent className="pt-2">
          {modelUsageQ.isLoading && <LoadingState label="Loading token chart…" />}
          {modelUsageQ.isError && (
            <ErrorState
              message={(modelUsageQ.error as Error).message}
              onRetry={() => modelUsageQ.refetch()}
            />
          )}
          {!modelUsageQ.isLoading && !modelUsageQ.isError && (
            <ModelTokenCompareChart hideTitle tagLabel={tagLabel} rows={modelUsageQ.data ?? []} />
          )}
        </CardContent>
      )}
    </Card>
  );
}

export function useModelTagLabel(tag: string) {
  const tagOptionsQ = useQuery({ queryKey: ["request-tag-options"], queryFn: listAIRequestTagOptions });
  if (tag === "") return "all tags";
  if (tag === tagOptionsQ.data?.no_tag_label) return tagOptionsQ.data.no_tag_label;
  return tag;
}

export function useModelApiKeyLabel(apiKeyID: string, projectID: string) {
  const { options } = useApiKeyFilterOptions(projectID);
  if (apiKeyID === "") return "all API keys";
  return options.find((k) => k.id === apiKeyID)?.name ?? apiKeyID;
}
