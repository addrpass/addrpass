"use client";

import Link from "next/link";
import Header from "../components/Header";
import Footer from "../components/Footer";

const plans = [
  {
    name: "Free",
    price: "$0",
    period: "forever",
    desc: "For individuals protecting their personal address.",
    cta: "Get started",
    ctaLink: "/register",
    highlight: false,
    features: [
      "3 addresses",
      "10 active shares",
      "QR code generation",
      "PIN protection",
      "Access monitoring",
      "Expiration & limits",
      "Community support",
    ],
  },
  {
    name: "Pro",
    price: "$9",
    period: "/month",
    desc: "For power users and freelancers who share addresses often.",
    cta: "Start free trial",
    ctaLink: "/register",
    highlight: true,
    features: [
      "Unlimited addresses",
      "Unlimited shares",
      "All Free features",
      "Webhook notifications",
      "API access",
      "Priority support",
      "Custom share branding",
    ],
  },
  {
    name: "Business",
    price: "$29",
    period: "/month",
    desc: "For companies integrating address privacy into their product.",
    cta: "Start free trial",
    ctaLink: "/register",
    highlight: false,
    features: [
      "Everything in Pro",
      "Business accounts",
      "API keys (OAuth2)",
      "Scoped access (full/delivery/zone)",
      "Delegation chains",
      "Shipping label API",
      "Team management",
      "Dedicated support",
    ],
  },
];

const faqs = [
  {
    q: "Is AddrPass really free for personal use?",
    a: "Yes. The Free plan includes 3 addresses and 10 active shares with full monitoring and QR codes. No credit card required, no time limit.",
  },
  {
    q: "Can I self-host instead of using the cloud?",
    a: "Absolutely. AddrPass is open source under AGPL-3.0. Clone the repo, run `docker compose up`, and you have the full platform on your own servers. Self-hosting is free and unlimited.",
  },
  {
    q: "What happens if I hit the free plan limits?",
    a: "You can upgrade to Pro anytime. Your existing addresses and shares stay intact. Or self-host for unlimited everything.",
  },
  {
    q: "How does pricing work for delivery companies?",
    a: "The Business plan includes API keys, delegation chains, and the shipping label API. For high-volume usage (100K+ resolutions/month), contact us for custom pricing.",
  },
  {
    q: "Is my data stored in the EU?",
    a: "Yes. Our cloud infrastructure is hosted in Germany. For self-hosted deployments, your data stays wherever you deploy it.",
  },
  {
    q: "Can I switch between plans?",
    a: "Yes. Upgrade or downgrade anytime. If you downgrade, you keep your data but new shares are limited to the lower plan.",
  },
];

function CheckIcon() {
  return (
    <svg width="16" height="16" viewBox="0 0 16 16" fill="none" className="flex-shrink-0 mt-0.5">
      <circle cx="8" cy="8" r="8" fill="#22D3EE" opacity="0.1" />
      <path d="M5 8l2.5 2.5L11 6" stroke="#22D3EE" strokeWidth="1.3" strokeLinecap="round" strokeLinejoin="round" />
    </svg>
  );
}

