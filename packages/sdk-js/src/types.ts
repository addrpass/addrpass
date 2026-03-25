export type ShareScope = "full" | "delivery" | "zone" | "verify";
export type ShareAccess = "public" | "authenticated";

export interface User {
  id: string;
  email: string;
  name: string;
  created_at: string;
  updated_at: string;
}

export interface Address {
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
}

export interface Share {
  id: string;
  address_id: string;
  user_id: string;
  token: string;
  access_type: ShareAccess;
  scope: ShareScope;
  pin: string;
  expires_at: string | null;
  max_accesses: number | null;
  access_count: number;
  active: boolean;
  created_at: string;
}

export interface AccessLog {
  id: string;
  share_id: string;
  ip: string;
  user_agent: string;
  country: string;
  business_id: string;
  business_name: string;
  access_at: string;
}

export interface ResolveResponse {
  address: Address;
  scope: string;
}

export interface CreateShareInput {
  address_id: string;
  access_type?: ShareAccess;
  scope?: ShareScope;
  pin?: string;
  expires_at?: string;
  max_accesses?: number;
}

export interface CreateShareResponse {
  share: Share;
  url: string;
}

export interface Business {
  id: string;
  name: string;
  owner_id: string;
  created_at: string;
}

export interface APIKey {
  id: string;
  business_id: string;
  client_id: string;
  name: string;
  scopes: string[];
  rate_limit_per_hour: number;
  active: boolean;
  created_at: string;
}

export interface OAuthTokenResponse {
  access_token: string;
  token_type: string;
  expires_in: number;
  scope: string;
}

export interface TokenExchangeResponse extends OAuthTokenResponse {
  share_token: string;
}

export interface Label {
  id: string;
  share_id: string;
  reference_code: string;
  zone_code: string;
  format: string;
  created_at: string;
}

export interface LabelResponse {
  label: Label;
  qr_code_url: string;
}

export interface Delegation {
  id: string;
  share_id: string;
  from_business_id: string;
  to_business_id: string;
  scope: ShareScope;
  expires_at: string | null;
  max_accesses: number | null;
  access_count: number;
  active: boolean;
  note: string;
  created_at: string;
}

export interface AddrPassConfig {
  baseURL?: string;
  apiKey?: string;
  clientId?: string;
  clientSecret?: string;
}

export interface AddrPassError {
  error: string;
  status: number;
}
