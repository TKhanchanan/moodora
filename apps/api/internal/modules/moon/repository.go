package moon

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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

func (r *Repository) SaveReport(ctx context.Context, userID string, birthDate *string, result PhaseResult, interpretation Interpretation, source Source) (ReportResponse, error) {
	interpretationSnapshot, err := json.Marshal(interpretation)
	if err != nil {
		return ReportResponse{}, fmt.Errorf("marshal interpretation snapshot: %w", err)
	}
	sourceSnapshot, err := json.Marshal(source)
	if err != nil {
		return ReportResponse{}, fmt.Errorf("marshal source snapshot: %w", err)
	}

	resultSnapshotBody := map[string]any{
		"targetDate":               result.TargetDate,
		"timezone":                 result.Timezone,
		"moonPhase":                result.MoonPhase,
		"illumination":             result.Illumination,
		"moonAge":                  result.MoonAge,
		"interpretation":           interpretation,
		"calculationMethodVersion": result.CalculationMethodVersion,
		"source":                   source,
	}
	resultSnapshot, err := json.Marshal(resultSnapshotBody)
	if err != nil {
		return ReportResponse{}, fmt.Errorf("marshal result snapshot: %w", err)
	}

	var user any
	if userID != "" {
		user = userID
	}
	var birth any
	if birthDate != nil && *birthDate != "" {
		birth = *birthDate
	}

	response := ReportResponse{
		TargetDate:               result.TargetDate,
		Timezone:                 result.Timezone,
		MoonPhase:                result.MoonPhase,
		Illumination:             result.Illumination,
		MoonAge:                  result.MoonAge,
		Interpretation:           interpretation,
		CalculationMethodVersion: result.CalculationMethodVersion,
		Source:                   source,
	}

	err = r.db.QueryRow(ctx, `
		INSERT INTO moon_reports (
			user_id, birth_date, target_date, timezone, moon_phase, illumination, moon_age,
			calculation_method_version, interpretation_snapshot, source_snapshot, result_snapshot
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		RETURNING id::text, image_url, created_at
	`, user, birth, result.TargetDate, result.Timezone, result.MoonPhase, result.Illumination, result.MoonAge,
		result.CalculationMethodVersion, interpretationSnapshot, sourceSnapshot, resultSnapshot).
		Scan(&response.ID, &response.ImageURL, &response.CreatedAt)
	if err != nil {
		return ReportResponse{}, fmt.Errorf("save moon report: %w", err)
	}

	return response, nil
}

func (r *Repository) GetReport(ctx context.Context, id string) (ReportResponse, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id::text, target_date::text, timezone, moon_phase, illumination::float8, moon_age::float8,
			image_url, calculation_method_version, interpretation_snapshot, source_snapshot, created_at
		FROM moon_reports
		WHERE id = $1
	`, id)

	var response ReportResponse
	var interpretationRaw []byte
	var sourceRaw []byte
	if err := row.Scan(
		&response.ID,
		&response.TargetDate,
		&response.Timezone,
		&response.MoonPhase,
		&response.Illumination,
		&response.MoonAge,
		&response.ImageURL,
		&response.CalculationMethodVersion,
		&interpretationRaw,
		&sourceRaw,
		&response.CreatedAt,
	); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ReportResponse{}, ErrNotFound
		}
		return ReportResponse{}, fmt.Errorf("get moon report: %w", err)
	}

	if err := json.Unmarshal(interpretationRaw, &response.Interpretation); err != nil {
		return ReportResponse{}, fmt.Errorf("decode interpretation snapshot: %w", err)
	}
	if err := json.Unmarshal(sourceRaw, &response.Source); err != nil {
		return ReportResponse{}, fmt.Errorf("decode source snapshot: %w", err)
	}
	return response, nil
}

func ParseDateInLocation(value string, location *time.Location) (time.Time, error) {
	parsed, err := time.ParseInLocation("2006-01-02", value, location)
	if err != nil {
		return time.Time{}, fmt.Errorf("date must use YYYY-MM-DD")
	}
	return parsed, nil
}
