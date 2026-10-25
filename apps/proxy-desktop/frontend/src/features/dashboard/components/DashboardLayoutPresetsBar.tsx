import { useEffect, useState } from "react";
import { Button } from "@/components/ui/button";
import { Input } from "@/components/ui/input";

type Props = {
  presets: { id: string; name: string }[];
  activePresetId: string;
  onDeletePreset: (id: string) => void;
  onRenamePreset: (id: string, name: string) => void;
};

export function DashboardLayoutPresetsBar({
  presets,
  activePresetId,
  onDeletePreset,
  onRenamePreset,
}: Props) {
  const activeName = presets.find((p) => p.id === activePresetId)?.name ?? "";
  const [layoutNameDraft, setLayoutNameDraft] = useState(activeName);

  useEffect(() => {
    setLayoutNameDraft(activeName);
  }, [activePresetId, activeName]);

  const saveLayoutName = () => {
    const trimmed = layoutNameDraft.trim();
    if (!trimmed || trimmed === activeName) return;
    onRenamePreset(activePresetId, trimmed);
  };

  return (
    <div className="flex flex-wrap items-end gap-2 rounded-lg border border-border bg-muted/20 p-3">
      <label className="flex flex-col gap-1 text-xs text-muted-foreground">
        Layout name
        <div className="flex flex-wrap items-center gap-1">
          <Input
            className="h-9 w-[min(100%,12rem)] text-sm"
            value={layoutNameDraft}
            onChange={(e) => setLayoutNameDraft(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter") saveLayoutName();
            }}
          />
          <Button
            type="button"
            size="sm"
            variant="default"
            disabled={!layoutNameDraft.trim() || layoutNameDraft.trim() === activeName}
            onClick={saveLayoutName}
          >
            Save
          </Button>
        </div>
      </label>
      <Button
        type="button"
        size="sm"
        variant="ghost"
        className="text-destructive hover:text-destructive"
        disabled={presets.length <= 1}
        onClick={() => onDeletePreset(activePresetId)}
      >
        Delete
      </Button>
    </div>
  );
}
