import { useState } from "react";

function PasswordForm({ onSuccess }: { onSuccess: () => void }) {
  const [value, setValue] = useState("");
  const [error, setError] = useState(false);
  const [loading, setLoading] = useState(false);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    setLoading(true);
    setError(false);

    try {
      const res = await fetch("/api/admin/login", {
        method: "POST",
        headers: { "Content-Type": "application/x-www-form-urlencoded" },
        body: `password=${encodeURIComponent(value)}`,
      });

      if (res.ok) {
        sessionStorage.setItem("admin_auth", "true");
        onSuccess();
      } else {
        setError(true);
      }
    } catch {
      setError(true);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="min-h-screen w-full flex items-center justify-center bg-white p-4">
      <form
        onSubmit={handleSubmit}
        className="w-full max-w-xs space-y-4 text-center"
      >
        <p className="font-['Plus_Jakarta_Sans'] text-xs uppercase tracking-[0.25em] text-stone-400">
          Admin Access
        </p>
        <input
          type="password"
          value={value}
          onChange={(e) => {
            setValue(e.target.value);
            setError(false);
          }}
          placeholder="Enter password"
          className="w-full rounded-full border border-stone-200 bg-white px-5 py-2.5 text-sm text-center placeholder:text-stone-300 focus:border-stone-400 focus:outline-none"
          autoFocus
        />
        {error && (
          <p className="text-xs text-red-400">Incorrect password</p>
        )}
        <button
          type="submit"
          disabled={loading}
          className="w-full rounded-full bg-[#1a1a2e] px-6 py-2.5 text-sm text-white transition-all hover:bg-[#2a2a4e] disabled:opacity-50 disabled:cursor-not-allowed"
        >
          {loading ? "Verifying..." : "Enter"}
        </button>
      </form>
    </div>
  );
}

export function AdminGuard({ children }: { children: React.ReactNode }) {
  const [authed, setAuthed] = useState(
    () => sessionStorage.getItem("admin_auth") === "true",
  );

  if (!authed) {
    return <PasswordForm onSuccess={() => setAuthed(true)} />;
  }

  return <>{children}</>;
}
