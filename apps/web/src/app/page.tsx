"use client";

import Link from "next/link";
import Header from "./components/Header";
import Footer from "./components/Footer";

/* ─── Custom Illustrations ────────────────────────────────── */

function HeroIllustration() {
  return (
    <svg viewBox="0 0 400 340" fill="none" className="w-full max-w-md">
      {/* Address card being locked */}
      <rect x="50" y="50" width="240" height="150" rx="14" fill="#F8FAFC" stroke="#E2E8F0" strokeWidth="1.5" />
      <rect x="72" y="78" width="100" height="7" rx="3.5" fill="#CBD5E1" />
      <rect x="72" y="94" width="160" height="5" rx="2.5" fill="#E2E8F0" />
      <rect x="72" y="107" width="130" height="5" rx="2.5" fill="#E2E8F0" />
      <rect x="72" y="120" width="90" height="5" rx="2.5" fill="#E2E8F0" />
      {/* Lock badge */}
      <rect x="72" y="145" width="52" height="26" rx="7" fill="#0F172A" />
      <circle cx="98" cy="155" r="3.5" fill="#22D3EE" />
      <rect x="96.5" y="157" width="3" height="5" rx="1" fill="#22D3EE" opacity="0.6" />

      {/* Token output */}
      <g className="float">
        <rect x="220" y="140" width="110" height="38" rx="10" fill="white" stroke="#22D3EE" strokeWidth="1.5" filter="url(#shadow)" />
        <rect x="234" y="152" width="32" height="14" rx="3" fill="#22D3EE" opacity="0.1" />
        <text x="238" y="163" fontSize="8" fontFamily="monospace" fontWeight="600" fill="#22D3EE">TOKEN</text>
        <rect x="272" y="154" width="44" height="4" rx="2" fill="#CBD5E1" />
        <rect x="272" y="162" width="30" height="3" rx="1.5" fill="#E2E8F0" />
      </g>

      {/* Arrow */}
      <path d="M195 160 L215 160" stroke="#22D3EE" strokeWidth="1.5" strokeDasharray="4 3" />
      <polygon points="213,156 220,160 213,164" fill="#22D3EE" />

      {/* QR floating */}
      <g className="float float-delay-1">
        <rect x="280" y="55" width="52" height="52" rx="8" fill="white" stroke="#A78BFA" strokeWidth="1.2" />
        <g transform="translate(290,65)">
          <rect width="7" height="7" rx="1" fill="#0F172A" />
          <rect x="10" width="7" height="7" rx="1" fill="#0F172A" />
          <rect x="20" width="7" height="7" rx="1" fill="#0F172A" />
          <rect y="10" width="7" height="7" rx="1" fill="#0F172A" />
          <rect x="10" y="10" width="7" height="7" rx="1" fill="#A78BFA" opacity="0.3" />
          <rect x="20" y="10" width="7" height="7" rx="1" fill="#0F172A" />
          <rect y="20" width="7" height="7" rx="1" fill="#0F172A" />
          <rect x="10" y="20" width="7" height="7" rx="1" fill="#0F172A" />
          <rect x="20" y="20" width="7" height="7" rx="1" fill="#0F172A" />
        </g>
      </g>

      {/* Monitoring pulse */}
      <g className="float float-delay-2">
        <circle cx="80" cy="250" r="20" fill="white" stroke="#FB7185" strokeWidth="1.2" />
        <path d="M70 250s4.5-6 10-6 10 6 10 6-4.5 6-10 6-10-6-10-6z" fill="none" stroke="#FB7185" strokeWidth="1.2" />
        <circle cx="80" cy="250" r="3" fill="#FB7185" />
      </g>

      {/* Notification badge */}
      <g className="float float-delay-1">
        <rect x="220" y="220" width="130" height="44" rx="10" fill="white" stroke="#E2E8F0" strokeWidth="1" />
        <circle cx="240" cy="242" r="8" fill="#34D399" opacity="0.15" />
        <path d="M236 242l3 3 5-5" stroke="#34D399" strokeWidth="1.2" strokeLinecap="round" />
        <rect x="254" y="236" width="80" height="4" rx="2" fill="#CBD5E1" />
        <rect x="254" y="244" width="56" height="3" rx="1.5" fill="#E2E8F0" />
      </g>

      {/* Dots */}
      <circle cx="360" cy="40" r="2.5" fill="#22D3EE" opacity="0.25" />
      <circle cx="30" cy="140" r="2" fill="#A78BFA" opacity="0.25" />
      <circle cx="375" cy="210" r="2" fill="#FB7185" opacity="0.25" />
      <circle cx="170" cy="280" r="1.5" fill="#34D399" opacity="0.3" />

      <defs>
        <filter id="shadow" x="-4" y="-2" width="120" height="48">
          <feDropShadow dx="0" dy="2" stdDeviation="4" floodOpacity="0.06" />
        </filter>
      </defs>
    </svg>
  );
}

