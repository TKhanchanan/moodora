package moon

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Service struct {
	repo       *Repository
	location   *time.Location
	calculator Calculator
	now        func() time.Time
}

func NewService(repo *Repository, location *time.Location) *Service {
	return &Service{
		repo:       repo,
		location:   location,
		calculator: Calculator{},
		now:        time.Now,
	}
}

func (s *Service) Today(ctx context.Context, userID string) (ReportResponse, error) {
	result := s.calculator.Calculate(s.now(), s.location)
	interpretation := BuildInterpretation(result.MoonPhase)
	return s.repo.SaveReport(ctx, userID, nil, result, interpretation, InternalCalculationSource())
}

func (s *Service) Birthday(ctx context.Context, userID string, req BirthdayRequest) (ReportResponse, error) {
	timezone := strings.TrimSpace(req.Timezone)
	location := s.location
	if timezone != "" {
		loaded, err := time.LoadLocation(timezone)
		if err != nil {
			return ReportResponse{}, fmt.Errorf("invalid timezone")
		}
		location = loaded
	}
	if location == nil {
		location = time.UTC
	}

	birthDate := strings.TrimSpace(req.BirthDate)
	if birthDate == "" {
		return ReportResponse{}, fmt.Errorf("birthDate is required")
	}
	target, err := ParseDateInLocation(birthDate, location)
	if err != nil {
		return ReportResponse{}, err
	}

	result := s.calculator.Calculate(target, location)
	interpretation := BuildInterpretation(result.MoonPhase)
	return s.repo.SaveReport(ctx, userID, &birthDate, result, interpretation, InternalCalculationSource())
}

func (s *Service) GetReport(ctx context.Context, id string) (ReportResponse, error) {
	return s.repo.GetReport(ctx, id)
}
