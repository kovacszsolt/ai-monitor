import { createContext, useContext, type ReactNode } from "react";

const DashboardExportContext = createContext(false);

export function DashboardExportProvider({
  forceExpanded,
  children,
}: {
  forceExpanded: boolean;
  children: ReactNode;
}) {
  return <DashboardExportContext.Provider value={forceExpanded}>{children}</DashboardExportContext.Provider>;
}

export function useDashboardExportForceExpanded(): boolean {
  return useContext(DashboardExportContext);
}
