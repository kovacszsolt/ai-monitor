export type ByteUnit = "B" | "KiB" | "MiB" | "GiB";

export const BYTE_UNIT_OPTIONS: { value: ByteUnit; label: string; factor: number }[] = [
  { value: "B", label: "B", factor: 1 },
  { value: "KiB", label: "KiB", factor: 1024 },
  { value: "MiB", label: "MiB", factor: 1024 ** 2 },
  { value: "GiB", label: "GiB", factor: 1024 ** 3 },
];

export function byteUnitFactor(unit: ByteUnit): number {
  return BYTE_UNIT_OPTIONS.find((o) => o.value === unit)?.factor ?? 1;
}

/** Pick a readable unit for a byte count (binary IEC). */
export function pickByteDisplayUnit(bytes: number): ByteUnit {
  if (!Number.isFinite(bytes) || bytes <= 0) return "MiB";
  if (bytes >= 1024 ** 3 && bytes % 1024 ** 3 === 0) return "GiB";
  if (bytes >= 1024 ** 2 && bytes % 1024 ** 2 === 0) return "MiB";
  if (bytes >= 1024 && bytes % 1024 === 0) return "KiB";
  if (bytes >= 1024 ** 2) return "MiB";
  if (bytes >= 1024) return "KiB";
  return "B";
}

export function bytesToUnitAmount(bytes: number, unit: ByteUnit): number {
  const factor = byteUnitFactor(unit);
  if (unit === "B") return bytes;
  const n = bytes / factor;
  const rounded = Math.round(n * 1000) / 1000;
  return Number.isInteger(rounded) ? Math.trunc(rounded) : rounded;
}

export function unitAmountToBytes(amount: number, unit: ByteUnit): number {
  if (!Number.isFinite(amount) || amount < 0) return 0;
  if (unit === "B") return Math.round(amount);
  return Math.round(amount * byteUnitFactor(unit));
}

export function formatBytesExact(bytes: number): string {
  if (!Number.isFinite(bytes)) return "—";
  return `${bytes.toLocaleString()} bytes`;
}
