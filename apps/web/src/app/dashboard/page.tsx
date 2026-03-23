"use client";

import { useEffect, useState, useCallback } from "react";
import { useRouter } from "next/navigation";
import Link from "next/link";
import { useAuth } from "@/lib/auth";
import * as api from "@/lib/api";

function AddressForm({ onCreated }: { onCreated: () => void }) {
  const { token } = useAuth();
  const [open, setOpen] = useState(false);
  const [loading, setLoading] = useState(false);
  const [form, setForm] = useState({ label: "", line1: "", line2: "", city: "", state: "", post_code: "", country: "", phone: "" });

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    if (!token) return;
    setLoading(true);
    try {
      await api.createAddress(token, form);
      setForm({ label: "", line1: "", line2: "", city: "", state: "", post_code: "", country: "", phone: "" });
      setOpen(false);
      onCreated();
    } finally {
      setLoading(false);
    }
  };

  if (!open) {
    return (
      <button onClick={() => setOpen(true)} className="btn-primary rounded-lg px-4 py-2 text-sm font-semibold text-white">
        + New Address
      </button>
    );
  }

  return (
    <form onSubmit={handleSubmit} className="bg-white border border-[#E5E7EB] rounded-xl p-5 space-y-3">
      <div className="grid grid-cols-2 gap-3">
        <input placeholder="Label (e.g. Home)" value={form.label} onChange={(e) => setForm({ ...form, label: e.target.value })} className="col-span-2 rounded-lg border border-[#E5E7EB] px-3 py-2 text-sm outline-none focus:border-[#0D9488]" />
        <input required placeholder="Address line 1 *" value={form.line1} onChange={(e) => setForm({ ...form, line1: e.target.value })} className="col-span-2 rounded-lg border border-[#E5E7EB] px-3 py-2 text-sm outline-none focus:border-[#0D9488]" />
        <input placeholder="Address line 2" value={form.line2} onChange={(e) => setForm({ ...form, line2: e.target.value })} className="col-span-2 rounded-lg border border-[#E5E7EB] px-3 py-2 text-sm outline-none focus:border-[#0D9488]" />
        <input required placeholder="City *" value={form.city} onChange={(e) => setForm({ ...form, city: e.target.value })} className="rounded-lg border border-[#E5E7EB] px-3 py-2 text-sm outline-none focus:border-[#0D9488]" />
        <input placeholder="State / Region" value={form.state} onChange={(e) => setForm({ ...form, state: e.target.value })} className="rounded-lg border border-[#E5E7EB] px-3 py-2 text-sm outline-none focus:border-[#0D9488]" />
        <input required placeholder="Post code *" value={form.post_code} onChange={(e) => setForm({ ...form, post_code: e.target.value })} className="rounded-lg border border-[#E5E7EB] px-3 py-2 text-sm outline-none focus:border-[#0D9488]" />
        <input required placeholder="Country (e.g. TR) *" value={form.country} onChange={(e) => setForm({ ...form, country: e.target.value })} className="rounded-lg border border-[#E5E7EB] px-3 py-2 text-sm outline-none focus:border-[#0D9488]" />
        <input placeholder="Phone" value={form.phone} onChange={(e) => setForm({ ...form, phone: e.target.value })} className="col-span-2 rounded-lg border border-[#E5E7EB] px-3 py-2 text-sm outline-none focus:border-[#0D9488]" />
      </div>
      <div className="flex gap-2">
        <button type="submit" disabled={loading} className="btn-primary rounded-lg px-4 py-2 text-sm font-semibold text-white disabled:opacity-50">
          {loading ? "Saving..." : "Save Address"}
        </button>
        <button type="button" onClick={() => setOpen(false)} className="rounded-lg px-4 py-2 text-sm text-[#6B7280] border border-[#E5E7EB] hover:bg-gray-50">
          Cancel
        </button>
      </div>
    </form>
  );
}

