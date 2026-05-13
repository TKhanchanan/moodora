WITH color_rows (code, name_th, name_en, hex, meaning) AS (
    VALUES
        ('soft_pink', 'ชมพูอ่อน', 'Soft Pink', '#F8BBD0', 'ความอ่อนโยน การเปิดใจ และความสัมพันธ์ที่นุ่มนวล'),
        ('sky_blue', 'ฟ้าใส', 'Sky Blue', '#90CAF9', 'ความสงบ การสื่อสาร และความคิดที่ปลอดโปร่ง'),
        ('sun_yellow', 'เหลืองสดใส', 'Sun Yellow', '#FFD54F', 'พลังบวก ความมั่นใจ และการเริ่มต้นอย่างสดใส'),
        ('mint_green', 'เขียวมิ้นต์', 'Mint Green', '#A5D6A7', 'ความสมดุล การฟื้นตัว และการดูแลตัวเอง'),
        ('lilac', 'ม่วงไลแลค', 'Lilac', '#CE93D8', 'แรงบันดาลใจ จินตนาการ และการฟังเสียงข้างใน'),
        ('cream_white', 'ขาวครีม', 'Cream White', '#FFF8E1', 'ความเรียบง่าย ความชัดเจน และพื้นที่ให้เริ่มใหม่'),
        ('charcoal', 'เทาชาร์โคล', 'Charcoal', '#424242', 'ความหนักแน่น การตั้งขอบเขต และการตัดสินใจอย่างสุขุม')
)
INSERT INTO lucky_colors (code, name_th, name_en, hex, meaning)
SELECT code, name_th, name_en, hex, meaning
FROM color_rows
ON CONFLICT (code) DO UPDATE
SET name_th = EXCLUDED.name_th,
    name_en = EXCLUDED.name_en,
    hex = EXCLUDED.hex,
    meaning = EXCLUDED.meaning,
    is_active = true,
    updated_at = now();

WITH rule_rows (day_of_week, birth_day_of_week, purpose, color_code, rule_type, weight) AS (
    VALUES
        (0, NULL, 'general', 'sun_yellow', 'lucky', 10),
        (1, NULL, 'general', 'cream_white', 'lucky', 10),
        (2, NULL, 'general', 'soft_pink', 'lucky', 10),
        (3, NULL, 'general', 'mint_green', 'lucky', 10),
        (4, NULL, 'general', 'sky_blue', 'lucky', 10),
        (5, NULL, 'general', 'lilac', 'lucky', 10),
        (6, NULL, 'general', 'charcoal', 'lucky', 10),
        (0, NULL, 'general', 'charcoal', 'avoid', 5),
        (1, NULL, 'general', 'soft_pink', 'avoid', 5),
        (2, NULL, 'general', 'sky_blue', 'avoid', 5),
        (3, NULL, 'general', 'lilac', 'avoid', 5),
        (4, NULL, 'general', 'sun_yellow', 'avoid', 5),
        (5, NULL, 'general', 'mint_green', 'avoid', 5),
        (6, NULL, 'general', 'cream_white', 'avoid', 5),
        (1, 1, 'career', 'sky_blue', 'lucky', 20),
        (5, 5, 'love', 'soft_pink', 'lucky', 20),
        (3, 3, 'study', 'mint_green', 'lucky', 20),
        (4, 4, 'interview', 'cream_white', 'lucky', 20),
        (2, 2, 'money', 'sun_yellow', 'lucky', 20)
)
INSERT INTO lucky_color_rules (
    day_of_week, birth_day_of_week, purpose, color_id, rule_type, weight
)
SELECT rule_rows.day_of_week, rule_rows.birth_day_of_week, rule_rows.purpose,
       lucky_colors.id, rule_rows.rule_type, rule_rows.weight
FROM rule_rows
JOIN lucky_colors ON lucky_colors.code = rule_rows.color_code
WHERE lucky_colors.is_active = true
ON CONFLICT DO NOTHING;

