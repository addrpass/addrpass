export default function Footer() {
  return (
    <footer className="border-t border-[#E5E7EB] bg-white">
      <div className="mx-auto max-w-6xl px-6 py-12">
        <div className="flex flex-col sm:flex-row items-center justify-between gap-4">
          <div className="flex items-center gap-2">
            <div className="w-6 h-6 rounded-md bg-gradient-to-br from-teal-500 to-teal-700 flex items-center justify-center">
              <svg
                width="12"
                height="12"
                viewBox="0 0 24 24"
                fill="none"
                stroke="white"
                strokeWidth="2.5"
                strokeLinecap="round"
                strokeLinejoin="round"
              >
                <rect x="3" y="11" width="18" height="11" rx="2" ry="2" />
                <path d="M7 11V7a5 5 0 0 1 10 0v4" />
              </svg>
            </div>
            <span className="text-sm font-semibold text-[#111827]">AddrPass</span>
          </div>
          <div className="flex items-center gap-6 text-sm text-[#6B7280]">
            <a
              href="https://github.com/addrpass/addrpass"
              target="_blank"
              rel="noopener noreferrer"
              className="hover:text-[#0D9488] transition-colors"
            >
              GitHub
            </a>
            <a
              href="https://github.com/addrpass/addrpass/blob/main/apps/api/openapi.yaml"
              target="_blank"
              rel="noopener noreferrer"
              className="hover:text-[#0D9488] transition-colors"
            >
              API Docs
            </a>
            <span>AGPL-3.0</span>
          </div>
        </div>
        <p className="mt-6 text-center text-xs text-[#9CA3AF]">
          Your address, your control. Open source and self-hostable.
        </p>
      </div>
    </footer>
  );
}
