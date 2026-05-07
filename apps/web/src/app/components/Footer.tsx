import Link from "next/link";

export default function Footer() {
  return (
    <footer className="bg-[#0F172A] text-white">
      <div className="mx-auto max-w-6xl px-6 py-16">
        <div className="grid md:grid-cols-4 gap-10">
          {/* Brand */}
          <div className="md:col-span-2">
            <div className="flex items-center gap-2.5 mb-4">
              <svg width="24" height="24" viewBox="0 0 1024 1024" fill="none">
                <rect x="0" y="0" width="1024" height="1024" rx="224" ry="224" fill="#F0EFE8" />
                <rect x="232" y="232" width="560" height="560" rx="80" fill="none" stroke="#1A1A1A" strokeWidth="48" />
                <circle cx="512" cy="430" r="76" fill="#1A1A1A" />
                <path d="M472 480L446 668L578 668L552 480Z" fill="#1A1A1A" />
                <circle cx="512" cy="430" r="28" fill="#E7B300" />
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
          <span>&copy; {new Date().getFullYear()} <a href="https://omelas.tech" target="_blank" rel="noopener noreferrer" className="hover:text-white transition-colors">Omelas</a>. Built in the EU.</span>
          <span className="flex items-center gap-3">
            <Link href="/privacy" className="hover:text-white transition-colors">Privacy Policy</Link>
            <span>Open source under AGPL-3.0</span>
          </span>
        </div>
      </div>
    </footer>
  );
}
