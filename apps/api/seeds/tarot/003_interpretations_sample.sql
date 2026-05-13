WITH sample_cards AS (
    SELECT id, source_code
    FROM tarot_cards
    WHERE source_code IN ('ar01', 'sw08', 'cu01')
),
translation_rows (source_code, language, name, description, meaning_upright, meaning_reversed) AS (
    VALUES
        (
            'ar01',
            'th',
            'นักมายากล',
            'ไพ่ใบนี้ชวนให้สำรวจพลัง ความตั้งใจ และสิ่งที่คุณมีอยู่ในมือ',
            'การตั้งใจลงมือทำอย่างชัดเจน',
            'การกระจายพลังหรือยังไม่มั่นใจในศักยภาพของตัวเอง'
        ),
        (
            'sw08',
            'th',
            'ดาบแปด',
            'ไพ่ใบนี้สะท้อนช่วงเวลาที่ใจอาจรู้สึกติดกรอบ และชวนให้มองหาทางเลือกที่ยังมีอยู่',
            'ความรู้สึกติดขัดหรือจำกัดตัวเอง',
            'การค่อย ๆ คลี่คลายจากความคิดที่กดทับ'
        ),
        (
            'cu01',
            'th',
            'ถ้วยหนึ่ง',
            'ไพ่ใบนี้สะท้อนพื้นที่ของความรู้สึก ความอ่อนโยน และการเริ่มต้นทางใจ',
            'การเปิดใจและเริ่มต้นความรู้สึกใหม่',
            'การดูแลใจตัวเองก่อนเปิดรับสิ่งใหม่'
        )
)
INSERT INTO tarot_card_translations (
    card_id, language, name, description, meaning_upright, meaning_reversed
)
SELECT sample_cards.id, translation_rows.language, translation_rows.name, translation_rows.description,
       translation_rows.meaning_upright, translation_rows.meaning_reversed
FROM translation_rows
JOIN sample_cards ON sample_cards.source_code = translation_rows.source_code
ON CONFLICT (card_id, language) DO UPDATE
SET name = EXCLUDED.name,
    description = EXCLUDED.description,
    meaning_upright = EXCLUDED.meaning_upright,
    meaning_reversed = EXCLUDED.meaning_reversed,
    updated_at = now();

WITH sample_cards AS (
    SELECT id, source_code
    FROM tarot_cards
    WHERE source_code IN ('ar01', 'sw08', 'cu01')
),
interpretation_rows (
    source_code, language, topic, orientation, short_meaning, full_meaning, advice, version
) AS (
    VALUES
        (
            'ar01',
            'th',
            'general',
            'upright',
            'ตั้งใจและลงมือ',
            'นักมายากลชวนให้คุณมองเห็นทรัพยากรที่มีอยู่ และใช้มันอย่างตั้งใจในวันนี้',
            'เลือกหนึ่งเรื่องที่สำคัญ แล้วลงมืออย่างชัดเจนโดยไม่ต้องรีบพิสูจน์ทุกอย่างพร้อมกัน',
            1
        ),
        (
            'ar01',
            'th',
            'love',
            'upright',
            'สื่อสารจากความตั้งใจ',
            'ในเรื่องความรัก ไพ่ใบนี้ชวนให้สื่อสารอย่างจริงใจและรู้ว่าคุณต้องการสร้างความสัมพันธ์แบบไหน',
            'ถามใจตัวเองให้ชัดก่อน แล้วค่อยพูดหรือทำด้วยความเคารพทั้งตัวเองและอีกฝ่าย',
            1
        ),
        (
            'ar01',
            'th',
            'general',
            'reversed',
            'รวมพลังกลับมา',
            'ไพ่กลับหัวชวนให้สังเกตว่าพลังของคุณกระจายไปหลายทางเกินไปหรือไม่',
            'พักสั้น ๆ แล้วจัดลำดับสิ่งที่ควรทำก่อนหลังอย่างอ่อนโยนกับตัวเอง',
            1
        ),
        (
            'sw08',
            'th',
            'general',
            'upright',
            'เห็นกรอบที่จำกัดใจ',
            'ดาบแปดสะท้อนความรู้สึกติดขัด แต่ก็ชวนให้แยกแยะว่าส่วนไหนคือข้อจำกัดจริง และส่วนไหนคือความกลัว',
            'เขียนทางเลือกเล็ก ๆ ที่ยังทำได้อย่างน้อยหนึ่งข้อ แล้วเริ่มจากตรงนั้น',
            1
        ),
        (
            'sw08',
            'th',
            'general',
            'reversed',
            'ค่อย ๆ คลี่คลาย',
            'ไพ่กลับหัวบอกถึงจังหวะที่คุณอาจเริ่มเห็นทางออกจากความคิดที่เคยกดทับ',
            'ให้เวลากับการเปลี่ยนมุมมอง และอย่าบังคับตัวเองให้ต้องหายกังวลทันที',
            1
        ),
        (
            'cu01',
            'th',
            'love',
            'upright',
            'เปิดพื้นที่ให้ความรู้สึก',
            'ถ้วยหนึ่งในเรื่องความรักชวนให้เปิดรับความรู้สึกอย่างอ่อนโยน โดยไม่ต้องสรุปอนาคตเร็วเกินไป',
            'ดูแลหัวใจของตัวเอง และสื่อสารความรู้สึกในจังหวะที่ปลอดภัย',
            1
        ),
        (
            'ar01',
            'en',
            'general',
            'upright',
            'Focused action',
            'The Magician invites you to notice what is already available and act with clear intention.',
            'Choose one meaningful step and take it with care.',
            1
        ),
        (
            'sw08',
            'en',
            'general',
            'upright',
            'Notice the limits',
            'Eight of Swords reflects a moment of feeling boxed in while still having room to reassess.',
            'Name one option that remains available, even if it is small.',
            1
        )
)
INSERT INTO tarot_card_interpretations (
    card_id, language, topic, orientation, short_meaning, full_meaning, advice, version
)
SELECT sample_cards.id, interpretation_rows.language, interpretation_rows.topic, interpretation_rows.orientation,
       interpretation_rows.short_meaning, interpretation_rows.full_meaning, interpretation_rows.advice,
       interpretation_rows.version
FROM interpretation_rows
JOIN sample_cards ON sample_cards.source_code = interpretation_rows.source_code
ON CONFLICT (card_id, language, topic, orientation, version) DO UPDATE
SET short_meaning = EXCLUDED.short_meaning,
    full_meaning = EXCLUDED.full_meaning,
    advice = EXCLUDED.advice,
    updated_at = now();
