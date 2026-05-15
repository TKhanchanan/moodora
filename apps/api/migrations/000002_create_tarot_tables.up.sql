CREATE TABLE tarot_cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_code TEXT NOT NULL,
    name_en TEXT NOT NULL,
    type TEXT NOT NULL,
    suit TEXT,
    meaning_up_en TEXT NOT NULL,
    meaning_rev_en TEXT NOT NULL,
    description_en TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT tarot_cards_source_code_not_empty CHECK (length(trim(source_code)) > 0),
    CONSTRAINT tarot_cards_name_en_not_empty CHECK (length(trim(name_en)) > 0),
    CONSTRAINT tarot_cards_type_valid CHECK (type IN ('major', 'minor')),
    CONSTRAINT tarot_cards_suit_valid CHECK (suit IS NULL OR suit IN ('wands', 'cups', 'swords', 'pentacles')),
    CONSTRAINT tarot_cards_major_has_no_suit CHECK (
        (type = 'major' AND suit IS NULL)
        OR (type = 'minor' AND suit IS NOT NULL)
    ),
    CONSTRAINT tarot_cards_meaning_up_en_not_empty CHECK (length(trim(meaning_up_en)) > 0),
    CONSTRAINT tarot_cards_meaning_rev_en_not_empty CHECK (length(trim(meaning_rev_en)) > 0),
    CONSTRAINT tarot_cards_description_en_not_empty CHECK (length(trim(description_en)) > 0)
);

CREATE UNIQUE INDEX tarot_cards_source_code_unique_idx ON tarot_cards (source_code);
CREATE INDEX tarot_cards_type_idx ON tarot_cards (type);
CREATE INDEX tarot_cards_suit_idx ON tarot_cards (suit);

CREATE TABLE tarot_card_translations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    card_id UUID NOT NULL REFERENCES tarot_cards (id) ON DELETE CASCADE,
    language TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    meaning_upright TEXT NOT NULL,
    meaning_reversed TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT tarot_card_translations_language_not_empty CHECK (length(trim(language)) > 0),
    CONSTRAINT tarot_card_translations_name_not_empty CHECK (length(trim(name)) > 0),
    CONSTRAINT tarot_card_translations_description_not_empty CHECK (length(trim(description)) > 0),
    CONSTRAINT tarot_card_translations_meaning_upright_not_empty CHECK (length(trim(meaning_upright)) > 0),
    CONSTRAINT tarot_card_translations_meaning_reversed_not_empty CHECK (length(trim(meaning_reversed)) > 0)
);

CREATE UNIQUE INDEX tarot_card_translations_card_language_unique_idx
    ON tarot_card_translations (card_id, language);
CREATE INDEX tarot_card_translations_language_idx ON tarot_card_translations (language);

CREATE TABLE tarot_card_interpretations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    card_id UUID NOT NULL REFERENCES tarot_cards (id) ON DELETE CASCADE,
    language TEXT NOT NULL,
    topic TEXT NOT NULL,
    orientation TEXT NOT NULL,
    short_meaning TEXT NOT NULL,
    full_meaning TEXT NOT NULL,
    advice TEXT NOT NULL,
    version INTEGER NOT NULL DEFAULT 1,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT tarot_card_interpretations_language_not_empty CHECK (length(trim(language)) > 0),
    CONSTRAINT tarot_card_interpretations_topic_valid CHECK (topic IN ('general', 'love', 'career', 'money')),
    CONSTRAINT tarot_card_interpretations_orientation_valid CHECK (orientation IN ('upright', 'reversed')),
    CONSTRAINT tarot_card_interpretations_short_meaning_not_empty CHECK (length(trim(short_meaning)) > 0),
    CONSTRAINT tarot_card_interpretations_full_meaning_not_empty CHECK (length(trim(full_meaning)) > 0),
    CONSTRAINT tarot_card_interpretations_advice_not_empty CHECK (length(trim(advice)) > 0),
    CONSTRAINT tarot_card_interpretations_version_positive CHECK (version > 0)
);

CREATE UNIQUE INDEX tarot_card_interpretations_lookup_unique_idx
    ON tarot_card_interpretations (card_id, language, topic, orientation, version);
CREATE INDEX tarot_card_interpretations_lookup_idx
    ON tarot_card_interpretations (language, topic, orientation);

CREATE TABLE tarot_spreads (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    card_count INTEGER NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT tarot_spreads_code_not_empty CHECK (length(trim(code)) > 0),
    CONSTRAINT tarot_spreads_name_not_empty CHECK (length(trim(name)) > 0),
    CONSTRAINT tarot_spreads_description_not_empty CHECK (length(trim(description)) > 0),
    CONSTRAINT tarot_spreads_card_count_positive CHECK (card_count > 0)
);

