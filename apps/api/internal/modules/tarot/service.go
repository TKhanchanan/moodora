package tarot

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"math/big"
	"strings"

	"github.com/jackc/pgx/v5"
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
	for i := range req.SelectedCardSourceCodes {
		req.SelectedCardSourceCodes[i] = strings.TrimSpace(req.SelectedCardSourceCodes[i])
	}

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

	drawn, err := s.selectReadingCards(ctx, tx, req, len(spread.Positions))
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
				SourceCode:     card.SourceCode,
				Name:           displayCardName(card, req.Language),
				NameEn:         card.NameEn,
				NameTh:         card.NameTh,
				Type:           card.Type,
				Suit:           card.Suit,
				Characteristic: cardCharacteristic(card, req.Language),
				Description:    displayCardDescription(card, req.Language),
				Assets:         card.Assets,
			},
			Orientation: orientation,
			Meaning:     found.Meaning,
			Advice:      found.Advice,
		})
	}

	summary := BuildSummary(req.Language, readingCards)
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

func (s *Service) selectReadingCards(ctx context.Context, tx pgx.Tx, req CreateReadingRequest, cardCount int) ([]Card, error) {
	if len(req.SelectedCardSourceCodes) == 0 {
		return s.repo.drawCards(ctx, tx, cardCount)
	}
	if len(req.SelectedCardSourceCodes) != cardCount {
		return nil, fmt.Errorf("selectedCardSourceCodes must contain exactly %d cards", cardCount)
	}

	seen := map[string]bool{}
	for _, sourceCode := range req.SelectedCardSourceCodes {
		if sourceCode == "" {
			return nil, fmt.Errorf("selectedCardSourceCodes contains an empty card")
		}
		if seen[sourceCode] {
			return nil, fmt.Errorf("selectedCardSourceCodes contains duplicate card %q", sourceCode)
		}
		seen[sourceCode] = true
	}
	return s.repo.findCardsBySourceCodes(ctx, tx, req.SelectedCardSourceCodes)
}

func displayCardName(card Card, language string) string {
	if language == "th" && card.NameTh != "" {
		return card.NameTh
	}
	return card.NameEn
}

func displayCardDescription(card Card, language string) string {
	if language == "th" && card.DescriptionTh != "" {
		return card.DescriptionTh
	}
	return card.DescriptionEn
}

func cardCharacteristic(card Card, language string) string {
	if language == "th" {
		if card.Type == "major" {
			return "เมเจอร์อาร์คานา: ไพ่บทเรียนสำคัญและจังหวะเปลี่ยนผ่านของชีวิต"
		}
		return "ไมเนอร์อาร์คานา ชุด" + suitLabel(card.Suit, language) + ": ไพ่สถานการณ์ประจำวันและพลังที่จับต้องได้"
	}

	if card.Type == "major" {
		return "Major Arcana: a broad life theme, lesson, or turning point for reflection"
	}
	return "Minor Arcana, suit of " + suitLabel(card.Suit, language) + ": a practical day-to-day influence"
}

func suitLabel(suit *string, language string) string {
	if suit == nil {
		if language == "th" {
			return "ไม่มีชุด"
		}
		return "none"
	}
	if language != "th" {
		return *suit
	}
	switch *suit {
	case "wands":
		return "ไม้เท้า"
	case "cups":
		return "ถ้วย"
	case "swords":
		return "ดาบ"
	case "pentacles":
		return "เหรียญ"
	default:
		return *suit
	}
}

func (s *Service) GetReading(ctx context.Context, id string) (ReadingResponse, error) {
	return s.repo.GetReading(ctx, id)
}

func BuildSummary(language string, cards []ReadingCard) string {
	fallback := "This reading is a quiet prompt for self-reflection."
	prefix := "Reflect on this pattern: "

	if language == "th" {
		fallback = "คำทำนายนี้เป็นข้อความเตือนใจให้คุณได้ใช้เวลาทบทวนตัวเองอย่างสงบ"
		prefix = "ภาพรวมคำทำนายของคุณ: "
	}

	if len(cards) == 0 {
		return fallback
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
		return fallback
	}
	return prefix + strings.Join(parts, " ")
}

func randomBool() (bool, error) {
	value, err := rand.Int(rand.Reader, big.NewInt(2))
	if err != nil {
		return false, fmt.Errorf("generate card orientation: %w", err)
	}
	return value.Int64() == 1, nil
}
