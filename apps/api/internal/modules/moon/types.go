package moon

import "time"

const CalculationMethodVersion = "moon_phase_v1"

type PhaseResult struct {
	TargetDate               string  `json:"targetDate"`
	Timezone                 string  `json:"timezone"`
	MoonPhase                string  `json:"moonPhase"`
	Illumination             float64 `json:"illumination"`
	MoonAge                  float64 `json:"moonAge"`
	CalculationMethodVersion string  `json:"calculationMethodVersion"`
}

type Interpretation struct {
	Title   string `json:"title"`
	Message string `json:"message"`
	Advice  string `json:"advice"`
}

type Source struct {
	SourceCode string `json:"sourceCode"`
	SourceName string `json:"sourceName"`
	SourceURL  string `json:"sourceUrl"`
}

type ReportResponse struct {
	ID                       string         `json:"id"`
	TargetDate               string         `json:"targetDate"`
	Timezone                 string         `json:"timezone"`
	MoonPhase                string         `json:"moonPhase"`
	Illumination             float64        `json:"illumination"`
	MoonAge                  float64        `json:"moonAge"`
	ImageURL                 *string        `json:"imageUrl"`
	Interpretation           Interpretation `json:"interpretation"`
	CalculationMethodVersion string         `json:"calculationMethodVersion"`
	Source                   Source         `json:"source"`
	CreatedAt                time.Time      `json:"createdAt"`
}

type BirthdayRequest struct {
	BirthDate string `json:"birthDate"`
	Timezone  string `json:"timezone"`
}

type APOD struct {
	Date        string `json:"date"`
	Title       string `json:"title"`
	Explanation string `json:"explanation"`
	URL         string `json:"url"`
	MediaType   string `json:"mediaType"`
}
