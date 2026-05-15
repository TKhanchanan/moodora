import { apiBaseUrl } from "@/lib/config/env";
import type {
  ApiErrorBody,
  CheckInResponse,
  CoinTransactionsResponse,
  CreateTarotReadingRequest,
  DailyInsight,
  HealthResponse,
  TarotCardsResponse,
  TarotReading,
  TarotSpreadsResponse,
  VersionResponse,
  Wallet
} from "./types";

export class ApiError extends Error {
  readonly status: number;
  readonly code: string;

  constructor(message: string, status: number, code = "api_error") {
    super(message);
    this.status = status;
    this.code = code;
  }
}

async function apiFetch<T>(path: string, init?: RequestInit): Promise<T> {
  const response = await fetch(`${apiBaseUrl}${path}`, {
    ...init,
    headers: {
      "Content-Type": "application/json",
      ...init?.headers
    }
  });

  if (!response.ok) {
    let body: ApiErrorBody | undefined;
    try {
      body = (await response.json()) as ApiErrorBody;
    } catch {
      body = undefined;
    }
    throw new ApiError(
      body?.error?.message ?? `API request failed with status ${response.status}`,
      response.status,
      body?.error?.code
    );
  }

  return (await response.json()) as T;
}

export const api = {
  health: () => apiFetch<HealthResponse>("/health"),
  version: () => apiFetch<VersionResponse>("/api/v1/version"),
  dailyInsight: () => apiFetch<DailyInsight>("/api/v1/daily-insights/today"),
  tarotCards: () => apiFetch<TarotCardsResponse>("/api/v1/tarot/cards"),
  tarotSpreads: () => apiFetch<TarotSpreadsResponse>("/api/v1/tarot/spreads"),
  tarotReading: (id: string) => apiFetch<TarotReading>(`/api/v1/tarot/readings/${id}`),
  createTarotReading: (body: CreateTarotReadingRequest) =>
    apiFetch<TarotReading>("/api/v1/tarot/readings", {
      method: "POST",
      body: JSON.stringify(body)
    }),
  wallet: () => apiFetch<Wallet>("/api/v1/wallet"),
  coinTransactions: () => apiFetch<CoinTransactionsResponse>("/api/v1/coin-transactions"),
  checkIn: () =>
    apiFetch<CheckInResponse>("/api/v1/check-ins", {
      method: "POST"
    })
};
