import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "AddrPass — Your Address, Your Control",
  description:
    "Open-source address management platform. Share your address via tokenized links and QR codes with full access control, monitoring, and expiration.",
  keywords: [
    "address privacy",
    "address sharing",
    "tokenized address",
    "QR code delivery",
    "contact privacy",
    "address vault",
    "open source",
    "GDPR",
    "data privacy",
  ],
  authors: [{ name: "AddrPass" }],
  openGraph: {
    title: "AddrPass — Your Address, Your Control",
    description:
      "Share your address via tokenized links and QR codes. Monitor access. Revoke anytime. Open source.",
    type: "website",
    url: "https://addrpass.com",
  },
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body className="bg-[#FAFBFC] text-[#111827] antialiased">
        {children}
      </body>
    </html>
  );
}
