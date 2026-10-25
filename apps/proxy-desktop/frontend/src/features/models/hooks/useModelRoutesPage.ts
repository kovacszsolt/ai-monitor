import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useMemo, useState } from "react";
import { toast } from "sonner";
import {
  type ModelSortColumn,
  sortModelRoutes,
  TEXT_SORT_COLUMNS,
  priceInputValue,
} from "@/features/models/components/modelRouteSort";
import { createModelRoute, listModelRoutes, setModelRouteEnabled, updateModelRoute } from "@/features/models/api/modelsApi";
import { listProviders } from "@/features/providers/api/providersApi";
import { service } from "../../../../wailsjs/go/models";

export function useModelRoutesPage() {
  const qc = useQueryClient();
  const { data, isLoading, error, refetch } = useQuery({ queryKey: ["models"], queryFn: listModelRoutes });
  const { data: providers } = useQuery({ queryKey: ["providers"], queryFn: listProviders });
  const [createOpen, setCreateOpen] = useState(false);
  const [edit, setEdit] = useState<service.ModelRouteDTO | null>(null);
  const [publicName, setPublicName] = useState("");
  const [providerID, setProviderID] = useState("");
  const [upstream, setUpstream] = useState("");
  const [protocol, setProtocol] = useState<string>("openai_chat");
  const [inputUSD, setInputUSD] = useState("");
  const [cachedUSD, setCachedUSD] = useState("");
  const [outputUSD, setOutputUSD] = useState("");
  const [nameQuery, setNameQuery] = useState("");
  const [providerFilter, setProviderFilter] = useState("");
  const [sortBy, setSortBy] = useState<ModelSortColumn>("publicName");
  const [sortDesc, setSortDesc] = useState(false);

  const filteredRoutes = useMemo(() => {
    if (!data) {
      return [];
    }
    const q = nameQuery.trim().toLowerCase();
    return data.filter((m) => {
      if (providerFilter && m.provider_id !== providerFilter) {
        return false;
      }
      if (!q) {
        return true;
      }
      return (
        m.public_model_name.toLowerCase().includes(q) ||
        m.upstream_model_name.toLowerCase().includes(q)
      );
    });
  }, [data, nameQuery, providerFilter]);

  const displayRoutes = useMemo(
    () => sortModelRoutes(filteredRoutes, sortBy, sortDesc),
    [filteredRoutes, sortBy, sortDesc],
  );

  function sortDirection(column: ModelSortColumn): "asc" | "desc" | false {
    if (sortBy !== column) {
      return false;
    }
    return sortDesc ? "desc" : "asc";
  }

  function toggleSort(column: ModelSortColumn) {
    if (sortBy === column) {
      setSortDesc((prev) => !prev);
      return;
    }
    setSortBy(column);
    setSortDesc(!TEXT_SORT_COLUMNS.has(column));
  }

  const createMut = useMutation({
    mutationFn: () =>
      createModelRoute(
        publicName.trim(),
        providerID,
        upstream.trim(),
        protocol,
        inputUSD.trim(),
        cachedUSD.trim(),
        outputUSD.trim(),
      ),
    onSuccess: () => {
      toast.success("Model route created");
      setCreateOpen(false);
      qc.invalidateQueries({ queryKey: ["models"] });
    },
    onError: (e: Error) => toast.error(e.message),
  });

  const updateMut = useMutation({
    mutationFn: () =>
      updateModelRoute(
        edit!.id,
        upstream.trim(),
        protocol,
        inputUSD.trim(),
        cachedUSD.trim(),
        outputUSD.trim(),
      ),
    onSuccess: () => {
      toast.success("Model route updated");
      setEdit(null);
      qc.invalidateQueries({ queryKey: ["models"] });
    },
    onError: (e: Error) => toast.error(e.message),
  });

  const toggleMut = useMutation({
    mutationFn: ({ id, enabled }: { id: string; enabled: boolean }) => setModelRouteEnabled(id, enabled),
    onSuccess: () => qc.invalidateQueries({ queryKey: ["models"] }),
    onError: (e: Error) => toast.error(e.message),
  });

  function openEdit(m: service.ModelRouteDTO) {
    setEdit(m);
    setUpstream(m.upstream_model_name);
    setProtocol(m.protocol);
    setInputUSD(priceInputValue(m.input_token_usd_per_1m));
    setCachedUSD(priceInputValue(m.cached_input_token_usd_per_1m));
    setOutputUSD(priceInputValue(m.output_token_usd_per_1m));
  }

  function resetPricingFields() {
    setInputUSD("");
    setCachedUSD("");
    setOutputUSD("");
  }

  function openCreateDialog() {
    setCreateOpen(true);
    setPublicName("");
    setUpstream("");
    resetPricingFields();
    setProviderID(providers?.[0]?.id ?? "");
  }

  return {
    data,
    isLoading,
    error,
    refetch,
    providers,
    createOpen,
    setCreateOpen,
    edit,
    setEdit,
    publicName,
    setPublicName,
    providerID,
    setProviderID,
    upstream,
    setUpstream,
    protocol,
    setProtocol,
    inputUSD,
    setInputUSD,
    cachedUSD,
    setCachedUSD,
    outputUSD,
    setOutputUSD,
    nameQuery,
    setNameQuery,
    providerFilter,
    setProviderFilter,
    displayRoutes,
    sortDirection,
    toggleSort,
    createMut,
    updateMut,
    toggleMut,
    openEdit,
    openCreateDialog,
  };
}
