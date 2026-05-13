package lifestyle

import "time"

const (
	PurposeGeneral   = "general"
	PurposeLove      = "love"
	PurposeCareer    = "career"
	PurposeMoney     = "money"
	PurposeStudy     = "study"
	PurposeInterview = "interview"

	RuleTypeLucky = "lucky"
	RuleTypeAvoid = "avoid"
)

var validPurposes = map[string]bool{
	PurposeGeneral:   true,
	PurposeLove:      true,
	PurposeCareer:    true,
	PurposeMoney:     true,
	PurposeStudy:     true,
	PurposeInterview: true,
}

type Color struct {
	ID      string `json:"id"`
	Code    string `json:"code"`
	NameTH  string `json:"nameTh"`
	NameEN  string `json:"nameEn"`
	Hex     string `json:"hex"`
	Meaning string `json:"meaning"`
}

type Food struct {
	ID          string   `json:"id"`
	Code        string   `json:"code"`
	NameTH      string   `json:"nameTh"`
	NameEN      string   `json:"nameEn"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
}

type Item struct {
	ID          string   `json:"id"`
	Code        string   `json:"code"`
	NameTH      string   `json:"nameTh"`
	NameEN      string   `json:"nameEn"`
	Category    string   `json:"category"`
	Tags        []string `json:"tags"`
	Description string   `json:"description"`
}

type Avoidance struct {
	ID       string `json:"id"`
	Code     string `json:"code"`
	Category string `json:"category"`
	TextTH   string `json:"textTh"`
	TextEN   string `json:"textEn"`
	MoodTag  string `json:"moodTag"`
}

type CheckInStatus struct {
	CheckedIn bool   `json:"checkedIn"`
	LocalDate string `json:"localDate"`
}

type DailyTarot struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

type DailyInsight struct {
	Date          string         `json:"date"`
	Timezone      string         `json:"timezone"`
	LuckyColors   []Color        `json:"luckyColors"`
	AvoidColors   []Color        `json:"avoidColors"`
	LuckyFoods    []Food         `json:"luckyFoods"`
	LuckyItems    []Item         `json:"luckyItems"`
	Avoidances    []Avoidance    `json:"avoidances"`
	WalletBalance *int64         `json:"walletBalance,omitempty"`
	CheckInStatus *CheckInStatus `json:"checkInStatus,omitempty"`
	DailyTarot    DailyTarot     `json:"dailyTarot"`
	Message       string         `json:"message"`
}

type TodayContext struct {
	UserID         string
	Date           time.Time
	Timezone       string
	DayOfWeek      int
	BirthDayOfWeek *int
	Purpose        string
}
