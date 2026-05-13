# Moodora Initial Project Plan

Moodora starts as a monorepo with the backend implemented first. The backend should be a Go modular monolith backed by PostgreSQL, Redis, and MinIO for local development.

This plan intentionally avoids application code. It defines the first structure, module boundaries, database direction, and development milestones.

## Project Structure

```text
moodora/
  apps/
    api/
      cmd/
        api/
      internal/
        platform/
        modules/
      migrations/
      tests/
    web/
    admin/
  docs/
    architecture/
      initial-plan.md
    api/
  infra/
    local/
  docker-compose.yml
  .env.example
  README.md
  AGENTS.md
```

### Root

- `AGENTS.md`: project rules and product direction.
- `README.md`: project overview, setup commands, and local development notes.
- `.env.example`: documented local environment variables only. Secrets must not be committed.
- `docker-compose.yml`: local PostgreSQL, Redis, and MinIO dependencies.

### `apps/api`

Go backend. This is the first implementation target and should be designed as a modular monolith.

Recommended shape:

- `cmd/api`: application entrypoint and process wiring.
- `internal/platform`: shared infrastructure adapters such as config, database, Redis, storage, logging, HTTP server, middleware, clock, and transactions.
- `internal/modules`: business modules with clear boundaries.
- `migrations`: PostgreSQL schema migrations.
- `tests`: integration test helpers and cross-module tests when needed.

Business logic belongs in services/usecases inside each module, not in HTTP handlers.

### `apps/web`

Future Next.js App Router PWA frontend using TypeScript, Tailwind CSS, and shadcn/ui. It should be mobile-first, minimal, soft, modern, and suitable for Gen Z users.

### `apps/admin`

Future admin dashboard for internal operations, content review, user support, wallet inspection, API key management, and operational reporting.

### `docs`

Architecture notes, API contracts, product decisions, and development plans.

### `infra`

Deployment and infrastructure configuration. Local infrastructure can start with Docker Compose at the root and later move environment-specific files under `infra`.

## Backend Modules

The backend should start with a small set of modules and expand only as behavior becomes real.

### Platform Layer

Shared technical capabilities used by modules:

- Configuration loading with `Asia/Bangkok` as the default timezone.
- PostgreSQL connection pool.
- Redis client for cache, rate limiting, and temporary state.
- MinIO/S3-compatible object storage client.
- HTTP router and middleware.
- Request validation and error responses.
- Logging.
- Transaction manager for database writes.
- Time provider for testable daily logic.
- Authentication primitives.

Platform code must not contain product business rules.

### Identity Module

Owns user identity and account state.

Initial responsibilities:

- User records.
- Sign-in identity model.
- User profile basics.
- Authentication session or token lifecycle.

Later responsibilities:

- Social login.
- Account deletion.
- User privacy controls.

### Wallet Module

Owns coin balance and ledger behavior.

Initial responsibilities:

- Coin account per user.
- Immutable coin ledger entries.
- Balance changes inside database transactions.
- Idempotency keys for coin operations.
- Wallet read model for current balance.

Rules:

- All coin operations must be transactional.
- Balance must be derived safely from ledger behavior or maintained with strict transactional updates.
- Handlers must call usecases; they must not mutate wallet state directly.

### Check-In Module

Owns daily check-in behavior.

Initial responsibilities:

- Daily check-in state by user and Bangkok-local date.
- Reward calculation through wallet usecases.
- Duplicate check-in prevention.
- Streak tracking if required by product.

Rules:

- Use `Asia/Bangkok` day boundaries.
- Reward grants must use wallet transactions.

### Insights Module

Owns daily user-facing insight generation and retrieval.

Initial responsibilities:

- Daily insight snapshot per user or per segment.
- Lucky colors.
- Lucky foods/items.
- Safe presentation language for interpretation-based guidance.

Rules:

- Store snapshots for user-facing results.
- Keep calculated facts separate from interpretation text.
- Do not claim scientific proof for fortune-telling.

### Tarot Module

Owns tarot readings.

Initial responsibilities:

- Reading request.
- Card draw result.
- Interpretation snapshot.
- Optional coin cost through wallet usecases.

Rules:

- Store the full user-facing reading snapshot.
- Keep card data, draw mechanics, and interpretation output distinct.
- Present readings as reflection or entertainment guidance.

### Astronomy Module

Owns calculated moon and astronomy data.

Initial responsibilities:

- Moon phase data.
- Moon-related daily context.
- External astronomy provider integration if needed.
- Cached calculations in Redis where appropriate.

Rules:

- Astronomy and moon data can be labeled as calculated data.
- External API responses should be normalized before reaching product modules.

### Assets Module

Owns file/object storage metadata.

Initial responsibilities:

- MinIO object metadata.
- Public asset URL mapping.
- Content-type and size metadata.

Later responsibilities:

- CDN migration.
- Upload policies.
- Admin media management.

### Public API Module

Future module for external API and white-label support.

Initial planning responsibilities:

- API clients.
- Hashed API keys only.
- Rate limit policy.
- Usage tracking.

Rules:

- Never store raw API keys.
- Design for tenant-aware access later, without adding multi-tenant complexity too early.

## Database Plan

PostgreSQL is the source of truth. Redis is for cache, rate limit, and temporary state. MinIO is for local object storage.

