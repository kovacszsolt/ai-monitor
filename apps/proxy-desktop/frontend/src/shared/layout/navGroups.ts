export type NavItem = {
  to: string;
  label: string;
  icon: string;
  end?: boolean;
};

export type NavGroup = {
  title: string;
  items: NavItem[];
};

export const navGroups: NavGroup[] = [
  {
    title: "Observe",
    items: [
      { to: "/", label: "Dashboard", icon: "dashboard", end: true },
      { to: "/requests", label: "Requests", icon: "data_exploration" },
    ],
  },
  {
    title: "Configure",
    items: [
      { to: "/projects", label: "Projects", icon: "folder_copy" },
      { to: "/providers", label: "Providers", icon: "hub" },
      { to: "/models", label: "Models", icon: "auto_awesome" },
    ],
  },
  {
    title: "Operate",
    items: [
      { to: "/chat", label: "Chat Testbed", icon: "terminal" },
      { to: "/proxy", label: "Proxy Server", icon: "dns" },
    ],
  },
  {
    title: "System",
    items: [
      { to: "/proxy-variables", label: "Runtime Vars", icon: "tune" },
      { to: "/settings", label: "Settings", icon: "settings" },
      { to: "/about", label: "About", icon: "info" },
    ],
  },
];
