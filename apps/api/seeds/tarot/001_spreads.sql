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
SELECT id, 1, 'general', 'ภาพรวมของวันนี้', 'ภาพรวมของวันนี้'
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
        (1, 'past', 'อดีต', 'อดีต'),
        (2, 'present', 'ปัจจุบัน', 'ปัจจุบัน'),
        (3, 'future', 'แนวโน้มอนาคต', 'แนวโน้มอนาคต')
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
        (1, 'current_situation', 'สถานการณ์ปัจจุบัน', 'สถานการณ์ปัจจุบัน'),
        (2, 'challenge', 'อุปสรรคหรือสิ่งที่ท้าทาย', 'อุปสรรคหรือสิ่งที่ท้าทาย'),
        (3, 'subconscious', 'สิ่งที่อยู่ลึกในใจ', 'สิ่งที่อยู่ลึกในใจ'),
        (4, 'past_influence', 'อิทธิพลจากอดีต', 'อิทธิพลจากอดีต'),
        (5, 'conscious_goal', 'สิ่งที่คาดหวังหรือเป้าหมาย', 'สิ่งที่คาดหวังหรือเป้าหมาย'),
        (6, 'near_future', 'อนาคตอันใกล้', 'อนาคตอันใกล้'),
        (7, 'self', 'ตัวคุณในสถานการณ์นี้', 'ตัวคุณในสถานการณ์นี้'),
        (8, 'environment', 'สภาพแวดล้อมหรือคนรอบตัว', 'สภาพแวดล้อมหรือคนรอบตัว'),
        (9, 'hopes_fears', 'ความหวังและความกังวล', 'ความหวังและความกังวล'),
        (10, 'final_outcome', 'แนวโน้มบทสรุป', 'แนวโน้มบทสรุป')
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
