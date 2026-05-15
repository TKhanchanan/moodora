import { Card } from "@/components/ui/card";

export default function ProfilePage() {
  return (
    <Card>
      <p className="text-sm text-ink/55">Profile</p>
      <h1 className="mt-3 text-3xl font-semibold">Your Moodora space</h1>
      <p className="mt-3 text-sm leading-6 text-ink/70">
        Authentication and profile editing are not implemented yet. Wallet, check-in, and daily insights use the backend development user when configured.
      </p>
    </Card>
  );
}
