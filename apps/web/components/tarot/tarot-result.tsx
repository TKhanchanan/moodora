"use client";

import Link from "next/link";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import { Card } from "@/components/ui/card";
import { ErrorState, LoadingState } from "@/components/ui/state";
import { CardArt } from "./card-art";

export function TarotResult({ id }: { id: string }) {
  const reading = useQuery({
    queryKey: ["tarot-reading", id],
    queryFn: () => api.tarotReading(id)
  });

  if (reading.isLoading) return <LoadingState label="Opening your reading..." />;
  if (reading.isError || !reading.data) return <ErrorState message="Could not load this reading from the backend." />;

  return (
    <div className="space-y-5">
      <Card>
        <p className="text-sm text-ink/55">{reading.data.spreadCode} · {reading.data.topic}</p>
        <h1 className="mt-3 text-3xl font-semibold">Your Tarot reflection</h1>
        {reading.data.question && <p className="mt-3 text-sm leading-6 text-ink/70">{reading.data.question}</p>}
        <p className="mt-4 rounded-2xl bg-mist px-4 py-3 text-sm leading-6 text-ink/75">{reading.data.summary}</p>
      </Card>

      <div className="grid gap-4">
        {reading.data.cards
          .slice()
          .sort((a, b) => a.positionNumber - b.positionNumber)
          .map((card) => (
            <Card key={`${card.positionNumber}-${card.card.sourceCode}`}>
              <div className="grid gap-4 sm:grid-cols-[140px_1fr]">
                <CardArt sourceCode={card.card.sourceCode} name={card.card.name} />
                <div>
                  <p className="text-sm text-ink/55">{card.positionNumber}. {card.positionName} · {card.orientation}</p>
                  <h2 className="mt-2 text-xl font-semibold">{card.card.name}</h2>
                  <p className="mt-3 text-sm leading-6 text-ink/70">{card.meaning}</p>
                  <p className="mt-3 rounded-2xl bg-mint/40 px-4 py-3 text-sm leading-6 text-ink/70">{card.advice}</p>
                </div>
              </div>
            </Card>
          ))}
      </div>

      <div className="flex flex-wrap gap-3">
        <Link className="rounded-full bg-ink px-5 py-3 text-sm font-semibold text-white" href="/tarot">Back to Tarot</Link>
        <Link className="rounded-full bg-white px-5 py-3 text-sm font-semibold text-ink shadow-sm" href="/daily">Daily</Link>
        <Link className="rounded-full bg-white px-5 py-3 text-sm font-semibold text-ink shadow-sm" href="/wallet">Wallet</Link>
      </div>
    </div>
  );
}
