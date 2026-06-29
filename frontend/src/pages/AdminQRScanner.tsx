import { useState, useCallback } from 'react'
import { api } from '../api/client'
import QRScanner from '../components/QRScanner'
import DebugPanel from '../components/DebugPanel'

export default function AdminQRScanner() {
  const [result, setResult] = useState<{ ok: boolean; message: string } | null>(null)

  const handleResult = useCallback(async (participantId: string) => {
    try {
      await api.markAttendance(participantId)
      setResult({ ok: true, message: '✓ Check-in recorded' })
    } catch (err) {
      setResult({ ok: false, message: err instanceof Error ? err.message : 'Failed to check in' })
    }
    setTimeout(() => setResult(null), 3500)
  }, [])

  return (
    <div className="admin-main">
      <DebugPanel />

      <h1 className="page-title">QR Scanner</h1>
      <p className="page-sub">Point the camera at a guest's invite QR to mark them as checked in</p>

      <div className="scanner-shell" style={{ maxWidth: 480 }}>
        <div className="scanner-inner">
          <QRScanner onResult={handleResult} />
          {result && (
            <div className={`scanner-result ${result.ok ? 'attended' : 'error'}`}>
              {result.message}
            </div>
          )}
        </div>
      </div>
    </div>
  )
}
