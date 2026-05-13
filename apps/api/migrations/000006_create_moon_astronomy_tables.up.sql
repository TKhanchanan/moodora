CREATE TABLE astronomy_sources (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    source_code TEXT NOT NULL,
    source_name TEXT NOT NULL,
    source_url TEXT NOT NULL,
    description TEXT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT astronomy_sources_code_not_empty CHECK (length(trim(source_code)) > 0),
    CONSTRAINT astronomy_sources_name_not_empty CHECK (length(trim(source_name)) > 0),
    CONSTRAINT astronomy_sources_url_not_empty CHECK (length(trim(source_url)) > 0),
    CONSTRAINT astronomy_sources_description_not_empty CHECK (length(trim(description)) > 0)
);

CREATE UNIQUE INDEX astronomy_sources_code_unique_idx ON astronomy_sources (source_code);
CREATE INDEX astronomy_sources_active_idx ON astronomy_sources (is_active);

CREATE TABLE moon_reports (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES users (id) ON DELETE SET NULL,
    birth_date DATE,
    target_date DATE NOT NULL,
    timezone TEXT NOT NULL DEFAULT 'Asia/Bangkok',
    moon_phase TEXT NOT NULL,
    illumination NUMERIC(5,2) NOT NULL,
    moon_age NUMERIC(5,2) NOT NULL,
    image_url TEXT,
    calculation_method_version TEXT NOT NULL,
    interpretation_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
    source_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
    result_snapshot JSONB NOT NULL DEFAULT '{}'::jsonb,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT now(),
    CONSTRAINT moon_reports_timezone_not_empty CHECK (length(trim(timezone)) > 0),
    CONSTRAINT moon_reports_phase_valid CHECK (
        moon_phase IN (
            'new_moon',
            'waxing_crescent',
            'first_quarter',
            'waxing_gibbous',
            'full_moon',
            'waning_gibbous',
            'last_quarter',
            'waning_crescent'
        )
    ),
    CONSTRAINT moon_reports_illumination_valid CHECK (illumination >= 0 AND illumination <= 100),
    CONSTRAINT moon_reports_moon_age_valid CHECK (moon_age >= 0 AND moon_age < 30),
    CONSTRAINT moon_reports_method_not_empty CHECK (length(trim(calculation_method_version)) > 0),
    CONSTRAINT moon_reports_interpretation_snapshot_object CHECK (jsonb_typeof(interpretation_snapshot) = 'object'),
    CONSTRAINT moon_reports_source_snapshot_object CHECK (jsonb_typeof(source_snapshot) = 'object'),
    CONSTRAINT moon_reports_result_snapshot_object CHECK (jsonb_typeof(result_snapshot) = 'object')
);

CREATE INDEX moon_reports_user_created_at_idx ON moon_reports (user_id, created_at DESC);
CREATE INDEX moon_reports_target_date_idx ON moon_reports (target_date);
CREATE INDEX moon_reports_birth_date_idx ON moon_reports (birth_date);
CREATE INDEX moon_reports_phase_idx ON moon_reports (moon_phase);

CREATE TRIGGER astronomy_sources_set_updated_at
BEFORE UPDATE ON astronomy_sources
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();

CREATE TRIGGER moon_reports_set_updated_at
BEFORE UPDATE ON moon_reports
FOR EACH ROW
EXECUTE FUNCTION set_updated_at();
