import Image from "next/image";
import type { TarotAsset } from "@/lib/api/types";

export function CardArt({
  sourceCode,
  name,
  assets
}: {
  sourceCode: string;
  name: string;
  assets?: TarotAsset[];
}) {
  const asset = assets?.find((item) => item.isDefault) ?? assets?.[0];

  if (asset) {
    return (
      <Image
        src={asset.url}
        alt={name}
        width={asset.width}
        height={asset.height}
        unoptimized
        className="aspect-[2/3] w-full rounded-2xl border border-white/70 object-cover shadow-sm"
      />
    );
  }

  return (
    <div className="grid aspect-[2/3] w-full place-items-center rounded-2xl border border-white/70 bg-gradient-to-br from-lilac via-mist to-sky p-4 text-center shadow-sm">
      <div>
        <p className="text-xs uppercase tracking-[0.25em] text-ink/45">{sourceCode}</p>
        <p className="mt-3 text-sm font-semibold text-ink">{name}</p>
      </div>
    </div>
  );
}
