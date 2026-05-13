package wallet

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var ErrInsufficientBalance = errors.New("insufficient wallet balance")

type Service struct {
	db *pgxpool.Pool
}

type Wallet struct {
	ID          string    `json:"id"`
	UserID      string    `json:"userId"`
	CoinBalance int64     `json:"coinBalance"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type Transaction struct {
	ID              string    `json:"id"`
	TransactionType string    `json:"transactionType"`
	Amount          int64     `json:"amount"`
	BalanceAfter    int64     `json:"balanceAfter"`
	Reason          string    `json:"reason"`
	CreatedAt       time.Time `json:"createdAt"`
}

func NewService(db *pgxpool.Pool) *Service {
	return &Service{db: db}
}

func (s *Service) GetWallet(ctx context.Context, userID string) (Wallet, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return Wallet{}, err
	}
	defer tx.Rollback(ctx)

	wallet, err := ensureWallet(ctx, tx, userID)
	if err != nil {
		return Wallet{}, err
	}
	if err := tx.Commit(ctx); err != nil {
		return Wallet{}, err
	}
	return wallet, nil
}

func (s *Service) ListTransactions(ctx context.Context, userID string) ([]Transaction, error) {
	rows, err := s.db.Query(ctx, `
		SELECT id::text, transaction_type, amount, balance_after, reason, created_at
		FROM coin_transactions
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT 100
	`, userID)
	if err != nil {
		return nil, fmt.Errorf("list coin transactions: %w", err)
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var transaction Transaction
		if err := rows.Scan(&transaction.ID, &transaction.TransactionType, &transaction.Amount, &transaction.BalanceAfter, &transaction.Reason, &transaction.CreatedAt); err != nil {
			return nil, fmt.Errorf("scan coin transaction: %w", err)
		}
		transactions = append(transactions, transaction)
	}
	return transactions, rows.Err()
}

func (s *Service) Grant(ctx context.Context, tx pgx.Tx, userID string, amount int64, reason string, idempotencyKey string) (Wallet, bool, error) {
	if amount <= 0 {
		return Wallet{}, false, fmt.Errorf("grant amount must be positive")
	}
	wallet, err := ensureWallet(ctx, tx, userID)
	if err != nil {
		return Wallet{}, false, err
	}

	if idempotencyKey != "" {
		var existingBalance int64
		err := tx.QueryRow(ctx, `
			SELECT balance_after
			FROM coin_transactions
			WHERE wallet_id = $1 AND idempotency_key = $2
		`, wallet.ID, idempotencyKey).Scan(&existingBalance)
		if err == nil {
			wallet.CoinBalance = existingBalance
			return wallet, false, nil
		}
		if !errors.Is(err, pgx.ErrNoRows) {
			return Wallet{}, false, err
		}
	}

	newBalance := wallet.CoinBalance + amount
	_, err = tx.Exec(ctx, `
		UPDATE wallets
		SET coin_balance = $1
		WHERE id = $2
	`, newBalance, wallet.ID)
	if err != nil {
		return Wallet{}, false, fmt.Errorf("update wallet balance: %w", err)
	}

	_, err = tx.Exec(ctx, `
		INSERT INTO coin_transactions (
			wallet_id, user_id, transaction_type, amount, balance_after, reason, idempotency_key
		)
		VALUES ($1, $2, 'check_in', $3, $4, $5, $6)
	`, wallet.ID, userID, amount, newBalance, reason, idempotencyKey)
	if err != nil {
		return Wallet{}, false, fmt.Errorf("insert coin transaction: %w", err)
	}

	wallet.CoinBalance = newBalance
	return wallet, true, nil
}

func ensureWallet(ctx context.Context, tx pgx.Tx, userID string) (Wallet, error) {
	_, err := tx.Exec(ctx, `
		INSERT INTO wallets (user_id, coin_balance)
		VALUES ($1, 0)
		ON CONFLICT (user_id) DO NOTHING
	`, userID)
	if err != nil {
		return Wallet{}, fmt.Errorf("ensure wallet: %w", err)
	}

	var wallet Wallet
	err = tx.QueryRow(ctx, `
		SELECT id::text, user_id::text, coin_balance, created_at, updated_at
		FROM wallets
		WHERE user_id = $1
		FOR UPDATE
	`, userID).Scan(&wallet.ID, &wallet.UserID, &wallet.CoinBalance, &wallet.CreatedAt, &wallet.UpdatedAt)
	if err != nil {
		return Wallet{}, fmt.Errorf("get wallet: %w", err)
	}
	return wallet, nil
}
