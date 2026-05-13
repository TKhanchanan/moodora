import { Card } from "./card";

export function LoadingState({ label = "Loading Moodora..." }: { label?: string }) {
  return (
    <Card>
      <div className="h-3 w-32 animate-pulse rounded-full bg-lilac/70" />
      <div className="mt-4 h-20 animate-pulse rounded-2xl bg-mist" />
      <p className="mt-4 text-sm text-ink/60">{label}</p>
    </Card>
  );
}

export function ErrorState({ message }: { message: string }) {
  return (
    <Card className="border-peach/60 bg-peach/20">
      <p className="text-sm font-semibold">Connection feels quiet</p>
      <p className="mt-2 text-sm text-ink/70">{message}</p>
    </Card>
  );
}

export function EmptyState({ message }: { message: string }) {
  return (
    <Card>
      <p className="text-sm text-ink/65">{message}</p>
    </Card>
  );
}
