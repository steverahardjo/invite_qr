import { useEffect, useRef } from 'react'
import { Html5Qrcode } from 'html5-qrcode'

interface QRScannerProps {
  onResult: (participantId: string) => void
}

export default function QRScanner({ onResult }: QRScannerProps) {
  const containerRef = useRef<HTMLDivElement>(null)
  const scannerRef = useRef<Html5Qrcode | null>(null)

  useEffect(() => {
    const containerId = 'qr-reader'
    if (!document.getElementById(containerId)) return

    const scanner = new Html5Qrcode(containerId)
    scannerRef.current = scanner

    scanner.start(
      { facingMode: 'environment' },
      { fps: 10, qrbox: { width: 250, height: 250 } },
      (decodedText) => {
        onResult(decodedText)
      },
      () => {},
    ).catch(() => {
      // Camera not available — silently degrade
    })

    return () => {
      scanner.stop().catch(() => {})
    }
  }, [onResult])

  return (
    <div ref={containerRef}>
      <div id="qr-reader" />
    </div>
  )
}
