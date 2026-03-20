import "./globals.css";
import type { Metadata } from "next";

export const metadata: Metadata = {
  title: "eth-pulse",
  description: "Realtime Ethereum whale movement feed",
};

export default function RootLayout({ children }: { children: React.ReactNode }) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
