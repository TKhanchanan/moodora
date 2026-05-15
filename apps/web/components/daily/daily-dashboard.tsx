"use client";

import Link from "next/link";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import { Card } from "@/components/ui/card";
import { EmptyState, ErrorState, LoadingState } from "@/components/ui/state";

export function DailyDashboard() {
  const query = useQuery({
    queryKey: ["daily-insight"],
    queryFn: api.dailyInsight
  });

  if (query.isLoading) return <LoadingState label="Loading today's guide..." />;
  if (query.isError) {
    return <ErrorState message="The API is offline or not seeded yet. Start apps/api and try again." />;
  }

  const insight = query.data;
  if (!insight) return <EmptyState message="No daily insight is available yet." />;

  return (
    <div className="space-y-5">
      <Card className="bg-white/80">
        <p className="text-sm text-ink/55">{insight.date} · {insight.timezone}</p>
        <h1 className="mt-3 text-3xl font-semibold leading-tight sm:text-4xl">Today feels softer with a little intention.</h1>
        <p className="mt-3 max-w-2xl text-sm leading-6 text-ink/70">
          Moodora offers lifestyle prompts for reflection, not guaranteed predictions.
        </p>
        <div className="mt-5 flex flex-wrap gap-3">
          <Link className="rounded-full bg-ink px-5 py-3 text-sm font-semibold text-white" href="/tarot">
            Start Tarot
          </Link>
          <Link className="rounded-full bg-white px-5 py-3 text-sm font-semibold text-ink shadow-sm" href="/wallet">
            Check in
          </Link>
        </div>
      </Card>

      <div className="grid gap-4 sm:grid-cols-2">
        <Card>
          <p className="text-sm text-ink/55">Wallet</p>
          <p className="mt-2 text-3xl font-semibold">{insight.walletBalance ?? "—"} coins</p>
          <p className="mt-2 text-sm text-ink/60">
            {insight.checkInStatus?.checkedIn ? "Checked in today." : "Check-in is available if your dev user is configured."}
          </p>
        </Card>
        <Card>
          <p className="text-sm text-ink/55">Daily Tarot</p>
          <p className="mt-2 text-base font-semibold">{insight.dailyTarot?.status ?? "Not connected"}</p>
          <p className="mt-2 text-sm leading-6 text-ink/65">{insight.dailyTarot?.message ?? "Create a Tarot reading when you feel ready."}</p>
        </Card>
      </div>

      <Card>
        <h2 className="text-lg font-semibold">Lucky colors</h2>
        <div className="mt-4 grid gap-3 sm:grid-cols-2">
          {insight.luckyColors?.map((color) => (
            <div key={color.code} className="rounded-2xl bg-mist p-4">
              <div className="flex items-center gap-3">
                <span className="h-10 w-10 rounded-full border border-white shadow-sm" style={{ backgroundColor: color.hex }} />
                <div>
                  <p className="font-semibold">{color.nameTh}</p>
                  <p className="text-xs text-ink/55">{color.nameEn}</p>
                </div>
              </div>
              <p className="mt-3 text-sm leading-6 text-ink/65">{color.meaning}</p>
            </div>
          ))}
        </div>
      </Card>

      <div className="grid gap-4 lg:grid-cols-3">
        <Card>
          <h2 className="text-lg font-semibold">Lucky foods</h2>
          <List values={insight.luckyFoods?.map((item) => `${item.nameTh} · ${item.description}`)} />
        </Card>
        <Card>
          <h2 className="text-lg font-semibold">Lucky items</h2>
          <List values={insight.luckyItems?.map((item) => `${item.nameTh} · ${item.description}`)} />
        </Card>
        <Card>
          <h2 className="text-lg font-semibold">Avoid today</h2>
          <List values={insight.avoidances?.map((item) => item.textTh)} />
        </Card>
      </div>

      {insight.avoidColors?.length > 0 && (
        <Card>
          <h2 className="text-lg font-semibold">Colors to keep light</h2>
          <div className="mt-3 flex flex-wrap gap-2">
            {insight.avoidColors.map((color) => (
              <span key={color.code} className="rounded-full bg-white px-3 py-2 text-sm text-ink/70">
                {color.nameTh}
              </span>
            ))}
          </div>
        </Card>
      )}
    </div>
  );
}

function List({ values }: { values?: string[] }) {
  if (!values?.length) return <p className="mt-3 text-sm text-ink/55">No recommendation yet.</p>;
  return (
    <ul className="mt-3 space-y-3">
      {values.map((value) => (
        <li key={value} className="rounded-2xl bg-mist px-4 py-3 text-sm leading-6 text-ink/70">
          {value}
        </li>
      ))}
    </ul>
  );
}
