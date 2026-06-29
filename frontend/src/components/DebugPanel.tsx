import { useState, useEffect } from 'react'
import { subscribe, getEntries, type LogEntry } from '../api/logger'

export default function DebugPanel() {
  const [open, setOpen] = useState(false)
  const [logs, setLogs] = useState<LogEntry[]>([])

  useEffect(() => {
    setLogs(getEntries())
    return subscribe(() => setLogs([...getEntries()]))
  }, [])

  const unread = logs.filter((l) => l.status !== 'pending').length

  return (
    <>
      <button
        onClick={() => setOpen(!open)}
        style={{
          position: 'fixed',
          bottom: 16,
          right: 16,
          zIndex: 9999,
          width: 44,
          height: 44,
          borderRadius: '50%',
          border: '1.5px solid rgba(201,169,110,0.4)',
          background: 'rgba(253,251,247,0.95)',
          backdropFilter: 'blur(12px)',
          color: '#5C3A2A',
          fontSize: 16,
          cursor: 'pointer',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          boxShadow: '0 4px 20px rgba(45,24,16,0.1)',
          fontFamily: 'Plus Jakarta Sans, sans-serif',
          fontWeight: 600,
          transition: 'all 0.3s cubic-bezier(0.32,0.72,0,1)',
        }}
        title="API Debug Log"
      >
        {unread > 0 ? unread : '∞'}
      </button>

      {open && (
        <div
          style={{
            position: 'fixed',
            bottom: 68,
            right: 16,
            zIndex: 9999,
            width: 420,
            maxHeight: 360,
            overflowY: 'auto',
            background: 'rgba(45,24,16,0.97)',
            backdropFilter: 'blur(16px)',
            borderRadius: 16,
            border: '1px solid rgba(201,169,110,0.2)',
            boxShadow: '0 12px 40px rgba(0,0,0,0.3)',
            padding: 16,
            fontFamily: 'Plus Jakarta Sans, sans-serif',
            fontSize: 12,
            color: '#E8DCC8',
          }}
        >
          <div
            style={{
              display: 'flex',
              justifyContent: 'space-between',
              alignItems: 'center',
              marginBottom: 12,
              paddingBottom: 8,
              borderBottom: '1px solid rgba(201,169,110,0.15)',
            }}
          >
            <span style={{ fontWeight: 600, fontSize: 13, color: '#C9A96E' }}>
              API Debug Log
            </span>
            <button
              onClick={() => setOpen(false)}
              style={{
                background: 'none',
                border: 'none',
                color: '#8B7D72',
                cursor: 'pointer',
                fontSize: 16,
                lineHeight: 1,
                padding: 2,
              }}
            >
              ✕
            </button>
          </div>

          {logs.length === 0 && (
            <div style={{ color: '#8B7D72', fontStyle: 'italic' }}>
              No API calls yet.
            </div>
          )}

          {logs
            .slice()
            .reverse()
            .map((log) => (
              <div
                key={log.id}
                style={{
                  display: 'flex',
                  alignItems: 'center',
                  gap: 8,
                  padding: '6px 8px',
                  borderRadius: 6,
                  marginBottom: 2,
                  background:
                    log.status === 'error'
                      ? 'rgba(183,110,121,0.12)'
                      : 'transparent',
                  fontSize: 11.5,
                }}
              >
                <span
                  style={{
                    width: 18,
                    height: 18,
                    borderRadius: '50%',
                    display: 'flex',
                    alignItems: 'center',
                    justifyContent: 'center',
                    fontSize: 10,
                    fontWeight: 700,
                    flexShrink: 0,
                    color: '#fff',
                    background:
                      log.status === 'success'
                        ? '#9CAF88'
                        : log.status === 'error'
                          ? '#B76E79'
                          : '#C9A96E',
                  }}
                >
                  {log.status === 'success'
                    ? '✓'
                    : log.status === 'error'
                      ? '✗'
                      : '…'}
                </span>
                <span style={{ fontWeight: 600, color: '#F5F0E8', minWidth: 44 }}>
                  {log.method}
                </span>
                <span style={{ color: '#E8DCC8', wordBreak: 'break-all', flex: 1 }}>
                  {log.path}
                </span>
                <span style={{ color: '#8B7D72', whiteSpace: 'nowrap' }}>
                  {log.timestamp.toLocaleTimeString([], {
                    hour: '2-digit',
                    minute: '2-digit',
                    second: '2-digit',
                  })}
                </span>
                {log.error && (
                  <span style={{ color: '#B76E79', fontStyle: 'italic', fontSize: 10 }}>
                    {log.error}
                  </span>
                )}
              </div>
            ))}
        </div>
      )}
    </>
  )
}