function PackageIllustration() {
  return (
    <svg viewBox="0 0 360 240" fill="none" className="w-full max-w-sm">
      <rect x="80" y="30" width="200" height="140" rx="8" fill="#F8FAFC" stroke="#E2E8F0" strokeWidth="1.5" />
      <line x1="80" y1="70" x2="280" y2="70" stroke="#E2E8F0" strokeWidth="1" />
      <line x1="180" y1="30" x2="180" y2="70" stroke="#E2E8F0" strokeWidth="1" />
      <rect x="110" y="90" width="50" height="50" rx="4" fill="white" stroke="#0F172A" strokeWidth="1" />
      <g transform="translate(118,98)">
        <rect width="7" height="7" rx="1" fill="#0F172A" /><rect x="9" width="7" height="7" rx="1" fill="#0F172A" /><rect x="18" width="7" height="7" rx="1" fill="#0F172A" />
        <rect y="9" width="7" height="7" rx="1" fill="#0F172A" /><rect x="9" y="9" width="7" height="7" rx="1" fill="#22D3EE" opacity="0.35" /><rect x="18" y="9" width="7" height="7" rx="1" fill="#0F172A" />
        <rect y="18" width="7" height="7" rx="1" fill="#0F172A" /><rect x="9" y="18" width="7" height="7" rx="1" fill="#0F172A" /><rect x="18" y="18" width="7" height="7" rx="1" fill="#0F172A" />
      </g>
      <rect x="178" y="95" width="76" height="18" rx="4" fill="#F1F5F9" />
      <text x="186" y="108" fontSize="9" fontFamily="monospace" fontWeight="600" fill="#0F172A">AP-7X3K-9M2P</text>
      <rect x="178" y="120" width="50" height="14" rx="4" fill="#F1F5F9" />
      <text x="184" y="131" fontSize="8" fontFamily="monospace" fill="#64748B">IST-340</text>
      <g opacity="0.25">
        <rect x="100" y="185" width="160" height="26" rx="4" fill="#FEE2E2" />
        <text x="130" y="202" fontSize="9" fill="#EF4444">No plaintext address</text>
      </g>
    </svg>
  );
}

/* ─── Icon Components ─────────────────────────────────────── */

function VaultIcon() {
  return (
    <svg width="36" height="36" viewBox="0 0 36 36" fill="none">
      <rect x="4" y="8" width="28" height="22" rx="5" fill="#0F172A" opacity="0.05" />
      <rect x="6" y="10" width="24" height="18" rx="4" fill="none" stroke="#0F172A" strokeWidth="1.3" />
      <circle cx="18" cy="18" r="4" fill="none" stroke="#22D3EE" strokeWidth="1.3" />
      <circle cx="18" cy="18" r="1.5" fill="#22D3EE" />
      <rect x="17" y="21" width="2" height="4" rx="1" fill="#22D3EE" opacity="0.5" />
    </svg>
  );
}

