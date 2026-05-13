package tarot

import (
	"strings"
	"testing"
)

func TestBuildSummaryUsesSelectedMeanings(t *testing.T) {
	summary := BuildSummary([]ReadingCard{
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
