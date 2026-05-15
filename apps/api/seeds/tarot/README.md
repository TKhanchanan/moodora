# Tarot Seeds

This directory contains local seed files for Tarot reference data.

Run core migrations before applying seeds:

```sh
migrate -path apps/api/migrations \
  -database "postgres://moodora:moodora@localhost:5432/moodora_db?sslmode=disable" \
  up
```

Seed configurable spreads:

```sh
psql "postgres://moodora:moodora@localhost:5432/moodora_db?sslmode=disable" \
  -v ON_ERROR_STOP=1 \
  -f apps/api/seeds/tarot/001_spreads.sql
```

The spread seed inserts or updates:

- `one_card` with one position: `general`
- `three_cards` with three positions: `past`, `present`, `future`
- `celtic_cross` with ten positions: `current_situation`, `challenge`, `subconscious`, `past_influence`, `conscious_goal`, `near_future`, `self`, `environment`, `hopes_fears`, `final_outcome`

Tarot card data from tarotapi.dev is intentionally not fetched by these seeds yet. `002_cards_template.sql` documents the future 78-card seed shape using `source_code` values such as `ar01` and `sw08`.

Import the built-in 78-card source set:

```sh
cd apps/api
go run ./cmd/import-tarot-cards
```

Seed sample Thai/English interpretations:

```sh
psql "postgres://moodora:moodora@localhost:5432/moodora_db?sslmode=disable" \
  -v ON_ERROR_STOP=1 \
  -f apps/api/seeds/tarot/003_interpretations_sample.sql
```

## Card Asset Storage

Tarot card image metadata belongs in `tarot_card_assets`. Actual image files should be stored in the configured S3-compatible object storage, such as MinIO for local development.

Expected object key convention:

```text
tarot/{deck_code}/webp/{size}/{source_code}.webp
tarot/{deck_code}/jpg/{size}/{source_code}.jpg
```

Example keys:

```text
tarot/rider_waite/webp/thumb/ar01.webp
tarot/rider_waite/jpg/large/sw08.jpg
```

Do not commit generated or downloaded Tarot image files to the repository.

## Asset Pipeline

Place original card files in:

```text
local-assets/tarot/rider_waite/originals/{source_code}.{png,jpg,jpeg}
```

Run the backend command from `apps/api`:

```sh
DATABASE_URL="postgres://moodora:moodora@localhost:5432/moodora_db?sslmode=disable" \
S3_ENDPOINT="http://localhost:9000" \
S3_REGION="auto" \
S3_BUCKET="moodora-assets" \
S3_ACCESS_KEY="moodora" \
S3_SECRET_KEY="moodora123" \
S3_PUBLIC_BASE_URL="http://localhost:9000/moodora-assets" \
go run ./cmd/process-tarot-assets
```

The command generates WebP and JPG derivatives under `processed-assets/tarot/rider_waite/`, uploads them to S3-compatible storage, and upserts `tarot_card_assets`.
