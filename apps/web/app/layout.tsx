import type { Metadata } from "next";

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
      <body>{children}</body>
    </html>
  );
}
