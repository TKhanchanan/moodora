"use client";

import { useMemo, useState } from "react";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { api } from "@/lib/api/client";
import type { CreateTarotReadingRequest } from "@/lib/api/types";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { EmptyState, ErrorState, LoadingState } from "@/components/ui/state";
import { CardArt } from "./card-art";

const topics = ["general", "love", "career", "money"] as const;

export function TarotSetup() {
  const router = useRouter();
  const spreads = useQuery({ queryKey: ["tarot-spreads"], queryFn: api.tarotSpreads });
  const cards = useQuery({ queryKey: ["tarot-cards"], queryFn: api.tarotCards });
  const [form, setForm] = useState<CreateTarotReadingRequest>({
    spreadCode: "three_cards",
    topic: "general",
    language: "th",
    allowReversed: true,
    question: ""
  });

  const selectedSpread = useMemo(
    () => spreads.data?.spreads.find((spread) => spread.code === form.spreadCode),
    [form.spreadCode, spreads.data]
  );

  const createReading = useMutation({
    mutationFn: api.createTarotReading,
    onSuccess: (reading) => router.push(`/tarot/result/${reading.id}`)
  });

  if (spreads.isLoading || cards.isLoading) return <LoadingState label="Loading Tarot table..." />;
  if (spreads.isError || cards.isError) return <ErrorState message="Tarot data is not available. Run backend migrations, imports, and seeds." />;

  const cardList = cards.data?.cards ?? [];
  const spreadList = spreads.data?.spreads ?? [];
  if (!spreadList.length) return <EmptyState message="No Tarot spreads are available yet." />;

  return (
    <div className="space-y-5">
      <Card>
        <p className="text-sm text-ink/55">Tarot</p>
        <h1 className="mt-3 text-3xl font-semibold">Choose a spread for reflection.</h1>
        <p className="mt-3 max-w-2xl text-sm leading-6 text-ink/70">
          Readings are stored by the backend and use seeded interpretations when available.
        </p>
      </Card>

      <div className="grid gap-4 lg:grid-cols-[1fr_0.9fr]">
        <Card>
          <h2 className="text-lg font-semibold">Reading setup</h2>
          <div className="mt-4 space-y-4">
            <label className="block text-sm font-medium">
              Spread
              <select
                className="mt-2 w-full rounded-2xl border border-white bg-white px-4 py-3"
                value={form.spreadCode}
                onChange={(event) => setForm((current) => ({ ...current, spreadCode: event.target.value }))}
              >
                {spreadList.map((spread) => (
                  <option key={spread.code} value={spread.code}>
                    {spread.name} · {spread.cardCount} cards
                  </option>
                ))}
              </select>
            </label>

            <label className="block text-sm font-medium">
              Topic
              <select
                className="mt-2 w-full rounded-2xl border border-white bg-white px-4 py-3"
                value={form.topic}
                onChange={(event) => setForm((current) => ({ ...current, topic: event.target.value as CreateTarotReadingRequest["topic"] }))}
              >
                {topics.map((topic) => (
                  <option key={topic} value={topic}>{topic}</option>
                ))}
              </select>
            </label>

            <label className="block text-sm font-medium">
              Language
              <select
                className="mt-2 w-full rounded-2xl border border-white bg-white px-4 py-3"
                value={form.language}
                onChange={(event) => setForm((current) => ({ ...current, language: event.target.value as "th" | "en" }))}
              >
                <option value="th">Thai</option>
                <option value="en">English</option>
              </select>
            </label>

            <label className="flex items-center justify-between rounded-2xl bg-mist px-4 py-3 text-sm font-medium">
              Allow reversed cards
              <input
                type="checkbox"
                checked={form.allowReversed}
                onChange={(event) => setForm((current) => ({ ...current, allowReversed: event.target.checked }))}
                className="h-5 w-5"
              />
            </label>

            <label className="block text-sm font-medium">
              Question
              <textarea
                className="mt-2 min-h-24 w-full rounded-2xl border border-white bg-white px-4 py-3"
                value={form.question}
                onChange={(event) => setForm((current) => ({ ...current, question: event.target.value }))}
                placeholder="วันนี้ควรโฟกัสเรื่องอะไร"
              />
            </label>

            {createReading.isError && (
              <p className="rounded-2xl bg-peach/30 px-4 py-3 text-sm text-ink/70">
                Could not create reading. Check that Tarot cards and interpretations are seeded.
              </p>
            )}

            <Button
              disabled={createReading.isPending}
              onClick={() => createReading.mutate(form)}
              className="w-full"
            >
              {createReading.isPending ? "Creating..." : "Create reading"}
            </Button>
          </div>
        </Card>

        <Card>
          <h2 className="text-lg font-semibold">{selectedSpread?.name ?? "Spread"}</h2>
          <p className="mt-2 text-sm leading-6 text-ink/65">{selectedSpread?.description}</p>
          <div className="mt-4 space-y-2">
            {selectedSpread?.positions.map((position) => (
              <div key={position.code} className="rounded-2xl bg-mist px-4 py-3">
                <p className="text-sm font-semibold">{position.positionNumber}. {position.name}</p>
                <p className="text-xs text-ink/55">{position.code}</p>
              </div>
            ))}
          </div>
        </Card>
      </div>

      <Card>
        <h2 className="text-lg font-semibold">Card library</h2>
        <div className="mt-4 grid grid-cols-2 gap-3 sm:grid-cols-4 lg:grid-cols-6">
          {cardList.slice(0, 24).map((card) => (
            <div key={card.sourceCode}>
              <CardArt sourceCode={card.sourceCode} name={card.nameEn} assets={card.assets} />
              <p className="mt-2 text-xs font-semibold">{card.nameEn}</p>
            </div>
          ))}
        </div>
      </Card>
    </div>
  );
}
