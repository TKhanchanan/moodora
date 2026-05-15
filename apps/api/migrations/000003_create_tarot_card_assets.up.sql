CREATE TABLE tarot_card_assets (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    card_id UUID NOT NULL REFERENCES tarot_cards (id) ON DELETE CASCADE,
    deck_code TEXT NOT NULL,
    size TEXT NOT NULL,
    format TEXT NOT NULL,
    url TEXT NOT NULL,
    width INTEGER NOT NULL,
    height INTEGER NOT NULL,
    file_size BIGINT NOT NULL,
    is_default BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT tarot_card_assets_deck_code_not_empty CHECK (length(trim(deck_code)) > 0),
    CONSTRAINT tarot_card_assets_size_valid CHECK (size IN ('thumb', 'medium', 'large')),
    CONSTRAINT tarot_card_assets_format_valid CHECK (format IN ('webp', 'jpg')),
    CONSTRAINT tarot_card_assets_url_not_empty CHECK (length(trim(url)) > 0),
    CONSTRAINT tarot_card_assets_width_positive CHECK (width > 0),
    CONSTRAINT tarot_card_assets_height_positive CHECK (height > 0),
    CONSTRAINT tarot_card_assets_file_size_positive CHECK (file_size > 0)
);

CREATE UNIQUE INDEX tarot_card_assets_variant_unique_idx
    ON tarot_card_assets (card_id, deck_code, size, format);
CREATE INDEX tarot_card_assets_card_id_idx ON tarot_card_assets (card_id);
CREATE INDEX tarot_card_assets_deck_code_idx ON tarot_card_assets (deck_code);
CREATE INDEX tarot_card_assets_size_idx ON tarot_card_assets (size);
CREATE INDEX tarot_card_assets_format_idx ON tarot_card_assets (format);
CREATE INDEX tarot_card_assets_is_default_idx ON tarot_card_assets (is_default);
CREATE UNIQUE INDEX tarot_card_assets_one_default_per_card_deck_idx
    ON tarot_card_assets (card_id, deck_code)
    WHERE is_default = true;

CREATE TRIGGER tarot_card_assets_set_updated_at
BEFORE UPDATE ON tarot_card_assets
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
