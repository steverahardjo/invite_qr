# Wedding QR Invite — Frontend

React + TypeScript frontend for the QR Invite system. Guests view their invitation and QR code via a unique link; admins manage the guest list, send invites, and scan QR codes at the door.

## Tech Stack

- **Runtime:** Bun
- **Framework:** React 19 + TypeScript
- **Bundler:** Vite
- **Routing:** React Router v7
- **QR Scanning:** html5-qrcode (browser camera)

## Setup

```bash
cd frontend
bun install
```

## Run (development)

```bash
bun run dev
```

Starts Vite on `http://localhost:5173`. The Vite dev server proxies `/api/*` requests to the Go backend at `http://localhost:8080` (configured in `vite.config.ts`).

## Build (production)

```bash
bun run build
```

Outputs static files to `dist/`.

## Routes

| Path                        | Auth | Description                                |
|-----------------------------|------|--------------------------------------------|
| `/invite/:external_id`      | No   | Guest's invitation page with QR code       |
| `/admin/login`              | No   | Admin sign-in                              |
| `/admin/dashboard`          | Yes  | Guest table, stats, scanner, add guest     |

## Project Structure

```
src/
  api/
    client.ts      — API client (fetch wrapper, auth helpers)
    logger.ts      — Debug logger (console + in-app panel)
  components/
    DebugPanel.tsx — Floating API log viewer (admin only)
    QRScanner.tsx  — Live camera QR scanner
  pages/
    InvitePage.tsx      — Public guest invitation
    AdminLogin.tsx      — Admin login form
    AdminDashboard.tsx  — Admin dashboard (table, stats, scanner)
  styles/
    global.css    — Wedding Editorial Luxury design system
  App.tsx         — Route definitions + admin layout
  main.tsx        — Entry point with BrowserRouter
```

## API Debug Logging

Every API call is logged to the browser console with a styled timestamp, method, path, and status. On the dashboard, a floating button (`∞`) opens an in-app debug panel showing the last 50 requests with color-coded success/error/pending states.

## Design

Wedding Editorial Luxury theme — Playfair Display (headings), Plus Jakarta Sans (body), warm ivory palette with champagne gold accents. Cards use a double-bezel nested architecture throughout.

## Backend

The Go backend lives one directory up. See `../README.md` for the full API reference and setup instructions.
