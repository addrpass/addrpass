import React, { useCallback, useState, useEffect } from "react";
import type { ShareScope, TokenExchangeResponse } from "./types";

const DEFAULT_APP_URL = "https://addrpass.com";
const DEFAULT_API_URL = "https://api.addrpass.com";

// ─── Types ──────────────────────────────────────────────────

export interface AddrPassProviderProps {
  clientId: string;
  clientSecret?: string;
  redirectUri: string;
  scope?: ShareScope;
  appURL?: string;
  apiURL?: string;
  children: React.ReactNode;
}

export interface AddrPassContextValue {
  authorize: () => void;
  exchangeCode: (code: string) => Promise<TokenExchangeResponse>;
  loading: boolean;
  error: string | null;
  result: TokenExchangeResponse | null;
  reset: () => void;
}

export interface AddrPassButtonProps {
  onToken?: (result: TokenExchangeResponse) => void;
  onError?: (error: string) => void;
  scope?: ShareScope;
  className?: string;
  style?: React.CSSProperties;
  children?: React.ReactNode;
}

// ─── Context ────────────────────────────────────────────────

const AddrPassContext = React.createContext<AddrPassContextValue | null>(null);

export function AddrPassProvider({
  clientId,
  clientSecret,
  redirectUri,
  scope = "delivery",
  appURL = DEFAULT_APP_URL,
  apiURL = DEFAULT_API_URL,
  children,
}: AddrPassProviderProps) {
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [result, setResult] = useState<TokenExchangeResponse | null>(null);

  const authorize = useCallback(() => {
    const state = Math.random().toString(36).substring(2, 14);
    const params = new URLSearchParams({
      client_id: clientId,
      redirect_uri: redirectUri,
      scope,
      state,
    });
    const url = `${appURL}/authorize?${params.toString()}`;

    const w = 600, h = 700;
    const left = (screen.width - w) / 2;
    const top = (screen.height - h) / 2;
    window.open(url, "addrpass_consent", `width=${w},height=${h},left=${left},top=${top}`);
  }, [clientId, redirectUri, scope, appURL]);

  const exchangeCode = useCallback(async (code: string): Promise<TokenExchangeResponse> => {
    setLoading(true);
    setError(null);
    try {
      const res = await fetch(`${apiURL}/api/v1/oauth/token`, {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          grant_type: "authorization_code",
          code,
          client_id: clientId,
          client_secret: clientSecret,
          redirect_uri: redirectUri,
        }),
      });
      if (!res.ok) {
        const err = await res.json().catch(() => ({ error: "Exchange failed" }));
        throw new Error(err.error);
      }
      const data: TokenExchangeResponse = await res.json();
      setResult(data);
      return data;
    } catch (e: unknown) {
      const msg = e instanceof Error ? e.message : "Unknown error";
      setError(msg);
      throw e;
    } finally {
      setLoading(false);
    }
  }, [apiURL, clientId, clientSecret, redirectUri]);

  const reset = useCallback(() => {
    setResult(null);
    setError(null);
    setLoading(false);
  }, []);

  // Listen for postMessage from consent popup
  useEffect(() => {
    function handler(event: MessageEvent) {
      if (event.data?.code) {
        exchangeCode(event.data.code).catch(() => {});
      }
      if (event.data?.error) {
        setError(event.data.error);
      }
    }
    window.addEventListener("message", handler);
    return () => window.removeEventListener("message", handler);
  }, [exchangeCode]);

  return (
    <AddrPassContext.Provider value={{ authorize, exchangeCode, loading, error, result, reset }}>
      {children}
    </AddrPassContext.Provider>
  );
}

// ─── Hook ───────────────────────────────────────────────────

export function useAddrPass(): AddrPassContextValue {
  const ctx = React.useContext(AddrPassContext);
  if (!ctx) {
    throw new Error("useAddrPass must be used within an <AddrPassProvider>");
  }
  return ctx;
}

// ─── Button Component ───────────────────────────────────────

const defaultButtonStyle: React.CSSProperties = {
  display: "inline-flex",
  alignItems: "center",
  gap: 8,
  padding: "10px 20px",
  borderRadius: 8,
  border: "1px solid #E2E8F0",
  background: "#fff",
  color: "#0F172A",
  fontFamily: "-apple-system, BlinkMacSystemFont, sans-serif",
  fontSize: 14,
  fontWeight: 600,
  cursor: "pointer",
  transition: "all 0.2s",
};

function ShieldIcon() {
  return (
    <svg width="16" height="16" viewBox="0 0 32 32" fill="none">
      <path d="M16 2L4 8v8c0 8.4 5.12 16.24 12 18 6.88-1.76 12-9.6 12-18V8L16 2z" fill="#0F172A" />
      <circle cx="16" cy="13" r="3" fill="#22D3EE" />
      <path d="M14 15.5L13 22h6l-1-6.5" fill="#22D3EE" opacity="0.7" />
    </svg>
  );
}

export function AddrPassButton({
  onToken,
  onError,
  className,
  style,
  children,
}: AddrPassButtonProps) {
  const { authorize, result, error } = useAddrPass();

  useEffect(() => {
    if (result && onToken) onToken(result);
  }, [result, onToken]);

  useEffect(() => {
    if (error && onError) onError(error);
  }, [error, onError]);

  return (
    <button
      type="button"
      onClick={authorize}
      className={className}
      style={className ? style : { ...defaultButtonStyle, ...style }}
    >
      <ShieldIcon />
      {children || "Share via AddrPass"}
    </button>
  );
}

// Re-export types
export type { ShareScope, TokenExchangeResponse } from "./types";
