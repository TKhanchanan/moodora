# Moodora API Migrations

Migrations are plain PostgreSQL SQL files under `apps/api/migrations`.

Use the `golang-migrate` CLI for local development and deployment automation:

```sh
brew install golang-migrate
```

From the repository root, start local dependencies:

```sh
docker compose up -d postgres redis minio
```

Run migrations:

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

Check the current migration version:

```sh
migrate -path apps/api/migrations \
  -database "postgres://moodora:moodora@localhost:5432/moodora_db?sslmode=disable" \
  version
```

Deployment environments should run the same migration files against their managed PostgreSQL database using a secret-provided `DATABASE_URL`.

## Seeds

Seed files live outside migrations because they are environment-controlled reference data.

Tarot spread seeds:

```sh
psql "postgres://moodora:moodora@localhost:5432/moodora_db?sslmode=disable" \
  -v ON_ERROR_STOP=1 \
  -f apps/api/seeds/tarot/001_spreads.sql
```

The future Tarot card seed structure is documented in:

```text
apps/api/seeds/tarot/002_cards_template.sql
```

Tarot card assets are stored as metadata in PostgreSQL and as files in object storage.
Use this object key convention:

```text
tarot/{deck_code}/webp/{size}/{source_code}.webp
tarot/{deck_code}/jpg/{size}/{source_code}.jpg
```
