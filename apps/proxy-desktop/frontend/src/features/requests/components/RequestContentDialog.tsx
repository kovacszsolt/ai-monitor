import {
  Dialog,
  DialogContent,
  DialogDescription,
  DialogHeader,
  DialogTitle,
} from "@/components/ui/dialog";
import { RequestContentPanel } from "@/features/requests/components/RequestContentPanel";
import { formatDateTime } from "@/lib/format";
import { service } from "../../../../wailsjs/go/models";

type RequestSummary = Pick<
  service.AIRequestDTO,
  "id" | "requested_model" | "started_at" | "project_name" | "project_id"
>;

export function RequestContentDialog({
  request,
  open,
  onOpenChange,
}: {
  request: RequestSummary | null;
  open: boolean;
  onOpenChange: (open: boolean) => void;
}) {
  if (!request) {
    return null;
  }

  const projectLabel = request.project_name || request.project_id;

  return (
    <Dialog open={open} onOpenChange={onOpenChange}>
      <DialogContent className="max-h-[90vh] max-w-3xl overflow-y-auto">
        <DialogHeader>
          <DialogTitle>Prompt &amp; response</DialogTitle>
          <DialogDescription>
            {projectLabel} · {request.requested_model} · {formatDateTime(request.started_at)}
          </DialogDescription>
        </DialogHeader>
        <RequestContentPanel requestId={request.id} contentClassName="max-h-[55vh]" />
      </DialogContent>
    </Dialog>
  );
}
