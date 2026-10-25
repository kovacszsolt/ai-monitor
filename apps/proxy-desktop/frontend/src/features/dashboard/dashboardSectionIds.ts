export const DASHBOARD_SECTION_IDS = {
  usageByHour: "usage-by-hour",
  usageByModel: "usage-by-model",
  statTotalRequests: "stat-total-requests",
  statRequests24h: "stat-requests-24h",
  statTokens: "stat-tokens",
  statAvgDuration: "stat-avg-duration",
  usageByProxyTag: "usage-by-proxy-tag",
  tokenMixChart: "token-mix-chart",
  byState: "by-state",
  topModels: "top-models",
  topProjects: "top-projects",
} as const;

/** Removed composite sections — expanded when loading saved layouts. */
export const LEGACY_DASHBOARD_SECTION_IDS = {
  summaryStats: "summary-stats",
  breakdown: "breakdown",
} as const;
