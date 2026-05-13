WITH one_card AS (
    INSERT INTO tarot_spreads (code, name, description, card_count)
    VALUES (
        'one_card',
        'One Card',
        'A single-card spread for a concise reflection.',
        1
    )
    ON CONFLICT (code) DO UPDATE
    SET name = EXCLUDED.name,
        description = EXCLUDED.description,
        card_count = EXCLUDED.card_count,
        updated_at = now()
    RETURNING id
)
INSERT INTO tarot_spread_positions (spread_id, position_number, code, name, description)
SELECT id, 1, 'message', 'Message', 'The main message or reflection for this reading.'
FROM one_card
ON CONFLICT (spread_id, position_number) DO UPDATE
SET code = EXCLUDED.code,
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    updated_at = now();

WITH three_cards AS (
    INSERT INTO tarot_spreads (code, name, description, card_count)
    VALUES (
        'three_cards',
        'Three Cards',
        'A three-card spread for past, present, and future reflection.',
        3
    )
    ON CONFLICT (code) DO UPDATE
    SET name = EXCLUDED.name,
        description = EXCLUDED.description,
        card_count = EXCLUDED.card_count,
        updated_at = now()
    RETURNING id
),
positions (position_number, code, name, description) AS (
    VALUES
        (1, 'past', 'Past', 'A past influence or context for the question.'),
        (2, 'present', 'Present', 'The current energy or situation.'),
        (3, 'future', 'Future', 'A possible direction or next reflection point.')
)
INSERT INTO tarot_spread_positions (spread_id, position_number, code, name, description)
SELECT three_cards.id, positions.position_number, positions.code, positions.name, positions.description
FROM three_cards
CROSS JOIN positions
ON CONFLICT (spread_id, position_number) DO UPDATE
SET code = EXCLUDED.code,
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    updated_at = now();

WITH celtic_cross AS (
    INSERT INTO tarot_spreads (code, name, description, card_count)
    VALUES (
        'celtic_cross',
        'Celtic Cross',
        'A ten-card spread for a deeper self-reflection reading.',
        10
    )
    ON CONFLICT (code) DO UPDATE
    SET name = EXCLUDED.name,
        description = EXCLUDED.description,
        card_count = EXCLUDED.card_count,
        updated_at = now()
    RETURNING id
),
positions (position_number, code, name, description) AS (
    VALUES
        (1, 'present', 'Present', 'The central situation or current energy.'),
        (2, 'challenge', 'Challenge', 'The crossing influence or immediate challenge.'),
        (3, 'foundation', 'Foundation', 'The root context beneath the situation.'),
        (4, 'past', 'Past', 'Recent past influence.'),
        (5, 'conscious', 'Conscious', 'What is visible, intended, or consciously held.'),
        (6, 'near_future', 'Near Future', 'A likely near-term development.'),
        (7, 'self', 'Self', 'The seeker''s current stance or self-view.'),
        (8, 'environment', 'Environment', 'External influences or surrounding context.'),
        (9, 'hopes_fears', 'Hopes and Fears', 'Inner hopes, worries, or expectations.'),
        (10, 'outcome', 'Outcome', 'A possible outcome or integration point.')
)
INSERT INTO tarot_spread_positions (spread_id, position_number, code, name, description)
SELECT celtic_cross.id, positions.position_number, positions.code, positions.name, positions.description
FROM celtic_cross
CROSS JOIN positions
ON CONFLICT (spread_id, position_number) DO UPDATE
SET code = EXCLUDED.code,
    name = EXCLUDED.name,
    description = EXCLUDED.description,
    updated_at = now();
