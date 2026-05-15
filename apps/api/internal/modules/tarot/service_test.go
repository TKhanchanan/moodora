package tarot

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestBuildSummaryUsesSelectedMeanings(t *testing.T) {
	summary := BuildSummary("en", []ReadingCard{
		{Meaning: "First reflection."},
		{Meaning: "Second reflection."},
		{Meaning: "Third reflection."},
	})

	if !strings.Contains(summary, "First reflection.") || !strings.Contains(summary, "Second reflection.") {
		t.Fatalf("summary did not include first two meanings: %q", summary)
	}
	if strings.Contains(summary, "Third reflection.") {
		t.Fatalf("summary should stay short: %q", summary)
	}
}

func TestFallbackInterpretationUsesRawCardMeaning(t *testing.T) {
	card := Card{
		MeaningUpEn:  "upright meaning",
		MeaningRevEn: "reversed meaning",
	}

	upright := fallbackInterpretation(card, OrientationUpright)
	if upright.Meaning != "upright meaning" {
		t.Fatalf("upright fallback = %q", upright.Meaning)
	}

	reversed := fallbackInterpretation(card, OrientationReversed)
	if reversed.Meaning != "reversed meaning" {
		t.Fatalf("reversed fallback = %q", reversed.Meaning)
	}
}

func TestReadingInfoIncludesAssetsInJSON(t *testing.T) {
	card := ReadingCard{
		Card: ReadingInfo{
			SourceCode: "ar01",
			Name:       "The Magician",
			Assets: []Asset{
				{
					DeckCode:  "rider_waite",
					Size:      "medium",
					Format:    "webp",
					URL:       "http://localhost:9000/moodora-assets/tarot/rider_waite/webp/medium/ar01.webp",
					Width:     480,
					Height:    840,
					FileSize:  12345,
					IsDefault: true,
				},
			},
		},
	}

	payload, err := json.Marshal(card)
	if err != nil {
		t.Fatalf("marshal reading card: %v", err)
	}
	if !strings.Contains(string(payload), `"assets"`) {
		t.Fatalf("reading card JSON should include assets: %s", payload)
	}
	if !strings.Contains(string(payload), "ar01.webp") {
		t.Fatalf("reading card JSON should include asset URL: %s", payload)
	}
}

func TestSelectReadingCardsValidatesSelectedCount(t *testing.T) {
	service := NewService(nil)
	_, err := service.selectReadingCards(t.Context(), nil, CreateReadingRequest{
		SelectedCardSourceCodes: []string{"ar00"},
	}, 3)
	if err == nil || !strings.Contains(err.Error(), "exactly 3 cards") {
		t.Fatalf("expected selected count error, got %v", err)
	}
}

func TestSelectReadingCardsRejectsDuplicateSelectedCards(t *testing.T) {
	service := NewService(nil)
	_, err := service.selectReadingCards(t.Context(), nil, CreateReadingRequest{
		SelectedCardSourceCodes: []string{"ar00", "ar00"},
	}, 2)
	if err == nil || !strings.Contains(err.Error(), "duplicate card") {
		t.Fatalf("expected duplicate card error, got %v", err)
	}
}
