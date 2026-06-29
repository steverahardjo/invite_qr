import { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { api, type Participant } from '../api/client'

export default function InvitePage() {
  const { external_id } = useParams<{ external_id: string }>()
  const [participant, setParticipant] = useState<Participant | null>(null)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    if (!external_id) {
      setError('No invite ID provided')
      return
    }
    api.getInvite(external_id)
      .then(setParticipant)
      .catch((err) => setError(err.message))
  }, [external_id])

  if (error) {
    return (
      <div className="invite-page">
        <div className="invite-card">
          <div className="ornament">✦ ✦ ✦</div>
          <h1>Invite Not Found</h1>
          <p className="subtitle">{error}</p>
        </div>
      </div>
    )
  }

  if (!participant) {
    return (
      <div className="invite-page">
        <div className="invite-card">
          <div className="ornament">✦ ✦ ✦</div>
          <p className="loading-state">Preparing your invitation…</p>
        </div>
      </div>
    )
  }

  return (
    <div className="invite-page">
      <div className="invite-card">
        <div className="ornament">✦ ✦ ✦</div>

        <h1>You're Invited</h1>
        <p className="subtitle">To celebrate our wedding</p>

        <div className="detail-row">
          <span className="label">Guest</span>
          <span className="value">{participant.name}</span>
        </div>
        {participant.email && (
          <div className="detail-row">
            <span className="label">Email</span>
            <span className="value">{participant.email}</span>
          </div>
        )}

        <div className="qr-shell">
          <div className="qr-inner">
            <img
              src={api.getQR(participant.external_id)}
              alt="Your personal QR code"
            />
          </div>
        </div>

        <p className="invite-footer">
          Present this QR code at the venue for check-in
        </p>

        <div style={{ marginTop: 20 }}>
          <span className={`status-chip ${participant.accessed ? 'attended' : 'pending'}`}>
            {participant.accessed ? '✓ Checked In' : 'Not Yet Checked In'}
          </span>
        </div>
      </div>
    </div>
  )
}
