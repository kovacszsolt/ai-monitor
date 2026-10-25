/** Stacked token segments — fixed colors for Recharts (SVG). */
export const TOKEN_STACK_SERIES = [
  { key: "input_tokens", label: "Input", fill: "#60a5fa" },
  { key: "cached_input_tokens", label: "Cached", fill: "#c084fc" },
  { key: "output_tokens", label: "Output", fill: "#4ade80" },
] as const;

export type TokenStackKey = (typeof TOKEN_STACK_SERIES)[number]["key"];

const TOKEN_FIELD_SEP = "::";

export function tokenStackDataKey(modelKey: string, segment: TokenStackKey): string {
  return `${modelKey}${TOKEN_FIELD_SEP}${segment}`;
}

export function parseTokenStackDataKey(
  dataKey: string,
): { modelKey: string; segment: TokenStackKey } | null {
  const sep = dataKey.indexOf(TOKEN_FIELD_SEP);
  if (sep <= 0) return null;
  const modelKey = dataKey.slice(0, sep);
  const segment = dataKey.slice(sep + TOKEN_FIELD_SEP.length) as TokenStackKey;
  if (!TOKEN_STACK_SERIES.some((s) => s.key === segment)) return null;
  return { modelKey, segment };
}

export function totalTokenCount(row: {
  input_tokens: number;
  cached_input_tokens: number;
  output_tokens: number;
}): number {
  return row.input_tokens + row.cached_input_tokens + row.output_tokens;
}

export function totalTokenCountFromPoint(p: {
  input_tokens?: number;
  cached_input_tokens?: number;
  output_tokens?: number;
}): number {
  return (p.input_tokens ?? 0) + (p.cached_input_tokens ?? 0) + (p.output_tokens ?? 0);
}
