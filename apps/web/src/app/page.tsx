"use client";

import Header from "./components/Header";
import Footer from "./components/Footer";

/* ─── Custom SVG Illustrations ────────────────────────────── */

function HeroIllustration() {
  return (
    <svg viewBox="0 0 400 320" fill="none" className="w-full max-w-md">
      {/* Address card */}
      <rect x="60" y="40" width="280" height="160" rx="16" fill="#F8FAFC" stroke="#E2E8F0" strokeWidth="1.5" />
      <rect x="80" y="65" width="120" height="8" rx="4" fill="#CBD5E1" />
      <rect x="80" y="83" width="180" height="6" rx="3" fill="#E2E8F0" />
      <rect x="80" y="97" width="140" height="6" rx="3" fill="#E2E8F0" />
      <rect x="80" y="111" width="100" height="6" rx="3" fill="#E2E8F0" />

      {/* Lock overlay */}
      <rect x="80" y="140" width="60" height="30" rx="8" fill="#0F172A" />
      <circle cx="110" cy="152" r="4" fill="#22D3EE" />
      <rect x="108" y="154" width="4" height="6" rx="1" fill="#22D3EE" opacity="0.7" />

      {/* Token flying out */}
      <g className="float">
        <rect x="230" y="135" width="90" height="36" rx="10" fill="white" stroke="#22D3EE" strokeWidth="1.5" />
        <text x="248" y="158" fontSize="10" fontFamily="monospace" fontWeight="600" fill="#0F172A">aX7k...9mP</text>
      </g>

      {/* Arrow from card to token */}
      <path d="M200 155 L225 155" stroke="#22D3EE" strokeWidth="1.5" strokeDasharray="4 3" />
      <polygon points="223,151 230,155 223,159" fill="#22D3EE" />

      {/* QR code floating */}
      <g className="float float-delay-1">
        <rect x="260" y="60" width="56" height="56" rx="8" fill="white" stroke="#A78BFA" strokeWidth="1.5" />
        <rect x="270" y="70" width="8" height="8" rx="1" fill="#0F172A" />
        <rect x="282" y="70" width="8" height="8" rx="1" fill="#0F172A" />
        <rect x="294" y="70" width="8" height="8" rx="1" fill="#0F172A" />
        <rect x="270" y="82" width="8" height="8" rx="1" fill="#0F172A" />
        <rect x="294" y="82" width="8" height="8" rx="1" fill="#0F172A" />
        <rect x="270" y="94" width="8" height="8" rx="1" fill="#0F172A" />
        <rect x="282" y="94" width="8" height="8" rx="1" fill="#0F172A" />
        <rect x="294" y="94" width="8" height="8" rx="1" fill="#0F172A" />
        <rect x="282" y="82" width="8" height="8" rx="1" fill="#A78BFA" opacity="0.3" />
      </g>

      {/* Eye / monitoring indicator */}
      <g className="float float-delay-2">
        <circle cx="90" cy="240" r="22" fill="white" stroke="#FB7185" strokeWidth="1.5" />
        <path d="M78 240s5.4-7 12-7 12 7 12 7-5.4 7-12 7-12-7-12-7z" fill="none" stroke="#FB7185" strokeWidth="1.5" />
        <circle cx="90" cy="240" r="3.5" fill="#FB7185" />
      </g>

      {/* Connection lines */}
      <path d="M110 170 Q110 210 90 220" stroke="#E2E8F0" strokeWidth="1" strokeDasharray="3 3" />
      <path d="M275 120 Q275 150 270 155" stroke="#E2E8F0" strokeWidth="1" strokeDasharray="3 3" />

      {/* Decorative dots */}
      <circle cx="350" cy="50" r="3" fill="#22D3EE" opacity="0.3" />
      <circle cx="40" cy="120" r="2" fill="#A78BFA" opacity="0.3" />
      <circle cx="370" cy="200" r="2.5" fill="#FB7185" opacity="0.3" />
    </svg>
  );
}

