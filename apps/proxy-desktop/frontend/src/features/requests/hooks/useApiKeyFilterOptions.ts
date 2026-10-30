import { useQuery } from "@tanstack/react-query";
import { useMemo } from "react";
import { listAPIKeys } from "@/features/projects/api/keysApi";
import { listAIRequestAPIKeyOptions } from "@/features/requests/api/requestsApi";

export type ApiKeyFilterOption = {
  id: string;
  name: string;
};

/** Project selected → all keys in project; otherwise keys seen in request logs. */
export function useApiKeyFilterOptions(projectID: string) {
  const fromProjectQ = useQuery({
    queryKey: ["keys", projectID],
    queryFn: () => listAPIKeys(projectID),
    enabled: projectID !== "",
  });
  const fromRequestsQ = useQuery({
    queryKey: ["request-api-key-options", projectID],
    queryFn: () => listAIRequestAPIKeyOptions(projectID),
    enabled: projectID === "",
  });

  const options: ApiKeyFilterOption[] = useMemo(() => {
    if (projectID !== "") {
      return (fromProjectQ.data ?? []).map((k) => ({
        id: k.id,
        name: k.name?.trim() ? k.name : k.id,
      }));
    }
    return fromRequestsQ.data ?? [];
  }, [projectID, fromProjectQ.data, fromRequestsQ.data]);

  const isLoading = projectID !== "" ? fromProjectQ.isLoading : fromRequestsQ.isLoading;
  const error = projectID !== "" ? fromProjectQ.error : fromRequestsQ.error;

  return { options, isLoading, error };
}
