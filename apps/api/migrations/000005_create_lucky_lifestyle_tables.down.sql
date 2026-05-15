DROP TRIGGER IF EXISTS daily_insights_set_updated_at ON daily_insights;
DROP TRIGGER IF EXISTS avoidance_recommendations_set_updated_at ON avoidance_recommendations;
DROP TRIGGER IF EXISTS lucky_items_set_updated_at ON lucky_items;
DROP TRIGGER IF EXISTS lucky_foods_set_updated_at ON lucky_foods;
DROP TRIGGER IF EXISTS lucky_color_rules_set_updated_at ON lucky_color_rules;
DROP TRIGGER IF EXISTS lucky_colors_set_updated_at ON lucky_colors;

DROP TABLE IF EXISTS daily_insights;
DROP TABLE IF EXISTS avoidance_recommendations;
DROP TABLE IF EXISTS lucky_items;
DROP TABLE IF EXISTS lucky_foods;
DROP TABLE IF EXISTS lucky_color_rules;
DROP TABLE IF EXISTS lucky_colors;
