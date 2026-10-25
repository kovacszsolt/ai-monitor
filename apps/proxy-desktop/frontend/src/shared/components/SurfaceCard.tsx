import type { ReactNode } from "react";
import { cn } from "@/lib/utils";

type Props = {
  children: ReactNode;
  className?: string;
  hover?: boolean;
};

export function SurfaceCard({ children, className, hover }: Props) {
  return (
    <div
      className={cn(
        "rounded-xl border border-surface-container-highest/30 bg-surface-container-low p-space-lg shadow-sm",
        hover && "transition-colors hover:bg-surface-container",
        className,
      )}
    >
      {children}
    </div>
  );
}
