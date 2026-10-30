import { useEffect, useState } from "react";
import { Input } from "@/components/ui/input";
import {
  bytesToUnitAmount,
  pickByteDisplayUnit,
  unitAmountToBytes,
  type ByteUnit,
  BYTE_UNIT_OPTIONS,
  formatBytesExact,
} from "@/lib/byteSize";
import { cn } from "@/lib/utils";

const selectClass =
  "flex h-9 min-w-[5.5rem] rounded border border-outline-variant bg-surface-container px-space-md py-1 text-body-sm transition-colors focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-primary/40 disabled:cursor-not-allowed disabled:opacity-50";

type Props = {
  id?: string;
  value: number;
  onChange: (bytes: number) => void;
  disabled?: boolean;
  minBytes?: number;
  className?: string;
};

export function ByteSizeInput({ id, value, onChange, disabled, minBytes = 1, className }: Props) {
  const [unit, setUnit] = useState<ByteUnit>(() => pickByteDisplayUnit(value));
  const [amount, setAmount] = useState(() => bytesToUnitAmount(value, pickByteDisplayUnit(value)));

  useEffect(() => {
    const nextUnit = pickByteDisplayUnit(value);
    setUnit(nextUnit);
    setAmount(bytesToUnitAmount(value, nextUnit));
  }, [value]);

  function commitAmount(nextAmount: number, nextUnit: ByteUnit) {
    setAmount(nextAmount);
    const bytes = Math.max(minBytes, unitAmountToBytes(nextAmount, nextUnit));
    onChange(bytes);
  }

  function handleUnitChange(nextUnit: ByteUnit) {
    setUnit(nextUnit);
    setAmount(bytesToUnitAmount(value, nextUnit));
  }

  const step = unit === "B" ? 1 : unit === "KiB" ? 1 : 0.001;

  return (
    <div className={cn("space-y-1", className)}>
      <div className="flex flex-wrap gap-2">
        <Input
          id={id}
          type="number"
          min={0}
          step={step}
          disabled={disabled}
          value={Number.isFinite(amount) ? amount : ""}
          onChange={(e) => {
            const parsed = e.target.value === "" ? 0 : Number(e.target.value);
            commitAmount(parsed, unit);
          }}
          className="min-w-[8rem] flex-1 font-mono text-xs"
        />
        <select
          aria-label="Size unit"
          className={selectClass}
          disabled={disabled}
          value={unit}
          onChange={(e) => handleUnitChange(e.target.value as ByteUnit)}
        >
          {BYTE_UNIT_OPTIONS.map((o) => (
            <option key={o.value} value={o.value}>
              {o.label}
            </option>
          ))}
        </select>
      </div>
      <p className="font-mono text-body-sm text-on-surface-variant">{formatBytesExact(value)}</p>
    </div>
  );
}
