import { useEffect, useRef, useState } from 'react'
import { Html5Qrcode } from 'html5-qrcode'

interface QRScannerProps {
  onResult: (participantId: string) => void
}

export default function QRScanner({ onResult }: QRScannerProps) {
  const [open, setOpen] = useState(false)
  const scannerRef = useRef<Html5Qrcode | null>(null)
  const containerRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    if (!open) return

    const id = 'qr-reader-modal'
    const el = document.getElementById(id)
    if (!el) return

    const scanner = new Html5Qrcode(id)
    scannerRef.current = scanner

    scanner.start(
      { facingMode: 'environment' },
      { fps: 10, qrbox: { width: 250, height: 250 } },
      (decodedText) => {
        onResult(decodedText)
        setOpen(false)
      },
      () => {},
    ).catch(() => {})

    return () => {
      scanner.stop().catch(() => {})
      scannerRef.current = null
    }
  }, [open, onResult])

  return (
    <>
      <button className="btn btn-gold" onClick={() => setOpen(true)}>
        Open Scanner
      </button>

      {open && (
        <div
          onClick={() => setOpen(false)}
          style={{
            position: 'fixed',
            inset: 0,
            zIndex: 9998,
            background: 'rgba(45, 24, 16, 0.6)',
            backdropFilter: 'blur(6px)',
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'center',
            padding: 24,
          }}
        >
          <div
            onClick={(e) => e.stopPropagation()}
            style={{
              position: 'relative',
              zIndex: 9999,
              background: '#fff',
              borderRadius: 20,
              padding: 32,
              width: '100%',
              maxWidth: 420,
              boxShadow: '0 0 0 1px rgba(201,169,110,0.3), 0 24px 80px rgba(0,0,0,0.25)',
            }}
          >
            <button
              onClick={() => setOpen(false)}
              style={{
                position: 'absolute',
                top: 12,
                right: 16,
                background: 'none',
                border: 'none',
                fontSize: 20,
                color: '#8B7D72',
                cursor: 'pointer',
                lineHeight: 1,
                padding: 4,
              }}
            >
              ✕
            </button>

            <div style={{ fontFamily: "'Playfair Display', serif", fontSize: 20, fontWeight: 700, marginBottom: 4 }}>
              Scan QR Code
            </div>
            <p style={{ fontFamily: "'Plus Jakarta Sans', sans-serif", fontSize: 13, color: '#8B7D72', marginBottom: 20 }}>
              Point camera at guest's invite QR
            </p>

            <div
              ref={containerRef}
              id="qr-reader-modal"
              style={{
                width: '100%',
                aspectRatio: '1 / 1',
                borderRadius: 14,
                overflow: 'hidden',
                background: '#F5F0E8',
              }}
            />
          </div>
        </div>
      )}
    </>
  )
}
