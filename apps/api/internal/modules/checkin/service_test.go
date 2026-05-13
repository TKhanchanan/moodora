package checkin

import "testing"

func TestRewardForStreakDay(t *testing.T) {
	tests := []struct {
		streakDay int
		want      int64
	}{
		{1, 5},
		{2, 5},
		{3, 8},
		{4, 8},
		{5, 10},
		{6, 15},
		{7, 25},
		{8, 25},
	}

	for _, tt := range tests {
		if got := RewardForStreakDay(tt.streakDay); got != tt.want {
			t.Fatalf("RewardForStreakDay(%d) = %d, want %d", tt.streakDay, got, tt.want)
		}
	}
}
