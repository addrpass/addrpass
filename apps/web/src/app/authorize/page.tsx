"use client";

import { useEffect, useState, Suspense } from "react";
import { useSearchParams, useRouter } from "next/navigation";
import Link from "next/link";
import { useAuth } from "@/lib/auth";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "https://api.addrpass.com";

type OAuthApp = {
  name: string;
  logo_url: string;
};

type Address = {
  id: string;
  label: string;
  line1: string;
  city: string;
  country: string;
};

type ConsentData = {
  app: OAuthApp;
  addresses: Address[];
  scope: string;
  state: string;
};

function ConsentContent() {
  const searchParams = useSearchParams();
  const router = useRouter();
  const { token, user, loading: authLoading } = useAuth();

  const clientId = searchParams.get("client_id") || "";
  const redirectUri = searchParams.get("redirect_uri") || "";
  const scope = searchParams.get("scope") || "delivery";
  const state = searchParams.get("state") || "";

  const [data, setData] = useState<ConsentData | null>(null);
  const [selectedAddress, setSelectedAddress] = useState("");
  const [selectedScope, setSelectedScope] = useState(scope);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");

  useEffect(() => {
    if (authLoading) return;
    if (!user || !token) {
      // Redirect to login with return URL
      const returnUrl = encodeURIComponent(window.location.href);
      router.push(`/login?return=${returnUrl}`);
      return;
    }

    // Fetch consent data
    fetch(`${API_URL}/api/v1/oauth/authorize?client_id=${clientId}&redirect_uri=${encodeURIComponent(redirectUri)}&scope=${scope}&state=${state}`, {
      headers: { Authorization: `Bearer ${token}` },
    })
      .then((r) => {
        if (!r.ok) throw new Error("Invalid application");
        return r.json();
      })
      .then((d) => {
        setData(d);
        if (d.addresses.length > 0) setSelectedAddress(d.addresses[0].id);
        setLoading(false);
      })
      .catch((e) => {
        setError(e.message);
        setLoading(false);
      });
  }, [authLoading, user, token, clientId, redirectUri, scope, state, router]);

  const handleApprove = async () => {
    if (!token || !selectedAddress) return;
    setSubmitting(true);
    try {
      const res = await fetch(`${API_URL}/api/v1/oauth/consent`, {
        method: "POST",
        headers: { "Content-Type": "application/json", Authorization: `Bearer ${token}` },
        body: JSON.stringify({
          client_id: clientId,
          redirect_uri: redirectUri,
          scope: selectedScope,
          state: state,
          address_id: selectedAddress,
        }),
      });
      if (!res.ok) {
        const err = await res.json();
        throw new Error(err.error || "Failed to authorize");
      }
      const result = await res.json();

      // If opened as popup, post message to parent and close
      if (window.opener) {
        window.opener.postMessage({ code: result.code, state: state }, "*");
        window.close();
      } else {
        // Redirect directly
        window.location.href = result.redirect_url;
      }
    } catch (e: unknown) {
      setError(e instanceof Error ? e.message : "Failed to authorize");
      setSubmitting(false);
    }
  };

  const handleDeny = () => {
    if (window.opener) {
      window.opener.postMessage({ error: "access_denied" }, "*");
      window.close();
    } else {
      const url = redirectUri + "?error=access_denied" + (state ? `&state=${state}` : "");
      window.location.href = url;
    }
  };

  if (authLoading || loading) {
    return (
      <div className="text-center py-20">
        <div className="w-8 h-8 border-[3px] border-[#22D3EE] border-t-transparent rounded-full animate-spin mx-auto mb-4" />
        <p className="text-[#64748B] text-sm">Loading...</p>
      </div>
    );
  }

  if (error) {
    return (
      <div className="max-w-sm mx-auto bg-white border border-red-200 rounded-2xl p-8 text-center">
        <div className="w-12 h-12 mx-auto mb-4 rounded-xl bg-red-50 flex items-center justify-center">
          <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#DC2626" strokeWidth="2"><circle cx="12" cy="12" r="10" /><line x1="15" y1="9" x2="9" y2="15" /><line x1="9" y1="9" x2="15" y2="15" /></svg>
        </div>
        <h2 className="font-bold text-lg mb-1">Authorization Failed</h2>
        <p className="text-sm text-[#64748B]">{error}</p>
      </div>
    );
  }

  if (!data) return null;

  const scopeLabels: Record<string, string> = {
    full: "Full address + phone",
    delivery: "Address only (no phone)",
    zone: "City + postal code only",
    verify: "Confirm address exists",
  };

  return (
    <div className="max-w-sm mx-auto">
      <div className="bg-white border border-[#E2E8F0] rounded-2xl p-8">
        {/* App info */}
        <div className="text-center mb-6">
          {data.app.logo_url ? (
            <img src={data.app.logo_url} alt={data.app.name} className="w-12 h-12 rounded-xl mx-auto mb-3" />
          ) : (
            <div className="w-12 h-12 rounded-xl bg-[#F1F5F9] flex items-center justify-center mx-auto mb-3 text-lg font-bold text-[#64748B]">
              {data.app.name[0]}
            </div>
          )}
          <h2 className="font-bold text-lg">{data.app.name}</h2>
          <p className="text-sm text-[#64748B] mt-1">wants access to your address</p>
        </div>

        {/* Address selector */}
        <div className="mb-5">
          <label className="block text-xs font-medium text-[#64748B] mb-2">Select an address</label>
          <div className="space-y-2">
            {data.addresses.map((addr) => (
              <label
                key={addr.id}
                className={`flex items-start gap-3 p-3 rounded-xl border cursor-pointer transition-all ${
                  selectedAddress === addr.id
                    ? "border-[#22D3EE] bg-[#22D3EE]/[0.04]"
                    : "border-[#E2E8F0] hover:border-[#CBD5E1]"
                }`}
              >
                <input
                  type="radio"
                  name="address"
                  value={addr.id}
                  checked={selectedAddress === addr.id}
                  onChange={() => setSelectedAddress(addr.id)}
                  className="mt-1 accent-[#22D3EE]"
                />
                <div>
                  {addr.label && <span className="text-xs font-medium text-[#22D3EE]">{addr.label}</span>}
                  <p className="text-sm text-[#0F172A]">{addr.line1}</p>
                  <p className="text-xs text-[#64748B]">{addr.city}, {addr.country}</p>
                </div>
              </label>
            ))}
          </div>
          {data.addresses.length === 0 && (
            <p className="text-sm text-[#64748B] text-center py-4">No addresses found. <Link href="/dashboard" className="text-[#22D3EE] hover:underline">Add one first</Link>.</p>
          )}
        </div>

        {/* Scope selector */}
        <div className="mb-6">
          <label className="block text-xs font-medium text-[#64748B] mb-2">Access level</label>
          <select
            value={selectedScope}
            onChange={(e) => setSelectedScope(e.target.value)}
            className="w-full rounded-lg border border-[#E2E8F0] px-3 py-2.5 text-sm outline-none focus:border-[#22D3EE] focus:ring-2 focus:ring-[#22D3EE]/20"
          >
            {Object.entries(scopeLabels).map(([key, label]) => (
              <option key={key} value={key}>{label}</option>
            ))}
          </select>
        </div>

        {/* Actions */}
        <div className="flex gap-3">
          <button
            onClick={handleDeny}
            className="flex-1 rounded-xl border border-[#E2E8F0] py-2.5 text-sm font-medium text-[#64748B] hover:bg-[#F8FAFC] transition-colors"
          >
            Deny
          </button>
          <button
            onClick={handleApprove}
            disabled={submitting || !selectedAddress}
            className="flex-1 rounded-xl bg-[#0F172A] py-2.5 text-sm font-semibold text-white hover:bg-[#1E293B] transition-colors disabled:opacity-50"
          >
            {submitting ? "Sharing..." : "Share Address"}
          </button>
        </div>

        <p className="text-[10px] text-[#94A3B8] text-center mt-4 leading-relaxed">
          You can revoke access anytime from your AddrPass dashboard.
        </p>
      </div>
    </div>
  );
}

export default function AuthorizePage() {
  return (
    <div className="min-h-screen bg-[#FAFBFD] flex flex-col">
      <header className="border-b border-[#E2E8F0] bg-white">
        <div className="mx-auto max-w-lg px-6 py-4 flex items-center justify-center">
          <Link href="/" className="flex items-center gap-2">
            <svg width="24" height="24" viewBox="0 0 32 32" fill="none">
              <path d="M16 2L4 8v8c0 8.4 5.12 16.24 12 18 6.88-1.76 12-9.6 12-18V8L16 2z" fill="#0F172A" />
              <circle cx="16" cy="13" r="3" fill="#22D3EE" />
              <path d="M14 15.5L13 22h6l-1-6.5" fill="#22D3EE" opacity="0.7" />
            </svg>
            <span className="text-base font-bold">AddrPass</span>
          </Link>
        </div>
      </header>
      <main className="flex-1 flex items-center justify-center px-4 py-12">
        <Suspense fallback={<div className="text-[#64748B]">Loading...</div>}>
          <ConsentContent />
        </Suspense>
      </main>
    </div>
  );
}
