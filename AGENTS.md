# Moodora Project Guidelines

Moodora is a cosmic lifestyle / fortune platform.

## Product Direction

Moodora provides:

- Tarot readings
- Lucky colors
- Lucky foods/items
- Moon and astronomy insights
- Daily insights
- Coin wallet
- Daily check-in
- Future public API and white-label support

## Architecture

Use a monorepo structure:

- apps/api: Go backend
- apps/web: Next.js PWA frontend
- apps/admin: admin dashboard
- docs: architecture and API docs
- infra: deployment and infrastructure config

Backend must be designed as a modular monolith first.

## Backend Rules

- Use Go.
- Use PostgreSQL as the main database.
- Use Redis for cache, rate limit, and temporary state.
- Keep business logic in services/usecases, not handlers.
- Use clear module boundaries.
- Use transactions for wallet and coin operations.
- Store snapshots for user-facing reading/report results.
- Never store raw API keys, only hashed keys.
- Never hardcode secrets.
- Use Asia/Bangkok as the default timezone.

## Frontend Rules

- Use Next.js App Router.
- Use TypeScript.
- Use Tailwind CSS and shadcn/ui.
- Build mobile-first.
- Moodora frontend must be PWA-ready.
- Keep UI minimal, soft, modern, and Gen Z friendly.

## Coding Style

- Prefer simple, readable code.
- Avoid unnecessary abstraction.
- Use meaningful names.
- Add tests for core business logic.
- Do not change unrelated files.
- Keep commits small and focused.

## Safety / Trust

- Separate calculated data from interpretation.
- Astronomy/moon data can be shown as calculated.
- Tarot/lucky guidance must be presented as interpretation or self-reflection.
- Do not claim scientific proof for fortune-telling.
