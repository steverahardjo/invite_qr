import { Routes, Route, Navigate } from "react-router-dom";
import { InvitePage } from "@/pages/InvitePage";
import { EventPage } from "@/pages/EventPage";
import { DebugPage } from "@/pages/DebugPage";
import "./index.css";

export function App() {
  return (
    <Routes>
      <Route path="/debug" element={<DebugPage />} />
      <Route path="/:id/:user" element={<InvitePage />} />
      <Route path="/:id" element={<EventPage />} />
      <Route path="/" element={<Navigate to="/debug" replace />} />
      <Route path="*" element={<Navigate to="/debug" replace />} />
    </Routes>
  );
}

export default App;
