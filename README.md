# Moodora

Moodora is a cosmic lifestyle and fortune platform built as a monorepo.

## Local Dependencies

Start PostgreSQL, Redis, and MinIO:

```sh
docker compose up -d postgres redis minio
```

## Backend Migrations

The Go backend stores PostgreSQL migrations in `apps/api/migrations`.

Install the migration CLI:

```sh
brew install golang-migrate
```

Run migrations locally:

```sh
migrate -path apps/api/migrations \
  -database "postgres://moodora:moodora@localhost:5432/moodora_db?sslmode=disable" \
  up
```

Rollback the latest migration:

```sh
migrate -path apps/api/migrations \
  -database "postgres://moodora:moodora@localhost:5432/moodora_db?sslmode=disable" \
  down 1
```

Check migration version:

```sh
migrate -path apps/api/migrations \
  -database "postgres://moodora:moodora@localhost:5432/moodora_db?sslmode=disable" \
  version
```

Future deployment should run the same migration files using the environment's secret-provided `DATABASE_URL`.

## Backend Seeds

Create a local development user for wallet and check-in endpoints:

```sh
psql "postgres://moodora:moodora@localhost:5432/moodora_db?sslmode=disable" \
  -v ON_ERROR_STOP=1 \
  -v dev_user_id="00000000-0000-0000-0000-000000000001" \
  -f apps/api/seeds/dev_user.sql
```

Import the built-in 78-card Tarot source data into `tarot_cards`:

```sh
cd apps/api
DATABASE_URL="postgres://moodora:moodora@localhost:5432/moodora_db?sslmode=disable" \
go run ./cmd/import-tarot-cards
```

The importer is idempotent and uses upsert behavior. It does not download images or fetch tarotapi.dev directly. `apps/api/seeds/tarot/cards/tarotapi.dev.sample.json` documents the expected JSON shape for a reviewed 78-card source file.

Run Tarot reference seeds after card import:

```sh
psql "postgres://moodora:moodora@localhost:5432/moodora_db?sslmode=disable" \
  -v ON_ERROR_STOP=1 \
  -f apps/api/seeds/tarot/001_spreads.sql
```

Run sample Tarot translation and interpretation seeds:

```sh
psql "postgres://moodora:moodora@localhost:5432/moodora_db?sslmode=disable" \
  -v ON_ERROR_STOP=1 \
  -f apps/api/seeds/tarot/003_interpretations_sample.sql
```

Run Lucky Lifestyle seeds:

```sh
psql "postgres://moodora:moodora@localhost:5432/moodora_db?sslmode=disable" \
  -v ON_ERROR_STOP=1 \
  -f apps/api/seeds/lifestyle/001_lucky_lifestyle.sql
```

Run Moon and Astronomy source seeds:

```sh
psql "postgres://moodora:moodora@localhost:5432/moodora_db?sslmode=disable" \
  -v ON_ERROR_STOP=1 \
  -f apps/api/seeds/moon/001_astronomy_sources.sql
```

The future 78-card tarotapi.dev seed structure is documented in `apps/api/seeds/tarot/002_cards_template.sql`. It does not fetch external data yet.

## Run API

From `apps/api`:

```sh
APP_NAME=Moodora \
APP_ENV=local \
APP_PORT=8080 \
APP_TIMEZONE=Asia/Bangkok \
DEV_USER_ID=00000000-0000-0000-0000-000000000001 \
DATABASE_URL="postgres://moodora:moodora@localhost:5432/moodora_db?sslmode=disable" \
REDIS_URL="redis://localhost:6379" \
S3_ENDPOINT="http://localhost:9000" \
S3_BUCKET="moodora-assets" \
S3_ACCESS_KEY="moodora" \
S3_SECRET_KEY="moodora123" \
go run ./cmd/api
```

## API Smoke Tests

Health and version:

```sh
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/version
```

API documentation:

```sh
curl http://localhost:8080/api/v1/openapi.json
open http://localhost:8080/docs
open http://localhost:8080/swagger
```