### Core Database Principles

- Use migrations from the start.
- Use UUID or ULID-style public identifiers for external references.
- Use `created_at` and `updated_at` timestamps on mutable tables.
- Use `Asia/Bangkok` for product day calculations, while storing timestamps consistently.
- Use transactions for wallet, coin, and check-in reward flows.
- Prefer explicit constraints over application-only validation.
- Store snapshots for user-facing readings, reports, and daily insights.

### Initial Table Groups

Identity:

- `users`
- `user_profiles`
- `auth_identities`
- `auth_sessions` or token/session equivalent

Wallet:

- `wallet_accounts`
- `wallet_ledger_entries`
- `wallet_idempotency_keys`

Check-in:

- `daily_checkins`
- `checkin_reward_rules`

Insights:

- `daily_insight_snapshots`
- `lucky_item_snapshots`

Tarot:

- `tarot_cards`
- `tarot_readings`
- `tarot_reading_cards`
- `tarot_reading_snapshots`

Astronomy:

- `moon_daily_snapshots`
- `astronomy_provider_cache` if database-backed cache is needed beyond Redis

Assets:

- `assets`

Future public API:

- `api_clients`
- `api_key_hashes`
- `api_usage_events`
- `api_rate_limit_policies`

### Snapshot Strategy

Snapshots should preserve exactly what the user saw at the time of generation.

Use snapshots for:

- Tarot reading interpretations.
- Daily insight text.
- Lucky colors, foods, and items.
- Moon/astronomy summaries shown in user-facing reports.

Recommended snapshot fields:

- Stable snapshot ID.
- User ID where applicable.
- Source calculation or generation version.
- Input parameters needed for audit/debugging.
- Render-ready result payload.
- Created timestamp.

### Redis Plan

Use Redis for:

- Rate limiting.
- Short-lived sessions or verification state if needed.
- Cached astronomy provider responses.
- Temporary generation locks.
- Idempotency locks around expensive or duplicate-prone actions.

Redis must not be the only source of truth for wallet balances, coin grants, paid readings, or check-in history.

### MinIO Plan

Use MinIO locally as an S3-compatible store for:

- Static tarot card images.
- Generated or curated visual assets.
- Public content assets.
- Future user-uploaded assets if product scope requires it.

Object metadata should be recorded in PostgreSQL when the application needs to reference, audit, or manage the asset.

## First Development Milestones

### Milestone 1: Repository Foundation

- Create the monorepo directories: `apps/api`, `apps/web`, `apps/admin`, `docs`, and `infra`.
- Keep `docker-compose.yml` for PostgreSQL, Redis, and MinIO local dependencies.
- Document local setup in `README.md`.
- Add backend environment configuration documentation.

Acceptance criteria:

- A developer can understand the intended repository layout.
- Local dependencies can be started with Docker Compose.
- No application behavior is implemented yet.

### Milestone 2: Backend Skeleton

- Initialize the Go module under `apps/api`.
- Add the API process entrypoint.
- Add configuration loading with default timezone set to `Asia/Bangkok`.
- Add database, Redis, and MinIO connection wiring.
- Add health endpoints for process, database, Redis, and storage checks.

Acceptance criteria:

- The API starts locally.
- Health checks confirm local dependencies.
- No business logic exists in handlers.

### Milestone 3: Database Migrations

- Add migration tooling.
- Create initial identity and wallet schema.
- Add constraints and indexes for user identity, wallet accounts, ledger entries, and idempotency.
- Add migration documentation.

Acceptance criteria:

- Migrations run cleanly against local PostgreSQL.
- Wallet tables support transactional ledger behavior.
- Schema choices are documented.

### Milestone 4: Wallet Core

- Implement wallet usecases.
- Add transactional coin grant and spend flows.
- Add idempotency handling.
- Add unit tests for core wallet behavior.

Acceptance criteria:

- Coin balance cannot be changed outside wallet usecases.
- Duplicate operations with the same idempotency key do not double-spend or double-grant.
- Tests cover successful grants, spends, insufficient balance, and duplicate requests.

### Milestone 5: Daily Check-In

- Implement Bangkok-local daily check-in.
- Grant check-in rewards through wallet usecases.
- Prevent duplicate check-ins for the same user and local date.
- Add tests for timezone boundaries and duplicate prevention.

Acceptance criteria:

- A user can check in once per Bangkok-local date.
- Reward grants are transactional.
- Tests prove date-boundary behavior.

### Milestone 6: First User-Facing Snapshot

- Implement a minimal daily insight or tarot reading flow.
- Store the user-facing result as a snapshot.
- Separate calculated data from interpretation text.
- Add tests for snapshot persistence.

Acceptance criteria:

- The API can create and retrieve a stored user-facing snapshot.
- Result language treats tarot/lucky guidance as interpretation or self-reflection.
- The system does not claim scientific proof for fortune-telling.

### Milestone 7: Frontend Preparation

- Initialize `apps/web` with Next.js App Router, TypeScript, Tailwind CSS, and shadcn/ui.
- Add PWA-ready foundations.
- Create mobile-first screens only after backend contracts are stable enough for the first workflow.

Acceptance criteria:

- Frontend setup follows the project rules.
- The first screen is tied to a real backend workflow.
- UI remains minimal, soft, modern, and mobile-first.

