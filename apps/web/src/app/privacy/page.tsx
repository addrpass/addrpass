import Header from "../components/Header";
import Footer from "../components/Footer";

export const metadata = {
  title: "Privacy Policy — AddrPass",
  description: "How AddrPass handles your data.",
};

export default function PrivacyPage() {
  return (
    <div className="min-h-screen bg-[#FAFBFD]">
      <Header />

      <main className="mx-auto max-w-3xl px-6 py-16 lg:py-20">
        <p className="text-xs font-semibold uppercase tracking-widest text-[#22D3EE] mb-3">Legal</p>
        <h1 className="text-4xl font-bold tracking-tight mb-4">Privacy Policy</h1>
        <p className="text-sm text-[#64748B] mb-12">Last updated: March 26, 2026</p>

        <div className="space-y-10 text-[15px] text-[#334155] leading-relaxed">
          <section>
            <h2 className="text-lg font-semibold text-[#0F172A] mb-3">Who we are</h2>
            <p>
              AddrPass is operated by Omelas (&ldquo;we&rdquo;, &ldquo;us&rdquo;, &ldquo;our&rdquo;). Our infrastructure
              is hosted in Germany (EU). We are committed to protecting your privacy and handling your data transparently.
            </p>
          </section>

          <section>
            <h2 className="text-lg font-semibold text-[#0F172A] mb-3">What data we collect</h2>
            <ul className="list-disc pl-5 space-y-2">
              <li><strong>Account information:</strong> Email address and hashed password when you register.</li>
              <li><strong>Addresses:</strong> The addresses you choose to store in your vault. These are stored encrypted at rest.</li>
              <li><strong>Access logs:</strong> IP address, user agent, and timestamp when someone resolves one of your share links. This is shown to you in your dashboard so you can monitor access.</li>
              <li><strong>Usage data:</strong> Monthly resolution counts for billing purposes (paid plans only).</li>
            </ul>
          </section>

          <section>
            <h2 className="text-lg font-semibold text-[#0F172A] mb-3">What data we do NOT collect</h2>
            <ul className="list-disc pl-5 space-y-2">
              <li>We do not use cookies for tracking or advertising.</li>
              <li>We do not use analytics services (no Google Analytics, no trackers).</li>
              <li>We do not sell, rent, or share your personal data with third parties.</li>
              <li>We do not use your data for advertising or profiling.</li>
            </ul>
          </section>

          <section>
            <h2 className="text-lg font-semibold text-[#0F172A] mb-3">How we use your data</h2>
            <ul className="list-disc pl-5 space-y-2">
              <li><strong>Authentication:</strong> Your email and password are used to sign you in.</li>
              <li><strong>Address sharing:</strong> When you create a share link, the recipient can view the address fields you authorized. You control the scope, expiration, and access limits.</li>
              <li><strong>Access monitoring:</strong> Access logs are collected so you can see who accessed your address and when.</li>
              <li><strong>Billing:</strong> Resolution counts are tracked to enforce plan limits and calculate usage-based billing.</li>
            </ul>
          </section>

          <section>
            <h2 className="text-lg font-semibold text-[#0F172A] mb-3">Browser extension</h2>
            <p>The AddrPass browser extension:</p>
            <ul className="list-disc pl-5 space-y-2 mt-2">
              <li>Stores your authentication token and cached addresses locally in your browser using <code className="bg-[#F1F5F9] px-1.5 py-0.5 rounded text-sm">chrome.storage.local</code>.</li>
              <li>Communicates only with <code className="bg-[#F1F5F9] px-1.5 py-0.5 rounded text-sm">api.addrpass.com</code> to authenticate, fetch addresses, and create share links.</li>
              <li>Detects address form fields on web pages using HTML attributes (autocomplete, name, id, placeholder). It does not read or transmit page content.</li>
              <li>Does not use remote code. All scripts are bundled locally in the extension package.</li>
              <li>Does not collect browsing history, analytics, or telemetry.</li>
            </ul>
          </section>

          <section>
            <h2 className="text-lg font-semibold text-[#0F172A] mb-3">Data storage and security</h2>
            <ul className="list-disc pl-5 space-y-2">
              <li>All data is stored on servers in Germany (EU).</li>
              <li>All connections use TLS encryption (HTTPS).</li>
              <li>Passwords are hashed with bcrypt before storage.</li>
              <li>Share tokens are generated using 144-bit CSPRNG and are not reversible.</li>
              <li>API keys are stored as bcrypt hashes; the plaintext is shown only once at creation.</li>
            </ul>
          </section>

          <section>
            <h2 className="text-lg font-semibold text-[#0F172A] mb-3">Your rights</h2>
            <p>Under GDPR and applicable data protection laws, you have the right to:</p>
            <ul className="list-disc pl-5 space-y-2 mt-2">
              <li><strong>Access:</strong> View all data we hold about you (visible in your dashboard).</li>
              <li><strong>Rectification:</strong> Update your addresses and account information at any time.</li>
              <li><strong>Deletion:</strong> Delete your account and all associated data.</li>
              <li><strong>Portability:</strong> Export your data via the API.</li>
              <li><strong>Revocation:</strong> Revoke any share link instantly, cutting off access.</li>
            </ul>
          </section>

          <section>
            <h2 className="text-lg font-semibold text-[#0F172A] mb-3">Self-hosted instances</h2>
            <p>
              If you self-host AddrPass, your data never touches our servers. You are responsible for
              your own data storage, backups, and compliance. This privacy policy applies only to the
              cloud service at addrpass.com.
            </p>
          </section>

          <section>
            <h2 className="text-lg font-semibold text-[#0F172A] mb-3">Third-party services</h2>
            <ul className="list-disc pl-5 space-y-2">
              <li><strong>Cloudflare:</strong> DNS and static site hosting for addrpass.com. Cloudflare may process request metadata (IP, headers) per their privacy policy.</li>
              <li><strong>Stripe:</strong> Payment processing for paid plans. We do not store credit card information; Stripe handles this directly.</li>
            </ul>
          </section>

          <section>
            <h2 className="text-lg font-semibold text-[#0F172A] mb-3">Changes to this policy</h2>
            <p>
              We may update this policy from time to time. Changes will be posted on this page with an
              updated date. For significant changes, we will notify registered users by email.
            </p>
          </section>

          <section>
            <h2 className="text-lg font-semibold text-[#0F172A] mb-3">Contact</h2>
            <p>
              For privacy questions or data requests, open an issue on{" "}
              <a href="https://github.com/addrpass/addrpass/issues" target="_blank" rel="noopener noreferrer" className="text-[#22D3EE] hover:underline">
                GitHub
              </a>{" "}
              or email us at privacy@addrpass.com.
            </p>
          </section>
        </div>
      </main>

      <Footer />
    </div>
  );
}
