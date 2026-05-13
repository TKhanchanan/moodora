export type ApiErrorBody = {
  error?: {
    code?: string;
    message?: string;
  };
};

export type HealthResponse = {
  status: string;
  checks: Record<string, string>;
  timezone: string;
};

export type VersionResponse = {
  name: string;
  env: string;
  version: string;
  timezone: string;
};

export type TarotAsset = {
  id: string;
  deckCode: string;
  size: "thumb" | "medium" | "large";
  format: "webp" | "jpg";
  url: string;
  width: number;
  height: number;
  fileSize: number;
  isDefault: boolean;
};

export type TarotCard = {
  id: string;
  sourceCode: string;
  nameEn: string;
  type: "major" | "minor";
  suit: string | null;
  meaningUpEn: string;
  meaningRevEn: string;
  descriptionEn: string;
  assets: TarotAsset[];
};

export type TarotCardsResponse = {
  cards: TarotCard[];
};

export type TarotSpreadPosition = {
  positionNumber: number;
  code: string;
  name: string;
  description: string;
};

export type TarotSpread = {
  code: string;
  name: string;
  description: string;
  cardCount: number;
  positions: TarotSpreadPosition[];
};

export type TarotSpreadsResponse = {
  spreads: TarotSpread[];
};

export type CreateTarotReadingRequest = {
  spreadCode: string;
  topic: "general" | "love" | "career" | "money";
  language: "th" | "en";
  allowReversed: boolean;
  question?: string;
};

export type TarotReading = {
  id: string;
  spreadCode: string;
  topic: string;
  language: string;
  question: string;
  cards: Array<{
    positionNumber: number;
    positionCode: string;
    positionName: string;
    card: {
      sourceCode: string;
      name: string;
    };
    orientation: "upright" | "reversed";
    meaning: string;
    advice: string;
  }>;
  summary: string;
  createdAt: string;
};

export type LuckyColor = {
  id: string;
  code: string;
  nameTh: string;
  nameEn: string;
  hex: string;
  meaning: string;
};

export type LuckyThing = {
  id: string;
  code: string;
  nameTh: string;
  nameEn: string;
  category: string;
  tags: string[];
  description: string;
};

export type Avoidance = {
  id: string;
  code: string;
  category: string;
  textTh: string;
  textEn: string;
  moodTag: string;
};

export type DailyInsight = {
  date: string;
  timezone: string;
  luckyColors: LuckyColor[];
  avoidColors: LuckyColor[];
  luckyFoods: LuckyThing[];
  luckyItems: LuckyThing[];
  avoidances: Avoidance[];
  walletBalance?: number;
  checkInStatus?: {
    checkedIn: boolean;
    localDate: string;
  };
  dailyTarot?: {
    status: string;
    message: string;
  };
  message?: string;
};

export type Wallet = {
  id: string;
  userId: string;
  coinBalance: number;
  createdAt: string;
  updatedAt: string;
};

export type CoinTransaction = {
  id: string;
  transactionType: string;
  amount: number;
  balanceAfter: number;
  reason: string;
  createdAt: string;
};

export type CoinTransactionsResponse = {
  transactions: CoinTransaction[];
};

export type CheckInResponse = {
  checkedIn: boolean;
  rewardCoins: number;
  streakDay: number;
  walletBalance: number;
  timezone: string;
  nextCheckInAt: string;
  alreadyChecked: boolean;
};
