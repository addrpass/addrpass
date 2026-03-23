const API_URL = process.env.NEXT_PUBLIC_API_URL || "http://localhost:8080";

type FetchOptions = {
  method?: string;
  body?: unknown;
  token?: string;
};

async function apiFetch<T>(path: string, opts: FetchOptions = {}): Promise<T> {
  const headers: Record<string, string> = {
    "Content-Type": "application/json",
  };
  if (opts.token) {
    headers["Authorization"] = `Bearer ${opts.token}`;
  }

  const res = await fetch(`${API_URL}${path}`, {
    method: opts.method || "GET",
    headers,
    body: opts.body ? JSON.stringify(opts.body) : undefined,
  });

  if (!res.ok) {
    const err = await res.json().catch(() => ({ error: res.statusText }));
    throw new Error(err.error || "Request failed");
  }

  if (res.status === 204) return {} as T;
  return res.json();
}

// Auth
export type User = {
  id: string;
  email: string;
  name: string;
  created_at: string;
};

export type AuthResponse = {
  token: string;
  user: User;
};

export function register(email: string, password: string, name: string) {
  return apiFetch<AuthResponse>("/api/v1/auth/register", {
    method: "POST",
    body: { email, password, name },
  });
}

export function login(email: string, password: string) {
  return apiFetch<AuthResponse>("/api/v1/auth/login", {
    method: "POST",
    body: { email, password },
  });
}

export function getMe(token: string) {
  return apiFetch<User>("/api/v1/auth/me", { token });
}

// Addresses
export type Address = {
  id: string;
  user_id: string;
  label: string;
  line1: string;
  line2: string;
  city: string;
  state: string;
  post_code: string;
  country: string;
  phone: string;
  created_at: string;
  updated_at: string;
};

export type CreateAddressInput = {
  label: string;
  line1: string;
  line2?: string;
  city: string;
  state?: string;
  post_code: string;
  country: string;
  phone?: string;
};

export function createAddress(token: string, data: CreateAddressInput) {
  return apiFetch<Address>("/api/v1/addresses", {
    method: "POST",
    body: data,
    token,
  });
}

export function listAddresses(token: string) {
  return apiFetch<Address[]>("/api/v1/addresses", { token });
}

export function deleteAddress(token: string, id: string) {
  return apiFetch<void>(`/api/v1/addresses/${id}`, {
    method: "DELETE",
    token,
  });
}

// Shares
export type Share = {
  id: string;
  address_id: string;
  token: string;
  access_type: string;
  pin: string;
  expires_at: string | null;
  max_accesses: number | null;
  access_count: number;
  active: boolean;
  created_at: string;
};

export type CreateShareInput = {
  address_id: string;
  access_type: "public" | "authenticated";
  pin?: string;
  expires_at?: string;
  max_accesses?: number;
};

export type CreateShareResponse = {
  share: Share;
  url: string;
};

export function createShare(token: string, data: CreateShareInput) {
  return apiFetch<CreateShareResponse>("/api/v1/shares", {
    method: "POST",
    body: data,
    token,
  });
}

export function listShares(token: string) {
  return apiFetch<Share[]>("/api/v1/shares", { token });
}

export function revokeShare(token: string, id: string) {
  return apiFetch<void>(`/api/v1/shares/${id}/revoke`, {
    method: "PATCH",
    token,
  });
}

export function deleteShare(token: string, id: string) {
  return apiFetch<void>(`/api/v1/shares/${id}`, {
    method: "DELETE",
    token,
  });
}

// Access logs
export type AccessLog = {
  id: string;
  share_id: string;
  ip: string;
  user_agent: string;
  country: string;
  access_at: string;
};

export function getAccessLogs(token: string, shareId: string) {
  return apiFetch<AccessLog[]>(`/api/v1/shares/${shareId}/accesses`, { token });
}

// QR Code URL
export function getQRCodeURL(shareToken: string) {
  return `${API_URL}/api/v1/qr/${shareToken}`;
}
