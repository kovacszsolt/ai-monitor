import type { ReactNode } from "react";
import { cn } from "@/lib/utils";

type Props = {
  section: string;
  accent?: string;
  title: string;
  description?: ReactNode;
  actions?: ReactNode;
  className?: string;
};

export function PageHeader({ section, accent, title, description, actions, className }: Props) {
  return (
    <div className={cn("flex flex-col gap-space-md md:flex-row md:items-end md:justify-between", className)}>
      <div className="flex flex-col gap-1">
        <div className="flex items-center gap-space-xs font-label-sm text-label-sm uppercase tracking-wider text-on-surface-variant">
          <span>{section}</span>
          {accent && (
            <>
              <span className="text-outline-variant">/</span>
              <span className="text-primary">{accent}</span>
            </>
          )}
        </div>
        <h1 className="font-headline-lg text-headline-lg tracking-tight text-on-surface">{title}</h1>
        {description && <p className="text-body-md text-on-surface-variant">{description}</p>}
      </div>
      {actions && <div className="flex flex-wrap items-center gap-space-sm">{actions}</div>}
    </div>
  );
}
