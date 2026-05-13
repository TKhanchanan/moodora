import { TarotResult } from "@/components/tarot/tarot-result";

export default async function TarotResultPage({
  params
}: {
  params: Promise<{ id: string }>;
}) {
  const { id } = await params;
  return <TarotResult id={id} />;
}
