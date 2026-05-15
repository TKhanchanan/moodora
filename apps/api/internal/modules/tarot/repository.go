package tarot

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

func (r *Repository) UpsertCards(ctx context.Context, cards []SourceCard) error {
	for _, card := range cards {
		var suit any
		if card.Suit != "" {
			suit = card.Suit
		}

		_, err := r.db.Exec(ctx, `
			INSERT INTO tarot_cards (
				source_code, name_en, type, suit, meaning_up_en, meaning_rev_en, description_en
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7)
			ON CONFLICT (source_code) DO UPDATE
			SET name_en = EXCLUDED.name_en,
				type = EXCLUDED.type,
				suit = EXCLUDED.suit,
				meaning_up_en = EXCLUDED.meaning_up_en,
				meaning_rev_en = EXCLUDED.meaning_rev_en,
				description_en = EXCLUDED.description_en,
				updated_at = now()
		`, card.SourceCode, card.NameEn, card.Type, suit, card.MeaningUpEn, card.MeaningRevEn, card.DescriptionEn)
		if err != nil {
			return fmt.Errorf("upsert tarot card %s: %w", card.SourceCode, err)
		}
	}

	return nil
}

func (r *Repository) ListCards(ctx context.Context) ([]Card, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id::text, source_code, name_en, type, suit, meaning_up_en, meaning_rev_en, description_en
		FROM tarot_cards
		ORDER BY source_code
	`)
	if err != nil {
		return nil, fmt.Errorf("list tarot cards: %w", err)
	}
	defer rows.Close()

	cards, err := scanCards(rows)
	if err != nil {
		return nil, err
	}
	return r.attachCardDetails(ctx, r.db, cards)
}

func (r *Repository) GetCardBySourceCode(ctx context.Context, sourceCode string) (Card, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id::text, source_code, name_en, type, suit, meaning_up_en, meaning_rev_en, description_en
		FROM tarot_cards
		WHERE source_code = $1
	`, sourceCode)

	card, err := scanCard(row)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Card{}, ErrNotFound
		}
		return Card{}, fmt.Errorf("get tarot card: %w", err)
	}

	cards, err := r.attachAssets(ctx, []Card{card})
	if err != nil {
		return Card{}, err
	}
	cards, err = attachTranslations(ctx, r.db, cards)
	if err != nil {
		return Card{}, err
	}
	return cards[0], nil
}

