package lifestyle

import (
	"reflect"
	"testing"
	"time"
)

func TestPickDeterministicStableForSameDateSeed(t *testing.T) {
	values := []string{"a", "b", "c", "d", "e"}
	first := pickDeterministic(values, "user-1|2026-05-13|foods", 2)
	second := pickDeterministic(values, "user-1|2026-05-13|foods", 2)

	if !reflect.DeepEqual(first, second) {
		t.Fatalf("selection changed for same seed: %v != %v", first, second)
	}
	if len(first) != 2 {
		t.Fatalf("selection length = %d, want 2", len(first))
	}
}

func TestPickDeterministicChangesWithDateSeed(t *testing.T) {
	values := []string{"a", "b", "c", "d", "e", "f", "g"}
	first := pickDeterministic(values, "user-1|2026-05-13|items", 3)
	second := pickDeterministic(values, "user-1|2026-05-14|items", 3)

	if reflect.DeepEqual(first, second) {
		t.Fatalf("selection should usually change when date seed changes: %v", first)
	}
}

func TestLocalTodayUsesBangkokDate(t *testing.T) {
	location, err := time.LoadLocation("Asia/Bangkok")
	if err != nil {
		t.Fatal(err)
	}
	now := time.Date(2026, 5, 13, 18, 30, 0, 0, time.UTC)
	today := LocalToday(now, location)

	if got := today.Format("2006-01-02"); got != "2026-05-14" {
		t.Fatalf("Bangkok local date = %s, want 2026-05-14", got)
	}
}
