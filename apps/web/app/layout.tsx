import type { Metadata } from "next";

import { Providers } from "@/app/providers";

import "./globals.css";

export const metadata: Metadata = {
  title: "YLX - The Marketplace for Developers",
  description:
    "A developer-first marketplace for buying and selling tech gear.",
};

export default function RootLayout({
  children,
}: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="en">
      <body>
        <a
          href="#main-content"
          className="absolute top-2 left-2 z-50 -translate-y-20 rounded-md bg-primary px-3 py-2 text-sm text-primary-foreground transition-transform focus:translate-y-0"
        >
          Skip to content
        </a>
        <Providers>{children}</Providers>
      </body>
    </html>
  );
}
