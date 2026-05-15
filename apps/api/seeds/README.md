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

Run sample Tarot interpretation seeds after importing cards:

```sh
psql "postgres://moodora:moodora@localhost:5432/moodora_db?sslmode=disable" \
  -v ON_ERROR_STOP=1 \
  -f apps/api/seeds/tarot/003_interpretations_sample.sql
```

Create the local development user used by `DEV_USER_ID`:

```sh
psql "postgres://moodora:moodora@localhost:5432/moodora_db?sslmode=disable" \
  -v ON_ERROR_STOP=1 \
  -v dev_user_id="00000000-0000-0000-0000-000000000001" \
  -f apps/api/seeds/dev_user.sql
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

`apps/api/seeds/tarot/002_cards_template.sql` documents the future 78-card tarotapi.dev import shape but does not insert card data yet.
