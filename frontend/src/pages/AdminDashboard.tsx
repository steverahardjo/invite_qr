import { useEffect, useState, useCallback } from "react";
import { api, type Participant, type ParticipantInput } from "../api/client";
import DebugPanel from "../components/DebugPanel";

interface Stats {
  total: number;
  attended: number;
  pending: number;
  sent: number;
}

function computeStats(participants: Participant[] | null | undefined): Stats {
  const list = participants ?? [];
  return {
    total: list.length,
    attended: list.filter((p) => p.accessed).length,
    pending: list.filter((p) => !p.accessed).length,
    sent: list.filter((p) => p.sent).length,
  };
}

export default function AdminDashboard() {
  const [participants, setParticipants] = useState<Participant[]>([]);
  const [stats, setStats] = useState<Stats>({
    total: 0,
    attended: 0,
    pending: 0,
    sent: 0,
  });
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [sendingId, setSendingId] = useState<number | null>(null);
  const [sendError, setSendError] = useState<string | null>(null);
  const [sendSuccess, setSendSuccess] = useState<string | null>(null);
  const [showForm, setShowForm] = useState(false);
  const [form, setForm] = useState<ParticipantInput>({
    name: "",
    email: "",
    wa_number: "",
  });
  const [formError, setFormError] = useState<string | null>(null);
  const [formSuccess, setFormSuccess] = useState<string | null>(null);
  const [submitting, setSubmitting] = useState(false);

  const fetchParticipants = useCallback(async () => {
    try {
      const data = await api.getParticipants();
      setParticipants(data ?? []);
      setStats(computeStats(data));
    } catch (err) {
      setError(
        err instanceof Error ? err.message : "Failed to load participants",
      );
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchParticipants();
  }, [fetchParticipants]);

  async function handleAddParticipant(e: React.FormEvent) {
    e.preventDefault();
    setFormError(null);
    setFormSuccess(null);

    if (!form.name.trim()) {
      setFormError("Name is required");
      return;
    }

    setSubmitting(true);
    try {
      await api.addParticipant(form);
      setFormSuccess(`${form.name} added successfully`);
      setForm({ name: "", email: "", wa_number: "" });
      setShowForm(false);
      await fetchParticipants();
    } catch (err) {
      setFormError(
        err instanceof Error ? err.message : "Failed to add participant",
      );
    } finally {
      setSubmitting(false);
    }
  }

  async function handleSendInvite(p: Participant) {
    setSendingId(p.id);
    setSendError(null);
    setSendSuccess(null);
    try {
      await api.sendInvite(p);
      setSendSuccess(`Invite sent to ${p.name}`);
      setTimeout(() => setSendSuccess(null), 3000);
      await fetchParticipants();
    } catch (err) {
      setSendError(
        err instanceof Error ? err.message : "Failed to send invite",
      );
    } finally {
      setSendingId(null);
    }
  }

  if (loading) {
    return (
      <div className="admin-main">
        <DebugPanel />
        <h1 className="page-title">Dashboard</h1>
        <p className="page-sub">Loading your guest list…</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="admin-main">
        <DebugPanel />
        <h1 className="page-title">Dashboard</h1>
        <div className="msg-error">{error}</div>
      </div>
    );
  }

  return (
    <div className="admin-main">
      <DebugPanel />

      <h1 className="page-title">Dashboard</h1>
      <p className="page-sub">Manage your wedding guests and check-ins</p>

      {/* Stats */}
      <div className="stats-grid">
        <div className="stat-card">
          <span className="stat-icon">👥</span>
          <div className="stat-label">Total Guests</div>
          <div className="stat-value">{stats.total}</div>
        </div>
        <div className="stat-card">
          <span className="stat-icon">✓</span>
          <div className="stat-label">Checked In</div>
          <div className="stat-value sage">{stats.attended}</div>
        </div>
        <div className="stat-card">
          <span className="stat-icon">⏳</span>
          <div className="stat-label">Pending</div>
          <div className="stat-value gold">{stats.pending}</div>
        </div>
        <div className="stat-card">
          <span className="stat-icon">✉</span>
          <div className="stat-label">Invites Sent</div>
          <div className="stat-value">{stats.sent}</div>
        </div>
      </div>

      {/* Add Participant */}
      <div className="form-section">
        <div className="section-header">
          <h2 className="section-title">Guests</h2>
          <button
            className={`btn ${showForm ? "btn-ghost" : "btn-gold"} btn-sm`}
            onClick={() => {
              setShowForm(!showForm);
              setFormError(null);
              setFormSuccess(null);
            }}
          >
            {showForm ? "Cancel" : "+ Add Guest"}
          </button>
        </div>

        {showForm && (
          <div className="form-card">
            <form onSubmit={handleAddParticipant}>
              <div className="form-group">
                <label className="label">Full Name</label>
                <input
                  className="input"
                  type="text"
                  value={form.name}
                  onChange={(e) => setForm({ ...form, name: e.target.value })}
                  placeholder="Guest's full name"
                  required
                />
              </div>
              <div className="form-row">
                <div className="form-group">
                  <label className="label">Email</label>
                  <input
                    className="input"
                    type="email"
                    value={form.email}
                    onChange={(e) =>
                      setForm({ ...form, email: e.target.value })
                    }
                    placeholder="email@example.com"
                  />
                </div>
                <div className="form-group">
                  <label className="label">WhatsApp</label>
                  <input
                    className="input"
                    type="text"
                    value={form.wa_number}
                    onChange={(e) =>
                      setForm({ ...form, wa_number: e.target.value })
                    }
                    placeholder="+1234567890"
                  />
                </div>
              </div>

              {formError && <div className="msg-error">{formError}</div>}
              {formSuccess && <div className="msg-success">{formSuccess}</div>}

              <button
                type="submit"
                className="btn btn-primary"
                disabled={submitting}
                style={{ width: "100%" }}
              >
                {submitting ? "Adding…" : "Add Guest"}
              </button>
            </form>
          </div>
        )}
      </div>

      {sendError && <div className="msg-error">{sendError}</div>}
      {sendSuccess && <div className="msg-success">{sendSuccess}</div>}

      {/* Participants Table */}
      <div className="table-shell">
        <div className="table-inner">
          <table>
            <thead>
              <tr>
                <th>Name</th>
                <th>Email</th>
                <th>WhatsApp</th>
                <th>Check-in</th>
                <th>Invite</th>
              </tr>
            </thead>
            <tbody>
              {participants.length === 0 ? (
                <tr>
                  <td
                    colSpan={5}
                    style={{
                      textAlign: "center",
                      color: "var(--text-muted)",
                      padding: 40,
                    }}
                  >
                    No guests yet. Add your first one above.
                  </td>
                </tr>
              ) : (
                participants.map((p) => (
                  <tr key={p.id}>
                    <td style={{ fontWeight: 600 }}>{p.name}</td>
                    <td>{p.email || "—"}</td>
                    <td>{p.wa_number || "—"}</td>
                    <td>
                      <span
                        className={`badge ${p.accessed ? "badge-ok" : "badge-pending"}`}
                      >
                        {p.accessed ? "Checked In" : "Pending"}
                      </span>
                    </td>
                    <td>
                      {p.sent ? (
                        <span className="badge badge-ok">Sent</span>
                      ) : (
                        <button
                          className="btn btn-sm"
                          onClick={() => handleSendInvite(p)}
                          disabled={sendingId === p.id}
                          style={{
                            background: "var(--gold)",
                            color: "#fff",
                            border: "none",
                            borderRadius: 6,
                            padding: "4px 12px",
                            cursor: sendingId === p.id ? "wait" : "pointer",
                            fontSize: "0.8rem",
                            fontWeight: 600,
                          }}
                        >
                          {sendingId === p.id ? "Sending…" : "Send"}
                        </button>
                      )}
                    </td>
                  </tr>
                ))
              )}
            </tbody>
          </table>
        </div>
      </div>
    </div>
  );
}
