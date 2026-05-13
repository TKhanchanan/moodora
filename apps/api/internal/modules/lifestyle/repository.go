package lifestyle

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrNotFound = errors.New("not found")

type Repository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) *Repository {
	return &Repository{db: db}
}

func (r *Repository) ListColorsByRule(ctx context.Context, dayOfWeek int, birthDayOfWeek *int, purpose string, ruleType string) ([]Color, error) {
	birth := any(nil)
	if birthDayOfWeek != nil {
		birth = *birthDayOfWeek
	}

	rows, err := r.db.Query(ctx, `
		SELECT lc.id::text, lc.code, lc.name_th, lc.name_en, lc.hex, lc.meaning
		FROM lucky_color_rules lcr
		JOIN lucky_colors lc ON lc.id = lcr.color_id
		WHERE lcr.is_active = true
			AND lc.is_active = true
			AND lcr.day_of_week = $1
			AND lcr.purpose = $2
			AND lcr.rule_type = $3
			AND (lcr.birth_day_of_week IS NULL OR lcr.birth_day_of_week = $4)
		ORDER BY
			CASE WHEN lcr.birth_day_of_week = $4 THEN 0 ELSE 1 END,
			lcr.weight DESC,
			lc.code
		LIMIT 3
	`, dayOfWeek, purpose, ruleType, birth)
	if err != nil {
		return nil, fmt.Errorf("list lucky color rules: %w", err)
	}
	defer rows.Close()

	var colors []Color
	for rows.Next() {
		var color Color
		if err := rows.Scan(&color.ID, &color.Code, &color.NameTH, &color.NameEN, &color.Hex, &color.Meaning); err != nil {
			return nil, fmt.Errorf("scan lucky color: %w", err)
		}
		colors = append(colors, color)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return colors, nil
}

func (r *Repository) ListFoods(ctx context.Context) ([]Food, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id::text, code, name_th, name_en, category, array_to_string(tags, ','), description
		FROM lucky_foods
		WHERE is_active = true
		ORDER BY code
	`)
	if err != nil {
		return nil, fmt.Errorf("list lucky foods: %w", err)
	}
	defer rows.Close()

	var foods []Food
	for rows.Next() {
		var food Food
		var tags string
		if err := rows.Scan(&food.ID, &food.Code, &food.NameTH, &food.NameEN, &food.Category, &tags, &food.Description); err != nil {
			return nil, fmt.Errorf("scan lucky food: %w", err)
		}
		food.Tags = splitTags(tags)
		foods = append(foods, food)
	}
	return foods, rows.Err()
}

func (r *Repository) ListItems(ctx context.Context) ([]Item, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id::text, code, name_th, name_en, category, array_to_string(tags, ','), description
		FROM lucky_items
		WHERE is_active = true
		ORDER BY code
	`)
	if err != nil {
		return nil, fmt.Errorf("list lucky items: %w", err)
	}
	defer rows.Close()

	var items []Item
	for rows.Next() {
		var item Item
		var tags string
		if err := rows.Scan(&item.ID, &item.Code, &item.NameTH, &item.NameEN, &item.Category, &tags, &item.Description); err != nil {
			return nil, fmt.Errorf("scan lucky item: %w", err)
		}
		item.Tags = splitTags(tags)
		items = append(items, item)
	}
	return items, rows.Err()
}

