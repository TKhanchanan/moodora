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

Run Tarot reference seeds after migrations:

```sh
psql "postgres://moodora:moodora@localhost:5432/moodora_db?sslmode=disable" \
  -v ON_ERROR_STOP=1 \
  -f apps/api/seeds/tarot/001_spreads.sql
```

The future 78-card tarotapi.dev seed structure is documented in `apps/api/seeds/tarot/002_cards_template.sql`. It does not fetch external data yet.
