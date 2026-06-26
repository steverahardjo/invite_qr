import { serve } from "bun";
import index from "./index.html";

function allParticipants() {
  const names = ["Alice Johnson", "Bob Smith", "Carol Davis", "David Lee", "Eva Martinez", "Frank Wilson", "Grace Kim", "Henry Brown", "Iris Chen", "Jack Taylor", "Karen White", "Leo Harris", "Mia Clark", "Noah Lewis", "Olivia Walker", "Peter Hall", "Quinn Allen", "Rachel Young", "Sam King", "Tina Wright", "Uma Scott", "Victor Adams", "Wendy Baker", "Xander Green", "Yara Nelson", "Zach Hill"];
  const all = [];
  for (let i = 0; i < 86; i++) {
    const name = names[i % names.length] ?? `Guest ${i + 1}`;
    const emailRequired = i % 3 !== 0;
    const waRequired = i % 4 !== 0;
    all.push({
      id: i + 1,
      name,
      email: emailRequired || !waRequired ? `${name.toLowerCase().replace(/\s+/g, ".")}@email.com` : null,
      wa_number: waRequired || !emailRequired ? `+628${String(1000000000 + i).slice(0, 9)}` : null,
      accessed: i % 3 === 0,
      sent: i % 5 !== 0,
    });
  }
  return all;
}

const server = serve({
  routes: {
    "/*": index,

    "/api/hello": {
      async GET() {
        return Response.json({ message: "Hello, world!", method: "GET" });
      },
      async PUT() {
        return Response.json({ message: "Hello, world!", method: "PUT" });
      },
    },

    "/api/admin/login": {
      async POST(req) {
        const text = await req.text();
        const params = new URLSearchParams(text);
        const password = params.get("password");
        if (password === "admin123") {
          return Response.json({ ok: true });
        }
        return Response.json({ ok: false }, { status: 401 });
      },
    },

    "/api/admin/summary": {
      async GET(req) {
        const url = new URL(req.url);
        const search = url.searchParams.get("search")?.toLowerCase();
        const all = allParticipants();
        const filtered = search
          ? all.filter((p: any) => p.name.toLowerCase().includes(search) || (p.email && p.email.toLowerCase().includes(search)))
          : all;
        return Response.json({
          total: filtered.length,
          seen: filtered.filter((p: any) => p.accessed).length,
          sent: filtered.filter((p: any) => p.sent).length,
        });
      },
    },

    "/api/admin/resend": {
      async POST(req) {
        const { id } = await req.json();
        return Response.json({ ok: true, id, message: `Invite resent to participant #${id}` });
      },
    },

    "/api/admin/participants": {
      async GET(req) {
        const url = new URL(req.url);
        const page = parseInt(url.searchParams.get("page") ?? "1");
        const limit = parseInt(url.searchParams.get("limit") ?? "20");
        const search = url.searchParams.get("search")?.toLowerCase();
        const sentFilter = url.searchParams.get("sent");
        const accessedFilter = url.searchParams.get("accessed");

        let all = allParticipants();

        if (search) all = all.filter((p: any) => p.name.toLowerCase().includes(search) || (p.email && p.email.toLowerCase().includes(search)));
        if (sentFilter === "true") all = all.filter((p: any) => p.sent);
        if (sentFilter === "false") all = all.filter((p: any) => !p.sent);
        if (accessedFilter === "true") all = all.filter((p: any) => p.accessed);
        if (accessedFilter === "false") all = all.filter((p: any) => !p.accessed);

        const offset = (page - 1) * limit;
        return Response.json({
          data: all.slice(offset, offset + limit),
          has_next: offset + limit < all.length,
        });
      },
    },

    "/api/hello/:name": async req => {
      const name = req.params.name;
      return Response.json({ message: `Hello, ${name}!` });
    },
  },

  development: process.env.NODE_ENV !== "production" && {
    hmr: true,
    console: true,
  },
});

console.log(`🚀 Server running at ${server.url}`);
