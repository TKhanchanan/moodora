import { CardArt } from "./card-art";
import type { TarotReading } from "@/lib/api/types";

export function TarotVisualSpread({ reading }: { reading: TarotReading }) {
  const cards = [...reading.cards].sort((a, b) => a.positionNumber - b.positionNumber);

  const getCard = (pos: number) => cards.find(c => c.positionNumber === pos);

  if (reading.spreadCode === "three_cards") {
    return (
      <div className="my-8 rounded-2xl bg-white/60 p-6 shadow-sm border border-white backdrop-blur-sm">
        <div className="flex justify-center gap-4 sm:gap-8">
          {[1, 2, 3].map((pos) => {
            const card = getCard(pos);
            if (!card) return null;
            const isReversed = card.orientation === "reversed";
            const label = pos === 1 ? "อดีต" : pos === 2 ? "ปัจจุบัน" : "อนาคต";
            return (
              <div key={pos} className="flex flex-col items-center">
                <div className={`w-[90px] sm:w-[130px] transition-transform ${isReversed ? 'rotate-180' : ''}`}>
                  <CardArt sourceCode={card.card.sourceCode} name={card.card.name} assets={card.card.assets} />
                </div>
                <div className="mt-4 flex flex-col items-center">
                  <span className="grid h-6 w-6 place-items-center rounded-full bg-peach text-xs font-bold text-ink mb-2">{pos}</span>
                  <span className="rounded-full bg-mist px-3 py-1 text-xs font-semibold text-ink">{label}</span>
                </div>
              </div>
            );
          })}
        </div>
      </div>
    );
  }

  if (reading.spreadCode === "celtic_cross") {
    const w = 64;
    const h = 106;
    const gap = 12;

    const renderCard = (pos: number, label: string, isCross: boolean = false) => {
      const card = getCard(pos);
      if (!card) return null;
      const isReversed = card.orientation === "reversed";
      
      let transform = "";
      if (isCross) {
        transform = isReversed ? "rotate(270deg)" : "rotate(90deg)";
      } else if (isReversed) {
        transform = "rotate(180deg)";
      }

      return (
        <div className="flex flex-col items-center group relative" style={{ width: w }}>
          <div 
            className="relative shadow-md rounded-md bg-white/50 backdrop-blur-sm" 
            style={{ width: w, height: h, transform, transition: "transform 0.3s" }}
          >
            <CardArt sourceCode={card.card.sourceCode} name={card.card.name} assets={card.card.assets} />
            <div className="absolute -top-2 -left-2 z-30 grid h-6 w-6 place-items-center rounded-full bg-peach text-xs font-bold text-ink shadow-sm">
              {pos}
            </div>
          </div>
          <div className="absolute -bottom-8 opacity-0 group-hover:opacity-100 transition-opacity z-50 whitespace-nowrap rounded-full bg-ink px-2 py-1 text-[10px] text-white pointer-events-none">
            {pos}. {label}
          </div>
        </div>
      );
    };

    return (
      <div className="my-8 rounded-2xl bg-white/30 p-4 sm:p-8 shadow-sm border border-white backdrop-blur-sm">
        <div className="flex flex-wrap sm:flex-nowrap justify-between items-center gap-6 w-full max-w-3xl mx-auto min-h-[460px]">
          {/* Column 1: Past */}
          <div className="flex h-full items-center justify-center">
            {renderCard(4, "อดีต")}
          </div>

          {/* Column 2: Central Cross */}
          <div className="relative flex items-center justify-center" style={{ width: w, height: 3 * h + 2 * gap }}>
            <div className="absolute top-0">{renderCard(5, "จุดมุ่งหมาย")}</div>
            <div className="absolute" style={{ top: h + gap }}>{renderCard(1, "ปัจจุบัน")}</div>
            <div className="absolute z-20" style={{ top: h + gap }}>{renderCard(2, "อุปสรรค", true)}</div>
            <div className="absolute bottom-0">{renderCard(3, "จิตใต้สำนึก")}</div>
          </div>

          {/* Column 3: Future */}
          <div className="flex h-full items-center justify-center">
            {renderCard(6, "อนาคต")}
          </div>

          {/* Column 4: Staff */}
          <div className="flex flex-col justify-between items-center" style={{ height: 4 * h + 1.5 * gap }}>
            {renderCard(10, "ภาพรวม")}
            {renderCard(9, "ความรู้สึกในใจ")}
            {renderCard(8, "สภาพแวดล้อม")}
            {renderCard(7, "ตัวตนภายใน")}
          </div>
        </div>
      </div>
    );
  }

  return null;
}
