import type { ReactNode } from "react";
import { Card, CardContent, CardHeader } from "@/components/ui/card";
import { CollapsibleCardHeader } from "@/shared/components/CollapsibleCardHeader";
import type { DashboardSectionId } from "@/features/dashboard/dashboardLayout";
import { useDashboardSectionCollapsed } from "@/features/dashboard/hooks/useDashboardSectionCollapsed";

export function DashboardCollapsibleSection({
  sectionId,
  title,
  children,
}: {
  sectionId: DashboardSectionId;
  title: string;
  children: ReactNode;
}) {
  const { collapsed, toggle } = useDashboardSectionCollapsed(sectionId);

  return (
    <Card className="h-full">
      <CardHeader className="pb-2">
        <CollapsibleCardHeader title={title} collapsed={collapsed} onToggle={toggle} />
      </CardHeader>
      {!collapsed && <CardContent>{children}</CardContent>}
    </Card>
  );
}
