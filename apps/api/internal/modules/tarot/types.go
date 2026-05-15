package tarot

import "time"

const (
	TopicGeneral = "general"
	TopicLove    = "love"
	TopicCareer  = "career"
	TopicMoney   = "money"

	OrientationUpright  = "upright"
	OrientationReversed = "reversed"
)

var validTopics = map[string]bool{
	TopicGeneral: true,
	TopicLove:    true,
	TopicCareer:  true,
	TopicMoney:   true,
}

var validLanguages = map[string]bool{
	"en": true,
	"th": true,
}

type SourceCard struct {
	SourceCode    string `json:"name_short"`
	NameEn        string `json:"name"`
	Type          string `json:"type"`
	Suit          string `json:"suit"`
	MeaningUpEn   string `json:"meaning_up"`
	MeaningRevEn  string `json:"meaning_rev"`
	DescriptionEn string `json:"desc"`
}

type Card struct {
	ID            string  `json:"id"`
	SourceCode    string  `json:"sourceCode"`
	NameEn        string  `json:"nameEn"`
	NameTh        string  `json:"nameTh"`
	Type          string  `json:"type"`
	Suit          *string `json:"suit"`
	MeaningUpEn   string  `json:"meaningUpEn"`
	MeaningRevEn  string  `json:"meaningRevEn"`
	DescriptionEn string  `json:"descriptionEn"`
	DescriptionTh string  `json:"descriptionTh"`
	Assets        []Asset `json:"assets"`
}

type Asset struct {
	ID        string `json:"id"`
	DeckCode  string `json:"deckCode"`
	Size      string `json:"size"`
	Format    string `json:"format"`
	URL       string `json:"url"`
	Width     int    `json:"width"`
	Height    int    `json:"height"`
	FileSize  int64  `json:"fileSize"`
	IsDefault bool   `json:"isDefault"`
}

type Spread struct {
	ID          string           `json:"-"`
	Code        string           `json:"code"`
	Name        string           `json:"name"`
	Description string           `json:"description"`
	CardCount   int              `json:"cardCount"`
	Positions   []SpreadPosition `json:"positions"`
}

type SpreadPosition struct {
	ID             string `json:"-"`
	PositionNumber int    `json:"positionNumber"`
	Code           string `json:"code"`
	Name           string `json:"name"`
	Description    string `json:"description"`
}

type CreateReadingRequest struct {
	SpreadCode              string   `json:"spreadCode"`
	Topic                   string   `json:"topic"`
	Language                string   `json:"language"`
	AllowReversed           bool     `json:"allowReversed"`
	Question                string   `json:"question"`
	SelectedCardSourceCodes []string `json:"selectedCardSourceCodes"`
}

type ReadingResponse struct {
	ID         string        `json:"id"`
	SpreadCode string        `json:"spreadCode"`
	Topic      string        `json:"topic"`
	Language   string        `json:"language"`
	Question   string        `json:"question"`
	Cards      []ReadingCard `json:"cards"`
	Summary    string        `json:"summary"`
	CreatedAt  time.Time     `json:"createdAt"`
}

type ReadingCard struct {
	PositionNumber int         `json:"positionNumber"`
	PositionCode   string      `json:"positionCode"`
	PositionName   string      `json:"positionName"`
	Card           ReadingInfo `json:"card"`
	Orientation    string      `json:"orientation"`
	Meaning        string      `json:"meaning"`
	Advice         string      `json:"advice"`
}

type ReadingInfo struct {
	SourceCode     string  `json:"sourceCode"`
	Name           string  `json:"name"`
	NameEn         string  `json:"nameEn"`
	NameTh         string  `json:"nameTh"`
	Type           string  `json:"type"`
	Suit           *string `json:"suit"`
	Characteristic string  `json:"characteristic"`
	Description    string  `json:"description"`
	Assets         []Asset `json:"assets"`
}

type interpretation struct {
	Meaning string
	Advice  string
}
