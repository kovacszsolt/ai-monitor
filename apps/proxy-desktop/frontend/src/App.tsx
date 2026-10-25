import { lazy, Suspense } from "react";
import { BrowserRouter, Navigate, Route, Routes } from "react-router-dom";
import { Toaster } from "sonner";
import { AppLayout } from "@/shared/layout/AppLayout";
import { LoadingState } from "@/shared/components/PageStates";

const DashboardPage = lazy(() =>
  import("@/features/dashboard/pages/DashboardPage").then((m) => ({ default: m.DashboardPage })),
);
const RequestsPage = lazy(() =>
  import("@/features/requests/pages/RequestsPage").then((m) => ({ default: m.RequestsPage })),
);
const ProjectsPage = lazy(() =>
  import("@/features/projects/pages/ProjectsPage").then((m) => ({ default: m.ProjectsPage })),
);
const ProjectDetailPage = lazy(() =>
  import("@/features/projects/pages/ProjectDetailPage").then((m) => ({ default: m.ProjectDetailPage })),
);
const ProvidersPage = lazy(() =>
  import("@/features/providers/pages/ProvidersPage").then((m) => ({ default: m.ProvidersPage })),
);
const ModelsPage = lazy(() =>
  import("@/features/models/pages/ModelsPage").then((m) => ({ default: m.ModelsPage })),
);
const ProxyPage = lazy(() =>
  import("@/features/proxy/pages/ProxyPage").then((m) => ({ default: m.ProxyPage })),
);
const ChatPage = lazy(() =>
  import("@/features/chat/pages/ChatPage").then((m) => ({ default: m.ChatPage })),
);
const SettingsPage = lazy(() =>
  import("@/features/settings/pages/SettingsPage").then((m) => ({ default: m.SettingsPage })),
);
const AboutPage = lazy(() =>
  import("@/features/about/pages/AboutPage").then((m) => ({ default: m.AboutPage })),
);
const ProxyVariablesPage = lazy(() =>
  import("@/features/settings/pages/ProxyVariablesPage").then((m) => ({ default: m.ProxyVariablesPage })),
);

function RouteFallback() {
  return <LoadingState label="Loading…" />;
}

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route element={<AppLayout />}>
          <Route
            path="/"
            element={
              <Suspense fallback={<RouteFallback />}>
                <DashboardPage />
              </Suspense>
            }
          />
          <Route
            path="/requests"
            element={
              <Suspense fallback={<RouteFallback />}>
                <RequestsPage />
              </Suspense>
            }
          />
          <Route
            path="/projects"
            element={
              <Suspense fallback={<RouteFallback />}>
                <ProjectsPage />
              </Suspense>
            }
          />
          <Route
            path="/projects/:id"
            element={
              <Suspense fallback={<RouteFallback />}>
                <ProjectDetailPage />
              </Suspense>
            }
          />
          <Route
            path="/providers"
            element={
              <Suspense fallback={<RouteFallback />}>
                <ProvidersPage />
              </Suspense>
            }
          />
          <Route
            path="/models"
            element={
              <Suspense fallback={<RouteFallback />}>
                <ModelsPage />
              </Suspense>
            }
          />
          <Route
            path="/proxy"
            element={
              <Suspense fallback={<RouteFallback />}>
                <ProxyPage />
              </Suspense>
            }
          />
          <Route
            path="/chat"
            element={
              <Suspense fallback={<RouteFallback />}>
                <ChatPage />
              </Suspense>
            }
          />
          <Route
            path="/proxy-variables"
            element={
              <Suspense fallback={<RouteFallback />}>
                <ProxyVariablesPage />
              </Suspense>
            }
          />
          <Route
            path="/settings"
            element={
              <Suspense fallback={<RouteFallback />}>
                <SettingsPage />
              </Suspense>
            }
          />
          <Route
            path="/about"
            element={
              <Suspense fallback={<RouteFallback />}>
                <AboutPage />
              </Suspense>
            }
          />
          <Route path="/upstream" element={<Navigate to="/providers" replace />} />
        </Route>
      </Routes>
      <Toaster theme="dark" richColors position="top-right" />
    </BrowserRouter>
  );
}

export default App;