export default function PricingPage() {
  return (
    <div className="min-h-screen bg-[#FAFBFD]">
      <Header />

      <main className="mx-auto max-w-5xl px-6">
        {/* Header */}
        <section className="py-16 lg:py-20 text-center">
          <p className="text-xs font-semibold uppercase tracking-widest text-[#22D3EE] mb-3">Pricing</p>
          <h1 className="text-4xl lg:text-5xl font-bold tracking-tight">
            Simple, transparent pricing.
          </h1>
          <p className="mt-4 text-[#64748B] max-w-lg mx-auto leading-relaxed">
            Free for personal use. Paid plans for teams and businesses. Self-host for free if you want full control.
          </p>
        </section>

        {/* Plans */}
        <section className="grid md:grid-cols-3 gap-5 pb-20">
          {plans.map((plan) => (
            <div
              key={plan.name}
              className={`rounded-2xl p-7 flex flex-col ${
                plan.highlight
                  ? "bg-[#0F172A] text-white ring-2 ring-[#22D3EE]/30"
                  : "bg-white border border-[#E2E8F0]"
              }`}
            >
              <div>
                <h3 className={`font-semibold text-lg ${plan.highlight ? "text-white" : "text-[#0F172A]"}`}>{plan.name}</h3>
                <div className="mt-3 flex items-baseline gap-1">
                  <span className={`text-4xl font-bold ${plan.highlight ? "text-white" : "text-[#0F172A]"}`}>{plan.price}</span>
                  <span className={`text-sm ${plan.highlight ? "text-[#94A3B8]" : "text-[#64748B]"}`}>{plan.period}</span>
                </div>
                <p className={`mt-3 text-sm leading-relaxed ${plan.highlight ? "text-[#94A3B8]" : "text-[#64748B]"}`}>{plan.desc}</p>
              </div>

              <div className="mt-6 flex-1">
                <ul className="space-y-2.5">
                  {plan.features.map((f) => (
                    <li key={f} className="flex items-start gap-2.5">
                      {plan.highlight ? (
                        <svg width="16" height="16" viewBox="0 0 16 16" fill="none" className="flex-shrink-0 mt-0.5">
                          <circle cx="8" cy="8" r="8" fill="#22D3EE" opacity="0.2" />
                          <path d="M5 8l2.5 2.5L11 6" stroke="#22D3EE" strokeWidth="1.3" strokeLinecap="round" strokeLinejoin="round" />
                        </svg>
                      ) : (
                        <CheckIcon />
                      )}
                      <span className={`text-sm ${plan.highlight ? "text-[#CBD5E1]" : "text-[#334155]"}`}>{f}</span>
                    </li>
                  ))}
                </ul>
              </div>

              <Link
                href={plan.ctaLink}
                className={`mt-8 block text-center rounded-full py-2.5 text-sm font-semibold transition-all ${
                  plan.highlight
                    ? "bg-[#22D3EE] text-[#0F172A] hover:bg-[#06B6D4]"
                    : "border border-[#E2E8F0] text-[#0F172A] hover:border-[#CBD5E1] hover:bg-[#F8FAFC]"
                }`}
              >
                {plan.cta}
              </Link>
            </div>
          ))}
        </section>

        <div className="section-divider" />

        {/* Self-host callout */}
        <section className="py-16 text-center">
          <div className="inline-flex items-center gap-2 pill text-[#64748B] text-xs mb-4">
            <svg width="14" height="14" viewBox="0 0 24 24" fill="currentColor"><path d="M12 0c-6.626 0-12 5.373-12 12 0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23.957-.266 1.983-.399 3.003-.404 1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576 4.765-1.589 8.199-6.086 8.199-11.386 0-6.627-5.373-12-12-12z" /></svg>
            Open source
          </div>
          <h3 className="text-2xl font-bold tracking-tight mb-3">Prefer to self-host?</h3>
          <p className="text-[#64748B] max-w-md mx-auto mb-6 leading-relaxed text-sm">
            The entire AddrPass platform is open source under AGPL-3.0. Run it on your own infrastructure with Docker. No limits, no cost, full control.
          </p>
          <div className="inline-flex items-center gap-2 bg-[#0F172A] text-[#94A3B8] rounded-xl px-5 py-3 font-mono text-sm">
            <span className="text-[#64748B]">$</span>
            <span className="text-[#E2E8F0]">docker compose up</span>
          </div>
          <div className="mt-4">
            <a href="https://github.com/addrpass/addrpass" target="_blank" rel="noopener noreferrer" className="text-sm text-[#22D3EE] hover:underline font-medium">
              View on GitHub &rarr;
            </a>
          </div>
        </section>

        <div className="section-divider" />

        {/* FAQ */}
        <section className="py-16 pb-24">
          <h3 className="text-2xl font-bold tracking-tight text-center mb-12">Frequently asked questions</h3>
          <div className="grid md:grid-cols-2 gap-x-12 gap-y-8 max-w-3xl mx-auto">
            {faqs.map((faq) => (
              <div key={faq.q}>
                <h4 className="font-semibold text-[14px] text-[#0F172A] mb-2">{faq.q}</h4>
                <p className="text-[13px] text-[#64748B] leading-relaxed">{faq.a}</p>
              </div>
            ))}
          </div>
        </section>
      </main>

      <Footer />
    </div>
  );
}
