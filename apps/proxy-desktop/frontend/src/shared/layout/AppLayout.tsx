import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { useState } from "react";
import { Outlet } from "react-router-dom";
import { toast } from "sonner";
import { Button } from "@/components/ui/button";
import { startProxy } from "@/features/proxy/api/proxyApi";
import { formatBuildTime } from "@/lib/format";
import { MaterialIcon } from "@/shared/components/MaterialIcon";
import { ConnectionBanner } from "@/shared/layout/ConnectionBanner";
import { SidebarNav } from "@/shared/layout/SidebarNav";
import { WailsTitleBarDrag } from "@/shared/layout/WailsTitleBarDrag";
import { cn } from "@/lib/utils";
import { GetConnectionInfo } from "../../../wailsjs/go/service/DesktopService";

const SIDEBAR_COLLAPSED_KEY = "proxy-desktop.sidebar-collapsed";

function readSidebarCollapsed(): boolean {
  try {
    return localStorage.getItem(SIDEBAR_COLLAPSED_KEY) === "1";
  } catch {
    return false;
  }
}

export function AppLayout() {
  const [sidebarCollapsed, setSidebarCollapsed] = useState(readSidebarCollapsed);
  const qc = useQueryClient();

  const connection = useQuery({
    queryKey: ["connection"],
    queryFn: () => GetConnectionInfo(),
    refetchInterval: 5000,
  });

  const startProxyMut = useMutation({
    mutationFn: startProxy,
    onSuccess: (result) => {
      const failed = result.unhealthy_providers?.length ?? 0;
      if (failed > 0) {
        toast.warning(`Proxy started; ${failed} provider(s) failed the health check`);
      } else {
        toast.success("Proxy started");
      }
      qc.invalidateQueries({ queryKey: ["connection"] });
      qc.invalidateQueries({ queryKey: ["proxy-status"] });
    },
    onError: (e: Error) => toast.error(e.message),
  });

  const proxyRunning = connection.data?.proxy_running ?? false;

  return (
    <div className="flex h-screen flex-col overflow-hidden bg-background">
      <header
        className={cn(
          "relative z-50 flex h-14 shrink-0 items-center border-b border-outline-variant/40 bg-surface-container-lowest/90 px-space-md shadow-[0_1px_8px_rgba(0,0,0,0.4)] backdrop-blur-md wails-app-titlebar",
        )}
      >
        <WailsTitleBarDrag className="absolute inset-0" aria-hidden />
        <div className="relative z-10 flex w-full min-w-0 items-center justify-between gap-space-md pointer-events-none">
          <div className="flex min-w-0 shrink-0 items-center gap-space-sm wails-app-titlebar-brand">
            <div className="flex items-baseline gap-space-xs">
              <span className="whitespace-nowrap font-headline-sm text-headline-sm text-on-surface">
                AI Proxy Manager
              </span>
              <span
                className="whitespace-nowrap font-label-sm text-label-sm text-on-surface-variant"
                title={connection.data?.build_time ? `Built ${connection.data.build_time} (UTC)` : undefined}
              >
                {formatBuildTime(connection.data?.build_time)}
              </span>
            </div>
          </div>
          <div className="min-w-0 flex-1" aria-hidden />
          <div className="flex shrink-0 items-center gap-space-sm pointer-events-auto wails-no-drag">
            {!proxyRunning && (
              <Button
                type="button"
                size="sm"
                variant="secondary"
                disabled={startProxyMut.isPending}
                onClick={() => startProxyMut.mutate()}
              >
                <MaterialIcon name="play_circle" size={14} />
                Start proxy
              </Button>
            )}
            <ConnectionBanner info={connection.data} />
          </div>
        </div>
      </header>

      <div className="flex min-h-0 min-w-0 flex-1">
        <aside
          className={cn(
            "flex shrink-0 flex-col border-r border-outline-variant/40 bg-surface-container-lowest py-space-sm transition-[width] duration-200",
            sidebarCollapsed ? "w-14 px-1" : "w-60 px-space-xs",
          )}
          aria-expanded={!sidebarCollapsed}
        >
          <div className="flex-1 overflow-y-auto">
            <SidebarNav collapsed={sidebarCollapsed} />
          </div>
          <div className="mt-auto shrink-0 px-space-xs pt-space-xs">
            <Button
              type="button"
              variant="ghost"
              size="sm"
              className={cn(
                "w-full justify-between font-label-sm text-label-sm text-on-surface-variant",
                sidebarCollapsed && "justify-center px-0",
              )}
              aria-label={sidebarCollapsed ? "Expand menu" : "Collapse menu"}
              onClick={() => {
                setSidebarCollapsed((prev) => {
                  const next = !prev;
                  try {
                    localStorage.setItem(SIDEBAR_COLLAPSED_KEY, next ? "1" : "0");
                  } catch {
                    /* ignore */
                  }
                  return next;
                });
              }}
            >
              <span className="flex items-center gap-1.5">
                <MaterialIcon name="dock_to_right" size={16} />
                {!sidebarCollapsed && <span>Collapse Rail</span>}
              </span>
            </Button>
          </div>
        </aside>
        <main className="min-h-0 min-w-0 flex-1 overflow-y-auto bg-background">
          <Outlet />
        </main>
      </div>
    </div>
  );
}
