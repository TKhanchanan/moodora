WITH thai_names (source_code, name_th) AS (
    VALUES
        ('ar00', 'คนโง่'), ('ar01', 'นักมายากล'), ('ar02', 'มหาปุโรหิตหญิง'), ('ar03', 'จักรพรรดินี'),
        ('ar04', 'จักรพรรดิ'), ('ar05', 'สังฆราช'), ('ar06', 'คู่รัก'), ('ar07', 'รถศึก'),
        ('ar08', 'พละกำลัง'), ('ar09', 'ฤๅษี'), ('ar10', 'กงล้อแห่งโชคชะตา'), ('ar11', 'ความยุติธรรม'),
        ('ar12', 'คนห้อยหัว'), ('ar13', 'ความตาย'), ('ar14', 'ความพอดี'), ('ar15', 'ปีศาจ'),
        ('ar16', 'หอคอย'), ('ar17', 'ดวงดาว'), ('ar18', 'ดวงจันทร์'), ('ar19', 'ดวงอาทิตย์'),
        ('ar20', 'การพิพากษา'), ('ar21', 'โลก'),
        ('wa01', 'ไม้เท้าหนึ่ง'), ('wa02', 'ไม้เท้าสอง'), ('wa03', 'ไม้เท้าสาม'), ('wa04', 'ไม้เท้าสี่'),
        ('wa05', 'ไม้เท้าห้า'), ('wa06', 'ไม้เท้าหก'), ('wa07', 'ไม้เท้าเจ็ด'), ('wa08', 'ไม้เท้าแปด'),
        ('wa09', 'ไม้เท้าเก้า'), ('wa10', 'ไม้เท้าสิบ'), ('wa11', 'มหาดเล็กไม้เท้า'), ('wa12', 'อัศวินไม้เท้า'),
        ('wa13', 'ราชินีไม้เท้า'), ('wa14', 'ราชาไม้เท้า'),
        ('cu01', 'ถ้วยหนึ่ง'), ('cu02', 'ถ้วยสอง'), ('cu03', 'ถ้วยสาม'), ('cu04', 'ถ้วยสี่'),
        ('cu05', 'ถ้วยห้า'), ('cu06', 'ถ้วยหก'), ('cu07', 'ถ้วยเจ็ด'), ('cu08', 'ถ้วยแปด'),
        ('cu09', 'ถ้วยเก้า'), ('cu10', 'ถ้วยสิบ'), ('cu11', 'มหาดเล็กถ้วย'), ('cu12', 'อัศวินถ้วย'),
        ('cu13', 'ราชินีถ้วย'), ('cu14', 'ราชาถ้วย'),
        ('sw01', 'ดาบหนึ่ง'), ('sw02', 'ดาบสอง'), ('sw03', 'ดาบสาม'), ('sw04', 'ดาบสี่'),
        ('sw05', 'ดาบห้า'), ('sw06', 'ดาบหก'), ('sw07', 'ดาบเจ็ด'), ('sw08', 'ดาบแปด'),
        ('sw09', 'ดาบเก้า'), ('sw10', 'ดาบสิบ'), ('sw11', 'มหาดเล็กดาบ'), ('sw12', 'อัศวินดาบ'),
        ('sw13', 'ราชินีดาบ'), ('sw14', 'ราชาดาบ'),
        ('pe01', 'เหรียญหนึ่ง'), ('pe02', 'เหรียญสอง'), ('pe03', 'เหรียญสาม'), ('pe04', 'เหรียญสี่'),
        ('pe05', 'เหรียญห้า'), ('pe06', 'เหรียญหก'), ('pe07', 'เหรียญเจ็ด'), ('pe08', 'เหรียญแปด'),
        ('pe09', 'เหรียญเก้า'), ('pe10', 'เหรียญสิบ'), ('pe11', 'มหาดเล็กเหรียญ'), ('pe12', 'อัศวินเหรียญ'),
        ('pe13', 'ราชินีเหรียญ'), ('pe14', 'ราชาเหรียญ')
),
card_rows AS (
    SELECT
        tc.id,
        tc.source_code,
        tc.name_en,
        tc.type,
        tc.suit,
        tc.meaning_up_en,
        tc.meaning_rev_en,
        tc.description_en,
        thai_names.name_th,
        CASE tc.suit
            WHEN 'wands' THEN 'ไม้เท้า'
            WHEN 'cups' THEN 'ถ้วย'
            WHEN 'swords' THEN 'ดาบ'
            WHEN 'pentacles' THEN 'เหรียญ'
            ELSE NULL
        END AS suit_th
    FROM tarot_cards tc
    JOIN thai_names ON thai_names.source_code = tc.source_code
),
translation_rows AS (
    SELECT
        id AS card_id,
        'en' AS language,
        name_en AS name,
        description_en AS description,
        meaning_up_en AS meaning_upright,
        meaning_rev_en AS meaning_reversed
    FROM card_rows
    UNION ALL
    SELECT
        id AS card_id,
        'th' AS language,
        name_th AS name,
        CASE
            WHEN type = 'major' THEN name_th || ' เป็นไพ่เมเจอร์อาร์คานาที่สะท้อนบทเรียนสำคัญ จังหวะเปลี่ยนผ่าน และหัวข้อใหญ่ของชีวิตเพื่อการใคร่ครวญ'
            ELSE name_th || ' เป็นไพ่ไมเนอร์อาร์คานาชุด' || suit_th || ' ที่สะท้อนสถานการณ์ประจำวัน การตัดสินใจ และพลังที่จับต้องได้'
        END AS description,
        'เมื่อหัวตั้ง ' || name_th || ' ชวนให้มองพลังด้านสร้างสรรค์ของไพ่ใบนี้ และใช้เป็นแนวทางใคร่ครวญอย่างมีสติ' AS meaning_upright,
        'เมื่อกลับหัว ' || name_th || ' ชวนให้สังเกตจุดติดขัด ความกลัว หรือพลังที่ยังไม่สมดุล โดยไม่ตัดสินตัวเองเร็วเกินไป' AS meaning_reversed
    FROM card_rows
)
INSERT INTO tarot_card_translations (
    card_id, language, name, description, meaning_upright, meaning_reversed
)
SELECT card_id, language, name, description, meaning_upright, meaning_reversed
FROM translation_rows
ON CONFLICT (card_id, language) DO UPDATE
SET name = EXCLUDED.name,
    description = EXCLUDED.description,
    meaning_upright = EXCLUDED.meaning_upright,
    meaning_reversed = EXCLUDED.meaning_reversed,
    updated_at = now();

