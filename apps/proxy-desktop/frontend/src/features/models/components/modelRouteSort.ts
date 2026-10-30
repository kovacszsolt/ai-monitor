import { service } from "../../../../wailsjs/go/models";

export type ModelSortColumn =
  | "publicName"
  | "provider"
  | "upstream"
  | "protocol"
  | "modality"
  | "inputUSD"
  | "cachedUSD"
  | "outputUSD"
  | "enabled";

export const TEXT_SORT_COLUMNS = new Set<ModelSortColumn>([
  "publicName",
  "provider",
  "upstream",
  "protocol",
  "modality",
]);

function compareNullableNumber(a?: number, b?: number): number {
  const aMissing = a === undefined || a === null || Number.isNaN(a);
  const bMissing = b === undefined || b === null || Number.isNaN(b);
  if (aMissing && bMissing) {
    return 0;
  }
  if (aMissing) {
    return 1;
  }
  if (bMissing) {
    return -1;
  }
  return a - b;
}

export function sortModelRoutes(
  rows: service.ModelRouteDTO[],
  sortBy: ModelSortColumn,
  sortDesc: boolean,
): service.ModelRouteDTO[] {
  const dir = sortDesc ? -1 : 1;
  return [...rows].sort((a, b) => {
    let cmp = 0;
    switch (sortBy) {
      case "publicName":
        cmp = a.public_model_name.localeCompare(b.public_model_name, undefined, { sensitivity: "base" });
        break;
      case "provider":
        cmp = (a.provider_name || a.provider_id).localeCompare(b.provider_name || b.provider_id, undefined, {
          sensitivity: "base",
        });
        break;
      case "upstream":
        cmp = a.upstream_model_name.localeCompare(b.upstream_model_name, undefined, { sensitivity: "base" });
        break;
      case "protocol":
        cmp = a.protocol.localeCompare(b.protocol);
        break;
      case "modality":
        cmp = (a.modality || "text-text").localeCompare(b.modality || "text-text");
        break;
      case "inputUSD":
        cmp = compareNullableNumber(a.input_token_usd_per_1m, b.input_token_usd_per_1m);
        break;
      case "cachedUSD":
        cmp = compareNullableNumber(a.cached_input_token_usd_per_1m, b.cached_input_token_usd_per_1m);
        break;
      case "outputUSD":
        cmp = compareNullableNumber(a.output_token_usd_per_1m, b.output_token_usd_per_1m);
        break;
      case "enabled":
        cmp = Number(a.enabled) - Number(b.enabled);
        break;
    }
    return cmp * dir;
  });
}

export function formatUSDPer1M(v?: number) {
  if (v === undefined || v === null || Number.isNaN(v)) {
    return "—";
  }
  return String(v);
}

export function priceInputValue(v?: number) {
  if (v === undefined || v === null || Number.isNaN(v)) {
    return "";
  }
  return String(v);
}
