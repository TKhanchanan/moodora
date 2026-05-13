# Lucky Lifestyle Seeds

Run migrations first:

```sh
migrate -path apps/api/migrations \
  -database "postgres://moodora:moodora@localhost:5432/moodora_db?sslmode=disable" \
  up
```

Seed lucky colors, rules, foods, items, and avoidance recommendations:

```sh
psql "postgres://moodora:moodora@localhost:5432/moodora_db?sslmode=disable" \
  -v ON_ERROR_STOP=1 \
  -f apps/api/seeds/lifestyle/001_lucky_lifestyle.sql
```

Seeds are idempotent and safe to rerun in local development.
