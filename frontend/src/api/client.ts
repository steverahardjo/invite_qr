import { logApiCall } from "./logger";

const BASE = "/api";

const TOKEN_KEY = "admin_token";

export function getToken(): string | null {
  return localStorage.getItem(TOKEN_KEY);
}

export function setToken(token: string): void {
  localStorage.setItem(TOKEN_KEY, token);
}

export function clearToken(): void {
  localStorage.removeItem(TOKEN_KEY);
}

export function isAuthenticated(): boolean {
  return !!getToken();
}

async function request<T>(path: string, options: RequestInit = {}): Promise<T> {
  const token = getToken();
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
    ...((options.headers as Record<string, string>) || {}),
  };
  if (token) {
    headers["Authorization"] = `Bearer ${token}`;
  }

  const method = (options.method || "GET").toUpperCase();
  logApiCall(method, path, "pending");

  try {
    const res = await fetch(`${BASE}${path}`, { ...options, headers });

    if (!res.ok) {
      const body = await res.json().catch(() => ({ error: res.statusText }));
      const msg = body.error || `Request failed: ${res.status}`;
      logApiCall(method, path, "error", msg);
      throw new Error(msg);
    }

    const ct = res.headers.get("content-type");
    if (ct?.includes("application/json")) {
      const data = await res.json();
      logApiCall(method, path, "success");
      return data;
    }
    logApiCall(method, path, "success");
    return undefined as T;
  } catch (err) {
    if (!(err instanceof Error && err.message.startsWith("Request"))) {
      logApiCall(method, path, "error", (err as Error).message);
    }
    throw err;
  }
}

export interface Participant {
  id: number;
  external_id: string;
  name: string;
  email: string;
  wa_number: string;
  accessed: boolean;
  sent: boolean;
}

export interface LoginResponse {
  token: string;
}

export interface ParticipantInput {
  name: string;
  email: string;
  wa_number: string;
}

export const api = {
  login: (username: string, password: string) =>
    request<LoginResponse>("/admin/login", {
      method: "POST",
      body: JSON.stringify({ username, password }),
    }),

  getParticipants: () => request<Participant[]>("/admin/participants"),

  addParticipant: (data: ParticipantInput) =>
    request<void>("/admin/participants", {
      method: "POST",
      body: JSON.stringify(data),
    }),

  markAttendance: (participantId: string) =>
    request<void>("/admin/attendance", {
      method: "POST",
      body: JSON.stringify({ participant_id: participantId }),
    }),

  getInvite: (externalId: string) =>
    request<Participant>(`/invite/${externalId}`),

  sendInvite: (participant: { id: number; email: string; wa_number: string; name: string }) => {
    const params = new URLSearchParams({
      guest_id: String(participant.id),
      email: participant.email,
      wa_number: participant.wa_number,
      name: participant.name,
    });
    return request<void>(`/send-invite?${params.toString()}`);
  },

  getQR: (participantId: string) =>
    `${BASE}/qr?participant_id=${participantId}`,
};
