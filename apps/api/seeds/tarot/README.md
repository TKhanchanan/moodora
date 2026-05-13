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

Tarot card data from tarotapi.dev is intentionally not fetched by these seeds yet. Add a separate import command later so source data can be validated, normalized, and reviewed before insertion.

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
