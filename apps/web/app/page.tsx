import { ApiStatus } from "@/components/app/api-status";
import { DailyDashboard } from "@/components/daily/daily-dashboard";

export default function HomePage() {
  return (
    <div>
      <DailyDashboard />
      <ApiStatus />
    </div>
  );
}
