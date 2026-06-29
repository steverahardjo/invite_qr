# QR Invite

Backend system for managing event invitations via QR codes. Guests receive invitations via WhatsApp or email with a unique link, view their invite and QR code, and admins scan the QR code at the door to mark attendance.

## Tech Stack

- **Language:** Go 1.25
- **Database:** PostgreSQL (via `pgx` driver + sqlc code generation)
- **Auth:** JWT (HS256) + bcrypt
- **QR Code:** `skip2/go-qrcode`
- **Messaging:** Meta/Facebook Graph API (WhatsApp), Resend (Email)
- **Logging:** Uber Zap

## Architecture

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│  Guest       │     │  Admin       │     │  Admin       │
│  (browser)   │     │  (browser)   │     │  (scanner)   │
└──────┬───────┘     └──────┬───────┘     └──────┬───────┘
       │                    │                    │
       ▼                    ▼                    ▼
┌───────────────────────────────────────────────────────┐
│                    HTTP Server                         │
│  ┌─────────┐  ┌──────────┐  ┌────────┐  ┌─────────┐  │
│  │ Public  │  │  Auth    │  │ Admin  │  │ Invite  │  │
│  │ Handlers│  │  Handlers│  │Handlers│  │Handlers │  │
│  └────┬────┘  └────┬─────┘  └───┬────┘  └────┬────┘  │
│       │            │            │             │        │
│  ┌────▼────┐  ┌────▼─────┐  ┌───▼────┐  ┌────▼────┐  │
│  │ Public  │  │JwtService│  │ admin  │  │ Invite  │  │
│  │ Service │  │          │  │ Store  │  │ Service │  │
│  └────┬────┘  └──────────┘  └───┬────┘  └────┬────┘  │
│       │                         │             │        │
│       └──────────┬──────────────┴─────────────┘        │
│                  ▼                                     │
│           ┌──────────────┐                              │
│           │  PostgreSQL  │                              │
│           └──────────────┘                              │
└───────────────────────────────────────────────────────┘
```

## Flow

1. **Admin adds participants** via `POST /api/admin/participants`
2. **Admin sends invites** via `POST /api/bulk-invite` or `GET /api/send-invite`
3. **Guest receives invite** — WhatsApp message or email with a unique link containing their `external_id`
4. **Guest opens link** — `GET /api/invite/{external_id}` returns their details, page shows QR code
5. **Guest requests QR code** — `GET /api/qr?participant_id={external_id}` returns a QR code PNG
6. **Admin scans QR code** at door — `POST /api/admin/attendance` marks them as attended
7. **Admin monitors** — `GET /api/admin/participants` lists all guests and check-in status

## Prerequisites

- Go 1.22+ (uses `http.ServeMux` with method-based routing)
- PostgreSQL 14+
- (Optional) WhatsApp API key for WhatsApp delivery
- (Optional) Resend API key for email delivery

## Setup

### 1. Database

```sql
CREATE DATABASE qr_invite;
```

Then run the schema:

```bash
psql -d qr_invite -f db/schema.sql
```

Or use the CSV in `testdata/` to seed sample data:

```bash
psql -d qr_invite -c "\copy participants (name, email, wa_number) FROM 'testdata/participants.csv' WITH CSV HEADER;"
```

### 2. Environment Variables

| Variable | Required | Default | Description |
|---|---|---|---|
| `DB_NAME` | Yes | — | PostgreSQL database name |
| `DB_USER` | Yes | — | Database user |
| `DB_HOST` | Yes | `localhost` | Database host |
| `DB_PORT` | Yes | `5432` | Database port |
| `DB_PASSWORD` | Yes | — | Database password |
| `DB_SSLMODE` | No | `disable` | PostgreSQL SSL mode |
| `JWT_SECRET` | Yes | — | Secret key for signing JWT tokens |
| `ADMIN_PASSWORD` | Yes | — | Admin login password (bcrypt-hashed at startup) |
| `PORT` | No | `8080` | HTTP server port |
| `BASE_WEB_URL` | No | `http://localhost:8080/` | Base URL for invite links |
| `WA_API_KEY` | No* | — | Meta Graph API key for WhatsApp |
| `WA_PHONE` | No* | — | WhatsApp sender phone number |
| `RESEND_API_KEY` | No* | — | Resend API key for email |
| `RESEND_EMAIL` | No* | — | Sender email address |

