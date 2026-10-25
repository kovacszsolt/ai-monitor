import { cn } from "@/lib/utils";

type Props = {
  active?: boolean;
  variant?: "tertiary" | "error";
  size?: "sm" | "md";
  className?: string;
};

export function StatusPulse({ active = true, variant = "tertiary", size = "md", className }: Props) {
  const dot = size === "sm" ? "h-1.5 w-1.5" : "h-2 w-2";
  const color = variant === "tertiary" ? "bg-tertiary" : "bg-error";
  return (
    <span className={cn("relative flex", dot, className)}>
      {active && (
        <span
          className={cn("absolute inline-flex h-full w-full animate-ping rounded-full opacity-75", color)}
        />
      )}
      <span className={cn("relative inline-flex rounded-full", dot, color)} />
    </span>
  );
}
