"use client";

import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { api } from "@/lib/api/client";
import { Button } from "@/components/ui/button";
import { Card } from "@/components/ui/card";
import { ErrorState, LoadingState } from "@/components/ui/state";

export function WalletPanel() {
  const queryClient = useQueryClient();
  const wallet = useQuery({ queryKey: ["wallet"], queryFn: api.wallet });
  const transactions = useQuery({ queryKey: ["coin-transactions"], queryFn: api.coinTransactions });
  const checkIn = useMutation({
    mutationFn: api.checkIn,
    onSuccess: () => {
      void queryClient.invalidateQueries({ queryKey: ["wallet"] });
      void queryClient.invalidateQueries({ queryKey: ["coin-transactions"] });
      void queryClient.invalidateQueries({ queryKey: ["daily-insight"] });
    }
  });

  if (wallet.isLoading) return <LoadingState label="Loading wallet..." />;
  if (wallet.isError) {
    return <ErrorState message="Wallet needs apps/api running with DEV_USER_ID configured and seeded." />;
  }

  const result = checkIn.data;

  return (
    <div className="space-y-5">
      <Card>
        <p className="text-sm text-ink/55">Wallet</p>
        <h1 className="mt-3 text-4xl font-semibold">{wallet.data?.coinBalance ?? 0} coins</h1>
        <p className="mt-3 text-sm leading-6 text-ink/70">
          Coins are an in-app balance for Moodora features. They are not cash or a payment instrument.
        </p>
        <Button
          className="mt-5 w-full sm:w-auto"
          disabled={checkIn.isPending}
          onClick={() => checkIn.mutate()}
        >
          {checkIn.isPending ? "Checking in..." : "Daily check-in"}
        </Button>
      </Card>

      {result && (
        <Card className={result.alreadyChecked ? "bg-sky/30" : "bg-mint/35"}>
          <p className="text-sm font-semibold">
            {result.alreadyChecked ? "Already checked in today" : "Check-in complete"}
          </p>
          <div className="mt-3 grid gap-3 sm:grid-cols-3">
            <Metric label="Reward" value={`${result.rewardCoins} coins`} />
            <Metric label="Streak" value={`Day ${result.streakDay}`} />
            <Metric label="Balance" value={`${result.walletBalance} coins`} />
          </div>
          <p className="mt-3 text-sm text-ink/60">Timezone: {result.timezone}</p>
        </Card>
      )}

      {checkIn.isError && <ErrorState message="Could not check in. Try again after confirming the API is running." />}

      <Card>
        <h2 className="text-lg font-semibold">Recent coin transactions</h2>
        {transactions.isLoading && <p className="mt-3 text-sm text-ink/55">Loading transactions...</p>}
        {transactions.isError && <p className="mt-3 text-sm text-ink/55">Transactions are unavailable.</p>}
        <div className="mt-4 space-y-3">
          {transactions.data?.transactions?.length ? (
            transactions.data.transactions.map((item) => (
              <div key={item.id} className="flex items-center justify-between rounded-2xl bg-mist px-4 py-3">
                <div>
                  <p className="text-sm font-semibold">{item.reason}</p>
                  <p className="text-xs text-ink/50">{item.transactionType}</p>
                </div>
                <div className="text-right">
                  <p className="text-sm font-semibold">{item.amount > 0 ? "+" : ""}{item.amount}</p>
                  <p className="text-xs text-ink/50">{item.balanceAfter} balance</p>
                </div>
              </div>
            ))
          ) : (
            <p className="text-sm text-ink/55">No transactions yet.</p>
          )}
        </div>
      </Card>
    </div>
  );
}

function Metric({ label, value }: { label: string; value: string }) {
  return (
    <div className="rounded-2xl bg-white/75 px-4 py-3">
      <p className="text-xs text-ink/50">{label}</p>
      <p className="mt-1 text-lg font-semibold">{value}</p>
    </div>
  );
}
