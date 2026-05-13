CREATE TABLE lucky_colors (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL,
    name_th TEXT NOT NULL,
    name_en TEXT NOT NULL,
    hex TEXT NOT NULL,
    meaning TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT lucky_colors_code_not_empty CHECK (length(trim(code)) > 0),
    CONSTRAINT lucky_colors_name_th_not_empty CHECK (length(trim(name_th)) > 0),
    CONSTRAINT lucky_colors_name_en_not_empty CHECK (length(trim(name_en)) > 0),
    CONSTRAINT lucky_colors_hex_format CHECK (hex ~ '^#[0-9A-Fa-f]{6}$'),
    CONSTRAINT lucky_colors_meaning_not_empty CHECK (length(trim(meaning)) > 0)
);

CREATE UNIQUE INDEX lucky_colors_code_unique_idx ON lucky_colors (code);
CREATE INDEX lucky_colors_active_idx ON lucky_colors (is_active);

CREATE TABLE lucky_color_rules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    day_of_week INTEGER NOT NULL,
    birth_day_of_week INTEGER,
    purpose TEXT NOT NULL,
    color_id UUID NOT NULL REFERENCES lucky_colors (id) ON DELETE CASCADE,
    rule_type TEXT NOT NULL,
    weight INTEGER NOT NULL DEFAULT 1,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT lucky_color_rules_day_valid CHECK (day_of_week BETWEEN 0 AND 6),
    CONSTRAINT lucky_color_rules_birth_day_valid CHECK (birth_day_of_week IS NULL OR birth_day_of_week BETWEEN 0 AND 6),
    CONSTRAINT lucky_color_rules_purpose_valid CHECK (purpose IN ('general', 'love', 'career', 'money', 'study', 'interview')),
    CONSTRAINT lucky_color_rules_type_valid CHECK (rule_type IN ('lucky', 'avoid')),
    CONSTRAINT lucky_color_rules_weight_positive CHECK (weight > 0)
);

CREATE INDEX lucky_color_rules_lookup_idx
    ON lucky_color_rules (day_of_week, birth_day_of_week, purpose, rule_type, is_active, weight DESC);
CREATE INDEX lucky_color_rules_color_id_idx ON lucky_color_rules (color_id);
CREATE INDEX lucky_color_rules_active_idx ON lucky_color_rules (is_active);
CREATE UNIQUE INDEX lucky_color_rules_unique_idx
    ON lucky_color_rules (day_of_week, COALESCE(birth_day_of_week, -1), purpose, color_id, rule_type);

CREATE TABLE lucky_foods (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL,
    name_th TEXT NOT NULL,
    name_en TEXT NOT NULL,
    category TEXT NOT NULL,
    tags TEXT[] NOT NULL DEFAULT '{}',
    description TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT lucky_foods_code_not_empty CHECK (length(trim(code)) > 0),
    CONSTRAINT lucky_foods_name_th_not_empty CHECK (length(trim(name_th)) > 0),
    CONSTRAINT lucky_foods_name_en_not_empty CHECK (length(trim(name_en)) > 0),
    CONSTRAINT lucky_foods_category_not_empty CHECK (length(trim(category)) > 0),
    CONSTRAINT lucky_foods_description_not_empty CHECK (length(trim(description)) > 0)
);

CREATE UNIQUE INDEX lucky_foods_code_unique_idx ON lucky_foods (code);
CREATE INDEX lucky_foods_active_idx ON lucky_foods (is_active);
CREATE INDEX lucky_foods_category_idx ON lucky_foods (category);
CREATE INDEX lucky_foods_tags_idx ON lucky_foods USING GIN (tags);

CREATE TABLE lucky_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL,
    name_th TEXT NOT NULL,
    name_en TEXT NOT NULL,
    category TEXT NOT NULL,
    tags TEXT[] NOT NULL DEFAULT '{}',
    description TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT lucky_items_code_not_empty CHECK (length(trim(code)) > 0),
    CONSTRAINT lucky_items_name_th_not_empty CHECK (length(trim(name_th)) > 0),
    CONSTRAINT lucky_items_name_en_not_empty CHECK (length(trim(name_en)) > 0),
    CONSTRAINT lucky_items_category_not_empty CHECK (length(trim(category)) > 0),
    CONSTRAINT lucky_items_description_not_empty CHECK (length(trim(description)) > 0)
);

CREATE UNIQUE INDEX lucky_items_code_unique_idx ON lucky_items (code);
CREATE INDEX lucky_items_active_idx ON lucky_items (is_active);
CREATE INDEX lucky_items_category_idx ON lucky_items (category);
CREATE INDEX lucky_items_tags_idx ON lucky_items USING GIN (tags);

CREATE TABLE avoidance_recommendations (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code TEXT NOT NULL,
    category TEXT NOT NULL,
    text_th TEXT NOT NULL,
    text_en TEXT NOT NULL,
    mood_tag TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT avoidance_recommendations_code_not_empty CHECK (length(trim(code)) > 0),
    CONSTRAINT avoidance_recommendations_category_not_empty CHECK (length(trim(category)) > 0),
    CONSTRAINT avoidance_recommendations_text_th_not_empty CHECK (length(trim(text_th)) > 0),
    CONSTRAINT avoidance_recommendations_text_en_not_empty CHECK (length(trim(text_en)) > 0),
    CONSTRAINT avoidance_recommendations_mood_tag_not_empty CHECK (length(trim(mood_tag)) > 0)
);

CREATE UNIQUE INDEX avoidance_recommendations_code_unique_idx ON avoidance_recommendations (code);
CREATE INDEX avoidance_recommendations_active_idx ON avoidance_recommendations (is_active);
CREATE INDEX avoidance_recommendations_category_idx ON avoidance_recommendations (category);
CREATE INDEX avoidance_recommendations_mood_tag_idx ON avoidance_recommendations (mood_tag);

CREATE TABLE daily_insights (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users (id) ON DELETE CASCADE,
    insight_date DATE NOT NULL,
    timezone TEXT NOT NULL DEFAULT 'Asia/Bangkok',
    result_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT daily_insights_timezone_not_empty CHECK (length(trim(timezone)) > 0),
    CONSTRAINT daily_insights_result_snapshot_object CHECK (jsonb_typeof(result_snapshot) = 'object')
);

CREATE UNIQUE INDEX daily_insights_user_date_timezone_unique_idx
    ON daily_insights (user_id, insight_date, timezone)
    WHERE user_id IS NOT NULL;
CREATE UNIQUE INDEX daily_insights_anon_date_timezone_unique_idx
    ON daily_insights (insight_date, timezone)
    WHERE user_id IS NULL;
CREATE INDEX daily_insights_user_created_at_idx ON daily_insights (user_id, created_at DESC);
CREATE INDEX daily_insights_insight_date_idx ON daily_insights (insight_date);

CREATE TRIGGER lucky_colors_set_updated_at
BEFORE UPDATE ON lucky_colors
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER lucky_color_rules_set_updated_at
BEFORE UPDATE ON lucky_color_rules
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER lucky_foods_set_updated_at
BEFORE UPDATE ON lucky_foods
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER lucky_items_set_updated_at
BEFORE UPDATE ON lucky_items
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER avoidance_recommendations_set_updated_at
BEFORE UPDATE ON avoidance_recommendations
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER daily_insights_set_updated_at
BEFORE UPDATE ON daily_insights
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