\* Required only if using the respective delivery channel.

### 3. Run

```bash
# Install dependencies
go mod tidy

# Start the server (with env vars exported or in .env)
go run cmd/main.go
```

## API Reference

### Public Endpoints (no auth required)

#### `GET /api/invite/{token}`

Returns participant details by their external UUID.

**Response `200`**
```json
{
  "id": 1,
  "external_id": "a1b2c3d4-...",
  "name": "Alice Johnson",
  "email": "alice@example.com",
  "wa_number": "+628111111001",
  "accessed": false,
  "sent": true
}
```

**Response `400`**
```json
missing token
```

---

#### `GET /api/user?id={uuid}`

Returns participant details by `id` query parameter.

**Response `200`** — same participant JSON as above.

---

#### `GET /api/qr?participant_id={uuid}`

Generates a 256×256 PNG QR code containing the participant's external UUID. Returns `image/png`.

**Response `200`** — PNG image bytes.

**Response `400`**
```json
missing participant_id
```

---

#### `POST /api/admin/login`

Authenticates admin credentials and returns a JWT token.

**Request**
```json
{
  "username": "admin",
  "password": "your-password"
}
```

**Response `200`**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIs..."
}
```

**Response `401`**
```json
{
  "error": "invalid credentials"
}
```

---

### Admin Endpoints (JWT required)

All admin endpoints require `Authorization: Bearer <token>` header.

#### `GET /api/admin/participants`

Lists up to 100 participants ordered by ID.

**Response `200`**
```json
[
  {
    "id": 1,
    "external_id": "a1b2c3d4-...",
    "name": "Alice Johnson",
    "email": "alice@example.com",
    "wa_number": "+628111111001",
    "accessed": false,
    "sent": true
  }
]
```

---

#### `POST /api/admin/participants`

Adds a new participant.

**Request**
```json
{
  "name": "Charlie Brown",
  "email": "charlie@example.com",
  "wa_number": "+628111111101"
}
```

**Response `201`** — no body.

---

#### `POST /api/admin/attendance`

Marks a participant as attended by scanning their QR code (external UUID). This is the endpoint used by the door scanner.

**Request**
```json
{
  "participant_id": "a1b2c3d4-..."
}
```

**Response `200`** — no body.

**Response `400`**
```json
missing participant_id
```

---

### Invite Endpoints (no auth required)

#### `POST /api/bulk-invite`

Sends invitations to all participants who haven't received one yet. Sends in batches of 50 with a 500ms throttle, via WhatsApp or email depending on what contact info is available.

**Response `200`**

---

#### `GET /api/send-invite?guest_id={id}&email={email}&wa_number={wa}&name={name}`

Sends a one-time invitation to a specific participant.

| Parameter | Required | Description |
|---|---|---|
| `guest_id` | Yes | Participant's internal `id` (not UUID) |
| `email` | No* | Email address to send to |
| `wa_number` | No* | WhatsApp number to send to |
| `name` | No | Override display name |

\* At least one of `email` or `wa_number` is required.

**Response `200`**

**Response `400`**
```json
missing guest_id
```
or
```json
either email or wa_number is required
```

## Project Structure

```
cmd/
  main.go              ← Entry point, wiring, route registration
db/
  schema.sql           ← Database schema (participants table)
  query.sql            ← Named SQL queries for sqlc
  sqlc.yaml            ← sqlc configuration
  db_gen/              ← Auto-generated Go code from sqlc
internal/
  admin/               ← Admin handlers (participant CRUD, attendance)
  auth/                ← JWT auth, login handler, middleware
  config/              ← DB connection, config loading
  invite/              ← Invite service, queue, WhatsApp/email senders
  public/              ← Guest-facing handlers (invite page, QR code)
  server/              ← Logger, WebServer, context helpers
testdata/
  participants.csv     ← 100 sample participants for testing
```