function ShieldIcon() {
  return (
    <svg width="40" height="40" viewBox="0 0 40 40" fill="none">
      <path d="M20 3L5 10v10c0 10.5 6.4 20.3 15 22.5 8.6-2.2 15-12 15-22.5V10L20 3z" fill="#0F172A" opacity="0.06" />
      <path d="M20 6L8 12v8c0 8.4 5.12 16.24 12 18 6.88-1.76 12-9.6 12-18v-8L20 6z" fill="none" stroke="#0F172A" strokeWidth="1.5" />
      <circle cx="20" cy="18" r="3" fill="#22D3EE" />
      <rect x="18.5" y="20" width="3" height="5" rx="1" fill="#22D3EE" opacity="0.6" />
    </svg>
  );
}

function TokenIcon() {
  return (
    <svg width="40" height="40" viewBox="0 0 40 40" fill="none">
      <rect x="4" y="10" width="32" height="20" rx="6" fill="#0F172A" opacity="0.06" />
      <rect x="6" y="12" width="28" height="16" rx="4" fill="none" stroke="#0F172A" strokeWidth="1.5" />
      <rect x="10" y="17" width="6" height="6" rx="1.5" fill="#A78BFA" />
      <rect x="19" y="18" width="12" height="2" rx="1" fill="#CBD5E1" />
      <rect x="19" y="22" width="8" height="2" rx="1" fill="#E2E8F0" />
    </svg>
  );
}

function EyeIcon() {
  return (
    <svg width="40" height="40" viewBox="0 0 40 40" fill="none">
      <circle cx="20" cy="20" r="16" fill="#FB7185" opacity="0.06" />
      <path d="M8 20s5.4-9 12-9 12 9 12 9-5.4 9-12 9-12-9-12-9z" fill="none" stroke="#0F172A" strokeWidth="1.5" />
      <circle cx="20" cy="20" r="4" fill="none" stroke="#0F172A" strokeWidth="1.5" />
      <circle cx="20" cy="20" r="1.5" fill="#FB7185" />
    </svg>
  );
}

function RevokeIcon() {
  return (
    <svg width="40" height="40" viewBox="0 0 40 40" fill="none">
      <circle cx="20" cy="20" r="16" fill="#0F172A" opacity="0.06" />
      <circle cx="20" cy="20" r="12" fill="none" stroke="#0F172A" strokeWidth="1.5" />
      <path d="M15 15l10 10M25 15l-10 10" stroke="#FB7185" strokeWidth="1.5" strokeLinecap="round" />
    </svg>
  );
}

function APIIcon() {
  return (
    <svg width="40" height="40" viewBox="0 0 40 40" fill="none">
      <rect x="4" y="4" width="32" height="32" rx="8" fill="#22D3EE" opacity="0.06" />
      <path d="M12 16l-4 4 4 4" stroke="#0F172A" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" />
      <path d="M28 16l4 4-4 4" stroke="#0F172A" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" />
      <path d="M22 12l-4 16" stroke="#22D3EE" strokeWidth="1.5" strokeLinecap="round" />
    </svg>
  );
}

function SelfHostIcon() {
  return (
    <svg width="40" height="40" viewBox="0 0 40 40" fill="none">
      <rect x="8" y="6" width="24" height="28" rx="4" fill="#0F172A" opacity="0.06" />
      <rect x="10" y="8" width="20" height="10" rx="3" fill="none" stroke="#0F172A" strokeWidth="1.5" />
      <rect x="10" y="22" width="20" height="10" rx="3" fill="none" stroke="#0F172A" strokeWidth="1.5" />
      <circle cx="14" cy="13" r="1.5" fill="#34D399" />
      <circle cx="14" cy="27" r="1.5" fill="#34D399" />
      <rect x="18" y="12" width="8" height="2" rx="1" fill="#CBD5E1" />
      <rect x="18" y="26" width="8" height="2" rx="1" fill="#CBD5E1" />
    </svg>
  );
}

