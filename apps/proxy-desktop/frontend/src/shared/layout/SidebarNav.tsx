import { NavLink } from "react-router-dom";
import { MaterialIcon } from "@/shared/components/MaterialIcon";
import { HoverTooltip } from "@/shared/components/HoverTooltip";
import { navGroups, type NavGroup } from "@/shared/layout/navGroups";
import { cn } from "@/lib/utils";

type Props = {
  collapsed: boolean;
  groups?: NavGroup[];
};

export function SidebarNav({ collapsed, groups = navGroups }: Props) {
  return (
    <nav className={cn("flex flex-1 flex-col", collapsed ? "gap-0.5" : "gap-space-md")}>
      {groups.map((group) => (
        <div key={group.title} className="space-y-1">
          {!collapsed && (
            <div className="px-space-sm py-1 font-label-sm text-label-sm uppercase tracking-wider text-on-surface-variant">
              {group.title}
            </div>
          )}
          {group.items.map(({ to, label, icon, end }) => {
            const link = (
              <NavLink
                to={to}
                end={end}
                aria-label={label}
                className={({ isActive }) =>
                  cn(
                    "flex items-center rounded py-1.5 text-on-surface-variant transition-colors hover:bg-surface-container hover:text-on-surface",
                    collapsed ? "justify-center px-2" : "gap-space-sm px-space-sm",
                    isActive && "bg-surface-container-high font-semibold text-on-surface",
                  )
                }
              >
                <MaterialIcon name={icon} size={18} className="shrink-0" />
                {!collapsed && <span className="font-body-sm text-body-sm">{label}</span>}
              </NavLink>
            );
            if (!collapsed) {
              return <div key={to}>{link}</div>;
            }
            return (
              <HoverTooltip key={to} label={label} side="right" className="flex w-full justify-center">
                {link}
              </HoverTooltip>
            );
          })}
        </div>
      ))}
    </nav>
  );
}