func (r *Repository) ListSpreads(ctx context.Context) ([]Spread, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id::text, code, name, description, card_count
		FROM tarot_spreads
		WHERE is_active = true
		ORDER BY card_count, code
	`)
	if err != nil {
		return nil, fmt.Errorf("list tarot spreads: %w", err)
	}
	defer rows.Close()

	var spreads []Spread
	for rows.Next() {
		var spread Spread
		if err := rows.Scan(&spread.ID, &spread.Code, &spread.Name, &spread.Description, &spread.CardCount); err != nil {
			return nil, fmt.Errorf("scan tarot spread: %w", err)
		}
		spread.Positions, err = r.listSpreadPositions(ctx, spread.ID)
		if err != nil {
			return nil, err
		}
		spreads = append(spreads, spread)
	}

	return spreads, rows.Err()
}

func (r *Repository) GetSpreadByCode(ctx context.Context, code string) (Spread, error) {
	row := r.db.QueryRow(ctx, `
		SELECT id::text, code, name, description, card_count
		FROM tarot_spreads
		WHERE code = $1 AND is_active = true
	`, code)

	var spread Spread
	if err := row.Scan(&spread.ID, &spread.Code, &spread.Name, &spread.Description, &spread.CardCount); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return Spread{}, ErrNotFound
		}
		return Spread{}, fmt.Errorf("get tarot spread: %w", err)
	}

	positions, err := r.listSpreadPositions(ctx, spread.ID)
	if err != nil {
		return Spread{}, err
	}
	spread.Positions = positions
	return spread, nil
}

func (r *Repository) listSpreadPositions(ctx context.Context, spreadID string) ([]SpreadPosition, error) {
	rows, err := r.db.Query(ctx, `
		SELECT id::text, position_number, code, name, description
		FROM tarot_spread_positions
		WHERE spread_id = $1
		ORDER BY position_number
	`, spreadID)
	if err != nil {
		return nil, fmt.Errorf("list tarot spread positions: %w", err)
	}
	defer rows.Close()

	var positions []SpreadPosition
	for rows.Next() {
		var position SpreadPosition
		if err := rows.Scan(&position.ID, &position.PositionNumber, &position.Code, &position.Name, &position.Description); err != nil {
			return nil, fmt.Errorf("scan tarot spread position: %w", err)
		}
		positions = append(positions, position)
	}
	return positions, rows.Err()
}

func (r *Repository) attachAssets(ctx context.Context, cards []Card) ([]Card, error) {
	return attachAssets(ctx, r.db, cards)
}

func (r *Repository) attachCardDetails(ctx context.Context, db assetQuerier, cards []Card) ([]Card, error) {
	cards, err := attachAssets(ctx, db, cards)
	if err != nil {
		return nil, err
	}
	return attachTranslations(ctx, db, cards)
}

type assetQuerier interface {
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

func attachAssets(ctx context.Context, db assetQuerier, cards []Card) ([]Card, error) {
	for i := range cards {
		rows, err := db.Query(ctx, `
			SELECT id::text, deck_code, size, format, url, width, height, file_size, is_default
			FROM tarot_card_assets
			WHERE card_id = $1
			ORDER BY is_default DESC, deck_code, size, format
		`, cards[i].ID)
		if err != nil {
			return nil, fmt.Errorf("list tarot card assets: %w", err)
		}

		for rows.Next() {
			var asset Asset
			if err := rows.Scan(&asset.ID, &asset.DeckCode, &asset.Size, &asset.Format, &asset.URL, &asset.Width, &asset.Height, &asset.FileSize, &asset.IsDefault); err != nil {
				rows.Close()
				return nil, fmt.Errorf("scan tarot card asset: %w", err)
			}
			cards[i].Assets = append(cards[i].Assets, asset)
		}
		if err := rows.Err(); err != nil {
			rows.Close()
			return nil, err
		}
		rows.Close()
	}
	return cards, nil
}

func attachTranslations(ctx context.Context, db assetQuerier, cards []Card) ([]Card, error) {
	for i := range cards {
		rows, err := db.Query(ctx, `
			SELECT language, name, description, meaning_upright, meaning_reversed
			FROM tarot_card_translations
			WHERE card_id = $1
				AND language IN ('th', 'en')
		`, cards[i].ID)
		if err != nil {
			return nil, fmt.Errorf("list tarot card translations: %w", err)
		}

		for rows.Next() {
			var language string
			var name string
			var description string
			var meaningUpright string
			var meaningReversed string
			if err := rows.Scan(&language, &name, &description, &meaningUpright, &meaningReversed); err != nil {
				rows.Close()
				return nil, fmt.Errorf("scan tarot card translation: %w", err)
			}
			if language == "th" {
				cards[i].NameTh = name
				cards[i].DescriptionTh = description
			}
		}
		if err := rows.Err(); err != nil {
			rows.Close()
			return nil, err
		}
		rows.Close()
	}
	return cards, nil
}

func scanCards(rows pgx.Rows) ([]Card, error) {
	var cards []Card
	for rows.Next() {
		card, err := scanCard(rows)
		if err != nil {
			return nil, err
		}
		cards = append(cards, card)
	}
	return cards, rows.Err()
}

type cardScanner interface {
	Scan(dest ...any) error
}

func scanCard(row cardScanner) (Card, error) {
	var card Card
	if err := row.Scan(&card.ID, &card.SourceCode, &card.NameEn, &card.Type, &card.Suit, &card.MeaningUpEn, &card.MeaningRevEn, &card.DescriptionEn); err != nil {
		return Card{}, err
	}
	if card.Assets == nil {
		card.Assets = []Asset{}
	}
	return card, nil
}

func (r *Repository) drawCards(ctx context.Context, tx pgx.Tx, limit int) ([]Card, error) {
	rows, err := tx.Query(ctx, `
		SELECT id::text, source_code, name_en, type, suit, meaning_up_en, meaning_rev_en, description_en
		FROM tarot_cards
		ORDER BY random()
		LIMIT $1
	`, limit)
	if err != nil {
		return nil, fmt.Errorf("draw tarot cards: %w", err)
	}
	defer rows.Close()

	cards, err := scanCards(rows)
	if err != nil {
		return nil, err
	}
	return r.attachCardDetails(ctx, tx, cards)
}

func (r *Repository) findCardsBySourceCodes(ctx context.Context, tx pgx.Tx, sourceCodes []string) ([]Card, error) {
	rows, err := tx.Query(ctx, `
		SELECT id::text, source_code, name_en, type, suit, meaning_up_en, meaning_rev_en, description_en
		FROM tarot_cards
		WHERE source_code = ANY($1)
	`, sourceCodes)
	if err != nil {
		return nil, fmt.Errorf("find selected tarot cards: %w", err)
	}
	defer rows.Close()

	foundCards, err := scanCards(rows)
	if err != nil {
		return nil, err
	}

	bySourceCode := make(map[string]Card, len(foundCards))
	for _, card := range foundCards {
		bySourceCode[card.SourceCode] = card
	}

	cards := make([]Card, 0, len(sourceCodes))
	for _, sourceCode := range sourceCodes {
		card, ok := bySourceCode[sourceCode]
		if !ok {
			return nil, fmt.Errorf("selected tarot card %q not found", sourceCode)
		}
		cards = append(cards, card)
	}
	return r.attachCardDetails(ctx, tx, cards)
}

func (r *Repository) findInterpretation(ctx context.Context, tx pgx.Tx, card Card, language string, topic string, orientation string) (interpretation, error) {
	queries := []struct {
		language string
		topic    string
	}{
		{language: language, topic: topic},
		{language: language, topic: TopicGeneral},
		{language: "en", topic: TopicGeneral},
	}

	for _, query := range queries {
		row := tx.QueryRow(ctx, `
			SELECT full_meaning, advice
			FROM tarot_card_interpretations
			WHERE card_id = $1
				AND language = $2
				AND topic = $3
				AND orientation = $4
			ORDER BY version DESC
			LIMIT 1
		`, card.ID, query.language, query.topic, orientation)

		var found interpretation
		if err := row.Scan(&found.Meaning, &found.Advice); err == nil {
			return found, nil
		} else if !errors.Is(err, pgx.ErrNoRows) {
			return interpretation{}, fmt.Errorf("find tarot interpretation: %w", err)
		}
	}

	if orientation == OrientationReversed {
		return fallbackInterpretation(card, orientation), nil
	}
	return fallbackInterpretation(card, orientation), nil
}

func fallbackInterpretation(card Card, orientation string) interpretation {
	if orientation == OrientationReversed {
		return interpretation{Meaning: card.MeaningRevEn, Advice: "Use this as a reflection point and move with care."}
	}
	return interpretation{Meaning: card.MeaningUpEn, Advice: "Use this as a reflection point and choose your next step intentionally."}
}

func (r *Repository) createReading(ctx context.Context, req CreateReadingRequest, spread Spread, cards []ReadingCard, summary string) (ReadingResponse, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return ReadingResponse{}, fmt.Errorf("begin tarot reading transaction: %w", err)
	}
	defer tx.Rollback(ctx)

	resultSnapshot, err := json.Marshal(map[string]any{
		"summary": summary,
		"cards":   cards,
	})
	if err != nil {
		return ReadingResponse{}, fmt.Errorf("marshal reading result snapshot: %w", err)
	}

	var readingID string
	var createdAt timeScanner
	err = tx.QueryRow(ctx, `
		INSERT INTO tarot_readings (spread_id, language, topic, question, status, result_snapshot)
		VALUES ($1, $2, $3, $4, 'completed', $5)
		RETURNING id::text, created_at
	`, spread.ID, req.Language, req.Topic, req.Question, resultSnapshot).Scan(&readingID, &createdAt.Time)
	if err != nil {
		return ReadingResponse{}, fmt.Errorf("insert tarot reading: %w", err)
	}

	for _, card := range cards {
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
		`, readingID, card.PositionNumber, card.Orientation, card.Meaning, card.Advice, cardSnapshot, spread.ID, card.Card.SourceCode)
		if err != nil {
			return ReadingResponse{}, fmt.Errorf("insert tarot reading card: %w", err)
		}
	}

	if err := tx.Commit(ctx); err != nil {
		return ReadingResponse{}, fmt.Errorf("commit tarot reading transaction: %w", err)
	}

	return ReadingResponse{
		ID:         readingID,
		SpreadCode: spread.Code,
		Topic:      req.Topic,
		Language:   req.Language,
		Question:   req.Question,
		Cards:      cards,
		Summary:    summary,
		CreatedAt:  createdAt.Time,
	}, nil
}

func (r *Repository) GetReading(ctx context.Context, id string) (ReadingResponse, error) {
	row := r.db.QueryRow(ctx, `
		SELECT tr.id::text, ts.code, tr.topic, tr.language, COALESCE(tr.question, ''), tr.result_snapshot, tr.created_at
		FROM tarot_readings tr
		JOIN tarot_spreads ts ON ts.id = tr.spread_id
		WHERE tr.id = $1
	`, id)

	var response ReadingResponse
	var snapshot []byte
	if err := row.Scan(&response.ID, &response.SpreadCode, &response.Topic, &response.Language, &response.Question, &snapshot, &response.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ReadingResponse{}, ErrNotFound
		}
		return ReadingResponse{}, fmt.Errorf("get tarot reading: %w", err)
	}

	var decoded struct {
		Summary string        `json:"summary"`
		Cards   []ReadingCard `json:"cards"`
	}
	if err := json.Unmarshal(snapshot, &decoded); err != nil {
		return ReadingResponse{}, fmt.Errorf("decode tarot reading snapshot: %w", err)
	}
	response.Summary = decoded.Summary
	response.Cards = decoded.Cards
	return response, nil
}

type timeScanner struct {
	Time time.Time
}
