import { MaterialIcon } from "@/shared/components/MaterialIcon";
import { StatusPulse } from "@/shared/components/StatusPulse";
import { service } from "../../../wailsjs/go/models";

type ConnectionInfo = service.ConnectionInfoDTO;

function basename(path: string): string {
  const parts = path.replace(/\\/g, "/").split("/");
  return parts[parts.length - 1] || path;
}

export function ConnectionBanner({ info }: { info?: ConnectionInfo | null }) {
  if (!info) return null;

  const ok = info.connected && info.schema_ready;
  const dbName = basename(info.database_path ?? "proxy.db");

  return (
    <div className="flex max-w-full flex-wrap items-center justify-end gap-space-sm px-space-md py-1 text-label-sm wails-no-drag">
      <div className="flex items-center gap-space-xs rounded bg-surface-container px-space-sm py-1">
        <StatusPulse active={info.proxy_running} variant={info.proxy_running ? "tertiary" : "error"} size="sm" />
        <span className="font-label-sm text-label-sm text-on-surface">
          {info.proxy_running ? "Proxy on" : "Proxy off"}
        </span>
      </div>
      <div className="flex items-center gap-1.5 rounded bg-surface-container-low px-space-sm py-1 font-label-sm text-label-sm text-on-surface-variant">
        <MaterialIcon name="database" size={14} className="text-tertiary" />
        <span>{dbName}</span>
        <span className="text-outline-variant">·</span>
        <span className={ok ? "text-tertiary" : "text-error"}>{ok ? "Connected" : "Database issue"}</span>
      </div>
      {!info.config_bridge_ok && info.bridge_message && (
        <span className="max-w-[200px] truncate text-label-sm text-secondary" title={info.bridge_message}>
          Config mismatch
        </span>
      )}
      {info.message && <span className="text-label-sm text-error">{info.message}</span>}
    </div>
  );
}