CREATE UNIQUE INDEX tarot_spreads_code_unique_idx ON tarot_spreads (code);
CREATE INDEX tarot_spreads_active_idx ON tarot_spreads (is_active);

CREATE TABLE tarot_spread_positions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    spread_id UUID NOT NULL REFERENCES tarot_spreads (id) ON DELETE CASCADE,
    position_number INTEGER NOT NULL,
    code TEXT NOT NULL,
    name TEXT NOT NULL,
    description TEXT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT tarot_spread_positions_position_positive CHECK (position_number > 0),
    CONSTRAINT tarot_spread_positions_code_not_empty CHECK (length(trim(code)) > 0),
    CONSTRAINT tarot_spread_positions_name_not_empty CHECK (length(trim(name)) > 0),
    CONSTRAINT tarot_spread_positions_description_not_empty CHECK (length(trim(description)) > 0)
);

CREATE UNIQUE INDEX tarot_spread_positions_spread_position_unique_idx
    ON tarot_spread_positions (spread_id, position_number);
CREATE UNIQUE INDEX tarot_spread_positions_spread_code_unique_idx
    ON tarot_spread_positions (spread_id, code);
CREATE INDEX tarot_spread_positions_spread_id_idx ON tarot_spread_positions (spread_id);

CREATE TABLE tarot_readings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users (id) ON DELETE SET NULL,
    spread_id UUID NOT NULL REFERENCES tarot_spreads (id) ON DELETE RESTRICT,
    language TEXT NOT NULL DEFAULT 'en',
    topic TEXT NOT NULL DEFAULT 'general',
    status TEXT NOT NULL DEFAULT 'completed',
    result_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT tarot_readings_language_not_empty CHECK (length(trim(language)) > 0),
    CONSTRAINT tarot_readings_topic_valid CHECK (topic IN ('general', 'love', 'career', 'money')),
    CONSTRAINT tarot_readings_status_valid CHECK (status IN ('pending', 'completed', 'failed')),
    CONSTRAINT tarot_readings_result_snapshot_object CHECK (jsonb_typeof(result_snapshot) = 'object')
);

CREATE INDEX tarot_readings_user_created_at_idx ON tarot_readings (user_id, created_at DESC);
CREATE INDEX tarot_readings_spread_created_at_idx ON tarot_readings (spread_id, created_at DESC);
CREATE INDEX tarot_readings_topic_idx ON tarot_readings (topic);
CREATE INDEX tarot_readings_status_idx ON tarot_readings (status);

CREATE TABLE tarot_reading_cards (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    reading_id UUID NOT NULL REFERENCES tarot_readings (id) ON DELETE CASCADE,
    card_id UUID NOT NULL REFERENCES tarot_cards (id) ON DELETE RESTRICT,
    spread_position_id UUID REFERENCES tarot_spread_positions (id) ON DELETE SET NULL,
    position_number INTEGER NOT NULL,
    orientation TEXT NOT NULL,
    meaning_snapshot TEXT NOT NULL,
    advice_snapshot TEXT NOT NULL,
    result_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT tarot_reading_cards_position_positive CHECK (position_number > 0),
    CONSTRAINT tarot_reading_cards_orientation_valid CHECK (orientation IN ('upright', 'reversed')),
    CONSTRAINT tarot_reading_cards_meaning_snapshot_not_empty CHECK (length(trim(meaning_snapshot)) > 0),
    CONSTRAINT tarot_reading_cards_advice_snapshot_not_empty CHECK (length(trim(advice_snapshot)) > 0),
    CONSTRAINT tarot_reading_cards_result_snapshot_object CHECK (jsonb_typeof(result_snapshot) = 'object')
);

CREATE UNIQUE INDEX tarot_reading_cards_reading_position_unique_idx
    ON tarot_reading_cards (reading_id, position_number);
CREATE INDEX tarot_reading_cards_reading_id_idx ON tarot_reading_cards (reading_id);
CREATE INDEX tarot_reading_cards_card_id_idx ON tarot_reading_cards (card_id);
CREATE INDEX tarot_reading_cards_orientation_idx ON tarot_reading_cards (orientation);

CREATE TRIGGER tarot_cards_set_updated_at
BEFORE UPDATE ON tarot_cards
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER tarot_card_translations_set_updated_at
BEFORE UPDATE ON tarot_card_translations
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER tarot_card_interpretations_set_updated_at
BEFORE UPDATE ON tarot_card_interpretations
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER tarot_spreads_set_updated_at
BEFORE UPDATE ON tarot_spreads
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER tarot_spread_positions_set_updated_at
BEFORE UPDATE ON tarot_spread_positions
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER tarot_readings_set_updated_at
BEFORE UPDATE ON tarot_readings
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
