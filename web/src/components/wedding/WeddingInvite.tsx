import { useState, useEffect, useCallback } from "react";
import { Card, CardContent } from "@/components/ui/card";
import { Button } from "@/components/ui/button";
import QRCode from "qrcode";

interface WeddingInviteProps {
  id: string;
  user: string;
}

const QR_SEEN_KEY = (id: string, user: string) => `qr_seen_${id}_${user}`;

function QRModal({
  url,
  id,
  user,
  onClose,
}: {
  url: string;
  id: string;
  user: string;
  onClose: () => void;
}) {
  const [step, setStep] = useState<"warning" | "qr" | "spent">(
    sessionStorage.getItem(QR_SEEN_KEY(id, user)) === "true"
      ? "spent"
      : "warning",
  );
  const [qrDataUrl, setQrDataUrl] = useState<string>("");

  useEffect(() => {
    if (step !== "qr") return;
    QRCode.toDataURL(url, {
      width: 280,
      margin: 2,
      color: { dark: "#1a1a2e", light: "#ffffff" },
    })
      .then(setQrDataUrl)
      .catch(() => {});
  }, [step, url]);

  useEffect(() => {
    if (step !== "qr") return;
    sessionStorage.setItem(QR_SEEN_KEY(id, user), "true");
  }, [step, id, user]);

  const handleReveal = useCallback(() => setStep("qr"), []);

  return (
    <div
      className="fixed inset-0 z-50 flex items-center justify-center bg-black/60 backdrop-blur-sm"
      onClick={onClose}
    >
      <div
        className="w-full max-w-sm mx-4"
        onClick={(e) => e.stopPropagation()}
      >
        {/* ── Outer shell (double-bezel) ── */}
        <div className="bg-white/5 p-1.5 rounded-[2rem] ring-1 ring-black/5">
          <div className="bg-white rounded-[calc(2rem-0.375rem)] p-8 shadow-2xl text-center space-y-5">
            {step === "spent" && (
              <>
                <div className="mx-auto flex h-14 w-14 items-center justify-center rounded-full bg-stone-100">
                  <svg className="h-6 w-6 text-stone-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                    <path strokeLinecap="round" strokeLinejoin="round" d="M16.5 10.5V6.75a4.5 4.5 0 10-9 0v3.75m-.75 11.25h10.5a2.25 2.25 0 002.25-2.25v-6.75a2.25 2.25 0 00-2.25-2.25H6.75a2.25 2.25 0 00-2.25 2.25v6.75a2.25 2.25 0 002.25 2.25z" />
                  </svg>
                </div>
                <p className="text-lg font-serif text-[#1a1a2e]">
                  Access Expired
                </p>
                <p className="text-sm text-stone-400 leading-relaxed">
                  This QR code has already been viewed and is no longer available.
                </p>
                <Button
                  variant="outline"
                  size="sm"
                  className="w-full rounded-full"
                  onClick={onClose}
                >
                  Close
                </Button>
              </>
            )}

            {step === "warning" && (
              <>
                <div className="mx-auto flex h-14 w-14 items-center justify-center rounded-full bg-amber-50">
                  <svg className="h-6 w-6 text-amber-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                    <path strokeLinecap="round" strokeLinejoin="round" d="M12 9v3.75m-9.303 3.376c-.866 1.5.217 3.374 1.948 3.374h14.71c1.73 0 2.813-1.874 1.948-3.374L13.949 3.378c-.866-1.5-3.032-1.5-3.898 0L2.697 16.126zM12 15.75h.007v.008H12v-.008z" />
                  </svg>
                </div>
                <p className="text-lg font-serif text-[#1a1a2e]">
                  One-Time Access
                </p>
                <p className="text-sm text-stone-400 leading-relaxed">
                  This QR code can only be viewed once. After closing, it will no longer be accessible.
                </p>
                <div className="flex gap-3">
                  <Button
                    variant="outline"
                    size="sm"
                    className="flex-1 rounded-full"
                    onClick={onClose}
                  >
                    Cancel
                  </Button>
                  <Button
                    size="sm"
                    className="flex-1 rounded-full bg-[#1a1a2e] hover:bg-[#2a2a4e] text-white"
                    onClick={handleReveal}
                  >
                    Reveal QR Code
                  </Button>
                </div>
              </>
            )}

            {step === "qr" && (
              <>
                {qrDataUrl && (
                  <img src={qrDataUrl} alt="QR Code" className="w-full" />
                )}
                <p className="text-sm text-stone-400">
                  Scan to view your invitation
                </p>
                <p className="text-[10px] uppercase tracking-[0.2em] text-amber-500 font-medium">
                  This QR code has been viewed
                </p>
                <Button
                  variant="outline"
                  size="sm"
                  className="w-full rounded-full"
                  onClick={onClose}
                >
                  Close
                </Button>
              </>
            )}
          </div>
        </div>
      </div>
    </div>
  );
}