function ShareModal({ address, onClose, onCreated }: { address: api.Address; onClose: () => void; onCreated: () => void }) {
  const { token } = useAuth();
  const [loading, setLoading] = useState(false);
  const [result, setResult] = useState<api.CreateShareResponse | null>(null);
  const [pin, setPin] = useState("");
  const [maxAccesses, setMaxAccesses] = useState("");

  const handleCreate = async () => {
    if (!token) return;
    setLoading(true);
    try {
      const data: api.CreateShareInput = { address_id: address.id, access_type: "public" };
      if (pin) data.pin = pin;
      if (maxAccesses) data.max_accesses = parseInt(maxAccesses);
      const resp = await api.createShare(token, data);
      setResult(resp);
      onCreated();
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/30 backdrop-blur-sm" onClick={onClose}>
      <div className="bg-white rounded-2xl p-6 w-full max-w-md mx-4 shadow-xl" onClick={(e) => e.stopPropagation()}>
        <h3 className="text-lg font-bold mb-1">Share: {address.label || address.line1}</h3>
        <p className="text-sm text-[#6B7280] mb-4">{address.line1}, {address.city}, {address.country}</p>

        {!result ? (
          <div className="space-y-3">
            <div>
              <label className="block text-xs font-medium text-[#6B7280] mb-1">PIN (optional)</label>
              <input placeholder="e.g. 1234" value={pin} onChange={(e) => setPin(e.target.value)} className="w-full rounded-lg border border-[#E5E7EB] px-3 py-2 text-sm outline-none focus:border-[#0D9488]" />
            </div>
            <div>
              <label className="block text-xs font-medium text-[#6B7280] mb-1">Max accesses (optional)</label>
              <input type="number" min="1" placeholder="Unlimited" value={maxAccesses} onChange={(e) => setMaxAccesses(e.target.value)} className="w-full rounded-lg border border-[#E5E7EB] px-3 py-2 text-sm outline-none focus:border-[#0D9488]" />
            </div>
            <div className="flex gap-2">
              <button onClick={handleCreate} disabled={loading} className="btn-primary rounded-lg px-4 py-2 text-sm font-semibold text-white flex-1 disabled:opacity-50">
                {loading ? "Creating..." : "Create Share"}
              </button>
              <button onClick={onClose} className="rounded-lg px-4 py-2 text-sm text-[#6B7280] border border-[#E5E7EB]">Cancel</button>
            </div>
          </div>
        ) : (
          <div className="space-y-4">
            <div>
              <label className="block text-xs font-medium text-[#6B7280] mb-1">Share URL</label>
              <div className="flex gap-2">
                <input readOnly value={result.url} className="flex-1 rounded-lg border border-[#E5E7EB] bg-gray-50 px-3 py-2 text-sm" />
                <button onClick={() => navigator.clipboard.writeText(result.url)} className="rounded-lg px-3 py-2 text-sm border border-[#E5E7EB] hover:bg-gray-50">Copy</button>
              </div>
            </div>
            <div className="text-center">
              <p className="text-xs text-[#6B7280] mb-2">QR Code</p>
              <img src={api.getQRCodeURL(result.share.token)} alt="QR Code" className="w-48 h-48 mx-auto rounded-lg border border-[#E5E7EB]" />
            </div>
            <button onClick={onClose} className="w-full rounded-lg px-4 py-2 text-sm font-semibold border border-[#E5E7EB] hover:bg-gray-50">Done</button>
          </div>
        )}
      </div>
    </div>
  );
}

function AccessLogsModal({ share, onClose }: { share: api.Share; onClose: () => void }) {
  const { token } = useAuth();
  const [logs, setLogs] = useState<api.AccessLog[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!token) return;
    api.getAccessLogs(token, share.id).then(setLogs).finally(() => setLoading(false));
  }, [token, share.id]);

  return (
    <div className="fixed inset-0 z-50 flex items-center justify-center bg-black/30 backdrop-blur-sm" onClick={onClose}>
      <div className="bg-white rounded-2xl p-6 w-full max-w-lg mx-4 shadow-xl max-h-[80vh] overflow-auto" onClick={(e) => e.stopPropagation()}>
        <h3 className="text-lg font-bold mb-1">Access Logs</h3>
        <p className="text-sm text-[#6B7280] mb-4">Token: {share.token.slice(0, 12)}... ({share.access_count} accesses)</p>

        {loading ? (
          <p className="text-sm text-[#6B7280]">Loading...</p>
        ) : logs.length === 0 ? (
          <p className="text-sm text-[#6B7280]">No accesses yet.</p>
        ) : (
          <div className="space-y-2">
            {logs.map((log) => (
              <div key={log.id} className="border border-[#E5E7EB] rounded-lg p-3 text-sm">
                <div className="flex justify-between">
                  <span className="text-[#374151] font-mono">{log.ip}</span>
                  <span className="text-[#6B7280]">{new Date(log.access_at).toLocaleString()}</span>
                </div>
                <p className="text-xs text-[#9CA3AF] mt-1 truncate">{log.user_agent}</p>
              </div>
            ))}
          </div>
        )}
        <button onClick={onClose} className="w-full mt-4 rounded-lg px-4 py-2 text-sm font-semibold border border-[#E5E7EB] hover:bg-gray-50">Close</button>
      </div>
    </div>
  );
}

