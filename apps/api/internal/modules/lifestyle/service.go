package lifestyle

import (
	"context"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"fmt"
	"sort"
	"strings"
	"time"
)

type Service struct {
	repo     *Repository
	location *time.Location
	now      func() time.Time
}

func NewService(repo *Repository, location *time.Location) *Service {
	return &Service{
		repo:     repo,
		location: location,
		now:      time.Now,
	}
}

func (s *Service) TodayColors(ctx context.Context, userID string, purpose string) (DailyInsight, error) {
	today, err := s.todayContext(ctx, userID, purpose)
	if err != nil {
		return DailyInsight{}, err
	}
	lucky, avoid, err := s.colors(ctx, today)
	if err != nil {
		return DailyInsight{}, err
	}
	return DailyInsight{Date: today.Date.Format("2006-01-02"), Timezone: today.Timezone, LuckyColors: lucky, AvoidColors: avoid}, nil
}

func (s *Service) TodayFoods(ctx context.Context, userID string) (DailyInsight, error) {
	today, err := s.todayContext(ctx, userID, PurposeGeneral)
	if err != nil {
		return DailyInsight{}, err
	}
	foods, err := s.foods(ctx, today)
	if err != nil {
		return DailyInsight{}, err
	}
	return DailyInsight{Date: today.Date.Format("2006-01-02"), Timezone: today.Timezone, LuckyFoods: foods}, nil
}

func (s *Service) TodayItems(ctx context.Context, userID string) (DailyInsight, error) {
	today, err := s.todayContext(ctx, userID, PurposeGeneral)
	if err != nil {
		return DailyInsight{}, err
	}
	items, err := s.items(ctx, today)
	if err != nil {
		return DailyInsight{}, err
	}
	return DailyInsight{Date: today.Date.Format("2006-01-02"), Timezone: today.Timezone, LuckyItems: items}, nil
}

func (s *Service) TodayAvoidances(ctx context.Context, userID string) (DailyInsight, error) {
	today, err := s.todayContext(ctx, userID, PurposeGeneral)
	if err != nil {
		return DailyInsight{}, err
	}
	avoidances, err := s.avoidances(ctx, today)
	if err != nil {
		return DailyInsight{}, err
	}
	return DailyInsight{Date: today.Date.Format("2006-01-02"), Timezone: today.Timezone, Avoidances: avoidances}, nil
}

func (s *Service) DailyInsight(ctx context.Context, userID string, purpose string) (DailyInsight, error) {
	today, err := s.todayContext(ctx, userID, purpose)
	if err != nil {
		return DailyInsight{}, err
	}

	date := today.Date.Format("2006-01-02")
	existing, err := s.repo.GetDailyInsight(ctx, userID, date, today.Timezone)
	if err == nil {
		return existing, nil
	}
	if !errors.Is(err, ErrNotFound) {
		return DailyInsight{}, err
	}

	luckyColors, avoidColors, err := s.colors(ctx, today)
	if err != nil {
		return DailyInsight{}, err
	}
	foods, err := s.foods(ctx, today)
	if err != nil {
		return DailyInsight{}, err
	}
	items, err := s.items(ctx, today)
	if err != nil {
		return DailyInsight{}, err
	}
	avoidances, err := s.avoidances(ctx, today)
	if err != nil {
		return DailyInsight{}, err
	}
	walletBalance, err := s.repo.WalletBalance(ctx, userID)
	if err != nil {
		return DailyInsight{}, err
	}
	checkInStatus, err := s.repo.CheckInStatus(ctx, userID, date)
	if err != nil {
		return DailyInsight{}, err
	}

	insight := DailyInsight{
		Date:          date,
		Timezone:      today.Timezone,
		LuckyColors:   luckyColors,
		AvoidColors:   avoidColors,
		LuckyFoods:    foods,
		LuckyItems:    items,
		Avoidances:    avoidances,
		WalletBalance: walletBalance,
		CheckInStatus: checkInStatus,
		DailyTarot: DailyTarot{
			Status:  "placeholder",
			Message: "Daily tarot can be connected after a product rule chooses the spread and timing.",
		},
		Message: "Use these as lifestyle prompts for reflection, not guaranteed outcomes.",
	}
	if err := s.repo.SaveDailyInsight(ctx, userID, insight); err != nil {
		return DailyInsight{}, err
	}
	return insight, nil
}

