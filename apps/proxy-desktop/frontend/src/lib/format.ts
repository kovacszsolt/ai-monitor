export function formatNumber(n: number): string {
  return new Intl.NumberFormat().format(n);
}

/** Y-axis ticks: full count below 1k, compact (k / M) at larger magnitudes. */
export function formatChartAxisCount(n: number): string {
  if (!Number.isFinite(n)) return "—";
  const abs = Math.abs(n);
  if (abs < 1000) {
    return new Intl.NumberFormat(undefined, { maximumFractionDigits: 0 }).format(n);
  }
  const digits = abs >= 1_000_000 ? 1 : abs >= 10_000 ? 0 : 1;
  return new Intl.NumberFormat(undefined, {
    notation: "compact",
    compactDisplay: "short",
    maximumFractionDigits: digits,
  }).format(n);
}

/** Y-axis ticks for USD amounts. */
export function formatChartAxisUSD(amount: number): string {
  if (!Number.isFinite(amount)) return "—";
  const abs = Math.abs(amount);
  if (abs < 1000) {
    return new Intl.NumberFormat(undefined, {
      style: "currency",
      currency: "USD",
      maximumFractionDigits: abs < 1 ? 2 : 0,
    }).format(amount);
  }
  return new Intl.NumberFormat(undefined, {
    style: "currency",
    currency: "USD",
    notation: "compact",
    compactDisplay: "short",
    maximumFractionDigits: 1,
  }).format(amount);
}

/** Y-axis ticks for tok/s (no unit suffix — legend/metric already names it). */
export function formatChartAxisTokPerSec(v: number): string {
  if (!Number.isFinite(v) || v < 0) return "—";
  if (v < 1000) {
    if (v >= 100) return `${v.toFixed(0)}`;
    if (v >= 10) return `${v.toFixed(1)}`;
    return `${v.toFixed(2)}`;
  }
  return new Intl.NumberFormat(undefined, {
    notation: "compact",
    compactDisplay: "short",
    maximumFractionDigits: 1,
  }).format(v);
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
  if (!Number.isFinite(ms) || ms < 0) return "—";
  if (ms < 1000) return `${ms} ms`;
  const s = ms / 1000;
  if (s < 60) return `${s.toFixed(2)} s`;
  const m = s / 60;
  if (m < 60) return `${m.toFixed(1)} min`;
  const h = m / 60;
  return `${h.toFixed(1)} h`;
}

export function formatOptionalDuration(has: boolean, ms: number): string {
  if (!has) return "—";
  return formatDuration(ms);
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
