package habit_test

import (
	"habit"
	"testing"
	"time"
)

func TestFetchHabitSetsMessageCorrectlyForNewHabit(t *testing.T) {
	t.Parallel()
	ht := habit.Tracker{}
	h := ht.FetchHabit("piano")
	want := "Good luck with your new habit 'piano'! Don't forget to do it again tomorrow."
	got := h.String()
	if want != got {
		t.Errorf("For %d day streak: want the message to be:\n%s,\n got\n%s", h.Streak, want, got)
	}
}

func TestFetchHabitSetsMessageCorrectlyForStreakBrokenStreak(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		want  string
		habit *habit.Habit
	}{
		{want: "Nice work: you've done the habit 'surf' for 4 days in a row now. Keep it up!", habit: &habit.Habit{Name: "surf", Streak: 3, Period: time.Now()}},
		{want: "You last did the habit 'running' 10 days ago, so you're starting a new streak today. Good luck!", habit: &habit.Habit{Name: "running", Streak: 10, Period: time.Now().Add(-10 * 24 * time.Hour)}},
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

func TestFetchHabitSetsMessageCorrectlyForAlreadyIncreasedStreak(t *testing.T) {
	t.Parallel()
	ht := habit.Tracker{
		"piano": &habit.Habit{
			Name:   "piano",
			Streak: 2,
			Period: habit.Tomorrow(),
		},
	}
	h := ht.FetchHabit("piano")
	want := "Nice work: you've done the habit 'piano' for 2 days in a row now. Keep it up!"
	got := h.String()
	if want != got {
		t.Errorf("For %d day streak: want the message to be:\n%s,\n got\n%s", h.Streak, want, got)
	}
}
func TestTracker_FetchHabitReturnPtrIsMatchTheMapPtr(t *testing.T) {
	t.Parallel()
	tracker := habit.Tracker{
		"piano": &habit.Habit{
			Name: "piano",
		},
	}
	h := tracker.FetchHabit("piano")
	if h != tracker["piano"] {
		t.Error("want FetchHabit return ptr to be equal to the Map(Tracker type) ptr")
	}
}
func TestFetchHabitReturnsANewHabitWithZeroDaysStreakOnNewHabit(t *testing.T) {
	t.Parallel()
	tracker := habit.Tracker{}
	tracker.FetchHabit("piano")
	want := 0
	got := tracker["piano"].Streak

	if want != got {
		t.Errorf("For a new habit want %d,\n got %d", want, got)
	}
}

func TestFetchHabitIncreasesStreakOnExistingHabit(t *testing.T) {
	t.Parallel()
	tracker := habit.Tracker{
		"piano": &habit.Habit{
			Name:   "piano",
			Streak: 1,
			Period: time.Now(),
		},
	}
	tracker.FetchHabit("piano")
	want := 2
	got := tracker["piano"].Streak
	if want != got {
		t.Errorf("want streak to increase to %d, got %d", want, got)
	}
}

func TestFetchHabitIncreaseStreakOncePerDay(t *testing.T) {
	t.Parallel()
	tracker := habit.Tracker{
		"piano": &habit.Habit{
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
	t.Parallel()
	fiveDaysAgo := time.Now().Add(-5 * 24 * time.Hour)
	tracker := habit.Tracker{
		"piano": &habit.Habit{
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
