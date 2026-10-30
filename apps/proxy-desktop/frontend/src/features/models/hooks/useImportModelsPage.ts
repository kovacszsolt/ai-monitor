import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useEffect, useMemo, useRef, useState } from "react";
import { useSearchParams } from "react-router-dom";
import { toast } from "sonner";
import {
  importModelRoutes,
  listProviderUpstreamModels,
  type UpstreamModelOption,
} from "@/features/models/api/modelsApi";
import { listProviders } from "@/features/providers/api/providersApi";

export function useImportModelsPage() {
  const qc = useQueryClient();
  const [searchParams] = useSearchParams();
  const initialProviderFromUrl = searchParams.get("provider") ?? "";

  const providersQ = useQuery({ queryKey: ["providers"], queryFn: listProviders });

  const [providerId, setProviderId] = useState("");
  const providerInitDone = useRef(false);

  useEffect(() => {
    if (providerInitDone.current || !providersQ.data?.length) {
      return;
    }
    providerInitDone.current = true;
    const fromUrl =
      initialProviderFromUrl &&
      providersQ.data.some((p) => p.id === initialProviderFromUrl)
        ? initialProviderFromUrl
        : "";
    setProviderId(fromUrl || providersQ.data[0].id);
  }, [providersQ.data, initialProviderFromUrl]);

  const [nameQuery, setNameQuery] = useState("");
  const [publicPrefix, setPublicPrefix] = useState("");
  const [models, setModels] = useState<UpstreamModelOption[] | null>(null);
  const [selected, setSelected] = useState<Set<string>>(() => new Set());
  const [fetchError, setFetchError] = useState<string | null>(null);

  const fetchMut = useMutation({
    mutationFn: () => listProviderUpstreamModels(providerId),
    onSuccess: (rows) => {
      setModels(rows);
      setFetchError(null);
      setSelected(new Set());
    },
    onError: (e: Error) => {
      setModels(null);
      setFetchError(e.message);
    },
  });

  const importMut = useMutation({
    mutationFn: () =>
      importModelRoutes(providerId, Array.from(selected), publicPrefix),
    onSuccess: (res) => {
      const skippedNote =
        res.skipped.length > 0 ? ` (${res.skipped.length} skipped)` : "";
      toast.success(`Imported ${res.created} model route(s)${skippedNote}`);
      setSelected(new Set());
      qc.invalidateQueries({ queryKey: ["models"] });
      void fetchMut.mutateAsync();
    },
    onError: (e: Error) => toast.error(e.message),
  });

  const filtered = useMemo(() => {
    if (!models) return [];
    const q = nameQuery.trim().toLowerCase();
    if (!q) return models;
    return models.filter((m) => m.id.toLowerCase().includes(q));
  }, [models, nameQuery]);

  const importableFiltered = useMemo(
    () => filtered.filter((m) => !m.alreadyImported),
    [filtered],
  );

  const firstSelectedPreview = useMemo(() => {
    const first = Array.from(selected)[0];
    if (!first) return null;
    return `${publicPrefix}${first}`;
  }, [selected, publicPrefix]);

  function onProviderChange(nextId: string) {
    setProviderId(nextId);
    setModels(null);
    setSelected(new Set());
    setFetchError(null);
  }

  function toggleId(id: string, alreadyImported: boolean) {
    if (alreadyImported) return;
    setSelected((prev) => {
      const next = new Set(prev);
      if (next.has(id)) {
        next.delete(id);
      } else {
        next.add(id);
      }
      return next;
    });
  }

  function selectAllImportable() {
    setSelected(new Set(importableFiltered.map((m) => m.id)));
  }

  function clearSelection() {
    setSelected(new Set());
  }

  const canFetch = !!providerId;
  const canImport = selected.size > 0 && !!providerId && !importMut.isPending;

  return {
    providersQ,
    providerId,
    onProviderChange,
    nameQuery,
    setNameQuery,
    publicPrefix,
    setPublicPrefix,
    models,
    selected,
    fetchError,
    fetchMut,
    importMut,
    filtered,
    importableFiltered,
    firstSelectedPreview,
    toggleId,
    selectAllImportable,
    clearSelection,
    canFetch,
    canImport,
  };
}
