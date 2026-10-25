import { Copy } from "lucide-react";
import { Button } from "@/components/ui/button";
import { HoverTooltip } from "@/shared/components/HoverTooltip";

const selectClassName =
  "h-9 min-w-[9rem] max-w-[14rem] rounded-md border border-border bg-card px-3 text-sm text-foreground";

type Props = {
  presets: { id: string; name: string }[];
  activePresetId: string;
  onSwitchPreset: (id: string) => void;
  onCopyLayout: () => void;
};

export function DashboardLayoutHeaderControls({
  presets,
  activePresetId,
  onSwitchPreset,
  onCopyLayout,
}: Props) {
  return (
    <>
      <label className="flex items-center gap-2 text-xs text-muted-foreground">
        Layout
        <select
          className={selectClassName}
          value={activePresetId}
          onChange={(e) => onSwitchPreset(e.target.value)}
        >
          {presets.map((p) => (
            <option key={p.id} value={p.id}>{p.name}</option>
          ))}
        </select>
      </label>
      <Button type="button" size="sm" variant="outline" onClick={onCopyLayout}>
        <HoverTooltip label="Copy layout" side="bottom" className="me-1.5 inline-flex">
          <Copy className="h-3.5 w-3.5" aria-hidden />
        </HoverTooltip>
        Copy layout
      </Button>
    </>
  );
}
