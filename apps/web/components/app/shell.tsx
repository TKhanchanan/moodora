import Link from "next/link";
import type { ReactNode } from "react";
import { QueryProvider } from "@/lib/query/provider";
import { ServiceWorker } from "./service-worker";

const nav = [
  { href: "/daily", label: "Daily" },
  { href: "/tarot", label: "Tarot" },
  { href: "/wallet", label: "Wallet" },
  { href: "/profile", label: "Profile" }
];

export function AppShell({ children }: { children: ReactNode }) {
  return (
    <QueryProvider>
      <ServiceWorker />
      <div className="mx-auto flex min-h-screen w-full max-w-5xl flex-col px-4 pb-24 pt-4 sm:px-6">
        <header className="sticky top-3 z-20 mb-5 rounded-3xl border border-white/70 bg-white/75 px-4 py-3 shadow-soft backdrop-blur">
          <div className="flex items-center justify-between gap-3">
            <Link href="/" className="flex items-center gap-2">
              <span className="grid h-9 w-9 place-items-center rounded-2xl bg-lilac text-lg">☾</span>
              <span className="text-base font-semibold tracking-wide">Moodora</span>
            </Link>
            <nav className="flex gap-1 overflow-x-auto text-sm">
              {nav.map((item) => (
                <Link
                  key={item.href}
                  href={item.href}
                  className="rounded-full px-3 py-2 text-ink/70 transition hover:bg-mist hover:text-ink"
                >
                  {item.label}
                </Link>
              ))}
            </nav>
          </div>
        </header>
        <main className="flex-1">{children}</main>
      </div>
    </QueryProvider>
  );
}
