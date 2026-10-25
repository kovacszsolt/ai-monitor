import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useEffect, useMemo, useState } from "react";
import { toast } from "sonner";
import { listProjects } from "@/features/projects/api/projectsApi";
import {
  bulkDeleteAIRequests,
  bulkSetAIRequestTag,
  deleteAIRequest,
  listAIRequestModelOptions,
  listAIRequestTagOptions,
  listAIRequests,
  type RequestSortColumn,
} from "@/features/requests/api/requestsApi";
import { formatNumber } from "@/lib/format";
import { service } from "../../../../wailsjs/go/models";

export const REQUEST_PAGE_SIZE = 50;

export type RequestListData = {
  total: number;
  items: service.AIRequestDTO[];
  page: number;
  page_size: number;
};

export function useRequestsList() {
  const qc = useQueryClient();
  const [page, setPage] = useState(1);
  const [projectID, setProjectID] = useState("");
  const [state, setState] = useState("");
  const [tag, setTag] = useState("");
  const [model, setModel] = useState("");
  const [expandedId, setExpandedId] = useState<string | null>(null);
  const [deletingId, setDeletingId] = useState<string | null>(null);
  const [selectedIds, setSelectedIds] = useState<Set<string>>(() => new Set());
  const [moveTag, setMoveTag] = useState("");
  const [sortBy, setSortBy] = useState<RequestSortColumn>("started");
  const [sortDesc, setSortDesc] = useState(true);
  const [contentRequest, setContentRequest] = useState<service.AIRequestDTO | null>(null);

  const projectsQ = useQuery({ queryKey: ["projects"], queryFn: listProjects });
  const tagOptionsQ = useQuery({ queryKey: ["request-tag-options"], queryFn: listAIRequestTagOptions });
  const tagBulkEnabled = tagOptionsQ.data?.enabled ?? false;
  const modelOptionsQ = useQuery({ queryKey: ["request-model-options"], queryFn: listAIRequestModelOptions });

  const listQ = useQuery({
    queryKey: ["requests", page, projectID, state, tag, model, sortBy, sortDesc],
    queryFn: () =>
      listAIRequests({
        page,
        pageSize: REQUEST_PAGE_SIZE,
        projectID,
        state,
        tag,
        model,
        sortBy,
        sortDesc,
      }),
  });

  const listQueryKey = useMemo(
    () => ["requests", page, projectID, state, tag, model, sortBy, sortDesc] as const,
    [page, projectID, state, tag, model, sortBy, sortDesc],
  );

  function sortDirection(column: RequestSortColumn): "asc" | "desc" | false {
    if (sortBy !== column) {
      return false;
    }
    return sortDesc ? "desc" : "asc";
  }

  function toggleSort(column: RequestSortColumn) {
    setPage(1);
    if (sortBy === column) {
      setSortDesc((prev) => !prev);
      return;
    }
    setSortBy(column);
    setSortDesc(column !== "project" && column !== "model" && column !== "tag" && column !== "state");
  }

  const deleteMut = useMutation({
    mutationFn: deleteAIRequest,
    onMutate: async (id) => {
      setDeletingId(id);
      await qc.cancelQueries({ queryKey: listQueryKey });
      const prev = qc.getQueryData<RequestListData>(listQueryKey);
      if (prev) {
        qc.setQueryData<RequestListData>(listQueryKey, {
          ...prev,
          total: Math.max(0, prev.total - 1),
          items: prev.items.filter((item) => item.id !== id),
        });
      }
      if (expandedId === id) {
        setExpandedId(null);
      }
      return { prev };
    },
    onSuccess: () => {
      toast.success("Request deleted");
      void qc.invalidateQueries({ queryKey: ["requests"] });
      qc.invalidateQueries({ queryKey: ["dashboard"] });
      qc.invalidateQueries({ queryKey: ["request-model-options"] });
      qc.invalidateQueries({ queryKey: ["request-tag-options"] });
      qc.invalidateQueries({ queryKey: ["model-usage-stats"] });
    },
    onError: (e: Error, _id, ctx) => {
      if (ctx?.prev) {
        qc.setQueryData(listQueryKey, ctx.prev);
      }
      toast.error(e.message);
    },
    onSettled: () => {
      setDeletingId(null);
    },
  });

  const totalPages = useMemo(() => {
    if (!listQ.data) return 1;
    return Math.max(1, Math.ceil(listQ.data.total / listQ.data.page_size));
  }, [listQ.data]);

  const pageItemIds = useMemo(
    () => listQ.data?.items.map((item) => item.id) ?? [],
    [listQ.data?.items],
  );
  const allPageSelected =
    pageItemIds.length > 0 && pageItemIds.every((id) => selectedIds.has(id));

  useEffect(() => {
    setSelectedIds(new Set());
  }, [projectID, state, tag, model]);

  const moveTagMut = useMutation({
    mutationFn: () => bulkSetAIRequestTag(Array.from(selectedIds), moveTag),
    onSuccess: (result) => {
      toast.success(`Tag updated on ${formatNumber(result.updated)} request(s)`);
      setSelectedIds(new Set());
      setMoveTag("");
      void qc.invalidateQueries({ queryKey: ["requests"] });
      qc.invalidateQueries({ queryKey: ["dashboard"] });
      qc.invalidateQueries({ queryKey: ["request-tag-options"] });
      qc.invalidateQueries({ queryKey: ["model-usage-stats"] });
    },
    onError: (e: Error) => toast.error(e.message),
  });

  const bulkDeleteMut = useMutation({
    mutationFn: () => bulkDeleteAIRequests(Array.from(selectedIds)),
    onMutate: async () => {
      const ids = new Set(selectedIds);
      await qc.cancelQueries({ queryKey: listQueryKey });
      const prev = qc.getQueryData<RequestListData>(listQueryKey);
      if (prev) {
        const removed = prev.items.filter((item) => ids.has(item.id)).length;
        qc.setQueryData<RequestListData>(listQueryKey, {
          ...prev,
          total: Math.max(0, prev.total - removed),
          items: prev.items.filter((item) => !ids.has(item.id)),
        });
      }
      if (expandedId && ids.has(expandedId)) {
        setExpandedId(null);
      }
      return { prev };
    },
    onSuccess: (result) => {
      toast.success(`Deleted ${formatNumber(result.deleted)} request(s)`);
      setSelectedIds(new Set());
      void qc.invalidateQueries({ queryKey: ["requests"] });
      qc.invalidateQueries({ queryKey: ["dashboard"] });
      qc.invalidateQueries({ queryKey: ["request-model-options"] });
      qc.invalidateQueries({ queryKey: ["request-tag-options"] });
      qc.invalidateQueries({ queryKey: ["model-usage-stats"] });
    },
    onError: (e: Error, _v, ctx) => {
      if (ctx?.prev) {
        qc.setQueryData(listQueryKey, ctx.prev);
      }
      toast.error(e.message);
    },
  });

  function toggleRowSelected(id: string) {
    setSelectedIds((prev) => {
      const next = new Set(prev);
      if (next.has(id)) {
        next.delete(id);
      } else {
        next.add(id);
      }
      return next;
    });
  }

  function toggleSelectAllOnPage() {
    setSelectedIds((prev) => {
      const next = new Set(prev);
      if (allPageSelected) {
        for (const id of pageItemIds) {
          next.delete(id);
        }
      } else {
        for (const id of pageItemIds) {
          next.add(id);
        }
      }
      return next;
    });
  }

  return {
    page,
    setPage,
    projectID,
    setProjectID,
    state,
    setState,
    tag,
    setTag,
    model,
    setModel,
    expandedId,
    setExpandedId,
    deletingId,
    selectedIds,
    setSelectedIds,
    moveTag,
    setMoveTag,
    sortBy,
    sortDesc,
    contentRequest,
    setContentRequest,
    projectsQ,
    tagOptionsQ,
    tagBulkEnabled,
    modelOptionsQ,
    listQ,
    sortDirection,
    toggleSort,
    deleteMut,
    totalPages,
    pageItemIds,
    allPageSelected,
    moveTagMut,
    bulkDeleteMut,
    toggleRowSelected,
    toggleSelectAllOnPage,
  };
}
