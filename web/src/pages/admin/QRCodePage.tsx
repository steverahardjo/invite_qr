import { useState, useEffect, useRef } from "react";
import { ArrowLeft } from "lucide-react";
import { useNavigate } from "react-router-dom";
import jsQR from "jsqr";

export function QRCodePage() {
  const navigate = useNavigate();
  const [scanned, setScanned] = useState<string>("");
  const videoRef = useRef<HTMLVideoElement>(null);
  const canvasRef = useRef<HTMLCanvasElement>(null);

  useEffect(() => {
    let stream: MediaStream | null = null;
    let timer: ReturnType<typeof setInterval> | null = null;

    navigator.mediaDevices
      .getUserMedia({ video: { facingMode: "environment" } })
      .then((s) => {
        stream = s;
        if (videoRef.current) videoRef.current.srcObject = s;

        timer = setInterval(() => {
          const video = videoRef.current;
          const canvas = canvasRef.current;
          if (!video || !canvas || scanned) return;

          canvas.width = video.videoWidth;
          canvas.height = video.videoHeight;
          const ctx = canvas.getContext("2d");
          if (!ctx) return;

          ctx.drawImage(video, 0, 0);
          const imageData = ctx.getImageData(0, 0, canvas.width, canvas.height);
          const code = jsQR(imageData.data, imageData.width, imageData.height);

          if (code) {
            setScanned(code.data);
            timer && clearInterval(timer);
          }
        }, 500);
      })
      .catch(() => {});

    return () => {
      timer && clearInterval(timer);
      if (stream) {
        stream.getTracks().forEach((t) => t.stop());
      }
    };
  }, [scanned]);

  return (
    <div className="min-h-screen w-full bg-white">
      <div className="mx-auto max-w-lg px-4 py-8">
        <button
          onClick={() => navigate("/admin")}
          className="mb-4 flex items-center gap-1.5 text-xs text-stone-400 hover:text-stone-600 transition-colors"
        >
          <ArrowLeft className="h-3.5 w-3.5" />
          Back
        </button>

        <p className="font-['Plus_Jakarta_Sans'] text-xs uppercase tracking-[0.25em] text-stone-400 text-center mb-6">
          Scan QR Code
        </p>

        <div className="relative aspect-square max-w-sm mx-auto rounded-2xl overflow-hidden bg-stone-950 ring-1 ring-stone-200">
          <video
            ref={videoRef}
            autoPlay
            playsInline
            muted
            className={`w-full h-full object-cover ${scanned ? "opacity-40" : ""}`}
          />
          <canvas ref={canvasRef} className="hidden" />
          {scanned && (
            <div className="absolute inset-0 flex items-center justify-center">
              <div className="bg-white/90 backdrop-blur rounded-2xl px-6 py-4 text-center max-w-[80%]">
                <p className="font-['Plus_Jakarta_Sans'] text-[10px] uppercase tracking-[0.25em] text-stone-400 mb-1">
                  Scanned
                </p>
                <p className="text-sm font-medium text-stone-800 break-all">
                  {scanned}
                </p>
                <button
                  onClick={() => setScanned("")}
                  className="mt-3 text-xs text-stone-500 underline hover:text-stone-800"
                >
                  Scan again
                </button>
              </div>
            </div>
          )}
        </div>

        <div className="mt-8 text-center">
          <p className="text-[11px] text-stone-400">
            Point your camera at a QR code to scan
          </p>
        </div>
      </div>
    </div>
  );
}
