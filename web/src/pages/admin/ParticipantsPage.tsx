import { useState, useEffect, useRef, useCallback } from "react";
import { useNavigate } from "react-router-dom";
import { ArrowLeft, Eye, Send, UsersRound, Search, RotateCcw } from "lucide-react";

interface Participant {
  id: number;
  name: string;
  email: string | null;
  wa_number: string | null;
  accessed: boolean;
  sent: boolean;
}

const PAGE_SIZE = 20;

export function ParticipantsPage() {
  const navigate = useNavigate();
  const [items, setItems] = useState<Participant[]>([]);
  const [page, setPage] = useState(1);
  const [hasNext, setHasNext] = useState(true);
  const [loading, setLoading] = useState(false);
  const [stats, setStats] = useState<{ total: number; seen: number; sent: number } | null>(null);
  const sentinelRef = useRef<HTMLDivElement>(null);

  // filters
  const [search, setSearch] = useState("");
  const [sentFilter, setSentFilter] = useState("");
  const [accessedFilter, setAccessedFilter] = useState("");
  const [resending, setResending] = useState<number | null>(null);

  // debounced search
  const [debouncedSearch, setDebouncedSearch] = useState("");
  useEffect(() => {
    const t = setTimeout(() => setDebouncedSearch(search), 300);
    return () => clearTimeout(t);
  }, [search]);

  // reset page when filters change
  useEffect(() => {
    setPage(1);
    setItems([]);
  }, [debouncedSearch, sentFilter, accessedFilter]);

  // build query string
  const qs = useCallback(() => {
    const p = new URLSearchParams();
    p.set("page", String(page));
    p.set("limit", String(PAGE_SIZE));
    if (debouncedSearch) p.set("search", debouncedSearch);
    if (sentFilter) p.set("sent", sentFilter);
    if (accessedFilter) p.set("accessed", accessedFilter);
    return p.toString();
  }, [page, debouncedSearch, sentFilter, accessedFilter]);

  // fetch summary (with filters)
  useEffect(() => {
    const p = new URLSearchParams();
    if (debouncedSearch) p.set("search", debouncedSearch);
    fetch(`/api/admin/summary?${p}`)
      .then((r) => r.json())
      .then(setStats)
      .catch(() => {});
  }, [debouncedSearch]);

  const loadPage = useCallback(async (p: number) => {
    setLoading(true);
    try {
      const params = new URLSearchParams();
      params.set("page", String(p));
      params.set("limit", String(PAGE_SIZE));
      if (debouncedSearch) params.set("search", debouncedSearch);
      if (sentFilter) params.set("sent", sentFilter);
      if (accessedFilter) params.set("accessed", accessedFilter);

      const res = await fetch(`/api/admin/participants?${params}`);
      const json = await res.json();
      setItems((prev) => (p === 1 ? json.data : [...prev, ...json.data]));
      setHasNext(json.has_next);
    } catch {
      // ignore
    } finally {
      setLoading(false);
    }
  }, [debouncedSearch, sentFilter, accessedFilter]);

  useEffect(() => {
    loadPage(1);
  }, [loadPage]);

  useEffect(() => {
    const el = sentinelRef.current;
    if (!el) return;

    const observer = new IntersectionObserver(
      (entries) => {
        if (entries[0]?.isIntersecting && hasNext && !loading) {
          setPage((prev) => {
            const next = prev + 1;
            loadPage(next);
            return next;
          });
        }
      },
      { rootMargin: "200px" },
    );

    observer.observe(el);
    return () => observer.disconnect();
  }, [hasNext, loading, loadPage]);

  const handleResend = async (id: number) => {
    setResending(id);
    try {
      await fetch("/api/admin/resend", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ id }),
      });
      setItems((prev) => prev.map((p) => (p.id === id ? { ...p, sent: true } : p)));
    } catch {
      // ignore
    } finally {
      setResending(null);
    }
  };

  return (
    <div className="min-h-screen w-full bg-white">
      <div className="mx-auto max-w-5xl px-4 py-8">
        <button
          onClick={() => navigate("/admin")}
          className="mb-4 flex items-center gap-1.5 text-xs text-stone-400 hover:text-stone-600 transition-colors"
        >
          <ArrowLeft className="h-3.5 w-3.5" />
          Back
        </button>

        <p className="font-['Plus_Jakarta_Sans'] text-xs uppercase tracking-[0.25em] text-stone-400 text-center mb-6">
          Participants
        </p>

        {/* ── Stats ── */}
        <div className="grid grid-cols-3 gap-3 mb-6 max-w-sm mx-auto">
          <div className="bg-stone-50/80 p-1 rounded-2xl ring-1 ring-stone-100">
            <div className="bg-white rounded-[calc(1.5rem-0.25rem)] px-3 py-3 text-center space-y-1">
              <UsersRound className="h-4 w-4 text-stone-400 mx-auto" />
              <p className="font-['Plus_Jakarta_Sans'] text-lg font-medium text-[#1a1a2e]">
                {stats?.total ?? "—"}
              </p>
              <p className="font-['Plus_Jakarta_Sans'] text-[9px] uppercase tracking-[0.15em] text-stone-400">
                Total
              </p>
            </div>
          </div>
          <div className="bg-stone-50/80 p-1 rounded-2xl ring-1 ring-stone-100">
            <div className="bg-white rounded-[calc(1.5rem-0.25rem)] px-3 py-3 text-center space-y-1">
              <Eye className="h-4 w-4 text-emerald-400 mx-auto" />
              <p className="font-['Plus_Jakarta_Sans'] text-lg font-medium text-[#1a1a2e]">
                {stats?.seen ?? "—"}
              </p>
              <p className="font-['Plus_Jakarta_Sans'] text-[9px] uppercase tracking-[0.15em] text-stone-400">
                Scanned
              </p>
            </div>
          </div>
          <div className="bg-stone-50/80 p-1 rounded-2xl ring-1 ring-stone-100">
            <div className="bg-white rounded-[calc(1.5rem-0.25rem)] px-3 py-3 text-center space-y-1">
              <Send className="h-4 w-4 text-amber-400 mx-auto" />
              <p className="font-['Plus_Jakarta_Sans'] text-lg font-medium text-[#1a1a2e]">
                {stats?.sent ?? "—"}
              </p>
              <p className="font-['Plus_Jakarta_Sans'] text-[9px] uppercase tracking-[0.15em] text-stone-400">
                Sent
              </p>
            </div>
          </div>
        </div>

        {/* ── Filters ── */}
        <div className="flex flex-wrap items-center gap-2 mb-4">
          <div className="relative flex-1 min-w-[180px]">
            <Search className="absolute left-3 top-1/2 -translate-y-1/2 h-3.5 w-3.5 text-stone-300" />
            <input
              type="text"
              value={search}
              onChange={(e) => setSearch(e.target.value)}
              placeholder="Search name or email..."
              className="w-full rounded-full border border-stone-200 bg-white pl-8 pr-4 py-2 text-xs text-stone-600 placeholder:text-stone-300 focus:border-stone-400 focus:outline-none"
            />
          </div>

          <select
            value={sentFilter}
            onChange={(e) => setSentFilter(e.target.value)}
            className="rounded-full border border-stone-200 bg-white px-3 py-2 text-xs text-stone-500 focus:border-stone-400 focus:outline-none"
          >
            <option value="">All invites</option>
            <option value="true">Sent</option>
            <option value="false">Pending</option>
          </select>

          <select
            value={accessedFilter}
            onChange={(e) => setAccessedFilter(e.target.value)}
            className="rounded-full border border-stone-200 bg-white px-3 py-2 text-xs text-stone-500 focus:border-stone-400 focus:outline-none"
          >
            <option value="">All scans</option>
            <option value="true">Scanned</option>
            <option value="false">Not scanned</option>
          </select>
        </div>

        {/* ── Table ── */}
        <div className="overflow-x-auto rounded-2xl ring-1 ring-stone-100">
          <table className="w-full text-left text-sm">
            <thead>
              <tr className="bg-stone-50 text-[10px] uppercase tracking-[0.15em] text-stone-400">
                <th className="px-4 py-3 font-medium">Name</th>
                <th className="px-4 py-3 font-medium">Email</th>
                <th className="px-4 py-3 font-medium">WA Number</th>
                <th className="px-4 py-3 font-medium text-center">Invite</th>
                <th className="px-4 py-3 font-medium text-center">QR Scan</th>
              </tr>
            </thead>
            <tbody className="divide-y divide-stone-100">
              {items.map((p) => (
                <tr key={p.id} className="hover:bg-stone-50/50 transition-colors">
                  <td className="px-4 py-3 font-['Playfair_Display'] text-[#1a1a2e] capitalize whitespace-nowrap">
                    {p.name}
                  </td>
                  <td className="px-4 py-3 text-stone-500">
                    {p.email ?? <span className="italic text-stone-300">—</span>}
                  </td>
                  <td className="px-4 py-3 text-stone-500">
                    {p.wa_number ?? <span className="italic text-stone-300">—</span>}
                  </td>
                  <td className="px-4 py-3 text-center">
                    {p.sent ? (
                      <span className="inline-flex items-center gap-1 rounded-full bg-amber-50 px-2.5 py-0.5 text-[10px] font-medium text-amber-600">
                        <Send className="h-2.5 w-2.5" />
                        Sent
                      </span>
                    ) : (
                      <button
                        onClick={() => handleResend(p.id)}
                        disabled={resending === p.id}
                        className="inline-flex items-center gap-1 rounded-full bg-stone-100 px-2.5 py-0.5 text-[10px] font-medium text-stone-500 hover:bg-stone-200 hover:text-stone-700 transition-colors disabled:opacity-50"
                      >
                        <RotateCcw className={`h-2.5 w-2.5 ${resending === p.id ? "animate-spin" : ""}`} />
                        {resending === p.id ? "Sending..." : "Resend"}
                      </button>
                    )}
                  </td>
                  <td className="px-4 py-3 text-center">
                    {p.accessed ? (
                      <span className="inline-flex items-center gap-1 rounded-full bg-emerald-50 px-2.5 py-0.5 text-[10px] font-medium text-emerald-600">
                        Scanned
                      </span>
                    ) : (
                      <span className="inline-flex items-center gap-1 rounded-full bg-stone-100 px-2.5 py-0.5 text-[10px] font-medium text-stone-400">
                        Not scanned
                      </span>
                    )}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        {/* sentinel */}
        <div ref={sentinelRef} className="h-4" />

        {loading && (
          <div className="flex justify-center py-6">
            <div className="h-5 w-5 animate-spin rounded-full border-2 border-stone-200 border-t-[#1a1a2e]" />
          </div>
        )}

        {!hasNext && items.length > 0 && (
          <p className="text-center text-[10px] text-stone-300 uppercase tracking-[0.2em] pt-6">
            All {stats?.total} participants loaded
          </p>
        )}
      </div>
    </div>
  );
}
