"use client";

import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import { Card } from "@/components/ui/card";

export function ApiStatus() {
  const health = useQuery({ queryKey: ["health"], queryFn: api.health });
  const version = useQuery({ queryKey: ["version"], queryFn: api.version });

  return (
    <Card className="mt-5">
      <div className="flex flex-wrap items-center justify-between gap-3">
        <div>
          <p className="text-sm font-semibold">API connection</p>
          <p className="mt-1 text-sm text-ink/60">
            {health.isError ? "API offline" : health.data?.status ?? "Checking..."}
          </p>
        </div>
        <div className="rounded-full bg-mist px-4 py-2 text-sm text-ink/70">
          {version.data ? `${version.data.env} · ${version.data.version}` : "Version pending"}
        </div>
      </div>
    </Card>
  );
}
