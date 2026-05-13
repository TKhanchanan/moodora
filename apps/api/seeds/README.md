# API Seeds

Seed files are environment-controlled reference data for local development.
They are intentionally separate from schema migrations.

Run migrations before seeds:

```sh
migrate -path apps/api/migrations \
  -database "postgres://moodora:moodora@localhost:5432/moodora_db?sslmode=disable" \
  up
```

Run Tarot seeds:

```sh
psql "postgres://moodora:moodora@localhost:5432/moodora_db?sslmode=disable" \
  -v ON_ERROR_STOP=1 \
  -f apps/api/seeds/tarot/001_spreads.sql
```

`apps/api/seeds/tarot/002_cards_template.sql` documents the future 78-card tarotapi.dev import shape but does not insert card data yet.
