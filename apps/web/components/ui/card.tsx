import type { ReactNode } from "react";
import { cn } from "@/lib/utils/cn";

export function Card({
  children,
  className
}: {
  children: ReactNode;
  className?: string;
}) {
  return (
    <section className={cn("rounded-3xl border border-white/75 bg-white/72 p-5 shadow-soft backdrop-blur", className)}>
      {children}
    </section>
  );
}
