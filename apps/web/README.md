# Moodora Web

Next.js PWA frontend for Moodora.

## Setup

Install dependencies from the repository root:

```sh
pnpm install
```

Create local env:

```sh
cp apps/web/.env.example apps/web/.env.local
```

Local API dependency:

```env
NEXT_PUBLIC_API_BASE_URL=http://localhost:8080
```

Start the backend first, then run the web app:

```sh
pnpm --filter web dev
```

Open:

```text
http://localhost:3000
```

## Routes

- `/` and `/daily`: daily dashboard from `GET /api/v1/daily-insights/today`
- `/tarot`: Tarot cards, spreads, and reading creation
- `/tarot/result/[id]`: backend Tarot reading result
- `/wallet`: wallet, check-in, and coin transactions
- `/offline`: friendly offline shell page
- `/profile`: placeholder profile page until auth is implemented

## PWA Notes

The app includes:

- manifest metadata for Moodora
- placeholder SVG icons in `public/icons`
- `/offline` page
- a simple service worker at `public/sw.js`

The service worker is registered in production builds only. Offline mode is limited to a cached shell; creating Tarot readings, check-ins, and coin updates still require the API.

## API Handling

All API calls use `NEXT_PUBLIC_API_BASE_URL`. Do not hardcode localhost inside components.

If the API is unavailable, pages show soft error states instead of crashing.
