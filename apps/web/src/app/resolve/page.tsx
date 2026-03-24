"use client";

import { useEffect, useState } from "react";
import { useSearchParams } from "next/navigation";
import Link from "next/link";
import { Suspense } from "react";

const API_URL = process.env.NEXT_PUBLIC_API_URL || "https://api.addrpass.com";

type Address = {
  label: string;
  line1: string;
  line2: string;
  city: string;
  state: string;
  post_code: string;
  country: string;
  phone: string;
};

type ResolveState =
  | { status: "loading" }
  | { status: "no_token" }
  | { status: "pin_required" }
  | { status: "resolved"; address: Address }
  | { status: "error"; message: string };

function ResolveContent() {
  const searchParams = useSearchParams();
  const token = searchParams.get("t");
  const [state, setState] = useState<ResolveState>(token ? { status: "loading" } : { status: "no_token" });
  const [pin, setPin] = useState("");
  const [copied, setCopied] = useState(false);

  const resolve = async (pinValue?: string) => {
    if (!token) return;
    setState({ status: "loading" });
    try {
      const url = pinValue
        ? `${API_URL}/api/v1/resolve/${token}?pin=${pinValue}`
        : `${API_URL}/api/v1/resolve/${token}`;
      const res = await fetch(url);

      if (res.status === 403) {
        setState({ status: "pin_required" });
        return;
      }
      if (!res.ok) {
        const err = await res.json().catch(() => ({ error: res.statusText }));
        setState({ status: "error", message: err.error || "Failed to resolve" });
        return;
      }

      const data = await res.json();
      setState({ status: "resolved", address: data.address });
    } catch {
      setState({ status: "error", message: "Failed to connect to server" });
    }
  };

  useEffect(() => {
    if (token) resolve();
  }, [token]); // eslint-disable-line react-hooks/exhaustive-deps

  const handlePinSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    resolve(pin);
  };

  const handleCopy = () => {
    if (state.status !== "resolved") return;
    const a = state.address;
    const text = [a.line1, a.line2, `${a.city}${a.state ? `, ${a.state}` : ""} ${a.post_code}`, a.country, a.phone].filter(Boolean).join("\n");
    navigator.clipboard.writeText(text);
    setCopied(true);
    setTimeout(() => setCopied(false), 2000);
  };

  return (
    <div className="w-full max-w-md">
      {state.status === "loading" && (
        <div className="text-center">
          <div className="w-10 h-10 border-[3px] border-teal-500 border-t-transparent rounded-full animate-spin mx-auto mb-4" />
          <p className="text-[#6B7280]">Resolving address...</p>
        </div>
      )}

      {state.status === "no_token" && (
        <div className="bg-white border border-[#E5E7EB] rounded-2xl p-8 text-center">
          <h2 className="text-lg font-bold mb-2">No share token provided</h2>
          <p className="text-sm text-[#6B7280]">This page resolves shared address links.</p>
        </div>
      )}

      {state.status === "pin_required" && (
        <div className="bg-white border border-[#E5E7EB] rounded-2xl p-8 text-center">
          <div className="w-14 h-14 mx-auto mb-4 rounded-xl bg-violet-50 flex items-center justify-center">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#7C3AED" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
              <path d="M7 11V7a5 5 0 0 1 10 0v4" />
            </svg>
          </div>
          <h2 className="text-lg font-bold mb-1">PIN Required</h2>
          <p className="text-sm text-[#6B7280] mb-6">This address is protected. Enter the PIN to view it.</p>
          <form onSubmit={handlePinSubmit} className="flex gap-2">
            <input type="text" required value={pin} onChange={(e) => setPin(e.target.value)} placeholder="Enter PIN" className="flex-1 rounded-lg border border-[#E5E7EB] px-4 py-2.5 text-sm text-center tracking-widest outline-none focus:border-[#0D9488] focus:ring-2 focus:ring-[#0D9488]/20" autoFocus />
            <button type="submit" className="btn-primary rounded-lg px-5 py-2.5 text-sm font-semibold text-white">Unlock</button>
          </form>
        </div>
      )}

      {state.status === "resolved" && (
        <div className="bg-white border border-[#E5E7EB] rounded-2xl p-8">
          <div className="flex items-center gap-3 mb-6">
            <div className="w-10 h-10 rounded-xl bg-teal-50 flex items-center justify-center">
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="#0D9488" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
                <path d="M21 10c0 7-9 13-9 13s-9-6-9-13a9 9 0 0 1 18 0z" />
                <circle cx="12" cy="10" r="3" />
              </svg>
            </div>
            <div>
              {state.address.label && <span className="text-xs font-medium text-teal-700 bg-teal-50 rounded-full px-2 py-0.5">{state.address.label}</span>}
              <p className="text-xs text-[#9CA3AF] mt-0.5">Shared via AddrPass</p>
            </div>
          </div>
          <div className="space-y-1.5 text-sm">
            <p className="font-medium text-[#111827]">{state.address.line1}</p>
            {state.address.line2 && <p className="text-[#6B7280]">{state.address.line2}</p>}
            <p className="text-[#6B7280]">{state.address.city}{state.address.state ? `, ${state.address.state}` : ""} {state.address.post_code}</p>
            <p className="text-[#6B7280]">{state.address.country}</p>
            {state.address.phone && <p className="text-[#6B7280] pt-2 border-t border-[#E5E7EB] mt-3">{state.address.phone}</p>}
          </div>
          <div className="mt-6 pt-4 border-t border-[#E5E7EB] flex gap-2">
            <button onClick={handleCopy} className="flex-1 rounded-lg border border-[#E5E7EB] px-4 py-2 text-sm font-medium text-[#374151] hover:bg-gray-50 transition-colors">
              {copied ? "Copied!" : "Copy Address"}
            </button>
            <a href={`https://maps.google.com/?q=${encodeURIComponent(`${state.address.line1}, ${state.address.city}, ${state.address.country}`)}`} target="_blank" rel="noopener noreferrer" className="flex-1 rounded-lg border border-[#E5E7EB] px-4 py-2 text-sm font-medium text-[#374151] hover:bg-gray-50 transition-colors text-center">
              Open in Maps
            </a>
          </div>
        </div>
      )}

      {state.status === "error" && (
        <div className="bg-white border border-red-200 rounded-2xl p-8 text-center">
          <div className="w-14 h-14 mx-auto mb-4 rounded-xl bg-red-50 flex items-center justify-center">
            <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="#DC2626" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
              <circle cx="12" cy="12" r="10" />
              <line x1="15" y1="9" x2="9" y2="15" />
              <line x1="9" y1="9" x2="15" y2="15" />
            </svg>
          </div>
          <h2 className="text-lg font-bold mb-1">Unable to Access</h2>
          <p className="text-sm text-[#6B7280]">{state.message}</p>
        </div>
      )}
    </div>
  );
}

export default function ResolvePage() {
  return (
    <div className="min-h-screen bg-[#FAFBFC] flex flex-col">
      <header className="border-b border-[#E5E7EB] bg-white">
        <div className="mx-auto max-w-lg px-6 py-4 flex items-center justify-between">
          <Link href="/" className="flex items-center gap-2">
            <div className="w-7 h-7 rounded-md bg-gradient-to-br from-teal-500 to-teal-700 flex items-center justify-center">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="white" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
                <path d="M7 11V7a5 5 0 0 1 10 0v4" />
              </svg>
            </div>
            <span className="text-base font-bold">AddrPass</span>
          </Link>
          <span className="text-xs text-[#9CA3AF]">Shared Address</span>
        </div>
      </header>
      <main className="flex-1 flex items-center justify-center px-4 py-12">
        <Suspense fallback={<div className="text-[#6B7280]">Loading...</div>}>
          <ResolveContent />
        </Suspense>
      </main>
      <footer className="border-t border-[#E5E7EB] py-6 text-center">
        <p className="text-xs text-[#9CA3AF]">
          Powered by <Link href="/" className="text-teal-600 hover:underline">AddrPass</Link> — Your address, your control.
        </p>
      </footer>
    </div>
  );
}
