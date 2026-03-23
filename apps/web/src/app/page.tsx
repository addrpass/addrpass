"use client";

import { useState } from "react";
import Header from "./components/Header";
import Footer from "./components/Footer";

function WaitlistForm() {
  const [email, setEmail] = useState("");
  const [submitted, setSubmitted] = useState(false);

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault();
    // TODO: Connect to email service
    setSubmitted(true);
  };

  if (submitted) {
    return (
      <div className="text-center p-4 rounded-xl bg-teal-50 border border-teal-200">
        <p className="text-teal-700 font-medium">
          You&apos;re on the list! We&apos;ll notify you when AddrPass launches.
        </p>
      </div>
    );
  }

  return (
    <form onSubmit={handleSubmit} className="flex gap-3 max-w-md mx-auto lg:mx-0">
      <input
        type="email"
        required
        value={email}
        onChange={(e) => setEmail(e.target.value)}
        placeholder="you@email.com"
        className="flex-1 rounded-full border border-[#E5E7EB] bg-white px-5 py-3 text-sm outline-none focus:border-[#0D9488] focus:ring-2 focus:ring-[#0D9488]/20 transition-all"
      />
      <button
        type="submit"
        className="btn-primary rounded-full px-6 py-3 text-sm font-semibold text-white whitespace-nowrap"
      >
        Join Waitlist
      </button>
    </form>
  );
}

const features = [
  {
    icon: (
      <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
        <path d="M7 11V7a5 5 0 0 1 10 0v4" />
      </svg>
    ),
    title: "Encrypted Vault",
    desc: "Store multiple addresses in your encrypted vault. Home, work, PO box — all in one place.",
    color: "text-teal-600",
    bg: "bg-teal-50",
  },
  {
    icon: (
      <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <path d="M10 13a5 5 0 0 0 7.54.54l3-3a5 5 0 0 0-7.07-7.07l-1.72 1.71" />
        <path d="M14 11a5 5 0 0 0-7.54-.54l-3 3a5 5 0 0 0 7.07 7.07l1.71-1.71" />
      </svg>
    ),
    title: "Tokenized Sharing",
    desc: "Share via unique links, QR codes, or short codes. Never type your address into a form again.",
    color: "text-violet-600",
    bg: "bg-violet-50",
  },
  {
    icon: (
      <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <path d="M1 12s4-8 11-8 11 8 11 8-4 8-11 8-11-8-11-8z" />
        <circle cx="12" cy="12" r="3" />
      </svg>
    ),
    title: "Access Monitoring",
    desc: "See who accessed your address, when, and from what device. Get alerts on suspicious access.",
    color: "text-amber-600",
    bg: "bg-amber-50",
  },
  {
    icon: (
      <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <circle cx="12" cy="12" r="10" />
        <polyline points="12 6 12 12 16 14" />
      </svg>
    ),
    title: "Auto-Expiration",
    desc: "Set expiry dates on every share. One-time delivery? Link expires after use.",
    color: "text-rose-600",
    bg: "bg-rose-50",
  },
  {
    icon: (
      <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <polyline points="16 18 22 12 16 6" />
        <polyline points="8 6 2 12 8 18" />
      </svg>
    ),
    title: "Developer API",
    desc: "REST API for e-commerce, delivery, and enterprise systems. SDKs for JS, Python, Go.",
    color: "text-cyan-600",
    bg: "bg-cyan-50",
  },
  {
    icon: (
      <svg width="24" height="24" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" strokeLinecap="round" strokeLinejoin="round">
        <rect x="2" y="2" width="20" height="8" rx="2" ry="2" />
        <rect x="2" y="14" width="20" height="8" rx="2" ry="2" />
        <line x1="6" y1="6" x2="6.01" y2="6" />
        <line x1="6" y1="18" x2="6.01" y2="18" />
      </svg>
    ),
    title: "Self-Hostable",
    desc: "Open source under AGPL. Run on your infrastructure with Docker. Your data stays yours.",
    color: "text-emerald-600",
    bg: "bg-emerald-50",
  },
];