func (s *Service) todayContext(ctx context.Context, userID string, purpose string) (TodayContext, error) {
	if purpose == "" {
		purpose = PurposeGeneral
	}
	purpose = strings.TrimSpace(purpose)
	if !validPurposes[purpose] {
		return TodayContext{}, fmt.Errorf("invalid purpose")
	}

	today := LocalToday(s.now(), s.location)
	birthDay, err := s.repo.BirthDayOfWeek(ctx, userID)
	if err != nil {
		return TodayContext{}, err
	}
	day := int(today.Weekday())
	if birthDay == nil {
		birthDay = &day
	}

	return TodayContext{
		UserID:         userID,
		Date:           today,
		Timezone:       s.location.String(),
		DayOfWeek:      day,
		BirthDayOfWeek: birthDay,
		Purpose:        purpose,
	}, nil
}

func (s *Service) colors(ctx context.Context, today TodayContext) ([]Color, []Color, error) {
	lucky, err := s.repo.ListColorsByRule(ctx, today.DayOfWeek, today.BirthDayOfWeek, today.Purpose, RuleTypeLucky)
	if err != nil {
		return nil, nil, err
	}
	if len(lucky) == 0 && today.Purpose != PurposeGeneral {
		lucky, err = s.repo.ListColorsByRule(ctx, today.DayOfWeek, today.BirthDayOfWeek, PurposeGeneral, RuleTypeLucky)
		if err != nil {
			return nil, nil, err
		}
	}

	avoid, err := s.repo.ListColorsByRule(ctx, today.DayOfWeek, today.BirthDayOfWeek, today.Purpose, RuleTypeAvoid)
	if err != nil {
		return nil, nil, err
	}
	if len(avoid) == 0 && today.Purpose != PurposeGeneral {
		avoid, err = s.repo.ListColorsByRule(ctx, today.DayOfWeek, today.BirthDayOfWeek, PurposeGeneral, RuleTypeAvoid)
		if err != nil {
			return nil, nil, err
		}
	}
	return lucky, avoid, nil
}

func (s *Service) foods(ctx context.Context, today TodayContext) ([]Food, error) {
	foods, err := s.repo.ListFoods(ctx)
	if err != nil {
		return nil, err
	}
	return pickDeterministic(foods, today.UserID+"|"+today.Date.Format("2006-01-02")+"|foods", 2), nil
}

func (s *Service) items(ctx context.Context, today TodayContext) ([]Item, error) {
	items, err := s.repo.ListItems(ctx)
	if err != nil {
		return nil, err
	}
	return pickDeterministic(items, today.UserID+"|"+today.Date.Format("2006-01-02")+"|items", 2), nil
}

func (s *Service) avoidances(ctx context.Context, today TodayContext) ([]Avoidance, error) {
	avoidances, err := s.repo.ListAvoidances(ctx)
	if err != nil {
		return nil, err
	}
	return pickDeterministic(avoidances, today.UserID+"|"+today.Date.Format("2006-01-02")+"|avoidances", 2), nil
}

func pickDeterministic[T any](values []T, seed string, limit int) []T {
	if limit <= 0 || len(values) == 0 {
		return []T{}
	}
	type ranked struct {
		index int
		rank  uint64
	}
	rankedValues := make([]ranked, len(values))
	for i := range values {
		sum := sha256.Sum256([]byte(fmt.Sprintf("%s|%d", seed, i)))
		rankedValues[i] = ranked{index: i, rank: binary.BigEndian.Uint64(sum[:8])}
	}
	sort.SliceStable(rankedValues, func(i, j int) bool {
		return rankedValues[i].rank < rankedValues[j].rank
	})

	if limit > len(values) {
		limit = len(values)
	}
	selected := make([]T, 0, limit)
	for i := 0; i < limit; i++ {
		selected = append(selected, values[rankedValues[i].index])
	}
	return selected
}
