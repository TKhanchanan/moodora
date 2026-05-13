package tarot

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) ImportCards(ctx context.Context, provider CardProvider) (int, error) {
	cards, err := provider.LoadCards()
	if err != nil {
		return 0, err
	}
	if len(cards) != 78 {
		return 0, fmt.Errorf("expected 78 tarot cards, got %d", len(cards))
	}
	if err := s.repo.UpsertCards(ctx, cards); err != nil {
		return 0, err
	}
	return len(cards), nil
}

func (s *Service) ListCards(ctx context.Context) ([]Card, error) {
	return s.repo.ListCards(ctx)
}

func (s *Service) GetCard(ctx context.Context, sourceCode string) (Card, error) {
	return s.repo.GetCardBySourceCode(ctx, sourceCode)
}

func (s *Service) ListSpreads(ctx context.Context) ([]Spread, error) {
	return s.repo.ListSpreads(ctx)
}

func (s *Service) GetSpread(ctx context.Context, code string) (Spread, error) {
	return s.repo.GetSpreadByCode(ctx, code)
}

func (s *Service) CreateReading(ctx context.Context, req CreateReadingRequest) (ReadingResponse, error) {
	req.SpreadCode = strings.TrimSpace(req.SpreadCode)
	req.Topic = strings.TrimSpace(req.Topic)
	req.Language = strings.TrimSpace(req.Language)
	req.Question = strings.TrimSpace(req.Question)

	if req.SpreadCode == "" {
		return ReadingResponse{}, fmt.Errorf("spreadCode is required")
	}
	if req.Topic == "" {
		req.Topic = TopicGeneral
	}
	if req.Language == "" {
		req.Language = "en"
	}
	if !validTopics[req.Topic] {
		return ReadingResponse{}, fmt.Errorf("invalid topic")
	}
	if !validLanguages[req.Language] {
		return ReadingResponse{}, fmt.Errorf("invalid language")
	}

	spread, err := s.repo.GetSpreadByCode(ctx, req.SpreadCode)
	if err != nil {
		return ReadingResponse{}, err
	}
	if len(spread.Positions) == 0 {
		return ReadingResponse{}, fmt.Errorf("spread has no positions")
	}

	tx, err := s.repo.db.Begin(ctx)
	if err != nil {
		return ReadingResponse{}, fmt.Errorf("begin tarot reading transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	drawn, err := s.repo.drawCards(ctx, tx, len(spread.Positions))
	if err != nil {
		return ReadingResponse{}, err
	}
	if len(drawn) < len(spread.Positions) {
		return ReadingResponse{}, fmt.Errorf("not enough tarot cards to create reading")
	}

	readingCards := make([]ReadingCard, 0, len(spread.Positions))
	for i, position := range spread.Positions {
		card := drawn[i]
		orientation := OrientationUpright
		if req.AllowReversed {
			reversed, err := randomBool()
			if err != nil {
				return ReadingResponse{}, err
			}
			if reversed {
				orientation = OrientationReversed
			}
		}

		found, err := s.repo.findInterpretation(ctx, tx, card, req.Language, req.Topic, orientation)
		if err != nil {
			return ReadingResponse{}, err
		}

		readingCards = append(readingCards, ReadingCard{
			PositionNumber: position.PositionNumber,
			PositionCode:   position.Code,
			PositionName:   position.Name,
			Card: ReadingInfo{
				SourceCode: card.SourceCode,
				Name:       card.NameEn,
			},
			Orientation: orientation,
			Meaning:     found.Meaning,
			Advice:      found.Advice,
		})
	}

	summary := BuildSummary(readingCards)
	resultSnapshot, err := json.Marshal(map[string]any{
		"summary": summary,
		"cards":   readingCards,
	})
	if err != nil {
		return ReadingResponse{}, fmt.Errorf("marshal reading result snapshot: %w", err)
	}

	var response ReadingResponse
	err = tx.QueryRow(ctx, `
		INSERT INTO tarot_readings (spread_id, language, topic, question, status, result_snapshot)
		VALUES ($1, $2, $3, $4, 'completed', $5)
		RETURNING id::text, created_at
	`, spread.ID, req.Language, req.Topic, req.Question, resultSnapshot).Scan(&response.ID, &response.CreatedAt)
	if err != nil {
		return ReadingResponse{}, fmt.Errorf("insert tarot reading: %w", err)
	}

	for _, card := range readingCards {
		cardSnapshot, err := json.Marshal(card)
		if err != nil {
			return ReadingResponse{}, fmt.Errorf("marshal reading card snapshot: %w", err)
		}

		_, err = tx.Exec(ctx, `
			INSERT INTO tarot_reading_cards (
				reading_id, card_id, spread_position_id, position_number, orientation,
				meaning_snapshot, advice_snapshot, result_snapshot
			)
			SELECT $1, tc.id, tsp.id, $2, $3, $4, $5, $6
			FROM tarot_cards tc
			LEFT JOIN tarot_spread_positions tsp
				ON tsp.spread_id = $7 AND tsp.position_number = $2
			WHERE tc.source_code = $8
		`, response.ID, card.PositionNumber, card.Orientation, card.Meaning, card.Advice, cardSnapshot, spread.ID, card.Card.SourceCode)
		if err != nil {
			return ReadingResponse{}, fmt.Errorf("insert tarot reading card: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return ReadingResponse{}, fmt.Errorf("commit tarot reading transaction: %w", err)
	}

	response.SpreadCode = spread.Code
	response.Topic = req.Topic
	response.Language = req.Language
	response.Question = req.Question
	response.Cards = readingCards
	response.Summary = summary
	return response, nil
}

func (s *Service) GetReading(ctx context.Context, id string) (ReadingResponse, error) {
	return s.repo.GetReading(ctx, id)
}

func BuildSummary(cards []ReadingCard) string {
	if len(cards) == 0 {
		return "This reading is a quiet prompt for self-reflection."
	}

	parts := make([]string, 0, len(cards))
	for _, card := range cards {
		meaning := strings.TrimSpace(card.Meaning)
		if meaning == "" {
			continue
		}
		parts = append(parts, meaning)
		if len(parts) == 2 {
			break
		}
	}
	if len(parts) == 0 {
		return "This reading is a quiet prompt for self-reflection."
	}
	return "Reflect on this pattern: " + strings.Join(parts, " ")
}

func randomBool() (bool, error) {
	value, err := rand.Int(rand.Reader, big.NewInt(2))
	if err != nil {
		return false, fmt.Errorf("generate card orientation: %w", err)
	}
	return value.Int64() == 1, nil
}
