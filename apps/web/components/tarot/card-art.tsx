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
        className="aspect-[3/5] w-full rounded-md object-cover shadow-sm"
      />
    );
  }

  return (
    <div className="grid aspect-[3/5] w-full place-items-center rounded-md bg-gradient-to-br from-lilac via-mist to-sky p-4 text-center shadow-sm">
      <div>
        <p className="text-xs uppercase tracking-[0.25em] text-ink/45">{sourceCode}</p>
        <p className="mt-3 text-sm font-semibold text-ink">{name}</p>
      </div>
    </div>
  );
}