function LinkIcon() {
  return (
    <svg width="36" height="36" viewBox="0 0 36 36" fill="none">
      <circle cx="18" cy="18" r="14" fill="#A78BFA" opacity="0.06" />
      <path d="M15 21a5.5 5.5 0 007.78.56l2.5-2.5a5.5 5.5 0 00-7.78-7.78l-1.42 1.42" stroke="#0F172A" strokeWidth="1.3" strokeLinecap="round" />
      <path d="M21 15a5.5 5.5 0 00-7.78-.56l-2.5 2.5a5.5 5.5 0 007.78 7.78l1.42-1.42" stroke="#A78BFA" strokeWidth="1.3" strokeLinecap="round" />
    </svg>
  );
}

function MonitorIcon() {
  return (
    <svg width="36" height="36" viewBox="0 0 36 36" fill="none">
      <circle cx="18" cy="18" r="14" fill="#FB7185" opacity="0.06" />
      <path d="M9 18s4-7 9-7 9 7 9 7-4 7-9 7-9-7-9-7z" fill="none" stroke="#0F172A" strokeWidth="1.3" />
      <circle cx="18" cy="18" r="3" fill="none" stroke="#0F172A" strokeWidth="1.3" />
      <circle cx="18" cy="18" r="1.2" fill="#FB7185" />
    </svg>
  );
}

function ClockIcon() {
  return (
    <svg width="36" height="36" viewBox="0 0 36 36" fill="none">
      <circle cx="18" cy="18" r="14" fill="#34D399" opacity="0.06" />
      <circle cx="18" cy="18" r="10" fill="none" stroke="#0F172A" strokeWidth="1.3" />
      <polyline points="18,12 18,18 22,20" stroke="#34D399" strokeWidth="1.3" strokeLinecap="round" />
    </svg>
  );
}

function CodeIcon() {
  return (
    <svg width="36" height="36" viewBox="0 0 36 36" fill="none">
      <rect x="4" y="4" width="28" height="28" rx="6" fill="#22D3EE" opacity="0.06" />
      <path d="M13 14l-4 4 4 4" stroke="#0F172A" strokeWidth="1.3" strokeLinecap="round" strokeLinejoin="round" />
      <path d="M23 14l4 4-4 4" stroke="#0F172A" strokeWidth="1.3" strokeLinecap="round" strokeLinejoin="round" />
      <path d="M20 11l-4 14" stroke="#22D3EE" strokeWidth="1.3" strokeLinecap="round" />
    </svg>
  );
}

function ServerIcon() {
  return (
    <svg width="36" height="36" viewBox="0 0 36 36" fill="none">
      <rect x="6" y="5" width="24" height="26" rx="4" fill="#0F172A" opacity="0.05" />
      <rect x="8" y="7" width="20" height="9" rx="3" fill="none" stroke="#0F172A" strokeWidth="1.3" />
      <rect x="8" y="20" width="20" height="9" rx="3" fill="none" stroke="#0F172A" strokeWidth="1.3" />
      <circle cx="12" cy="11.5" r="1.2" fill="#34D399" />
      <circle cx="12" cy="24.5" r="1.2" fill="#34D399" />
      <rect x="16" y="10.5" width="8" height="2" rx="1" fill="#CBD5E1" />
      <rect x="16" y="23.5" width="8" height="2" rx="1" fill="#CBD5E1" />
    </svg>
  );
}

/* ─── Data ────────────────────────────────────────────────── */

