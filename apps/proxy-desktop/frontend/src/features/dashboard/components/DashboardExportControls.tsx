import { useState } from "react";
import { Loader2 } from "lucide-react";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import type { DashboardLayoutUsagePrefs } from "@/features/dashboard/dashboardLayout";
import { buildExportFilterSummary } from "@/features/dashboard/export/buildExportFilterSummary";
import {
  exportDashboard,
  type DashboardExportFormat,
} from "@/features/dashboard/export/dashboardExport";

type Props = {
  exportRootRef: React.RefObject<HTMLElement | null>;
  layoutName: string;
  usagePrefs: DashboardLayoutUsagePrefs;
  projectLabel: string;
  apiKeyLabel: string;
  tagLabel: string;
  onExportStart: () => void;
  onExportEnd: () => void;
};

export function DashboardExportControls({
  exportRootRef,
  layoutName,
  usagePrefs,
  projectLabel,
  apiKeyLabel,
  tagLabel,
  onExportStart,
  onExportEnd,
}: Props) {
  const [busyFormat, setBusyFormat] = useState<DashboardExportFormat | null>(null);

  const runExport = async (format: DashboardExportFormat) => {
    const root = exportRootRef.current;
    if (!root) {
      toast.error("Nothing to export yet.");
      return;
    }
    setBusyFormat(format);
    onExportStart();
    try {
      const saved = await exportDashboard(root, format, {
        layoutName,
        exportedAtIso: new Date().toISOString(),
        filterSummary: buildExportFilterSummary(usagePrefs, projectLabel, apiKeyLabel, tagLabel),
      });
      if (saved === null) {
        toast.message("Export cancelled");
      } else {
        toast.success(`Saved ${format.toUpperCase()}`, {
          description: saved.length > 80 ? `…${saved.slice(-72)}` : saved,
        });
      }
    } catch (err) {
      const message = err instanceof Error ? err.message : "Export failed";
      console.error("[dashboard export]", err);
      toast.error(message, { duration: 8000 });
    } finally {
      setBusyFormat(null);
      onExportEnd();
    }
  };

  const disabled = busyFormat !== null;

  return (
    <div className="flex flex-wrap items-center gap-1">
      <span className="text-xs text-muted-foreground">Export</span>
      {(["png", "html", "pdf"] as const).map((format) => (
        <Button
          key={format}
          type="button"
          size="sm"
          variant="outline"
          disabled={disabled}
          onClick={() => runExport(format)}
        >
          {busyFormat === format ? <Loader2 className="h-3.5 w-3.5 animate-spin" /> : format.toUpperCase()}
        </Button>
      ))}
    </div>
  );
}
