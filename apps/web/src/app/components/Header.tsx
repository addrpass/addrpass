"use client";

import Link from "next/link";
import { useState } from "react";

function Logo() {
  return (
    <svg width="28" height="28" viewBox="0 0 32 32" fill="none">
      {/* Shield shape with keyhole */}
      <path d="M16 2L4 8v8c0 8.4 5.12 16.24 12 18 6.88-1.76 12-9.6 12-18V8L16 2z" fill="#0F172A" />
      <path d="M16 4L6 9v7c0 7.36 4.48 14.24 10 15.8V4z" fill="#1E293B" />
      <path d="M16 4l10 5v7c0 7.36-4.48 14.24-10 15.8V4z" fill="#0F172A" />
      {/* Keyhole / address pin */}
      <circle cx="16" cy="13" r="3" fill="#22D3EE" />
      <path d="M14 15.5L13 22h6l-1-6.5" fill="#22D3EE" opacity="0.7" />
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
          <a href="https://github.com/addrpass/addrpass" target="_blank" rel="noopener noreferrer" className="hover:text-[#0F172A] transition-colors">GitHub</a>
        </nav>

        <div className="flex items-center gap-3">
          <Link href="/login" className="hidden sm:block text-[13px] font-medium text-[#64748B] hover:text-[#0F172A] transition-colors">
            Sign in
          </Link>
          <Link href="/register" className="btn-primary rounded-full px-5 py-2 text-[13px] font-semibold">
            <span>Get Started</span>
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
          <a href="https://github.com/addrpass/addrpass" className="block text-[#64748B]">GitHub</a>
        </nav>
      )}
    </header>
  );
}
