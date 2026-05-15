"use client";

import { useMemo, useState } from "react";
import { useMutation, useQuery } from "@tanstack/react-query";
import { useRouter } from "next/navigation";
import { api } from "@/lib/api/client";
import type { CreateTarotReadingRequest, TarotCard } from "@/lib/api/types";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { EmptyState, ErrorState, LoadingState } from "@/components/ui/state";

const topicOptions: Array<{
  code: CreateTarotReadingRequest["topic"];
  label: string;
  icon: string;
}> = [
  { code: "general", label: "ภาพรวม", icon: "✦" },
  { code: "love", label: "ความรัก", icon: "♡" },
  { code: "career", label: "การงาน", icon: "⌁" },
  { code: "money", label: "การเงิน", icon: "◈" }
];

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
  const [shuffleNonce, setShuffleNonce] = useState(() => Math.floor(Math.random() * 1_000_000));
  const [selectedSourceCodes, setSelectedSourceCodes] = useState<string[]>([]);
  const [deckFanned, setDeckFanned] = useState(true);
  const [shuffleStep, setShuffleStep] = useState<"idle" | "gathering" | "animating">("idle");
  const isShuffling = shuffleStep !== "idle";

  const selectedSpread = useMemo(
    () => spreads.data?.spreads.find((spread) => spread.code === form.spreadCode),
    [form.spreadCode, spreads.data]
  );
  const requiredCardCount = selectedSpread?.cardCount ?? 0;

  const createReading = useMutation({
    mutationFn: api.createTarotReading,
    onSuccess: (reading) => router.push(`/tarot/result/${reading.id}`)
  });

  const cardList = useMemo(() => cards.data?.cards ?? [], [cards.data]);
  const shuffledCards = useMemo(() => shuffleCards(cardList, `${form.spreadCode}:${shuffleNonce}`), [cardList, form.spreadCode, shuffleNonce]);

  if (spreads.isLoading || cards.isLoading) return <LoadingState label="กำลังเตรียมโต๊ะพยากรณ์..." />;
  if (spreads.isError || cards.isError) return <ErrorState message="ไม่สามารถโหลดข้อมูลไพ่ทาโรต์ได้ กรุณาตรวจสอบการเชื่อมต่อระบบ" />;

  const spreadList = spreads.data?.spreads ?? [];
  if (!spreadList.length) return <EmptyState message="ยังไม่มีรูปแบบการเปิดไพ่ในระบบ" />;
  const selectionComplete = selectedSourceCodes.length === requiredCardCount;

  return (
    <div className="space-y-5">
      <Card>
        <p className="text-sm text-ink/55">ไพ่ทาโรต์</p>
        <h1 className="mt-3 text-3xl font-semibold">เลือกไพ่ทำนายดวงชะตา</h1>
        <p className="mt-3 max-w-2xl text-sm leading-6 text-ink/70">
          ไพ่ถูกสับและวางเรียงไว้บนโต๊ะแล้ว ตั้งจิตอธิษฐานแล้วเลือกไพ่ด้วยตัวคุณเอง เมื่อครบจำนวนแล้วกดดูคำทำนายได้เลย
        </p>
      </Card>

      <Card>
        <div className="grid gap-5">
          <div>
            <p className="text-sm font-semibold text-ink/70">เรื่องที่อยากถาม</p>
            <div className="mt-3 grid grid-cols-2 gap-3 sm:grid-cols-4">
              {topicOptions.map((topic) => {
                const selected = form.topic === topic.code;
                return (
                  <button
                    key={topic.code}
                    type="button"
                    className={[
                      "min-h-28 rounded-2xl border px-4 py-4 text-left transition",
                      selected ? "border-ink bg-ink text-white shadow-soft" : "border-white bg-white/80 text-ink hover:bg-white"
                    ].join(" ")}
                    onClick={() => setForm((current) => ({ ...current, topic: topic.code }))}
                  >
                    <span className="grid h-10 w-10 place-items-center rounded-md bg-mist text-xl text-ink">{topic.icon}</span>
                    <span className="mt-4 block text-sm font-semibold">{topic.label}</span>
                  </button>
                );
              })}
            </div>
          </div>

          <div>
            <p className="text-sm font-semibold text-ink/70">จำนวนไพ่</p>
            <div className="mt-3 grid grid-cols-1 gap-3 sm:grid-cols-3">
              {spreadList.map((spread) => {
                const selected = form.spreadCode === spread.code;
                return (
                  <button
                    key={spread.code}
                    type="button"
                    className={[
                      "min-h-28 rounded-2xl border px-4 py-4 text-left transition",
                      selected ? "border-ink bg-ink text-white shadow-soft" : "border-white bg-white/80 text-ink hover:bg-white"
                    ].join(" ")}
                    onClick={() => {
                      setForm((current) => ({ ...current, spreadCode: spread.code }));
                      setSelectedSourceCodes([]);
                      setShuffleNonce((current) => current + 1);
                      setDeckFanned(true);
                      setShuffleStep("idle");
                    }}
                  >
                    <span className="grid h-10 w-10 place-items-center rounded-md bg-mist text-sm font-bold text-ink">
                      {spread.cardCount}
                    </span>
                    <span className="mt-4 block text-sm font-semibold">{spread.cardCount} ใบ</span>
                    <span className="mt-1 block text-xs text-current/65">{spreadLabel(spread.code, spread.name)}</span>
                  </button>
                );
              })}
            </div>
          </div>

          {createReading.isError && (
            <p className="rounded-2xl bg-peach/30 px-4 py-3 text-sm text-ink/70">
              ไม่สามารถสร้างคำทำนายได้ กรุณาลองใหม่อีกครั้ง
            </p>
          )}
        </div>
      </Card>

      <Card className="overflow-hidden">
        <div className="flex flex-col gap-3 sm:flex-row sm:items-end sm:justify-between">
          <div>
            <h2 className="text-lg font-semibold">เลือกไพ่ {requiredCardCount} ใบ</h2>
            <p className="mt-2 text-sm text-ink/60">เลือกแล้ว {selectedSourceCodes.length}/{requiredCardCount} ใบ</p>
          </div>
        </div>

        <div className="relative mt-6 h-[460px] sm:h-[720px] w-full [--radius:160px] sm:[--radius:360px] [--shuffle-dist:55px] sm:[--shuffle-dist:90px]">
          <div className="absolute left-1/2 top-0 h-[110px] w-[66px] -translate-x-1/2 sm:h-[200px] sm:w-[120px]">
            {shuffledCards.map((card, index) => {
              const selectedIndex = selectedSourceCodes.indexOf(card.sourceCode);
              const isSelected = selectedIndex >= 0;
              const isDisabled = !isSelected && selectedSourceCodes.length >= requiredCardCount;

              const reverseOffset = (shuffledCards.length - 1) - index;
              const isAnimatingCard = shuffleStep === "animating" && reverseOffset < 5;

              return (
                <button
                  key={card.sourceCode}
                  type="button"
                  aria-label={isSelected ? `Selected card ${selectedIndex + 1}` : "Select hidden tarot card"}
                  disabled={isDisabled || createReading.isPending || isShuffling}
                  onClick={() => {
                    setSelectedSourceCodes((current) => {
                      if (current.includes(card.sourceCode)) {
                        return current.filter((sourceCode) => sourceCode !== card.sourceCode);
                      }
                      if (current.length >= requiredCardCount) {
                        return current;
                      }
                      return [...current, card.sourceCode];
                    });
                  }}
                  className={[
                    "absolute left-0 top-0 h-[110px] w-[66px] rounded-md bg-[#050505] bg-cover bg-center shadow-[0_0_5px_rgba(0,0,0,0.5)] focus:outline-none focus:ring-2 focus:ring-ink disabled:cursor-not-allowed disabled:opacity-45 sm:h-[200px] sm:w-[120px]",
                    isAnimatingCard ? "" : "transition-transform duration-[350ms] ease-in-out hover:-translate-y-2"
                  ].join(" ")}
                  style={{
                    backgroundImage: "url('/tarot/card-back.png')",
                    transform: cardStackTransform(index, shuffledCards.length, deckFanned, isSelected),
                    transformOrigin: "center var(--radius)",
                    zIndex: isAnimatingCard ? undefined : (isSelected ? 200 + selectedIndex : index + 1),
                    animation: isAnimatingCard ? `custom-tarot-shuffle 500ms linear ${reverseOffset * 120}ms forwards` : undefined
                  }}
                >
                  {isSelected && (
                    <span className="absolute -right-2 -top-2 grid h-7 w-7 place-items-center rounded-full bg-peach text-xs font-bold text-ink shadow-sm">
                      {selectedIndex + 1}
                    </span>
                  )}
                </button>
              );
            })}
          </div>

          <div className="absolute left-1/2 top-[380px] sm:top-[360px] z-[300] flex w-full max-w-[280px] -translate-x-1/2 -translate-y-1/2 flex-col items-center gap-3">
            <div className="flex gap-2.5">
              <button
                type="button"
                className="flex h-10 px-5 cursor-pointer items-center justify-center rounded-full bg-black text-sm font-semibold text-white shadow-md transition hover:bg-ink active:scale-95 disabled:opacity-50"
                onClick={() => setDeckFanned((current) => !current)}
                disabled={createReading.isPending || isShuffling}
              >
                {deckFanned ? "เก็บไพ่" : "กรีดไพ่"}
              </button>

              <button
                type="button"
                className="flex h-10 px-5 cursor-pointer items-center justify-center rounded-full bg-white text-sm font-semibold text-ink shadow-md transition hover:bg-mist active:scale-95 disabled:opacity-50"
                onClick={() => {
                  if (createReading.isPending || isShuffling) return;
                  setShuffleStep("gathering");
                  setDeckFanned(false);
                  setSelectedSourceCodes([]);
                  setTimeout(() => {
                    setShuffleStep("animating");
                    setTimeout(() => {
                      setShuffleNonce((current) => current + 1);
                      setShuffleStep("idle");
                      setTimeout(() => {
                        setDeckFanned(true);
                      }, 50);
                    }, 1000);
                  }, 350);
                }}
                disabled={createReading.isPending || isShuffling}
              >
                สับไพ่ใหม่
              </button>
            </div>

            {selectedSourceCodes.length > 0 && (
              <button
                type="button"
                className="flex h-8 px-4 cursor-pointer items-center justify-center rounded-full bg-peach text-xs font-semibold text-ink shadow-sm transition hover:bg-peach/80 active:scale-95 disabled:opacity-50"
                onClick={() => setSelectedSourceCodes([])}
                disabled={createReading.isPending || isShuffling}
              >
                ล้างไพ่ที่เลือก ({selectedSourceCodes.length})
              </button>
            )}
          </div>
        </div>

        <div className="mt-4 flex justify-end">
          <Button
            disabled={!selectionComplete || createReading.isPending || isShuffling}
            onClick={() => createReading.mutate({ ...form, selectedCardSourceCodes: selectedSourceCodes })}
            className="w-full sm:w-auto"
          >
            {createReading.isPending ? "Opening..." : "ดูคำทำนาย"}
          </Button>
        </div>
      </Card>

      <style>{`
        @keyframes custom-tarot-shuffle {
          0% {
            transform: translate(0px, 0px) rotate(0deg);
            z-index: 400;
          }
          45% {
            transform: translate(var(--shuffle-dist), -8px) rotate(4deg);
            z-index: 400;
          }
          50% {
            transform: translate(var(--shuffle-dist), -8px) rotate(4deg);
            z-index: -1;
          }
          100% {
            transform: translate(0px, 0px) rotate(0deg);
            z-index: -1;
          }
        }
      `}</style>
    </div>
  );
}

