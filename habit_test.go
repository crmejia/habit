package habit_test

import (
	"habit"
	"testing"
	"time"
)

func TestHabitImplementsStringer(t *testing.T) {
	testCases := []struct {
		want  string
		habit habit.Habit
	}{
		{want: "Good luck with your new habit 'piano'! Don't forget to do it again tomorrow.", habit: habit.Habit{Name: "piano", Streak: 0}},
		{want: "Nice work: you've done the habit 'surf' for 4 days in a row now. Keep it up!", habit: habit.Habit{Name: "surf", Streak: 4}},
	}
	for _, tc := range testCases {
		got := tc.habit.String()
		if tc.want != got {
			t.Errorf("For %d day streak: want the message to be %s,\n got %s", tc.habit.Streak, tc.want, got)
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