func (r *Repository) ListAvoidances(ctx context.Context) ([]Avoidance, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id::text, code, category, text_th, text_en, mood_tag
		FROM avoidance_recommendations
		WHERE is_active = true
		ORDER BY code
	`)
	if err != nil {
		return nil, fmt.Errorf("list avoidances: %w", err)
	}
	defer rows.Close()

	var avoidances []Avoidance
	for rows.Next() {
		var avoidance Avoidance
		if err := rows.Scan(&avoidance.ID, &avoidance.Code, &avoidance.Category, &avoidance.TextTH, &avoidance.TextEN, &avoidance.MoodTag); err != nil {
			return nil, fmt.Errorf("scan avoidance: %w", err)
		}
		avoidances = append(avoidances, avoidance)
	}
	return avoidances, rows.Err()
}

func (r *Repository) BirthDayOfWeek(ctx context.Context, userID string) (*int, error) {
	if userID == "" {
		return nil, nil
	}

	var day int
	err := r.db.QueryRow(ctx, `
		SELECT EXTRACT(DOW FROM birth_date)::int
		FROM user_profiles
		WHERE user_id = $1 AND birth_date IS NOT NULL
	`, userID).Scan(&day)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get user birth day: %w", err)
	}
	return &day, nil
}

func (r *Repository) WalletBalance(ctx context.Context, userID string) (*int64, error) {
	if userID == "" {
		return nil, nil
	}

	var balance int64
	err := r.db.QueryRow(ctx, `
		SELECT coin_balance
		FROM wallets
		WHERE user_id = $1
	`, userID).Scan(&balance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, nil
		}
		return nil, fmt.Errorf("get wallet balance: %w", err)
	}
	return &balance, nil
}

func (r *Repository) CheckInStatus(ctx context.Context, userID string, localDate string) (*CheckInStatus, error) {
	if userID == "" {
		return nil, nil
	}

	var exists bool
	err := r.db.QueryRow(ctx, `
		SELECT EXISTS (
			SELECT 1
			FROM check_ins
			WHERE user_id = $1 AND local_date = $2
		)
	`, userID, localDate).Scan(&exists)
	if err != nil {
		return nil, fmt.Errorf("get check-in status: %w", err)
	}
	return &CheckInStatus{CheckedIn: exists, LocalDate: localDate}, nil
}

func (r *Repository) GetDailyInsight(ctx context.Context, userID string, date string, timezone string) (DailyInsight, error) {
	var row pgx.Row
	if userID == "" {
		row = r.db.QueryRow(ctx, `
			SELECT result_snapshot
			FROM daily_insights
			WHERE user_id IS NULL AND insight_date = $1 AND timezone = $2
		`, date, timezone)
	} else {
		row = r.db.QueryRow(ctx, `
			SELECT result_snapshot
			FROM daily_insights
			WHERE user_id = $1 AND insight_date = $2 AND timezone = $3
		`, userID, date, timezone)
	}

	var raw []byte
	if err := row.Scan(&raw); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return DailyInsight{}, ErrNotFound
		}
		return DailyInsight{}, fmt.Errorf("get daily insight: %w", err)
	}

	var insight DailyInsight
	if err := json.Unmarshal(raw, &insight); err != nil {
		return DailyInsight{}, fmt.Errorf("decode daily insight snapshot: %w", err)
	}
	return insight, nil
}

func (r *Repository) SaveDailyInsight(ctx context.Context, userID string, insight DailyInsight) error {
	raw, err := json.Marshal(insight)
	if err != nil {
		return fmt.Errorf("marshal daily insight snapshot: %w", err)
	}

	if userID == "" {
		_, err = r.db.Exec(ctx, `
			INSERT INTO daily_insights (user_id, insight_date, timezone, result_snapshot)
			VALUES (NULL, $1, $2, $3)
			ON CONFLICT (insight_date, timezone) WHERE user_id IS NULL DO UPDATE
			SET result_snapshot = daily_insights.result_snapshot
		`, insight.Date, insight.Timezone, raw)
	} else {
		_, err = r.db.Exec(ctx, `
			INSERT INTO daily_insights (user_id, insight_date, timezone, result_snapshot)
			VALUES ($1, $2, $3, $4)
			ON CONFLICT (user_id, insight_date, timezone) WHERE user_id IS NOT NULL DO UPDATE
			SET result_snapshot = daily_insights.result_snapshot
		`, userID, insight.Date, insight.Timezone, raw)
	}
	if err != nil {
		return fmt.Errorf("save daily insight: %w", err)
	}
	return nil
}

func splitTags(value string) []string {
	if value == "" {
		return []string{}
	}
	parts := strings.Split(value, ",")
	tags := make([]string, 0, len(parts))
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if part != "" {
			tags = append(tags, part)
		}
	}
	return tags
}

func LocalToday(now time.Time, location *time.Location) time.Time {
	local := now.In(location)
	return time.Date(local.Year(), local.Month(), local.Day(), 0, 0, 0, 0, location)
}
