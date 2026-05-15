package tarot

import "testing"

func TestBuiltInProviderReturnsUnique78Cards(t *testing.T) {
	cards, err := BuiltInProvider{}.LoadCards()
	if err != nil {
		t.Fatalf("LoadCards() error = %v", err)
	}
	if len(cards) != 78 {
		t.Fatalf("card count = %d, want 78", len(cards))
	}

	seen := map[string]bool{}
	for _, card := range cards {
		if card.SourceCode == "" {
			t.Fatal("source code must not be empty")
		}
		if seen[card.SourceCode] {
			t.Fatalf("duplicate source code %q", card.SourceCode)
		}
		seen[card.SourceCode] = true
	}

	for _, code := range []string{"ar01", "sw08"} {
		if !seen[code] {
			t.Fatalf("expected source code %q", code)
		}
	}
}
