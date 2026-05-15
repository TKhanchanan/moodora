package moon

import (
	"testing"
	"time"
)

func TestCalculateKnownNewMoon(t *testing.T) {
	location := time.UTC
	result := Calculator{}.Calculate(time.Date(2000, 1, 6, 18, 14, 0, 0, location), location)

	if result.MoonPhase != "new_moon" {
		t.Fatalf("MoonPhase = %s, want new_moon", result.MoonPhase)
	}
	if result.MoonAge < 27 && result.MoonAge > 1 {
		t.Fatalf("MoonAge = %.2f, want near new moon boundary", result.MoonAge)
	}
}

func TestCalculateDeterministicForSameDateAndTimezone(t *testing.T) {
	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		t.Fatal(err)
	}
	target := time.Date(2026, 5, 13, 8, 30, 0, 0, location)

	first := Calculator{}.Calculate(target, location)
	second := Calculator{}.Calculate(target, location)

	if first != second {
		t.Fatalf("calculation changed for same input: %+v != %+v", first, second)
	}
}

func TestBuildInterpretationIsSelfReflection(t *testing.T) {
	interpretation := BuildInterpretation("full_moon")
	if interpretation.Message == "" || interpretation.Advice == "" {
		t.Fatalf("interpretation should include message and advice: %+v", interpretation)
	}
}

func TestParseDateInLocation(t *testing.T) {
	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		t.Fatal(err)
	}
	parsed, err := ParseDateInLocation("2000-02-14", location)
	if err != nil {
		t.Fatal(err)
	}
	if got := parsed.Format("2006-01-02"); got != "2000-02-14" {
		t.Fatalf("date = %s, want 2000-02-14", got)
	}
}
