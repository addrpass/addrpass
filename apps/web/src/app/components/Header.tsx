"use client";

import Link from "next/link";
import { useState } from "react";

function Logo() {
  return (
    <svg width="28" height="28" viewBox="0 0 1024 1024" fill="none">
      <rect x="0" y="0" width="1024" height="1024" rx="224" ry="224" fill="#F0EFE8" />
      <rect x="232" y="232" width="560" height="560" rx="80" fill="none" stroke="#1A1A1A" strokeWidth="48" />
      <circle cx="512" cy="430" r="76" fill="#1A1A1A" />
      <path d="M472 480L446 668L578 668L552 480Z" fill="#1A1A1A" />
      <circle cx="512" cy="430" r="28" fill="#E7B300" />
    </svg>
  );
}

export default function Header() {
  const [mobileOpen, setMobileOpen] = useState(false);

  return (
    <header className="sticky top-0 z-50 backdrop-blur-lg bg-[#FAFBFD]/90 border-b border-[#E2E8F0]">
      <div className="mx-auto max-w-6xl px-6 py-4 flex items-center justify-between">
        <Link href="/" className="flex items-center gap-2.5">
          <Logo />
          <span className="text-[17px] font-bold tracking-tight text-[#0F172A]">AddrPass</span>
        </Link>

        <nav className="hidden md:flex items-center gap-7 text-[13px] font-medium text-[#64748B]">
          <a href="#how-it-works" className="hover:text-[#0F172A] transition-colors">How It Works</a>
          <a href="#features" className="hover:text-[#0F172A] transition-colors">Features</a>
          <a href="#delivery" className="hover:text-[#0F172A] transition-colors">For Delivery</a>
          <a href="#developers" className="hover:text-[#0F172A] transition-colors">Developers</a>
          <Link href="/pricing" className="hover:text-[#0F172A] transition-colors">Pricing</Link>
        </nav>

        <div className="flex items-center gap-3">
          <Link href="/login" className="hidden sm:block text-[13px] font-medium text-[#64748B] hover:text-[#0F172A] transition-colors">
            Sign in
          </Link>
          <Link href="/register" className="btn-primary rounded-full px-5 py-2 text-[13px] font-semibold">
            <span>Get started free</span>
          </Link>
          <button onClick={() => setMobileOpen(!mobileOpen)} className="md:hidden p-1 text-[#64748B]">
            <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              {mobileOpen
                ? <path d="M18 6L6 18M6 6l12 12" />
                : <><line x1="4" y1="7" x2="20" y2="7" /><line x1="4" y1="12" x2="20" y2="12" /><line x1="4" y1="17" x2="20" y2="17" /></>}
            </svg>
          </button>
        </div>
      </div>

      {mobileOpen && (
        <nav className="md:hidden border-t border-[#E2E8F0] bg-white px-6 py-4 space-y-3 text-sm">
          <a href="#how-it-works" className="block text-[#64748B]" onClick={() => setMobileOpen(false)}>How It Works</a>
          <a href="#features" className="block text-[#64748B]" onClick={() => setMobileOpen(false)}>Features</a>
          <a href="#delivery" className="block text-[#64748B]" onClick={() => setMobileOpen(false)}>For Delivery</a>
          <Link href="/pricing" className="block text-[#64748B]" onClick={() => setMobileOpen(false)}>Pricing</Link>
          <a href="https://github.com/addrpass/addrpass" className="block text-[#64748B]">GitHub</a>
        </nav>
      )}
    </header>
  );
}
