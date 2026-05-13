package moon

import (
	"math"
	"time"
)

const synodicMonth = 29.530588853

// Reference new moon: 2000-01-06 18:14 UTC.
var knownNewMoon = time.Date(2000, 1, 6, 18, 14, 0, 0, time.UTC)

type Calculator struct{}

func (Calculator) Calculate(target time.Time, location *time.Location) PhaseResult {
	local := target.In(location)
	noonLocal := time.Date(local.Year(), local.Month(), local.Day(), 12, 0, 0, 0, location)
	days := noonLocal.UTC().Sub(knownNewMoon).Hours() / 24
	age := math.Mod(days, synodicMonth)
	if age < 0 {
		age += synodicMonth
	}

	phase := phaseName(age)
	illumination := (1 - math.Cos(2*math.Pi*age/synodicMonth)) / 2 * 100

	return PhaseResult{
		TargetDate:               noonLocal.Format("2006-01-02"),
		Timezone:                 location.String(),
		MoonPhase:                phase,
		Illumination:             round2(illumination),
		MoonAge:                  round2(age),
		CalculationMethodVersion: CalculationMethodVersion,
	}
}

func phaseName(age float64) string {
	switch {
	case age < 1.84566:
		return "new_moon"
	case age < 5.53699:
		return "waxing_crescent"
	case age < 9.22831:
		return "first_quarter"
	case age < 12.91963:
		return "waxing_gibbous"
	case age < 16.61096:
		return "full_moon"
	case age < 20.30228:
		return "waning_gibbous"
	case age < 23.99361:
		return "last_quarter"
	case age < 27.68493:
		return "waning_crescent"
	default:
		return "new_moon"
	}
}

func round2(value float64) float64 {
	return math.Round(value*100) / 100
}

func BuildInterpretation(phase string) Interpretation {
	switch phase {
	case "new_moon":
		return Interpretation{Title: "New Moon", Message: "A quiet moment for setting intentions and making space.", Advice: "Choose one small reset that supports how you want to feel."}
	case "waxing_crescent":
		return Interpretation{Title: "Waxing Crescent", Message: "This phase can be used as a prompt for gentle momentum.", Advice: "Take a small step before asking yourself to have the full plan."}
	case "first_quarter":
		return Interpretation{Title: "First Quarter", Message: "A useful checkpoint for decisions, effort, and adjustment.", Advice: "Notice what needs action and what can be simplified."}
	case "waxing_gibbous":
		return Interpretation{Title: "Waxing Gibbous", Message: "A reflective phase for refinement before something feels complete.", Advice: "Improve one detail without chasing perfection."}
	case "full_moon":
		return Interpretation{Title: "Full Moon", Message: "A bright phase for awareness, release, and emotional honesty.", Advice: "Name what feels clear, then choose what you no longer need to carry."}
	case "waning_gibbous":
		return Interpretation{Title: "Waning Gibbous", Message: "A phase for sharing, learning, and integrating recent lessons.", Advice: "Turn one recent experience into a practical note for yourself."}
	case "last_quarter":
		return Interpretation{Title: "Last Quarter", Message: "A reflective checkpoint for letting go and choosing cleaner priorities.", Advice: "Release one task, expectation, or thought that is taking too much space."}
	default:
		return Interpretation{Title: "Waning Crescent", Message: "A soft phase for rest, closure, and quiet preparation.", Advice: "Protect your energy and make room for the next cycle."}
	}
}
