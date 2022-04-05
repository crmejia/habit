package habit_test

import (
	"habit"
	"testing"
	"time"
)

func TestHabitSetsMessageCorrectlyForNewHabit(t *testing.T) {
	ht := habit.Tracker{}
	h := ht.FetchHabit("piano")
	want := "Good luck with your new habit 'piano'! Don't forget to do it again tomorrow."
	got := h.String()
	if want != got {
		t.Errorf("For %d day streak: want the message to be:\n%s,\n got\n%s", h.Streak, want, got)
	}

}
func TestHabitSetsMessageCorrectlyForStreakBrokenStreak(t *testing.T) {
	testCases := []struct {
		want  string
		habit habit.Habit
	}{
		//{want: "Good luck with your new habit 'piano'! Don't forget to do it again tomorrow.", habit: habit.Habit{Name: "piano", Streak: 0}},
		{want: "Nice work: you've done the habit 'surf' for 4 days in a row now. Keep it up!", habit: habit.Habit{Name: "surf", Streak: 3, Period: time.Now()}},
		{want: "You last did the habit 'running' 10 days ago, so you're starting a new streak today. Good luck!", habit: habit.Habit{Name: "running", Streak: 10, Period: time.Now().Add(-10 * 24 * time.Hour)}},
	}
	ht := habit.Tracker{}
	for _, tc := range testCases {
		ht[tc.habit.Name] = tc.habit
		h := ht.FetchHabit(tc.habit.Name)
		got := h.String()
		if tc.want != got {
			t.Errorf("For %d day streak: want the message to be:\n%s,\n got\n%s", tc.habit.Streak, tc.want, got)
		}
	}
}

func TestFetchHabitReturnsANewHabitWithZeroDaysStreakOnNewHabit(t *testing.T) {
	tracker := habit.Tracker{}
	h := tracker.FetchHabit("piano")
	want := 0
	got := h.Streak

	if want != got {
		t.Errorf("For a new habit want %d,\n got %d", want, got)
	}
}

func TestFetchHabitIncreasesStreakOnExistingHabit(t *testing.T) {
	tracker := habit.Tracker{
		"piano": habit.Habit{
			Name:   "piano",
			Streak: 1,
			Period: time.Now(),
		},
	}
	h := tracker.FetchHabit("piano")
	want := 2
	got := h.Streak
	if want != got {
		t.Errorf("want streak to increase to %d, got %d", want, got)
	}
}

func TestFetchHabitIncreaseStreakOncePerDay(t *testing.T) {
	tracker := habit.Tracker{
		"piano": habit.Habit{
			Name:   "piano",
			Streak: 1,
			Period: habit.Tomorrow(),
		},
	}
	h := tracker.FetchHabit("piano")
	want := 1
	got := h.Streak
	if want != got {
		t.Errorf("want streak to increase to %d, got %d", want, got)
	}
}

func TestFetchHabitResetsStreak(t *testing.T) {
	fiveDaysAgo := time.Now().Add(-5 * 24 * time.Hour)
	tracker := habit.Tracker{
		"piano": habit.Habit{
			Name:   "piano",
			Streak: 8,
			Period: fiveDaysAgo,
		},
	}
	h := tracker.FetchHabit("piano")
	want := 0
	got := h.Streak
	if want != got {
		t.Errorf("want streak to reset to %d, got %d", want, got)
	}
}