WITH topics (topic, topic_en, topic_th, focus_en, focus_th, advice_en, advice_th) AS (
    VALUES
        ('general', 'general life', 'ภาพรวมชีวิต', 'your current rhythm, emotional clarity, and the next grounded step', 'จังหวะชีวิตตอนนี้ ความชัดเจนในใจ และก้าวเล็ก ๆ ที่ทำได้จริง', 'Choose one practical action today, then pause and observe how it changes your perspective.', 'เลือกลงมือทำหนึ่งเรื่องที่จับต้องได้ในวันนี้ แล้วค่อยสังเกตว่ามุมมองของคุณเปลี่ยนไปอย่างไร'),
        ('love', 'love and relationships', 'ความรักและความสัมพันธ์', 'communication, emotional boundaries, and the way connection is being built', 'การสื่อสาร ขอบเขตทางใจ และรูปแบบความสัมพันธ์ที่กำลังก่อตัว', 'Speak honestly, listen carefully, and do not rush another person or yourself into a fixed answer.', 'สื่อสารอย่างจริงใจ ฟังให้มากพอ และอย่าเร่งตัวเองหรืออีกฝ่ายให้ต้องมีคำตอบทันที'),
        ('career', 'career and work', 'การงานและเส้นทางอาชีพ', 'priorities, collaboration, timing, and the quality of your effort', 'ลำดับความสำคัญ การทำงานร่วมกัน จังหวะเวลา และคุณภาพของความพยายาม', 'Clarify the next useful task and protect your energy from work that only creates noise.', 'ทำให้ก้าวถัดไปชัดเจน และรักษาพลังของตัวเองจากงานที่สร้างแต่เสียงรบกวน'),
        ('money', 'money and resources', 'การเงินและทรัพยากร', 'resources, spending choices, stability, and the relationship between value and security', 'ทรัพยากร การใช้จ่าย ความมั่นคง และความสัมพันธ์ระหว่างคุณค่ากับความปลอดภัย', 'Review one financial choice with calm attention before committing more time, money, or energy.', 'ทบทวนการตัดสินใจด้านเงินหนึ่งเรื่องอย่างใจเย็น ก่อนเพิ่มเวลา เงิน หรือพลังลงไป')
),
orientations (orientation, label_en, label_th) AS (
    VALUES
        ('upright', 'upright', 'หัวตั้ง'),
        ('reversed', 'reversed', 'กลับหัว')
),
languages (language) AS (
    VALUES ('en'), ('th')
),
card_rows AS (
    SELECT
        tc.id,
        tc.source_code,
        tc.name_en,
        tc.type,
        tc.suit,
        tc.meaning_up_en,
        tc.meaning_rev_en,
        th.name AS name_th,
        th.meaning_upright AS meaning_up_th,
        th.meaning_reversed AS meaning_rev_th
    FROM tarot_cards tc
    JOIN tarot_card_translations th ON th.card_id = tc.id AND th.language = 'th'
),
interpretation_rows AS (
    SELECT
        card_rows.id AS card_id,
        languages.language,
        topics.topic,
        orientations.orientation,
        CASE
            WHEN languages.language = 'th' AND orientations.orientation = 'upright' THEN card_rows.name_th || 'กับ' || topics.topic_th
            WHEN languages.language = 'th' THEN card_rows.name_th || 'กลับหัวกับ' || topics.topic_th
            WHEN orientations.orientation = 'upright' THEN card_rows.name_en || ' for ' || topics.topic_en
            ELSE card_rows.name_en || ' reversed for ' || topics.topic_en
        END AS short_meaning,
        CASE
            WHEN languages.language = 'th' AND orientations.orientation = 'upright'
                THEN card_rows.name_th || ' ในเรื่อง' || topics.topic_th || ' ชวนให้คุณใช้ความหมายหัวตั้งของไพ่ใบนี้เป็นกระจกสะท้อน ' || topics.focus_th || ' ความหมายหลักคือ: ' || card_rows.meaning_up_th || ' คำทำนายนี้เป็นการตีความเพื่อการใคร่ครวญ ไม่ใช่ข้อพิสูจน์แน่นอนของอนาคต'
            WHEN languages.language = 'th'
                THEN card_rows.name_th || ' กลับหัวในเรื่อง' || topics.topic_th || ' ชวนให้สังเกตพลังที่ติดขัดหรือยังไม่สมดุลเกี่ยวกับ ' || topics.focus_th || ' ความหมายหลักคือ: ' || card_rows.meaning_rev_th || ' คำทำนายนี้เป็นการตีความเพื่อการใคร่ครวญ ไม่ใช่ข้อพิสูจน์แน่นอนของอนาคต'
            WHEN orientations.orientation = 'upright'
                THEN card_rows.name_en || ' in ' || topics.topic_en || ' invites you to reflect on ' || topics.focus_en || '. Core meaning: ' || card_rows.meaning_up_en || ' Treat this as an interpretation for self-reflection, not scientific proof of the future.'
            ELSE card_rows.name_en || ' reversed in ' || topics.topic_en || ' asks you to notice where energy may feel blocked, delayed, or unbalanced around ' || topics.focus_en || '. Core meaning: ' || card_rows.meaning_rev_en || ' Treat this as an interpretation for self-reflection, not scientific proof of the future.'
        END AS full_meaning,
        CASE
            WHEN languages.language = 'th' AND orientations.orientation = 'upright'
                THEN topics.advice_th || ' ใช้พลังของ ' || card_rows.name_th || ' เป็นคำชวนให้เลือกอย่างมีสติ'
            WHEN languages.language = 'th'
                THEN topics.advice_th || ' หากรู้สึกฝืน ให้ลดจังหวะลงและกลับมาดูแลใจตัวเองก่อน'
            WHEN orientations.orientation = 'upright'
                THEN topics.advice_en || ' Let ' || card_rows.name_en || ' guide a conscious, grounded choice.'
            ELSE topics.advice_en || ' If it feels forced, slow down and return to a steadier inner pace first.'
        END AS advice,
        1 AS version
    FROM card_rows
    CROSS JOIN topics
    CROSS JOIN orientations
    CROSS JOIN languages
)
INSERT INTO tarot_card_interpretations (
    card_id, language, topic, orientation, short_meaning, full_meaning, advice, version
)
SELECT card_id, language, topic, orientation, short_meaning, full_meaning, advice, version
FROM interpretation_rows
ON CONFLICT (card_id, language, topic, orientation, version) DO UPDATE
SET short_meaning = EXCLUDED.short_meaning,
    full_meaning = EXCLUDED.full_meaning,
    advice = EXCLUDED.advice,
    updated_at = now();