function useIsMobile() {
  const [isMobile, setIsMobile] = useState(false);

  useEffect(() => {
    const check = () => {
      const hasTouch = navigator.maxTouchPoints > 0;
      const isNarrow = window.innerWidth < 1024;
      setIsMobile(hasTouch && isNarrow);
    };
    check();
    window.addEventListener("resize", check);
    return () => window.removeEventListener("resize", check);
  }, []);

  return isMobile;
}

export function WeddingInvite({ id, user }: WeddingInviteProps) {
  const [showQR, setShowQR] = useState(false);
  const isMobile = useIsMobile();
  const inviteUrl = `${window.location.origin}/${id}/${user}`;

  return (
    <div className="min-h-screen w-full bg-white">
      <div className="mx-auto max-w-lg px-4 py-16 sm:py-24 space-y-14">

        {/* ── Hero ── */}
        <div className="text-center space-y-6">
          <div className="space-y-3">
            <p className="font-['Plus_Jakarta_Sans'] text-[10px] uppercase tracking-[0.3em] text-stone-400">
              Together with their families
            </p>
            <div className="flex items-center justify-center gap-3">
              <span className="block h-px w-8 bg-gradient-to-r from-transparent via-amber-300/60 to-transparent" />
              <span className="block h-2 w-2 rotate-45 border border-amber-300/60" />
              <span className="block h-px w-8 bg-gradient-to-r from-transparent via-amber-300/60 to-transparent" />
            </div>
          </div>

          <h1 className="font-['Playfair_Display'] text-5xl sm:text-6xl text-[#1a1a2e] leading-[1.1]">
            Sarah
            <span className="block font-['Cormorant_Garamond'] text-2xl sm:text-3xl font-light italic text-stone-400 mt-1">
              &amp;
            </span>
            James
          </h1>

          <p className="font-['Plus_Jakarta_Sans'] text-[10px] uppercase tracking-[0.3em] text-stone-400">
            Request the honor of your presence
          </p>
        </div>

        {/* ── Guest Detail ── */}
        <div className="bg-stone-50/80 p-1 rounded-2xl ring-1 ring-stone-100">
          <div className="bg-white rounded-[calc(1.5rem-0.25rem)] px-6 py-5 text-center space-y-1.5">
            <p className="font-['Plus_Jakarta_Sans'] text-[10px] uppercase tracking-[0.25em] text-stone-400">
              Guest of Honor
            </p>
            <p className="font-['Playfair_Display'] text-2xl text-[#1a1a2e] capitalize">
              {user}
            </p>
            <p className="font-['Plus_Jakarta_Sans'] text-xs text-stone-400">
              You are cordially invited to celebrate their union
            </p>
          </div>
        </div>

        {/* ── Event Details ── */}
        <div className="bg-stone-50/80 p-1 rounded-2xl ring-1 ring-stone-100">
          <div className="bg-white rounded-[calc(1.5rem-0.25rem)] px-6 py-5 space-y-5">
            <div className="flex items-start gap-4">
              <div className="mt-0.5 flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-amber-50">
                <svg className="h-4 w-4 text-amber-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M6.75 3v2.25M17.25 3v2.25M3 18.75V7.5a2.25 2.25 0 012.25-2.25h13.5A2.25 2.25 0 0121 7.5v11.25m-18 0A2.25 2.25 0 005.25 21h13.5A2.25 2.25 0 0021 18.75m-18 0v-7.5A2.25 2.25 0 015.25 9h13.5A2.25 2.25 0 0121 11.25v7.5" />
                </svg>
              </div>
              <div className="space-y-0.5">
                <p className="font-['Plus_Jakarta_Sans'] text-[10px] uppercase tracking-[0.25em] text-stone-400">
                  Date &amp; Time
                </p>
                <p className="font-['Playfair_Display'] text-lg text-[#1a1a2e]">Saturday, June 15, 2025</p>
                <p className="font-['Plus_Jakarta_Sans'] text-sm text-stone-400">Half past four in the afternoon</p>
              </div>
            </div>

            <div className="h-px bg-stone-100" />

            <div className="flex items-start gap-4">
              <div className="mt-0.5 flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-amber-50">
                <svg className="h-4 w-4 text-amber-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M15 10.5a3 3 0 11-6 0 3 3 0 016 0z" />
                  <path strokeLinecap="round" strokeLinejoin="round" d="M19.5 10.5c0 7.142-7.5 11.25-7.5 11.25S4.5 17.642 4.5 10.5a7.5 7.5 0 1115 0z" />
                </svg>
              </div>
              <div className="space-y-0.5">
                <p className="font-['Plus_Jakarta_Sans'] text-[10px] uppercase tracking-[0.25em] text-stone-400">
                  Venue
                </p>
                <p className="font-['Playfair_Display'] text-lg text-[#1a1a2e]">The Riverside Chapel</p>
                <p className="font-['Plus_Jakarta_Sans'] text-sm text-stone-400">42 River Lane, Willow Creek</p>
              </div>
            </div>
          </div>
        </div>

        {/* ── Dress Code ── */}
        <div className="bg-stone-50/80 p-1 rounded-2xl ring-1 ring-stone-100">
          <div className="bg-white rounded-[calc(1.5rem-0.25rem)] px-6 py-5">
            <div className="flex items-start gap-4">
              <div className="mt-0.5 flex h-9 w-9 shrink-0 items-center justify-center rounded-full bg-amber-50">
                <svg className="h-4 w-4 text-amber-400" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M9.813 15.904L9 18.75l-.813-2.846a4.5 4.5 0 00-3.09-3.09L2.25 12l2.846-.813a4.5 4.5 0 003.09-3.09L9 5.25l.813 2.846a4.5 4.5 0 003.09 3.09L15.75 12l-2.846.813a4.5 4.5 0 00-3.09 3.09zM18.259 8.715L18 9.75l-.259-1.035a3.375 3.375 0 00-2.455-2.456L14.25 6l1.036-.259a3.375 3.375 0 002.455-2.456L18 2.25l.259 1.035a3.375 3.375 0 002.455 2.456L21.75 6l-1.036.259a3.375 3.375 0 00-2.455 2.456z" />
                </svg>
              </div>
              <div className="space-y-0.5">
                <p className="font-['Plus_Jakarta_Sans'] text-[10px] uppercase tracking-[0.25em] text-stone-400">
                  Dress Code
                </p>
                <p className="font-['Playfair_Display'] text-lg text-[#1a1a2e]">Formal Attire</p>
                <p className="font-['Plus_Jakarta_Sans'] text-sm text-stone-400">Black tie optional</p>
              </div>
            </div>
          </div>
        </div>

        {/* ── Reception ── */}
        <div className="text-center">
          <div className="flex items-center justify-center gap-3">
            <span className="block h-px w-6 bg-stone-200" />
            <p className="font-['Cormorant_Garamond'] text-sm italic text-stone-400">
              Reception to follow at the same venue
            </p>
            <span className="block h-px w-6 bg-stone-200" />
          </div>
        </div>

        {/* ── View QR Code ── */}
        <div className="text-center space-y-3">
          <button
            onClick={() => setShowQR(true)}
            disabled={!isMobile}
            className="group relative inline-flex items-center gap-2 rounded-full bg-[#1a1a2e] px-8 py-3.5 text-sm font-medium text-white transition-all duration-500 ease-[cubic-bezier(0.32,0.72,0,1)] hover:bg-[#2a2a4e] active:scale-[0.97] disabled:cursor-not-allowed disabled:opacity-40 disabled:hover:bg-[#1a1a2e] disabled:active:scale-100"
          >
            <svg className="h-4 w-4" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
              <path strokeLinecap="round" strokeLinejoin="round" d="M3.75 4.875c0-.621.504-1.125 1.125-1.125h4.5c.621 0 1.125.504 1.125 1.125v4.5c0 .621-.504 1.125-1.125 1.125h-4.5A1.125 1.125 0 013.75 9.375v-4.5zM3.75 14.625c0-.621.504-1.125 1.125-1.125h4.5c.621 0 1.125.504 1.125 1.125v4.5c0 .621-.504 1.125-1.125 1.125h-4.5a1.125 1.125 0 01-1.125-1.125v-4.5zM13.5 4.875c0-.621.504-1.125 1.125-1.125h4.5c.621 0 1.125.504 1.125 1.125v4.5c0 .621-.504 1.125-1.125 1.125h-4.5A1.125 1.125 0 0113.5 9.375v-4.5z" />
              <path strokeLinecap="round" strokeLinejoin="round" d="M13.5 14.625c0-.621.504-1.125 1.125-1.125h4.5c.621 0 1.125.504 1.125 1.125v4.5c0 .621-.504 1.125-1.125 1.125h-4.5a1.125 1.125 0 01-1.125-1.125v-4.5z" />
            </svg>
            <span>View QR Code</span>
            <span className="flex h-7 w-7 items-center justify-center rounded-full bg-white/10 transition-transform duration-500 ease-[cubic-bezier(0.32,0.72,0,1)] group-hover:translate-x-0.5">
              <svg className="h-3 w-3" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={2}>
                <path strokeLinecap="round" strokeLinejoin="round" d="M4.5 12h15m0 0l-6.75-6.75M19.5 12l-6.75 6.75" />
              </svg>
            </span>
          </button>
          {!isMobile && (
            <p className="text-[11px] text-stone-400">
              Open this page on your phone to access the QR code
            </p>
          )}
        </div>

        {/* ── Footer ── */}
        <p className="text-center font-['Plus_Jakarta_Sans'] text-[9px] uppercase tracking-[0.35em] text-stone-300">
          {id}
        </p>
      </div>

      {showQR && (
        <QRModal
          url={inviteUrl}
          id={id}
          user={user}
          onClose={() => setShowQR(false)}
        />
      )}
    </div>
  );
}
