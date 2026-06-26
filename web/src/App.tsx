import { Routes, Route, Navigate } from "react-router-dom";
import { InvitePage } from "@/pages/InvitePage";
import { EventPage } from "@/pages/EventPage";
import { DebugPage } from "@/pages/DebugPage";
import { AdminGuard } from "@/pages/admin/AdminGuard";
import { AdminDashboard } from "@/pages/admin/AdminDashboard";
import { QRCodePage } from "@/pages/admin/QRCodePage";
import { ParticipantsPage } from "@/pages/admin/ParticipantsPage";
import "./index.css";

export function App() {
  return (
    <Routes>
      <Route path="/debug" element={<DebugPage />} />

      <Route
        path="/admin"
        element={
          <AdminGuard>
            <AdminDashboard />
          </AdminGuard>
        }
      />
      <Route
        path="/admin/qr_code"
        element={
          <AdminGuard>
            <QRCodePage />
          </AdminGuard>
        }
      />
      <Route
        path="/admin/participants"
        element={
          <AdminGuard>
            <ParticipantsPage />
          </AdminGuard>
        }
      />

      <Route path="/:id/:user" element={<InvitePage />} />
      <Route path="/:id" element={<EventPage />} />
      <Route path="/" element={<Navigate to="/debug" replace />} />
      <Route path="*" element={<Navigate to="/debug" replace />} />
    </Routes>
  );
}

export default App;