function PackageIllustration() {
  return (
    <svg viewBox="0 0 360 240" fill="none" className="w-full max-w-sm">
      {/* Package box */}
      <rect x="80" y="40" width="200" height="140" rx="8" fill="#F8FAFC" stroke="#E2E8F0" strokeWidth="1.5" />
      <line x1="80" y1="80" x2="280" y2="80" stroke="#E2E8F0" strokeWidth="1" />
      <line x1="180" y1="40" x2="180" y2="80" stroke="#E2E8F0" strokeWidth="1" />

      {/* QR on package (no address!) */}
      <rect x="120" y="100" width="50" height="50" rx="4" fill="white" stroke="#0F172A" strokeWidth="1" />
      <rect x="128" y="108" width="8" height="8" rx="1" fill="#0F172A" />
      <rect x="138" y="108" width="8" height="8" rx="1" fill="#0F172A" />
      <rect x="128" y="118" width="8" height="8" rx="1" fill="#0F172A" />
      <rect x="138" y="118" width="8" height="8" rx="1" fill="#22D3EE" opacity="0.4" />
      <rect x="148" y="108" width="8" height="8" rx="1" fill="#0F172A" />
      <rect x="148" y="118" width="8" height="8" rx="1" fill="#0F172A" />
      <rect x="128" y="128" width="8" height="8" rx="1" fill="#0F172A" />
      <rect x="138" y="128" width="8" height="8" rx="1" fill="#0F172A" />
      <rect x="148" y="128" width="8" height="8" rx="1" fill="#0F172A" />

      {/* Reference code */}
      <rect x="185" y="105" width="76" height="20" rx="4" fill="#F1F5F9" />
      <text x="193" y="119" fontSize="9" fontFamily="monospace" fontWeight="600" fill="#0F172A">AP-7X3K-9M2P</text>

      {/* Zone code */}
      <rect x="185" y="130" width="50" height="16" rx="4" fill="#F1F5F9" />
      <text x="192" y="142" fontSize="8" fontFamily="monospace" fill="#64748B">IST-340</text>

      {/* Crossed out address */}
      <g opacity="0.3">
        <rect x="100" y="190" width="160" height="30" rx="4" fill="#FEE2E2" />
        <text x="115" y="208" fontSize="9" fill="#EF4444">No plaintext address</text>
        <line x1="110" y1="209" x2="250" y2="209" stroke="#EF4444" strokeWidth="0.5" />
      </g>
    </svg>
  );
}

/* ─── Page ────────────────────────────────────────────────── */

const features = [
  { icon: <ShieldIcon />, title: "Encrypted Vault", desc: "Store multiple addresses in your private vault. Home, work, shipping — organized and secure.", accent: "cyan" },
  { icon: <TokenIcon />, title: "Tokenized Sharing", desc: "Share via unique links, QR codes, or short codes. Recipients never see your raw data until you allow it.", accent: "lavender" },
  { icon: <EyeIcon />, title: "Access Monitoring", desc: "See exactly who accessed your address, when, from what device. Get notified on every access.", accent: "coral" },
  { icon: <RevokeIcon />, title: "Instant Revocation", desc: "Revoke access anytime. Set expiry dates and access limits. One-time delivery? One-time access.", accent: "coral" },
  { icon: <APIIcon />, title: "Developer API", desc: "REST API with OAuth2 for businesses. Scoped access: full, delivery, zone, or verify-only.", accent: "cyan" },
  { icon: <SelfHostIcon />, title: "Self-Hostable", desc: "Run on your own infrastructure with Docker. AGPL-3.0 open source. Your data never leaves your servers.", accent: "mint" },
];

const steps = [
  { num: "01", title: "Store", desc: "Add your addresses to your encrypted vault. Home, work, anything.", color: "#22D3EE" },
  { num: "02", title: "Share", desc: "Generate a tokenized link or QR code with access rules and scope.", color: "#A78BFA" },
  { num: "03", title: "Monitor", desc: "Track every access. See who, when, and from where. Get notified.", color: "#FB7185" },
  { num: "04", title: "Revoke", desc: "Expire or revoke access instantly. The token dies, the address stays safe.", color: "#34D399" },
];

const stats = [
  { value: "3.1B+", label: "personal records leaked in 2024" },
  { value: "53%", label: "of breaches expose home addresses" },
  { value: "17%", label: "of identity theft from physical mail" },
  { value: "$56B", label: "annual identity theft cost in the US" },
];