const features = [
  { icon: <VaultIcon />, title: "Encrypted Vault", desc: "Store multiple addresses securely. Home, work, PO box — organized in one place under your control." },
  { icon: <LinkIcon />, title: "Tokenized Sharing", desc: "Share via unique links, QR codes, or short codes. Recipients get a reference, never your raw address." },
  { icon: <MonitorIcon />, title: "Access Monitoring", desc: "See who accessed your address, when, and from where. Know when DHL reads your data." },
  { icon: <ClockIcon />, title: "Expiration & Limits", desc: "Set expiry dates and access limits per share. One-time delivery? One-time access. Then it dies." },
  { icon: <CodeIcon />, title: "Developer API", desc: "REST API with OAuth2 for businesses. Scoped access levels: full, delivery, zone, or verify-only." },
  { icon: <ServerIcon />, title: "Self-Hostable", desc: "Docker Compose, one command. Run on your own servers. AGPL-3.0 — inspect every line of code." },
];

const steps = [
  { num: "01", title: "Store", desc: "Add your addresses to your encrypted vault.", color: "#22D3EE" },
  { num: "02", title: "Share", desc: "Generate a tokenized link or QR code with rules.", color: "#A78BFA" },
  { num: "03", title: "Monitor", desc: "Track every access in real time. Get notified.", color: "#FB7185" },
  { num: "04", title: "Revoke", desc: "Kill the token. Instantly. The address stays safe.", color: "#34D399" },
];

/* ─── Page ────────────────────────────────────────────────── */

