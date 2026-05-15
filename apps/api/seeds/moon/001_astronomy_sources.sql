WITH source_rows (source_code, source_name, source_url, description) AS (
    VALUES
        (
            'internal_moon_calculation',
            'Moodora Internal Moon Calculation',
            'internal://moon-phase-v1',
            'Deterministic internal lunar phase calculation used for local moon reports.'
        ),
        (
            'nasa_apod_future',
            'NASA Astronomy Picture of the Day',
            'https://api.nasa.gov/planetary/apod',
            'Future adapter target for astronomy imagery and descriptions. Not called by current tests or APIs.'
        ),
        (
            'nasa_svs_future',
            'NASA Scientific Visualization Studio',
            'https://svs.gsfc.nasa.gov/',
            'Future reference source for moon visualizations and astronomy media. Not called by current tests or APIs.'
        )
)
INSERT INTO astronomy_sources (source_code, source_name, source_url, description)
SELECT source_code, source_name, source_url, description
FROM source_rows
ON CONFLICT (source_code) DO UPDATE
SET source_name = EXCLUDED.source_name,
    source_url = EXCLUDED.source_url,
    description = EXCLUDED.description,
    is_active = true,
    updated_at = now();
