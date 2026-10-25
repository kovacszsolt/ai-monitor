import { cn } from "@/lib/utils";

type Props = {
  name: string;
  size?: number;
  className?: string;
  filled?: boolean;
};

export function MaterialIcon({ name, size = 18, className, filled }: Props) {
  return (
    <span
      className={cn("material-symbols-outlined inline-flex shrink-0 items-center justify-center", className)}
      style={{
        fontSize: size,
        width: size,
        height: size,
        fontVariationSettings: filled ? '"FILL" 1, "wght" 400, "GRAD" 0, "opsz" 20' : undefined,
      }}
      aria-hidden
    >
      {name}
    </span>
  );
}
