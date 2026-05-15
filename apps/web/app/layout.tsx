import type { Metadata, Viewport } from "next";
import "./globals.css";
import { AppShell } from "@/components/app/shell";

export const metadata: Metadata = {
  title: "Moodora",
  description: "Cosmic lifestyle guidance for daily self-reflection.",
  applicationName: "Moodora",
  appleWebApp: {
    capable: true,
    title: "Moodora",
    statusBarStyle: "default"
  }
};

export const viewport: Viewport = {
  themeColor: "#A78BFA",
  width: "device-width",
  initialScale: 1
};

export default function RootLayout({ children }: Readonly<{ children: React.ReactNode }>) {
  return (
    <html lang="en">
      <body>
        <AppShell>{children}</AppShell>
      </body>
    </html>
  );
}
