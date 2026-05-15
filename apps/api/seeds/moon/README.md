# Moon Seeds

Run migrations first:

```sh
migrate -path apps/api/migrations \
  -database "postgres://moodora:moodora@localhost:5432/moodora_db?sslmode=disable" \
  up
```

Seed astronomy sources:

```sh
psql "postgres://moodora:moodora@localhost:5432/moodora_db?sslmode=disable" \
  -v ON_ERROR_STOP=1 \
  -f apps/api/seeds/moon/001_astronomy_sources.sql
```

These seeds are idempotent. NASA sources are placeholders for future adapters and are not called by current tests or API handlers.