export default function Home() {
  return (
    <div className="min-h-screen bg-[#FAFBFD]">
      <Header />
      <main className="mx-auto max-w-6xl px-6">

        {/* ─── Hero ──────────────────────────────────────── */}
        <section className="py-20 lg:py-28 relative">
          <div className="hero-mesh" />
          <div className="grid lg:grid-cols-2 gap-12 items-center relative z-10">
            <div className="text-center lg:text-left">
              <h1 className="text-4xl sm:text-5xl lg:text-[54px] font-bold tracking-tight leading-[1.1] fade-in">
                Stop giving your address
                <br />
                <span className="brand-gradient">to every website.</span>
              </h1>
              <p className="mt-6 text-base lg:text-[17px] text-[#64748B] max-w-lg leading-relaxed fade-in fade-in-delay-1">
                AddrPass lets you store your address once and share tokenized links instead. You control who sees it, for how long, and you can revoke access instantly.
              </p>
              <div className="mt-8 flex flex-wrap gap-3 justify-center lg:justify-start fade-in fade-in-delay-2">
                <Link href="/register" className="btn-primary rounded-full px-7 py-3 text-sm font-semibold">
                  <span>Get started free</span>
                </Link>
                <Link href="/pricing" className="btn-secondary inline-flex items-center gap-2 rounded-full border border-[#E2E8F0] bg-white px-6 py-3 text-sm font-semibold text-[#0F172A]">
                  View pricing
                </Link>
              </div>
              <div className="mt-5 flex flex-wrap items-center gap-x-5 gap-y-2 justify-center lg:justify-start fade-in fade-in-delay-3 text-xs text-[#94A3B8]">
                <span className="flex items-center gap-1.5"><span className="w-1.5 h-1.5 rounded-full bg-[#34D399]" /> Free forever for personal use</span>
                <span className="flex items-center gap-1.5"><span className="w-1.5 h-1.5 rounded-full bg-[#34D399]" /> No credit card required</span>
                <span className="flex items-center gap-1.5"><span className="w-1.5 h-1.5 rounded-full bg-[#34D399]" /> Self-hostable</span>
              </div>
            </div>
            <div className="hidden lg:block fade-in fade-in-delay-3">
              <HeroIllustration />
            </div>
          </div>
        </section>

        {/* ─── Social Proof Bar ──────────────────────────── */}
        <section className="rounded-2xl bg-[#0F172A] p-8 lg:p-10 text-white -mx-2">
          <div className="grid grid-cols-2 lg:grid-cols-4 gap-6">
            {[
              { value: "1.7B+", label: "records compromised in 2024", source: "ITRC 2024 Annual Data Breach Report", href: "https://www.idtheftcenter.org/post/2024-annual-data-breach-report-near-record-compromises/" },
              { value: "53%", label: "of breaches expose customer PII", source: "IBM Cost of a Data Breach Report 2025", href: "https://www.ibm.com/reports/data-breach" },
              { value: "161%", label: "rise in mail theft (2020–2023)", source: "USPS Postal Inspection Service", href: "https://www.uspis.gov/news/press-releases" },
              { value: "$47B", label: "annual ID fraud losses (2024)", source: "Javelin / AARP 2025 Identity Fraud Study", href: "https://www.aarp.org/money/scams-fraud/javelin-identity-theft-report-2024/" },
            ].map((s) => (
              <div key={s.label} className="text-center">
                <div className="text-2xl lg:text-3xl font-bold brand-gradient">{s.value}</div>
                <div className="text-xs text-[#94A3B8] mt-1.5 leading-relaxed">{s.label}</div>
              </div>
            ))}
          </div>
          <div className="mt-6 pt-4 border-t border-white/10 flex flex-wrap justify-center gap-x-6 gap-y-1">
            {[
              { source: "ITRC 2024 Annual Data Breach Report", href: "https://www.idtheftcenter.org/post/2024-annual-data-breach-report-near-record-compromises/" },
              { source: "IBM Cost of a Data Breach Report 2025", href: "https://www.ibm.com/reports/data-breach" },
              { source: "USPS Postal Inspection Service", href: "https://www.uspis.gov/news/press-releases" },
              { source: "Javelin / AARP 2025 Identity Fraud Study", href: "https://www.aarp.org/money/scams-fraud/javelin-identity-theft-report-2024/" },
            ].map((r) => (
              <a key={r.source} href={r.href} target="_blank" rel="noopener noreferrer" className="text-[10px] text-[#64748B] hover:text-[#94A3B8] transition-colors underline underline-offset-2">
                {r.source}
              </a>
            ))}
          </div>
        </section>

        {/* ─── How It Works ──────────────────────────────── */}
        <section id="how-it-works" className="py-24">
          <div className="text-center mb-16">
            <p className="text-xs font-semibold uppercase tracking-widest text-[#22D3EE] mb-3">How it works</p>
            <h2 className="text-3xl lg:text-4xl font-bold tracking-tight">Four steps to take back control</h2>
          </div>
          <div className="grid sm:grid-cols-2 lg:grid-cols-4 gap-5">
            {steps.map((step) => (
              <div key={step.num} className="card p-6 text-center">
                <div className="text-4xl font-black mb-4 leading-none" style={{ color: step.color, opacity: 0.18 }}>{step.num}</div>
                <h3 className="font-semibold text-[17px] text-[#0F172A] mb-2">{step.title}</h3>
                <p className="text-sm text-[#64748B] leading-relaxed">{step.desc}</p>
              </div>
            ))}
          </div>
        </section>

        <div className="section-divider" />

        {/* ─── Features ──────────────────────────────────── */}
        <section id="features" className="py-24">
          <div className="text-center mb-16">
            <p className="text-xs font-semibold uppercase tracking-widest text-[#A78BFA] mb-3">Features</p>
            <h2 className="text-3xl lg:text-4xl font-bold tracking-tight">Everything you need.<br /><span className="brand-gradient">Nothing you don&apos;t.</span></h2>
          </div>
          <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-5">
            {features.map((f) => (
              <div key={f.title} className="card p-6">
                <div className="mb-4">{f.icon}</div>
                <h3 className="font-semibold text-[15px] text-[#0F172A] mb-2">{f.title}</h3>
                <p className="text-[13px] text-[#64748B] leading-relaxed">{f.desc}</p>
              </div>
            ))}
          </div>
        </section>

        <div className="section-divider" />

        {/* ─── Delivery Use Case ─────────────────────────── */}
        <section id="delivery" className="py-24">
          <div className="grid lg:grid-cols-2 gap-16 items-center">
            <div>
              <p className="text-xs font-semibold uppercase tracking-widest text-[#FB7185] mb-3">For delivery companies</p>
              <h2 className="text-3xl lg:text-4xl font-bold tracking-tight leading-tight">
                QR code on the label.
                <br />
                <span className="brand-gradient">No address printed.</span>
              </h2>
              <p className="mt-4 text-[#64748B] leading-relaxed max-w-lg">
                Integrate via API. The shipping label shows a QR code and zone code &mdash; no plaintext address. Drivers scan to resolve. Customers see every access.
              </p>
              <div className="mt-8 space-y-3">
                {[
                  "E-commerce stores a token, not the address",
                  "Warehouse gets zone-only scope for sorting",
                  "Driver gets delivery scope (no phone number)",
                  "Customer gets notified on every resolution",
                  "Revoke after delivery \u2014 the token dies",
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
              <div className="mt-8">
                <Link href="/register" className="btn-primary rounded-full px-6 py-3 text-sm font-semibold inline-block">
                  <span>Start integrating</span>
                </Link>
              </div>
            </div>
            <div className="flex justify-center">
              <PackageIllustration />
            </div>
          </div>
        </section>

        <div className="section-divider" />

        {/* ─── SDK / Developers ──────────────────────────── */}
        <section id="developers" className="py-24">
          <div className="text-center mb-12">
            <p className="text-xs font-semibold uppercase tracking-widest text-[#22D3EE] mb-3">For developers</p>
            <h2 className="text-3xl lg:text-4xl font-bold tracking-tight">Integrate in minutes,<br /><span className="brand-gradient">not months.</span></h2>
            <p className="mt-3 text-[#64748B] max-w-md mx-auto">npm install, add a few lines, and your app supports tokenized address sharing.</p>
          </div>

          <div className="grid lg:grid-cols-2 gap-5">
            {/* Backend */}
            <div className="card p-6">
              <div className="flex items-center gap-2 mb-4">
                <span className="text-xs font-semibold uppercase tracking-wider text-[#64748B]">Backend</span>
                <span className="text-[10px] px-2 py-0.5 rounded-full bg-[#0F172A] text-[#94A3B8]">Node.js</span>
              </div>
              <pre className="bg-[#0F172A] rounded-xl p-5 text-[13px] leading-relaxed overflow-x-auto">
                <code className="text-[#E2E8F0]">
{`import { `}<span className="text-[#22D3EE]">AddrPassClient</span>{` } from `}<span className="text-[#34D399]">&quot;@addrpass/sdk&quot;</span>{`;

const addrpass = new `}<span className="text-[#22D3EE]">AddrPassClient</span>{`({
  clientId: `}<span className="text-[#34D399]">&quot;ap_your_id&quot;</span>{`,
  clientSecret: `}<span className="text-[#34D399]">&quot;aps_your_secret&quot;</span>{`,
});

`}<span className="text-[#64748B]">// Resolve a customer&apos;s shared address</span>{`
const { address } = await addrpass.`}<span className="text-[#A78BFA]">resolve</span>{`(token);
console.log(address.line1, address.city);`}
                </code>
              </pre>
            </div>

            {/* Frontend */}
            <div className="card p-6">
              <div className="flex items-center gap-2 mb-4">
                <span className="text-xs font-semibold uppercase tracking-wider text-[#64748B]">Frontend</span>
                <span className="text-[10px] px-2 py-0.5 rounded-full bg-[#0F172A] text-[#94A3B8]">React</span>
              </div>
              <pre className="bg-[#0F172A] rounded-xl p-5 text-[13px] leading-relaxed overflow-x-auto">
                <code className="text-[#E2E8F0]">
{`import { `}<span className="text-[#22D3EE]">AddrPassProvider</span>{`, `}<span className="text-[#22D3EE]">AddrPassButton</span>{` }
  from `}<span className="text-[#34D399]">&quot;@addrpass/sdk/react&quot;</span>{`;

<`}<span className="text-[#22D3EE]">AddrPassProvider</span>{`
  clientId=`}<span className="text-[#34D399]">&quot;ap_your_id&quot;</span>{`
  redirectUri=`}<span className="text-[#34D399]">&quot;/callback&quot;</span>{`
  scope=`}<span className="text-[#34D399]">&quot;delivery&quot;</span>{`>
  <`}<span className="text-[#22D3EE]">AddrPassButton</span>{`
    onToken={(r) => ship(r.`}<span className="text-[#A78BFA]">share_token</span>{`)} />
</`}<span className="text-[#22D3EE]">AddrPassProvider</span>{`>`}
                </code>
              </pre>
            </div>
          </div>

          <div className="mt-6 flex flex-wrap items-center justify-center gap-4">
            <div className="bg-[#0F172A] rounded-xl px-5 py-3 inline-flex items-center gap-3">
              <span className="text-[#64748B] text-sm font-mono">$</span>
              <span className="text-[#E2E8F0] text-sm font-mono">npm install @addrpass/sdk</span>
            </div>
            <a href="https://www.npmjs.com/package/@addrpass/sdk" target="_blank" rel="noopener noreferrer" className="pill text-[#64748B] text-xs hover:border-[#CBD5E1]">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="#CB3837"><path d="M0 7.334v8h6.666v1.332H12v-1.332h12v-8H0zm6.666 6.664H5.334v-4H3.999v4H1.335V8.667h5.331v5.331zm4 0v1.336H8.001V8.667h5.334v5.331h-2.669zm12.001 0h-1.33v-4h-1.336v4h-1.335v-4h-1.33v4h-2.671V8.667h8.002v5.331zM10.665 10H12v2.667h-1.335V10z" /></svg>
              @addrpass/sdk
            </a>
            <a href="https://github.com/addrpass/addrpass/tree/main/packages/sdk-js" target="_blank" rel="noopener noreferrer" className="pill text-[#64748B] text-xs hover:border-[#CBD5E1]">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z" /></svg>
              SDK Source
            </a>
          </div>
        </section>

        <div className="section-divider" />

        {/* ─── Open Source + CTA ──────────────────────────── */}
        <section className="py-24">
          <div className="rounded-2xl brand-gradient-subtle border border-[#E2E8F0] p-10 lg:p-16 text-center">
            <div className="flex items-center justify-center gap-3 mb-6">
              <a href="https://github.com/addrpass/addrpass" target="_blank" rel="noopener noreferrer" className="pill text-[#64748B] text-xs hover:border-[#CBD5E1]">
                <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z" /></svg>
                Open source on GitHub
              </a>
              <span className="pill text-[#64748B] text-xs">
                <svg width="12" height="12" viewBox="0 0 16 16" fill="none"><path d="M8 1L2 4.5v3.5c0 4.7 2.56 9.07 6 10 3.44-.93 6-5.3 6-10V4.5L8 1z" fill="#0F172A" opacity="0.15" /><path d="M5.5 8l2 2 3.5-3.5" stroke="#34D399" strokeWidth="1.2" strokeLinecap="round" /></svg>
                AGPL-3.0
              </span>
              <span className="pill text-[#64748B] text-xs hidden sm:inline-flex">
                Built in the EU
              </span>
            </div>
            <h2 className="text-3xl lg:text-4xl font-bold tracking-tight mb-4">
              Ready to take back control?
            </h2>
            <p className="text-[#64748B] max-w-md mx-auto mb-8 leading-relaxed">
              Free for personal use. Paid plans for teams and businesses. Self-host if you want full ownership.
            </p>
            <div className="flex flex-wrap gap-3 justify-center">
              <Link href="/register" className="btn-primary rounded-full px-7 py-3 text-sm font-semibold">
                <span>Get started free</span>
              </Link>
              <Link href="/pricing" className="btn-secondary rounded-full border border-[#E2E8F0] bg-white px-6 py-3 text-sm font-semibold text-[#0F172A]">
                View pricing
              </Link>
            </div>
          </div>
        </section>
      </main>
      <Footer />
    </div>
  );
}