export default function Home() {
  return (
    <div className="min-h-screen bg-[#FAFBFD]">
      <Header />

      <main className="mx-auto max-w-6xl px-6">
        {/* ─── Hero ────────────────────────────────────────── */}
        <section className="py-20 lg:py-28 relative">
          <div className="hero-mesh" />
          <div className="grid lg:grid-cols-2 gap-12 items-center relative z-10">
            <div className="text-center lg:text-left">
              <div className="inline-flex items-center gap-2 rounded-full bg-[#0F172A]/[0.03] border border-[#E2E8F0] px-4 py-1.5 text-[12px] font-medium text-[#64748B] mb-6 fade-in">
                <span className="w-1.5 h-1.5 rounded-full bg-[#34D399]" />
                Open source &middot; Self-hostable &middot; EU hosted
              </div>
              <h1 className="text-4xl sm:text-5xl lg:text-6xl font-bold tracking-tight leading-[1.1] fade-in fade-in-delay-1">
                Share your address,
                <br />
                <span className="brand-gradient">not your privacy.</span>
              </h1>
              <p className="mt-6 text-base lg:text-lg text-[#64748B] max-w-lg leading-relaxed fade-in fade-in-delay-2">
                Store addresses once. Share tokenized links with access control,
                real-time monitoring, and instant revocation. For people who care
                where their data goes.
              </p>
              <div className="mt-8 flex flex-wrap gap-3 justify-center lg:justify-start fade-in fade-in-delay-3">
                <a href="/register" className="btn-primary rounded-full px-7 py-3 text-sm font-semibold">
                  <span>Start for free</span>
                </a>
                <a href="https://github.com/addrpass/addrpass" target="_blank" rel="noopener noreferrer" className="btn-secondary inline-flex items-center gap-2 rounded-full border border-[#E2E8F0] bg-white px-6 py-3 text-sm font-semibold text-[#0F172A]">
                  <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z" /></svg>
                  View source
                </a>
              </div>
              <div className="mt-5 flex items-center gap-4 justify-center lg:justify-start fade-in fade-in-delay-4">
                <div className="flex items-center gap-1.5 text-xs text-[#94A3B8]">
                  <svg width="12" height="12" viewBox="0 0 16 16" fill="#34D399"><circle cx="8" cy="8" r="8" /></svg>
                  No credit card required
                </div>
                <div className="flex items-center gap-1.5 text-xs text-[#94A3B8]">
                  <svg width="12" height="12" viewBox="0 0 16 16" fill="#34D399"><circle cx="8" cy="8" r="8" /></svg>
                  AGPL-3.0 open source
                </div>
              </div>
            </div>
            <div className="hidden lg:block fade-in fade-in-delay-3">
              <HeroIllustration />
            </div>
          </div>
        </section>

        {/* ─── Stats ───────────────────────────────────────── */}
        <section className="brand-gradient-subtle rounded-2xl p-8 lg:p-10 -mx-2">
          <div className="grid grid-cols-2 lg:grid-cols-4 gap-6">
            {stats.map((s) => (
              <div key={s.label} className="text-center">
                <div className="text-2xl lg:text-3xl font-bold text-[#0F172A]">{s.value}</div>
                <div className="text-xs text-[#64748B] mt-1 leading-relaxed">{s.label}</div>
              </div>
            ))}
          </div>
        </section>

        {/* ─── How It Works ────────────────────────────────── */}
        <section id="how-it-works" className="py-24">
          <div className="text-center mb-16">
            <h2 className="text-3xl lg:text-4xl font-bold tracking-tight">How it works</h2>
            <p className="mt-3 text-[#64748B] max-w-md mx-auto">Four steps to take back control of your address data.</p>
          </div>
          <div className="grid sm:grid-cols-2 lg:grid-cols-4 gap-5">
            {steps.map((step) => (
              <div key={step.num} className="card p-6 text-center group">
                <div className="text-3xl font-bold mb-3" style={{ color: step.color, opacity: 0.25 }}>{step.num}</div>
                <h3 className="font-semibold text-lg text-[#0F172A] mb-2">{step.title}</h3>
                <p className="text-sm text-[#64748B] leading-relaxed">{step.desc}</p>
              </div>
            ))}
          </div>
        </section>

        <div className="section-divider" />

        {/* ─── Features ────────────────────────────────────── */}
        <section id="features" className="py-24">
          <div className="text-center mb-16">
            <h2 className="text-3xl lg:text-4xl font-bold tracking-tight">Built for privacy,<br /><span className="brand-gradient">designed for control.</span></h2>
            <p className="mt-3 text-[#64748B] max-w-md mx-auto">Everything you need to share your address without losing ownership.</p>
          </div>
          <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-5">
            {features.map((f) => (
              <div key={f.title} className="card p-6">
                <div className="mb-4">{f.icon}</div>
                <h3 className="font-semibold text-[15px] text-[#0F172A] mb-2">{f.title}</h3>
                <p className="text-sm text-[#64748B] leading-relaxed">{f.desc}</p>
              </div>
            ))}
          </div>
        </section>

        <div className="section-divider" />

        {/* ─── Delivery Use Case ───────────────────────────── */}
        <section id="delivery" className="py-24">
          <div className="grid lg:grid-cols-2 gap-16 items-center">
            <div>
              <div className="inline-flex items-center gap-2 pill text-[#64748B] mb-6">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><rect x="1" y="3" width="15" height="13" /><polygon points="16 8 20 8 23 11 23 16 16 16 16 8" /><circle cx="5.5" cy="18.5" r="2.5" /><circle cx="18.5" cy="18.5" r="2.5" /></svg>
                For delivery companies
              </div>
              <h2 className="text-3xl lg:text-4xl font-bold tracking-tight leading-tight">
                No address on the label.
                <br />
                <span className="brand-gradient">Just a QR code.</span>
              </h2>
              <p className="mt-4 text-[#64748B] leading-relaxed max-w-lg">
                Delivery companies integrate via API. The shipping label shows only a QR code and zone code &mdash; no plaintext address. Drivers scan to resolve. Customers see every access in real time.
              </p>
              <div className="mt-8 space-y-3">
                {[
                  "E-commerce stores token, not address (reduced liability)",
                  "Warehouse gets zone-only scope for sorting",
                  "Driver gets delivery scope (address, no phone)",
                  "Customer gets real-time access notifications",
                  "Revoke after delivery \u2014 token dies forever",
                ].map((item) => (
                  <div key={item} className="flex items-start gap-3">
                    <svg width="18" height="18" viewBox="0 0 20 20" fill="none" className="mt-0.5 flex-shrink-0">
                      <circle cx="10" cy="10" r="10" fill="#22D3EE" opacity="0.1" />
                      <path d="M6 10l3 3 5-5" stroke="#22D3EE" strokeWidth="1.5" strokeLinecap="round" strokeLinejoin="round" />
                    </svg>
                    <span className="text-sm text-[#334155]">{item}</span>
                  </div>
                ))}
              </div>
            </div>
            <div className="flex justify-center">
              <PackageIllustration />
            </div>
          </div>
        </section>

        <div className="section-divider" />

        {/* ─── Open Source CTA ──────────────────────────────── */}
        <section className="py-24">
          <div className="text-center max-w-2xl mx-auto">
            <div className="w-16 h-16 mx-auto mb-6 rounded-2xl bg-[#0F172A] flex items-center justify-center">
              <svg width="28" height="28" viewBox="0 0 24 24" fill="white"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z" /></svg>
            </div>
            <h2 className="text-3xl lg:text-4xl font-bold tracking-tight mb-4">
              Fully <span className="brand-gradient">open source</span>
            </h2>
            <p className="text-[#64748B] leading-relaxed mb-8">
              AddrPass is open source under AGPL-3.0. Self-host on your infrastructure, audit every line of code, or contribute.
              Built with Go, PostgreSQL, and Next.js. Deployed with Docker.
            </p>
            <div className="flex flex-wrap gap-3 justify-center">
              <a href="https://github.com/addrpass/addrpass" target="_blank" rel="noopener noreferrer" className="btn-secondary inline-flex items-center gap-2 rounded-full border border-[#E2E8F0] bg-white px-6 py-3 text-sm font-semibold text-[#0F172A]">
                <svg width="16" height="16" viewBox="0 0 24 24" fill="currentColor"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z" /></svg>
                View on GitHub
              </a>
              <a href="/register" className="btn-primary rounded-full px-6 py-3 text-sm font-semibold">
                <span>Start for free</span>
              </a>
            </div>
          </div>
        </section>
      </main>

      <Footer />
    </div>
  );
}
