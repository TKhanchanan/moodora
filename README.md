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
