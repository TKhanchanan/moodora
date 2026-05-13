import Link from "next/link";
import { Card } from "@/components/ui/card";

export default function OfflinePage() {
  return (
    <Card>
      <p className="text-sm text-ink/55">Offline</p>
      <h1 className="mt-3 text-3xl font-semibold">Moodora is in quiet mode.</h1>
      <p className="mt-3 max-w-xl text-sm leading-6 text-ink/70">
        You can browse the cached shell, but creating readings, checking in, and coin updates need a connection.
      </p>
      <Link className="mt-5 inline-flex rounded-full bg-ink px-5 py-3 text-sm font-semibold text-white" href="/">
        Back home
      </Link>
    </Card>
  );
}
