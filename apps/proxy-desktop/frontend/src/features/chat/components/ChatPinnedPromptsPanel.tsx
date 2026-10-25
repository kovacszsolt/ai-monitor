import { Pin, Play, X } from "lucide-react";
import { Button } from "@/components/ui/button";
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from "@/components/ui/card";
import { LoadingState } from "@/shared/components/PageStates";
import { service } from "../../../../wailsjs/go/models";

type Props = {
  pins: service.ChatPinnedPromptDTO[] | undefined;
  isLoading: boolean;
  runDisabled: boolean;
  onRun: (body: string) => void;
  onUnpin: (id: string) => void;
};

export function ChatPinnedPromptsPanel({ pins, isLoading, runDisabled, onRun, onUnpin }: Props) {
  return (
    <Card className="flex min-h-0 flex-col lg:max-h-[calc(100vh-14rem)]">
      <CardHeader className="pb-2">
        <CardTitle className="flex items-center gap-2 text-base">
          <Pin className="h-4 w-4" />
          Pinned questions
        </CardTitle>
        <CardDescription>Re-run with another provider or model.</CardDescription>
      </CardHeader>
      <CardContent className="min-h-0 flex-1 overflow-y-auto">
        {isLoading && <LoadingState label="Loading pins…" />}
        {!isLoading && (!pins || pins.length === 0) && (
          <p className="text-sm text-muted-foreground">
            No pinned questions yet. Use Pin on a message you sent.
          </p>
        )}
        {!isLoading && pins && pins.length > 0 && (
          <ul className="space-y-2">
            {pins.map((p) => (
              <li
                key={p.id}
                className="rounded-md border border-border p-2 text-sm"
              >
                <p className="font-medium leading-snug">{p.title || "Pinned question"}</p>
                <p className="mt-1 line-clamp-3 text-muted-foreground whitespace-pre-wrap">{p.body}</p>
                <div className="mt-2 flex flex-wrap gap-1">
                  <Button
                    type="button"
                    size="sm"
                    variant="outline"
                    disabled={runDisabled}
                    onClick={() => onRun(p.body)}
                  >
                    <Play className="mr-1 h-3.5 w-3.5" />
                    Run
                  </Button>
                  <Button
                    type="button"
                    size="sm"
                    variant="ghost"
                    onClick={() => onUnpin(p.id)}
                  >
                    <X className="mr-1 h-3.5 w-3.5" />
                    Unpin
                  </Button>
                </div>
              </li>
            ))}
          </ul>
        )}
      </CardContent>
    </Card>
  );
}