export default function DashboardPage() {
  const { user, token, loading: authLoading, logout } = useAuth();
  const router = useRouter();
  const [addresses, setAddresses] = useState<api.Address[]>([]);
  const [shares, setShares] = useState<api.Share[]>([]);
  const [sharingAddress, setSharingAddress] = useState<api.Address | null>(null);
  const [viewingLogs, setViewingLogs] = useState<api.Share | null>(null);

  const loadData = useCallback(async () => {
    if (!token) return;
    const [addrs, shrs] = await Promise.all([
      api.listAddresses(token),
      api.listShares(token),
    ]);
    setAddresses(addrs);
    setShares(shrs);
  }, [token]);

  useEffect(() => {
    if (!authLoading && !user) {
      router.push("/login");
    }
  }, [authLoading, user, router]);

  useEffect(() => {
    loadData();
  }, [loadData]);

  if (authLoading || !user) return null;

  const handleDelete = async (id: string) => {
    if (!token || !confirm("Delete this address? All shares will be removed too.")) return;
    await api.deleteAddress(token, id);
    loadData();
  };

  const handleRevoke = async (id: string) => {
    if (!token) return;
    await api.revokeShare(token, id);
    loadData();
  };

  return (
    <div className="min-h-screen bg-[#FAFBFC]">
      {/* Header */}
      <header className="sticky top-0 z-40 backdrop-blur-md bg-[#FAFBFC]/80 border-b border-[#E5E7EB]">
        <div className="mx-auto max-w-5xl px-6 py-3 flex items-center justify-between">
          <Link href="/" className="flex items-center gap-2">
            <div className="w-7 h-7 rounded-md bg-gradient-to-br from-teal-500 to-teal-700 flex items-center justify-center">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="white" strokeWidth="2.5" strokeLinecap="round" strokeLinejoin="round">
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
                <path d="M7 11V7a5 5 0 0 1 10 0v4" />
              </svg>
            </div>
            <span className="text-base font-bold">AddrPass</span>
          </Link>
          <div className="flex items-center gap-4">
            <span className="text-sm text-[#6B7280]">{user.email}</span>
            <button onClick={logout} className="text-sm text-[#6B7280] hover:text-red-600">Sign out</button>
          </div>
        </div>
      </header>

      <main className="mx-auto max-w-5xl px-6 py-8">
        {/* Addresses */}
        <div className="flex items-center justify-between mb-6">
          <h2 className="text-xl font-bold">Addresses</h2>
          <AddressForm onCreated={loadData} />
        </div>

        {addresses.length === 0 ? (
          <div className="text-center py-12 border border-dashed border-[#E5E7EB] rounded-xl">
            <p className="text-[#6B7280]">No addresses yet. Add your first one.</p>
          </div>
        ) : (
          <div className="grid gap-4 sm:grid-cols-2">
            {addresses.map((addr) => (
              <div key={addr.id} className="bg-white border border-[#E5E7EB] rounded-xl p-5 hover:border-[#0D9488] transition-colors">
                <div className="flex items-start justify-between">
                  <div>
                    {addr.label && <span className="inline-block text-xs font-medium text-teal-700 bg-teal-50 rounded-full px-2.5 py-0.5 mb-2">{addr.label}</span>}
                    <p className="font-medium text-sm">{addr.line1}</p>
                    {addr.line2 && <p className="text-sm text-[#6B7280]">{addr.line2}</p>}
                    <p className="text-sm text-[#6B7280]">{addr.city}{addr.state ? `, ${addr.state}` : ""} {addr.post_code}</p>
                    <p className="text-sm text-[#6B7280]">{addr.country}</p>
                    {addr.phone && <p className="text-xs text-[#9CA3AF] mt-1">{addr.phone}</p>}
                  </div>
                </div>
                <div className="flex gap-2 mt-4 border-t border-[#E5E7EB] pt-3">
                  <button onClick={() => setSharingAddress(addr)} className="text-xs font-medium text-teal-600 hover:underline">Share</button>
                  <button onClick={() => handleDelete(addr.id)} className="text-xs font-medium text-red-500 hover:underline">Delete</button>
                </div>
              </div>
            ))}
          </div>
        )}

        {/* Shares */}
        {shares.length > 0 && (
          <>
            <h2 className="text-xl font-bold mt-12 mb-6">Active Shares</h2>
            <div className="space-y-3">
              {shares.map((share) => {
                const addr = addresses.find((a) => a.id === share.address_id);
                return (
                  <div key={share.id} className={`bg-white border rounded-xl p-4 flex items-center justify-between ${share.active ? "border-[#E5E7EB]" : "border-red-200 bg-red-50/30"}`}>
                    <div>
                      <div className="flex items-center gap-2">
                        <span className="font-mono text-sm">{share.token.slice(0, 16)}...</span>
                        {!share.active && <span className="text-xs text-red-600 font-medium">Revoked</span>}
                        {share.pin && <span className="text-xs text-violet-600 font-medium">PIN</span>}
                      </div>
                      <p className="text-xs text-[#6B7280] mt-1">
                        {addr?.label || addr?.line1 || "Unknown"} — {share.access_count} accesses
                        {share.max_accesses ? ` / ${share.max_accesses} max` : ""}
                      </p>
                    </div>
                    <div className="flex gap-2">
                      <button onClick={() => setViewingLogs(share)} className="text-xs font-medium text-teal-600 hover:underline">Logs</button>
                      {share.active && (
                        <button onClick={() => handleRevoke(share.id)} className="text-xs font-medium text-red-500 hover:underline">Revoke</button>
                      )}
                    </div>
                  </div>
                );
              })}
            </div>
          </>
        )}
      </main>

      {/* Modals */}
      {sharingAddress && (
        <ShareModal address={sharingAddress} onClose={() => setSharingAddress(null)} onCreated={loadData} />
      )}
      {viewingLogs && (
        <AccessLogsModal share={viewingLogs} onClose={() => setViewingLogs(null)} />
      )}
    </div>
  );
}