const steps = [
  {
    num: "1",
    title: "Store",
    desc: "Add your addresses to your encrypted vault",
    color: "from-teal-500 to-teal-600",
  },
  {
    num: "2",
    title: "Share",
    desc: "Generate a link, QR code, or short code with access rules",
    color: "from-violet-500 to-violet-600",
  },
  {
    num: "3",
    title: "Monitor",
    desc: "Track who accessed your address and when",
    color: "from-amber-500 to-amber-600",
  },
  {
    num: "4",
    title: "Revoke",
    desc: "Expire or revoke access anytime, instantly",
    color: "from-rose-500 to-rose-600",
  },
];

const useCases = [
  { label: "Online Shopping", emoji: "🛒" },
  { label: "Package Delivery", emoji: "📦" },
  { label: "Moving", emoji: "🏠" },
  { label: "Freelancing", emoji: "💼" },
  { label: "Healthcare", emoji: "🏥" },
  { label: "E-Commerce API", emoji: "🔌" },
  { label: "Real Estate", emoji: "🏗️" },
  { label: "Personal Safety", emoji: "🛡️" },
];

const stats = [
  { value: "3.1B+", label: "records leaked in 2024" },
  { value: "53%", label: "of breaches expose addresses" },
  { value: "17%", label: "of ID theft from physical mail" },
  { value: "$56B", label: "annual identity theft cost (US)" },
];

