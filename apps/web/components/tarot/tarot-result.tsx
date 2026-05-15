"use client";

import Link from "next/link";
import { useQuery } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import { Card } from "@/components/ui/card";
import { ErrorState, LoadingState } from "@/components/ui/state";
import { CardArt } from "./card-art";
import { TarotVisualSpread } from "./tarot-visual-spread";

const topicMap: Record<string, string> = {
  general: "ภาพรวม",
  love: "ความรัก",
  career: "การงาน",
  money: "การเงิน"
};

const spreadMap: Record<string, string> = {
  one_card: "ไพ่ 1 ใบ",
  three_cards: "ไพ่ 3 ใบ",
  celtic_cross: "เซลติกครอส (10 ใบ)"
};

export function TarotResult({ id }: { id: string }) {
  const reading = useQuery({
    queryKey: ["tarot-reading", id],
    queryFn: () => api.tarotReading(id)
  });

  if (reading.isLoading) return <LoadingState label="กำลังเปิดไพ่ทำนาย..." />;
  if (reading.isError || !reading.data) return <ErrorState message="ไม่สามารถโหลดข้อมูลคำทำนายได้ กรุณาลองใหม่อีกครั้ง" />;

  const spreadLabel = spreadMap[reading.data.spreadCode] || reading.data.spreadCode;
  const topicLabel = topicMap[reading.data.topic] || reading.data.topic;

  return (
    <div className="space-y-5">
      <Card>
        <p className="text-sm font-semibold text-ink/60">{spreadLabel} · {topicLabel}</p>
        <h1 className="mt-3 text-3xl font-semibold">ผลการทำนายของคุณ</h1>
        <p className="mt-4 text-sm leading-7 text-ink/80">
          {reading.data.summary.replace("Reflect on this pattern:", "ภาพรวมคำทำนายของคุณ:")}
        </p>
      </Card>

      <TarotVisualSpread reading={reading.data} />

      <div className="grid gap-4">
        {reading.data.cards
          .slice()
          .sort((a, b) => a.positionNumber - b.positionNumber)
          .map((card) => {
            const isReversed = card.orientation === "reversed";
            const orientationLabel = isReversed ? "ไพ่กลับหัว (Reversed)" : "ไพ่หัวตั้ง (Upright)";

            return (
              <Card key={`${card.positionNumber}-${card.card.sourceCode}`}>
                <div className="grid gap-5 sm:grid-cols-[140px_1fr]">
                  <CardArt sourceCode={card.card.sourceCode} name={card.card.name} assets={card.card.assets} />
                  <div>
                    <p className="text-sm font-semibold text-ink/60">ใบที่ {card.positionNumber} · {card.positionName} · {orientationLabel}</p>
                    <h2 className="mt-2 text-xl font-bold">{card.card.nameTh || card.card.name} {`(${card.card.nameEn})`}</h2>
                    <p className="mt-3 text-sm leading-6 text-ink/75">
                      <strong className="text-ink">ลักษณะไพ่:</strong> {card.card.characteristic}
                    </p>
                    {card.card.description && (
                      <p className="mt-3 text-sm leading-7 text-ink/70">{card.card.description}</p>
                    )}
                    <p className="mt-3 text-sm leading-7 text-ink/80"><strong className="text-ink">ความหมาย:</strong> {card.meaning}</p>
                    <p className="mt-3 rounded-md bg-mint/30 px-3 py-2 text-sm leading-7 text-ink/80"><strong className="text-ink">คำแนะนำ:</strong> {card.advice}</p>
                  </div>
                </div>
              </Card>
            );
          })}
      </div>

      <div className="flex flex-wrap gap-3">
        <Link className="rounded-full bg-ink px-5 py-3 text-sm font-semibold text-white" href="/tarot">เปิดไพ่ใหม่</Link>
        <Link className="rounded-full bg-white px-5 py-3 text-sm font-semibold text-ink shadow-sm" href="/daily">ดวงรายวัน</Link>
        <Link className="rounded-full bg-white px-5 py-3 text-sm font-semibold text-ink shadow-sm" href="/wallet">กระเป๋าเหรียญ</Link>
      </div>
    </div>
  );
}