WITH food_rows (code, name_th, name_en, category, tags, description) AS (
    VALUES
        ('jasmine_tea', 'ชามะลิ', 'Jasmine Tea', 'drink', ARRAY['calm', 'focus'], 'เครื่องดื่มหอมอ่อน ๆ สำหรับเริ่มวันอย่างนุ่มนวล'),
        ('fruit_bowl', 'ผลไม้รวม', 'Fruit Bowl', 'snack', ARRAY['fresh', 'light'], 'ตัวเลือกสดชื่นที่ช่วยให้วันดูเบาขึ้น'),
        ('rice_soup', 'ข้าวต้มอุ่น ๆ', 'Warm Rice Soup', 'meal', ARRAY['comfort', 'grounding'], 'อาหารเรียบง่ายที่เหมาะกับวันที่ต้องการความสบายใจ'),
        ('green_salad', 'สลัดผักเขียว', 'Green Salad', 'meal', ARRAY['balance', 'refresh'], 'เมนูเบา ๆ ที่ชวนให้กลับมาดูแลตัวเอง'),
        ('dark_chocolate', 'ดาร์กช็อกโกแลต', 'Dark Chocolate', 'snack', ARRAY['mood', 'reward'], 'ของหวานชิ้นเล็กสำหรับให้กำลังใจตัวเองอย่างพอดี')
)
INSERT INTO lucky_foods (code, name_th, name_en, category, tags, description)
SELECT code, name_th, name_en, category, tags, description
FROM food_rows
ON CONFLICT (code) DO UPDATE
SET name_th = EXCLUDED.name_th,
    name_en = EXCLUDED.name_en,
    category = EXCLUDED.category,
    tags = EXCLUDED.tags,
    description = EXCLUDED.description,
    is_active = true,
    updated_at = now();

WITH item_rows (code, name_th, name_en, category, tags, description) AS (
    VALUES
        ('small_notebook', 'สมุดโน้ตเล่มเล็ก', 'Small Notebook', 'stationery', ARRAY['focus', 'planning'], 'เหมาะสำหรับจดความคิดหรือสิ่งที่อยากตั้งใจในวันนี้'),
        ('silver_ring', 'แหวนสีเงิน', 'Silver Ring', 'accessory', ARRAY['clarity', 'boundary'], 'เครื่องประดับเรียบ ๆ ที่ช่วยเตือนให้ตั้งขอบเขตกับตัวเอง'),
        ('canvas_bag', 'กระเป๋าผ้า', 'Canvas Bag', 'daily', ARRAY['light', 'ready'], 'ของใช้ประจำวันที่ชวนให้เตรียมตัวอย่างเรียบง่าย'),
        ('mint_lip_balm', 'ลิปบาล์มกลิ่นมิ้นต์', 'Mint Lip Balm', 'self-care', ARRAY['fresh', 'gentle'], 'ของชิ้นเล็กสำหรับดูแลตัวเองระหว่างวัน'),
        ('blue_pen', 'ปากกาสีน้ำเงิน', 'Blue Pen', 'stationery', ARRAY['communication', 'work'], 'เหมาะกับวันที่ต้องเขียน คิด หรือคุยเรื่องสำคัญ')
)
INSERT INTO lucky_items (code, name_th, name_en, category, tags, description)
SELECT code, name_th, name_en, category, tags, description
FROM item_rows
ON CONFLICT (code) DO UPDATE
SET name_th = EXCLUDED.name_th,
    name_en = EXCLUDED.name_en,
    category = EXCLUDED.category,
    tags = EXCLUDED.tags,
    description = EXCLUDED.description,
    is_active = true,
    updated_at = now();

WITH avoidance_rows (code, category, text_th, text_en, mood_tag) AS (
    VALUES
        ('avoid_rushing_reply', 'communication', 'หลีกเลี่ยงการตอบทันทีตอนอารมณ์ยังไม่นิ่ง', 'Avoid replying immediately while emotions are unsettled.', 'calm'),
        ('avoid_overplanning', 'mindset', 'อย่าวางแผนแน่นเกินจนไม่มีพื้นที่หายใจ', 'Avoid planning so tightly that there is no room to breathe.', 'balance'),
        ('avoid_impulse_buy', 'money', 'ชะลอการซื้อของตามอารมณ์ แล้วกลับมาดูอีกครั้งภายหลัง', 'Pause impulse purchases and revisit them later.', 'grounding'),
        ('avoid_self_blame', 'self-care', 'หลีกเลี่ยงการโทษตัวเองกับเรื่องที่ยังต้องใช้เวลา', 'Avoid blaming yourself for things that still need time.', 'gentle'),
        ('avoid_multitask', 'focus', 'ลดการทำหลายอย่างพร้อมกัน ถ้าวันนี้ต้องการความชัดเจน', 'Avoid multitasking if today needs clarity.', 'focus')
)
INSERT INTO avoidance_recommendations (code, category, text_th, text_en, mood_tag)
SELECT code, category, text_th, text_en, mood_tag
FROM avoidance_rows
ON CONFLICT (code) DO UPDATE
SET category = EXCLUDED.category,
    text_th = EXCLUDED.text_th,
    text_en = EXCLUDED.text_en,
    mood_tag = EXCLUDED.mood_tag,
    is_active = true,
    updated_at = now();
