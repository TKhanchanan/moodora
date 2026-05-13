package checkin

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/moodora/moodora/apps/api/internal/modules/wallet"
)

type Service struct {
	db       *pgxpool.Pool
	wallets  *wallet.Service
	location *time.Location
	now      func() time.Time
}

type Response struct {
	CheckedIn      bool   `json:"checkedIn"`
	RewardCoins    int64  `json:"rewardCoins"`
	StreakDay      int    `json:"streakDay"`
	WalletBalance  int64  `json:"walletBalance"`
	Timezone       string `json:"timezone"`
	NextCheckInAt  string `json:"nextCheckInAt"`
	AlreadyChecked bool   `json:"alreadyChecked"`
}

func NewService(db *pgxpool.Pool, wallets *wallet.Service, location *time.Location) *Service {
	return &Service{
		db:       db,
		wallets:  wallets,
		location: location,
		now:      time.Now,
	}
}

func RewardForStreakDay(streakDay int) int64 {
	rewards := []int64{5, 5, 8, 8, 10, 15, 25}
	if streakDay <= 0 {
		return rewards[0]
	}
	if streakDay > len(rewards) {
		return rewards[len(rewards)-1]
	}
	return rewards[streakDay-1]
}

func (s *Service) CheckIn(ctx context.Context, userID string) (Response, error) {
	now := s.now().In(s.location)
	localDate := now.Format("2006-01-02")
	nextCheckIn := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, s.location)

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return Response{}, err
	}
	defer tx.Rollback(ctx)

	streakDay, err := s.nextStreakDay(ctx, tx, userID, now)
	if err != nil {
		return Response{}, err
	}

	var existingReward int64
	err = tx.QueryRow(ctx, `
		SELECT reward_coins
		FROM check_ins
		WHERE user_id = $1 AND local_date = $2
	`, userID, localDate).Scan(&existingReward)
	if err == nil {
		var walletBalance int64
		if walletErr := tx.QueryRow(ctx, `
			SELECT coin_balance
			FROM wallets
			WHERE user_id = $1
		`, userID).Scan(&walletBalance); walletErr != nil {
			return Response{}, walletErr
		}
		return Response{
			CheckedIn:      false,
			RewardCoins:    0,
			StreakDay:      streakDay,
			WalletBalance:  walletBalance,
			Timezone:       s.location.String(),
			NextCheckInAt:  nextCheckIn.Format(time.RFC3339),
			AlreadyChecked: true,
		}, tx.Commit(ctx)
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return Response{}, err
	}

	reward := RewardForStreakDay(streakDay)
	idempotencyKey := fmt.Sprintf("check_in:%s:%s", userID, localDate)
	updatedWallet, granted, err := s.wallets.Grant(ctx, tx, userID, reward, "daily check-in", idempotencyKey)
	if err != nil {
		return Response{}, err
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO check_ins (user_id, local_date, timezone, reward_coins)
		VALUES ($1, $2, $3, $4)
		ON CONFLICT (user_id, local_date) DO NOTHING
	`, userID, localDate, s.location.String(), reward)
	if err != nil {
		return Response{}, fmt.Errorf("insert check-in: %w", err)
	}

	if err := tx.Commit(ctx); err != nil {
		return Response{}, err
	}

	return Response{
		CheckedIn:      granted,
		RewardCoins:    reward,
		StreakDay:      streakDay,
		WalletBalance:  updatedWallet.CoinBalance,
		Timezone:       s.location.String(),
		NextCheckInAt:  nextCheckIn.Format(time.RFC3339),
		AlreadyChecked: false,
	}, nil
}

func (s *Service) nextStreakDay(ctx context.Context, tx pgx.Tx, userID string, now time.Time) (int, error) {
	yesterday := now.AddDate(0, 0, -1)
	rows, err := tx.Query(ctx, `
		SELECT local_date
		FROM check_ins
		WHERE user_id = $1 AND local_date <= $2
		ORDER BY local_date DESC
		LIMIT 7
	`, userID, yesterday.Format("2006-01-02"))
	if err != nil {
		return 0, err
	}
	defer rows.Close()

	streak := 1
	expected := yesterday
	for rows.Next() {
		var localDate time.Time
		if err := rows.Scan(&localDate); err != nil {
			return 0, err
		}
		if localDate.Format("2006-01-02") != expected.Format("2006-01-02") {
			break
		}
		streak++
		expected = expected.AddDate(0, 0, -1)
	}
	return streak, rows.Err()
}