`/docs` and `/swagger` serve Swagger UI for the embedded OpenAPI document. The UI loads Swagger UI assets from a CDN, while `/api/v1/openapi.json` is served directly by the Go API.

Tarot read endpoints:

```sh
curl http://localhost:8080/api/v1/tarot/cards
curl http://localhost:8080/api/v1/tarot/cards/ar01
curl http://localhost:8080/api/v1/tarot/spreads
curl http://localhost:8080/api/v1/tarot/spreads/three_cards
```

Create and read a Tarot reading:

```sh
curl -X POST http://localhost:8080/api/v1/tarot/readings \
  -H "Content-Type: application/json" \
  -d '{"spreadCode":"three_cards","topic":"love","language":"th","allowReversed":true,"question":"วันนี้ควรโฟกัสเรื่องความสัมพันธ์อย่างไร"}'
```

Wallet and check-in endpoints use `DEV_USER_ID` until real auth exists:

```sh
curl http://localhost:8080/api/v1/wallet
curl -X POST http://localhost:8080/api/v1/check-ins
curl http://localhost:8080/api/v1/coin-transactions
```

Check-in rewards use Bangkok-local dates. The reward cycle is day 1: 5, day 2: 5, day 3: 8, day 4: 8, day 5: 10, day 6: 15, day 7: 25. After day 7, continued streak days keep the day 7 reward.

Lucky Lifestyle endpoints:

```sh
curl http://localhost:8080/api/v1/lucky-colors/today
curl http://localhost:8080/api/v1/lucky-colors/today?purpose=career
curl http://localhost:8080/api/v1/lucky-foods/today
curl http://localhost:8080/api/v1/lucky-items/today
curl http://localhost:8080/api/v1/avoidance/today
curl http://localhost:8080/api/v1/daily-insights/today
```

Daily insight snapshots are stored by local date and timezone so the same user/date response stays stable. Recommendations are lifestyle and self-reflection prompts, not guaranteed outcomes.

Moon endpoints:

```sh
curl http://localhost:8080/api/v1/moon/today

curl -X POST http://localhost:8080/api/v1/moon/birthday \
  -H "Content-Type: application/json" \
  -d '{"birthDate":"2000-02-14","timezone":"Asia/Bangkok"}'

curl http://localhost:8080/api/v1/moon/reports/{id}
```

Moon calculations use Moodora's internal deterministic `moon_phase_v1` method. NASA/APOD and NASA SVS are seeded as future astronomy sources only; the API does not call live NASA services yet.

## Tarot Asset Pipeline

Source Tarot images are local-only and must not be committed. Place Rider-Waite originals here:

```text
local-assets/tarot/rider_waite/originals/{source_code}.{png,jpg,jpeg}
```

Each `{source_code}` must match `tarot_cards.source_code`, such as `ar01` or `sw08`.

Install the WebP encoder:

```sh
brew install webp
```

Run migrations, import Tarot cards, and make sure the S3-compatible bucket exists. Then process and upload assets:

```sh
cd apps/api
DATABASE_URL="postgres://moodora:moodora@localhost:5432/moodora_db?sslmode=disable" \
S3_ENDPOINT="http://localhost:9000" \
S3_REGION="auto" \
S3_BUCKET="moodora-assets" \
S3_ACCESS_KEY="moodora" \
S3_SECRET_KEY="moodora123" \
S3_PUBLIC_BASE_URL="http://localhost:9000/moodora-assets" \
go run ./cmd/process-tarot-assets
```

Generated files are written locally under:

```text
processed-assets/tarot/rider_waite/
```

Uploaded object keys use:

```text
tarot/rider_waite/webp/thumb/{source_code}.webp
tarot/rider_waite/webp/medium/{source_code}.webp
tarot/rider_waite/webp/large/{source_code}.webp
tarot/rider_waite/jpg/medium/{source_code}.jpg
```

The command upserts `tarot_card_assets` records and is safe to rerun. `local-assets/` and `processed-assets/` are gitignored.
