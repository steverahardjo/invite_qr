import { useNavigate } from "react-router-dom";
import { QrCode, Users } from "lucide-react";

export function AdminDashboard() {
  const navigate = useNavigate();

  const cards = [
    {
      label: "QR Code Scanner",
      desc: "Scan participant QR codes to mark attendance",
      icon: QrCode,
      href: "/admin/qr_code",
    },
    {
      label: "Participants",
      desc: "View and manage all invited guests",
      icon: Users,
      href: "/admin/participants",
    },
  ];

  return (
    <div className="min-h-screen w-full bg-white">
      <div className="mx-auto max-w-lg px-4 py-12">
        <p className="font-['Plus_Jakarta_Sans'] text-xs uppercase tracking-[0.25em] text-stone-400 text-center mb-8">
          Admin Dashboard
        </p>

        <div className="space-y-3">
          {cards.map((c) => (
            <button
              key={c.href}
              onClick={() => navigate(c.href)}
              className="group w-full text-left bg-stone-50/80 p-1 rounded-2xl ring-1 ring-stone-100 transition-all duration-300 ease-[cubic-bezier(0.32,0.72,0,1)] hover:ring-amber-300/40 hover:shadow-md active:scale-[0.99]"
            >
              <div className="bg-white rounded-[calc(1.5rem-0.25rem)] px-5 py-4 flex items-center gap-4">
                <div className="flex h-10 w-10 shrink-0 items-center justify-center rounded-full bg-amber-50 group-hover:bg-amber-100 transition-colors duration-300">
                  <c.icon className="h-4 w-4 text-amber-400" />
                </div>
                <div className="flex-1 min-w-0">
                  <p className="font-['Playfair_Display'] text-lg text-[#1a1a2e]">
                    {c.label}
                  </p>
                  <p className="font-['Plus_Jakarta_Sans'] text-xs text-stone-400 mt-0.5">
                    {c.desc}
                  </p>
                </div>
                <svg className="h-4 w-4 shrink-0 text-stone-300 transition-transform duration-300 ease-[cubic-bezier(0.32,0.72,0,1)] group-hover:translate-x-0.5 group-hover:text-stone-500" fill="none" viewBox="0 0 24 24" stroke="currentColor" strokeWidth={1.5}>
                  <path strokeLinecap="round" strokeLinejoin="round" d="M8.25 4.5l7.5 7.5-7.5 7.5" />
                </svg>
              </div>
            </button>
          ))}
        </div>
      </div>
    </div>
  );
}
