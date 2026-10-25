export function formatNumber(n: number): string {
  return new Intl.NumberFormat().format(n);
}

/** Request list: show token count when captured (has flag) or usage_status is ok. */
export function formatRequestTokenCount(has: boolean, value: number, usageStatus: string): string {
  if (has || usageStatus === "ok") {
    return formatNumber(value ?? 0);
  }
  return "—";
}

export function formatUSD(has: boolean, amount: number): string {
  if (!has || !Number.isFinite(amount)) return "—";
  return new Intl.NumberFormat(undefined, {
    style: "currency",
    currency: "USD",
    minimumFractionDigits: 2,
    maximumFractionDigits: amount < 0.01 && amount > 0 ? 4 : 2,
  }).format(amount);
}

export function formatDuration(ms: number): string {
  if (ms < 1000) return `${ms} ms`;
  return `${(ms / 1000).toFixed(2)} s`;
}

export function formatTokensPerSec(v: number): string {
  if (!Number.isFinite(v) || v <= 0) return "—";
  if (v >= 100) return `${v.toFixed(0)} tok/s`;
  if (v >= 10) return `${v.toFixed(1)} tok/s`;
  return `${v.toFixed(2)} tok/s`;
}

export function formatOptionalTokensPerSec(has: boolean, v: number, estimated?: boolean): string {
  if (!has) return "—";
  const base = formatTokensPerSec(v);
  if (base === "—") return base;
  return estimated ? `~${base}` : base;
}

export function formatCacheHitRate(has: boolean, pct: number): string {
  if (!has || !Number.isFinite(pct)) return "—";
  return `${pct.toFixed(1)}%`;
}

export function formatDateTime(iso: string): string {
  if (!iso) return "—";
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return iso;
  return d.toLocaleString();
}

/** Human-readable build stamp for About and detail views. */
export function formatBuildTimeLabel(iso?: string): string {
  if (!iso) return "—";
  if (iso === "dev") return "Development build";
  const d = new Date(iso);
  if (Number.isNaN(d.getTime())) return iso;
  return d.toLocaleString(undefined, { dateStyle: "medium", timeStyle: "short" });
}

/** Desktop app link-time build stamp (UTC ISO from Go). */
export function formatBuildTime(iso?: string): string {
  if (!iso) return "Build: …";
  if (iso === "dev") return "Build: dev";
  const label = formatBuildTimeLabel(iso);
  return label === "—" ? "Build: …" : `Build: ${label}`;
}