function shuffleCards(cards: TarotCard[], salt: string) {
  const shuffled = [...cards];
  const nextRandom = seededRandom(hashString(salt));
  for (let i = shuffled.length - 1; i > 0; i -= 1) {
    const j = Math.floor(nextRandom() * (i + 1));
    [shuffled[i], shuffled[j]] = [shuffled[j], shuffled[i]];
  }
  return shuffled;
}

function hashString(value: string) {
  let hash = 2166136261;
  for (let i = 0; i < value.length; i += 1) {
    hash ^= value.charCodeAt(i);
    hash = Math.imul(hash, 16777619);
  }
  return hash >>> 0;
}

function seededRandom(seed: number) {
  let value = seed || 1;
  return () => {
    value = Math.imul(value ^ (value >>> 15), 1 | value);
    value ^= value + Math.imul(value ^ (value >>> 7), 61 | value);
    return ((value ^ (value >>> 14)) >>> 0) / 4294967296;
  };
}

function spreadLabel(code: string, fallback: string) {
  switch (code) {
    case "one_card":
      return "เร็ว กระชับ";
    case "three_cards":
      return "อดีต ปัจจุบัน อนาคต";
    case "celtic_cross":
      return "อ่านลึกแบบเต็ม";
    default:
      return fallback;
  }
}

function cardStackTransform(index: number, total: number, fanned: boolean, selected: boolean) {
  if (total <= 1) {
    return selected ? "translateY(-24px)" : "rotate(0deg)";
  }

  if (!fanned) {
    const offset = (index % 7) * 0.35;
    return `translate(${offset}px, ${offset}px) rotate(0deg)`;
  }

  const middleOffset = -1 * Math.trunc(total / 2);
  const angle = (360 / total) * (middleOffset + index);
  const lift = selected ? -34 : 0;
  return `rotate(${angle}deg) translateY(${lift}px)`;
}