export default function Home() {
  return (
    <div className="min-h-screen relative">
      <Header />

      <main className="relative z-10 mx-auto max-w-6xl px-6">
        {/* Hero */}
        <section className="py-20 lg:py-28 relative">
          <div className="hero-glow bg-teal-400 top-0 -left-40" />
          <div className="hero-glow bg-violet-400 top-20 right-0" />

          <div className="text-center lg:text-left max-w-3xl relative z-10">
            <h1 className="text-5xl font-bold tracking-tight sm:text-6xl lg:text-7xl fade-in">
              Your address,
              <br />
              <span className="gradient-text-hero">your control.</span>
            </h1>
            <p className="mt-6 text-lg lg:text-xl text-[#6B7280] max-w-2xl fade-in fade-in-delay-1">
              Stop copying your address into every website. Store it once, share
              tokenized links with access control, monitoring, and expiration.
              Open source and self-hostable.
            </p>
            <div className="mt-10 fade-in fade-in-delay-2" id="waitlist">
              <WaitlistForm />
            </div>
            <p className="mt-4 text-xs text-[#9CA3AF] fade-in fade-in-delay-3">
              Free and open source. No credit card required.
            </p>
          </div>
        </section>

        {/* Stats bar */}
        <section className="border-t border-b border-[#E5E7EB] py-10 -mx-6 px-6 bg-white/50">
          <div className="grid grid-cols-2 lg:grid-cols-4 gap-8">
            {stats.map((s) => (
              <div key={s.label} className="text-center">
                <div className="text-2xl lg:text-3xl font-bold gradient-text">
                  {s.value}
                </div>
                <div className="text-sm text-[#6B7280] mt-1">{s.label}</div>
              </div>
            ))}
          </div>
        </section>

        {/* How it works */}
        <section id="how-it-works" className="py-20 lg:py-28">
          <h2 className="text-center text-3xl lg:text-4xl font-bold mb-4">
            How it <span className="gradient-text">works</span>
          </h2>
          <p className="text-center text-[#6B7280] mb-16 max-w-xl mx-auto">
            Four steps to take back control of your address data.
          </p>
          <div className="grid sm:grid-cols-2 lg:grid-cols-4 gap-6">
            {steps.map((step) => (
              <div key={step.num} className="step-card rounded-2xl p-6 text-center">
                <div
                  className={`w-12 h-12 mx-auto rounded-full bg-gradient-to-br ${step.color} flex items-center justify-center text-white font-bold text-lg mb-4`}
                >
                  {step.num}
                </div>
                <h3 className="font-semibold text-lg mb-2">{step.title}</h3>
                <p className="text-sm text-[#6B7280]">{step.desc}</p>
              </div>
            ))}
          </div>
        </section>

        {/* Features */}
        <section id="features" className="border-t border-[#E5E7EB] py-20 lg:py-28">
          <h2 className="text-center text-3xl lg:text-4xl font-bold mb-4">
            <span className="gradient-text">Features</span>
          </h2>
          <p className="text-center text-[#6B7280] mb-16 max-w-xl mx-auto">
            Everything you need to share your address without losing control.
          </p>
          <div className="grid sm:grid-cols-2 lg:grid-cols-3 gap-6">
            {features.map((f) => (
              <div key={f.title} className="feature-card rounded-2xl p-6">
                <div
                  className={`w-12 h-12 rounded-xl ${f.bg} ${f.color} flex items-center justify-center mb-4`}
                >
                  {f.icon}
                </div>
                <h3 className="font-semibold text-lg mb-2">{f.title}</h3>
                <p className="text-sm text-[#6B7280] leading-relaxed">{f.desc}</p>
              </div>
            ))}
          </div>
        </section>

        {/* Use cases */}
        <section id="use-cases" className="border-t border-[#E5E7EB] py-20 lg:py-28">
          <h2 className="text-center text-3xl lg:text-4xl font-bold mb-4">
            Built for <span className="gradient-text">everyone</span>
          </h2>
          <p className="text-center text-[#6B7280] mb-12 max-w-xl mx-auto">
            From everyday online shopping to enterprise delivery APIs.
          </p>
          <div className="flex flex-wrap justify-center gap-3">
            {useCases.map((uc) => (
              <div
                key={uc.label}
                className="use-case-pill flex items-center gap-2 rounded-full border border-[#E5E7EB] bg-white px-5 py-2.5 text-sm font-medium text-[#374151]"
              >
                <span>{uc.emoji}</span>
                {uc.label}
              </div>
            ))}
          </div>
        </section>

        {/* Open source section */}
        <section className="border-t border-[#E5E7EB] py-20 lg:py-28">
          <div className="text-center max-w-2xl mx-auto">
            <div className="w-16 h-16 mx-auto mb-6 rounded-2xl bg-[#111827] flex items-center justify-center">
              <svg width="28" height="28" viewBox="0 0 24 24" fill="white">
                <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z" />
              </svg>
            </div>
            <h2 className="text-3xl lg:text-4xl font-bold mb-4">
              Fully <span className="gradient-text">open source</span>
            </h2>
            <p className="text-[#6B7280] mb-8">
              AddrPass is open source under AGPL-3.0. Self-host on your
              infrastructure, audit the code, or contribute. Your address data
              never has to leave your servers.
            </p>
            <div className="flex flex-wrap gap-4 justify-center">
              <a
                href="https://github.com/addrpass/addrpass"
                target="_blank"
                rel="noopener noreferrer"
                className="btn-secondary inline-flex items-center gap-2 rounded-full border border-[#E5E7EB] bg-white px-6 py-3 font-semibold text-[#111827] text-sm"
              >
                <svg width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
                  <path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z" />
                </svg>
                View on GitHub
              </a>
              <a
                href="#waitlist"
                className="btn-primary inline-flex items-center gap-2 rounded-full px-6 py-3 font-semibold text-white text-sm"
              >
                Get Early Access
              </a>
            </div>
          </div>
        </section>

        {/* Final CTA */}
        <section className="border-t border-[#E5E7EB] py-20 lg:py-28">
          <div className="text-center max-w-xl mx-auto">
            <h2 className="text-3xl lg:text-4xl font-bold mb-4">
              Ready to take back <span className="gradient-text">control</span>?
            </h2>
            <p className="text-[#6B7280] mb-8">
              Join the waitlist and be the first to know when AddrPass launches.
            </p>
            <WaitlistForm />
          </div>
        </section>
      </main>

      <Footer />
    </div>
  );
}
