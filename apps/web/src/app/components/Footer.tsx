import Link from "next/link";

export default function Footer() {
  return (
    <footer className="bg-[#0F172A] text-white">
      <div className="mx-auto max-w-6xl px-6 py-16">
        <div className="grid md:grid-cols-4 gap-10">
          {/* Brand */}
          <div className="md:col-span-2">
            <div className="flex items-center gap-2.5 mb-4">
              <svg width="24" height="24" viewBox="0 0 32 32" fill="none">
                <path d="M16 2L4 8v8c0 8.4 5.12 16.24 12 18 6.88-1.76 12-9.6 12-18V8L16 2z" fill="#1E293B" />
                <circle cx="16" cy="13" r="3" fill="#22D3EE" />
                <path d="M14 15.5L13 22h6l-1-6.5" fill="#22D3EE" opacity="0.7" />
              </svg>
              <span className="text-base font-bold">AddrPass</span>
            </div>
            <p className="text-sm text-[#94A3B8] max-w-sm leading-relaxed">
              Your address, your control. Open-source address management with tokenized sharing, access monitoring, and instant revocation.
            </p>
          </div>

          {/* Links */}
          <div>
            <h4 className="text-xs font-semibold uppercase tracking-wider text-[#64748B] mb-4">Product</h4>
            <div className="space-y-2.5 text-sm text-[#94A3B8]">
              <a href="#features" className="block hover:text-white transition-colors">Features</a>
              <a href="#delivery" className="block hover:text-white transition-colors">For Delivery</a>
              <Link href="/register" className="block hover:text-white transition-colors">Get Started</Link>
              <Link href="/login" className="block hover:text-white transition-colors">Sign In</Link>
            </div>
          </div>

          <div>
            <h4 className="text-xs font-semibold uppercase tracking-wider text-[#64748B] mb-4">Developers</h4>
            <div className="space-y-2.5 text-sm text-[#94A3B8]">
              <a href="https://www.npmjs.com/package/@addrpass/sdk" target="_blank" rel="noopener noreferrer" className="block hover:text-white transition-colors">npm SDK</a>
              <a href="https://github.com/addrpass/addrpass" target="_blank" rel="noopener noreferrer" className="block hover:text-white transition-colors">GitHub</a>
              <a href="https://github.com/addrpass/addrpass/blob/main/apps/api/openapi.yaml" target="_blank" rel="noopener noreferrer" className="block hover:text-white transition-colors">API Reference</a>
              <a href="https://github.com/addrpass/addrpass/blob/main/ARCHITECTURE.md" target="_blank" rel="noopener noreferrer" className="block hover:text-white transition-colors">Architecture</a>
            </div>
          </div>
        </div>

        <div className="section-divider mt-12 mb-6 opacity-20" />

        <div className="flex flex-col sm:flex-row items-center justify-between gap-4 text-xs text-[#64748B]">
          <span>Built in the EU. Your data never has to leave your servers.</span>
          <span>Open source under AGPL-3.0</span>
        </div>
      </div>
    </footer>
  );
}
