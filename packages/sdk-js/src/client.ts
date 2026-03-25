import type {
  AddrPassConfig,
  Address,
  Share,
  CreateShareInput,
  CreateShareResponse,
  ResolveResponse,
  AccessLog,
  OAuthTokenResponse,
  TokenExchangeResponse,
  Label,
  LabelResponse,
  Delegation,
  ShareScope,
} from "./types";

const DEFAULT_BASE_URL = "https://api.addrpass.com";

export class AddrPassClient {
  private baseURL: string;
  private token: string | null = null;

  constructor(config: AddrPassConfig = {}) {
    this.baseURL = (config.baseURL || DEFAULT_BASE_URL).replace(/\/$/, "");

    // If client credentials provided, authenticate on first use
    if (config.clientId && config.clientSecret) {
      this._clientId = config.clientId;
      this._clientSecret = config.clientSecret;
    }
    if (config.apiKey) {
      this.token = config.apiKey;
    }
  }

  private _clientId?: string;
  private _clientSecret?: string;
  private _tokenExpiry?: number;

  private async getToken(): Promise<string> {
    if (this.token && (!this._tokenExpiry || Date.now() < this._tokenExpiry)) {
      return this.token;
    }
    if (this._clientId && this._clientSecret) {
      const resp = await this.authenticateClient();
      this.token = resp.access_token;
      this._tokenExpiry = Date.now() + (resp.expires_in - 60) * 1000;
      return this.token;
    }
    throw new Error("No authentication configured. Provide apiKey or clientId+clientSecret.");
  }

  private async request<T>(method: string, path: string, body?: unknown, requiresAuth = true): Promise<T> {
    const headers: Record<string, string> = { "Content-Type": "application/json" };
    if (requiresAuth) {
      headers["Authorization"] = `Bearer ${await this.getToken()}`;
    }

    const res = await fetch(`${this.baseURL}${path}`, {
      method,
      headers,
      body: body ? JSON.stringify(body) : undefined,
    });

    if (!res.ok) {
      const err = await res.json().catch(() => ({ error: res.statusText }));
      throw Object.assign(new Error(err.error || "Request failed"), { status: res.status });
    }

    if (res.status === 204) return {} as T;
    return res.json();
  }

  // ─── Authentication ─────────────────────────────────────

  /** Authenticate using client credentials (OAuth2). */
  async authenticateClient(): Promise<OAuthTokenResponse> {
    return this.request<OAuthTokenResponse>("POST", "/api/v1/oauth/token", {
      grant_type: "client_credentials",
      client_id: this._clientId,
      client_secret: this._clientSecret,
    }, false);
  }

  /** Exchange an authorization code for a token + share token. */
  async exchangeCode(code: string, redirectUri: string): Promise<TokenExchangeResponse> {
    return this.request<TokenExchangeResponse>("POST", "/api/v1/oauth/token", {
      grant_type: "authorization_code",
      code,
      client_id: this._clientId,
      client_secret: this._clientSecret,
      redirect_uri: redirectUri,
    }, false);
  }

  // ─── Resolve (Public) ───────────────────────────────────

  /** Resolve a share token to an address. No auth required for public shares. */
  async resolve(token: string, pin?: string): Promise<ResolveResponse> {
    const query = pin ? `?pin=${encodeURIComponent(pin)}` : "";
    return this.request<ResolveResponse>("GET", `/api/v1/resolve/${token}${query}`, undefined, false);
  }

  /** Get QR code URL for a share token. */
  qrCodeURL(token: string): string {
    return `${this.baseURL}/api/v1/qr/${token}`;
  }

  // ─── Addresses ──────────────────────────────────────────

  async listAddresses(): Promise<Address[]> {
    return this.request<Address[]>("GET", "/api/v1/addresses");
  }

  async createAddress(data: Omit<Address, "id" | "user_id" | "created_at" | "updated_at">): Promise<Address> {
    return this.request<Address>("POST", "/api/v1/addresses", data);
  }

  async getAddress(id: string): Promise<Address> {
    return this.request<Address>("GET", `/api/v1/addresses/${id}`);
  }

  async updateAddress(id: string, data: Partial<Address>): Promise<Address> {
    return this.request<Address>("PUT", `/api/v1/addresses/${id}`, data);
  }

  async deleteAddress(id: string): Promise<void> {
    await this.request<void>("DELETE", `/api/v1/addresses/${id}`);
  }

  // ─── Shares ─────────────────────────────────────────────

  async createShare(data: CreateShareInput): Promise<CreateShareResponse> {
    return this.request<CreateShareResponse>("POST", "/api/v1/shares", data);
  }

  async listShares(): Promise<Share[]> {
    return this.request<Share[]>("GET", "/api/v1/shares");
  }

  async revokeShare(id: string): Promise<void> {
    await this.request<void>("PATCH", `/api/v1/shares/${id}/revoke`);
  }

  async deleteShare(id: string): Promise<void> {
    await this.request<void>("DELETE", `/api/v1/shares/${id}`);
  }

  async getAccessLogs(shareId: string): Promise<AccessLog[]> {
    return this.request<AccessLog[]>("GET", `/api/v1/shares/${shareId}/accesses`);
  }

  // ─── Labels ─────────────────────────────────────────────

  async createLabel(shareId: string): Promise<LabelResponse> {
    return this.request<LabelResponse>("POST", "/api/v1/labels", { share_id: shareId });
  }

  labelImageURL(referenceCode: string): string {
    return `${this.baseURL}/api/v1/labels/${referenceCode}/image`;
  }

  // ─── Delegations ────────────────────────────────────────

  async createDelegation(data: {
    share_id: string;
    to_business_id: string;
    scope?: ShareScope;
    expires_at?: string;
    max_accesses?: number;
    note?: string;
  }): Promise<Delegation> {
    return this.request<Delegation>("POST", "/api/v1/delegations", data);
  }

  async listDelegations(shareId: string): Promise<Delegation[]> {
    return this.request<Delegation[]>("GET", `/api/v1/shares/${shareId}/delegations`);
  }

  async revokeDelegation(id: string): Promise<void> {
    await this.request<void>("PATCH", `/api/v1/delegations/${id}/revoke`);
  }

  // ─── OAuth Helpers ──────────────────────────────────────

  /** Generate the authorization URL for the consent flow. */
  getAuthorizationURL(params: {
    redirectUri: string;
    scope?: ShareScope;
    state?: string;
  }): string {
    const appURL = this.baseURL.replace("api.", "");
    const qs = new URLSearchParams({
      client_id: this._clientId || "",
      redirect_uri: params.redirectUri,
      scope: params.scope || "delivery",
      ...(params.state ? { state: params.state } : {}),
    });
    return `${appURL}/authorize?${qs.toString()}`;
  }
}
