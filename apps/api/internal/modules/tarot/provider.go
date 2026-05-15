package tarot

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type CardProvider interface {
	LoadCards() ([]SourceCard, error)
}

type FileProvider struct {
	Path string
}

func (p FileProvider) LoadCards() ([]SourceCard, error) {
	data, err := os.ReadFile(p.Path)
	if err != nil {
		return nil, fmt.Errorf("read tarot source file: %w", err)
	}
	return decodeCards(data)
}

type BuiltInProvider struct{}

func (BuiltInProvider) LoadCards() ([]SourceCard, error) {
	return buildDefaultCards(), nil
}

func decodeCards(data []byte) ([]SourceCard, error) {
	var wrapped struct {
		Cards []SourceCard `json:"cards"`
	}
	if err := json.Unmarshal(data, &wrapped); err == nil && len(wrapped.Cards) > 0 {
		return normalizeSourceCards(wrapped.Cards)
	}

	var cards []SourceCard
	if err := json.Unmarshal(data, &cards); err != nil {
		return nil, fmt.Errorf("decode tarot source cards: %w", err)
	}
	return normalizeSourceCards(cards)
}

func normalizeSourceCards(cards []SourceCard) ([]SourceCard, error) {
	for i := range cards {
		cards[i].SourceCode = strings.TrimSpace(cards[i].SourceCode)
		cards[i].NameEn = strings.TrimSpace(cards[i].NameEn)
		cards[i].Type = strings.TrimSpace(cards[i].Type)
		cards[i].Suit = strings.TrimSpace(cards[i].Suit)
		cards[i].MeaningUpEn = strings.TrimSpace(cards[i].MeaningUpEn)
		cards[i].MeaningRevEn = strings.TrimSpace(cards[i].MeaningRevEn)
		cards[i].DescriptionEn = strings.TrimSpace(cards[i].DescriptionEn)

		if cards[i].SourceCode == "" || cards[i].NameEn == "" || cards[i].Type == "" {
			return nil, fmt.Errorf("card %d is missing required source fields", i)
		}
	}
	return cards, nil
}

func buildDefaultCards() []SourceCard {
	major := []string{
		"The Fool", "The Magician", "The High Priestess", "The Empress", "The Emperor", "The Hierophant",
		"The Lovers", "The Chariot", "Strength", "The Hermit", "Wheel of Fortune", "Justice",
		"The Hanged Man", "Death", "Temperance", "The Devil", "The Tower", "The Star",
		"The Moon", "The Sun", "Judgement", "The World",
	}

	cards := make([]SourceCard, 0, 78)
	for i, name := range major {
		cards = append(cards, SourceCard{
			SourceCode:    fmt.Sprintf("ar%02d", i),
			NameEn:        name,
			Type:          "major",
			MeaningUpEn:   fmt.Sprintf("%s invites clear reflection, openness, and intentional action.", name),
			MeaningRevEn:  fmt.Sprintf("%s reversed asks for patience, grounding, and a closer look at assumptions.", name),
			DescriptionEn: fmt.Sprintf("%s is part of the major arcana and represents a broad life theme for self-reflection.", name),
		})
	}

	suits := []struct {
		code string
		name string
	}{
		{"wa", "wands"},
		{"cu", "cups"},
		{"sw", "swords"},
		{"pe", "pentacles"},
	}
	ranks := []string{"Ace", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Page", "Knight", "Queen", "King"}

	for _, suit := range suits {
		for i, rank := range ranks {
			name := fmt.Sprintf("%s of %s", rank, strings.Title(suit.name))
			cards = append(cards, SourceCard{
				SourceCode:    fmt.Sprintf("%s%02d", suit.code, i+1),
				NameEn:        name,
				Type:          "minor",
				Suit:          suit.name,
				MeaningUpEn:   fmt.Sprintf("%s suggests a practical reflection connected to %s energy.", name, suit.name),
				MeaningRevEn:  fmt.Sprintf("%s reversed suggests slowing down and reassessing the %s theme.", name, suit.name),
				DescriptionEn: fmt.Sprintf("%s is a minor arcana card in the suit of %s.", name, suit.name),
			})
		}
	}

	return cards
}
