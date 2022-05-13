package habit_test

import (
	"habit"
	"testing"
	"time"
)

func TestTracker_FetchHabitReturnPtrMatchesMapPtr(t *testing.T) {
	t.Parallel()
	tracker := habit.Tracker{
		"piano": &habit.Habit{
			Name: "piano",
		},
	}
	h, _ := tracker.FetchHabit("piano")
	if h != tracker["piano"] {
		t.Error("want FetchHabit return ptr to be equal to the Map(Tracker type) ptr")
	}
}

func TestTracker_FetchHabitReturnsFalseOnNonExistentHabit(t *testing.T) {
	t.Parallel()
	tracker := habit.Tracker{}
	_, got := tracker.FetchHabit("piano")
	want := false

	if want != got {
		t.Errorf("For a new habit want %t,\n got %t", want, got)
	}
}

func TestTracker_FetchHabitIncreasesStreakOnExistingDailyHabit(t *testing.T) {
	t.Parallel()
	tracker := habit.Tracker{
		"piano": &habit.Habit{
			Name:    "piano",
			Streak:  1,
			DueDate: time.Now(),
		},
	}
	tracker.FetchHabit("piano")
	want := 2
	got := tracker["piano"].Streak
	if want != got {
		t.Errorf("want streak to increase to %d, got %d", want, got)
	}
}

func TestTracker_FetchHabitSetsCorrectDueDateOnExistingWeeklyHabit(t *testing.T) {
	t.Parallel()
	tracker := habit.Tracker{
		"piano": &habit.Habit{
			Name:     "piano",
			Streak:   1,
			DueDate:  time.Now(),
			Interval: habit.WeeklyInterval,
		},
	}
	tracker.FetchHabit("piano")
	want := time.Now().Add(habit.WeeklyInterval)
	got := tracker["piano"].DueDate
	if !habit.SameDay(want, got) {
		t.Errorf("want DueDate to be set to %q, got %q", want, got)
	}
}

func TestTracker_FetchHabitIncreaseStreakOncePerDay(t *testing.T) {
	t.Parallel()
	tracker := habit.Tracker{
		"piano": &habit.Habit{
			Name:     "piano",
			Interval: habit.DailyInterval,
			Streak:   1,
			DueDate:  time.Now().Add(habit.DailyInterval),
		},
	}
	h, _ := tracker.FetchHabit("piano")
	want := 1
	got := h.Streak
	if want != got {
		t.Errorf("want streak to increase to %d, got %d", want, got)
	}
}

func TestTracker_FetchHabitIncreaseWeeklyStreakOncePerWeeks(t *testing.T) {
	t.Parallel()
	tracker := habit.Tracker{
		"piano": &habit.Habit{
			Name:     "piano",
			Interval: habit.WeeklyInterval,
			Streak:   1,
			DueDate:  time.Now().Add(habit.WeeklyInterval),
		},
	}
	h, _ := tracker.FetchHabit("piano")
	want := 1
	got := h.Streak
	if want != got {
		t.Errorf("want streak to increase to %d, got %d", want, got)
	}
}

func TestTracker_FetchHabitResetsStreak(t *testing.T) {
	t.Parallel()
	fiveDaysAgo := time.Now().Add(-5 * 24 * time.Hour)
	tracker := habit.Tracker{
		"piano": &habit.Habit{
			Name:    "piano",
			Streak:  8,
			DueDate: fiveDaysAgo,
		},
	}
	h, _ := tracker.FetchHabit("piano")
	want := 0
	got := h.Streak
	if want != got {
		t.Errorf("want streak to reset to %d, got %d", want, got)
	}
}

func TestTracker_FetchHabitResetsStreakOnWeeklyHabit(t *testing.T) {
	t.Parallel()
	twoWeeksAgo := time.Now().Add(-14 * 24 * time.Hour)
	tracker := habit.Tracker{
		"piano": &habit.Habit{
			Name:     "piano",
			Interval: habit.WeeklyInterval,
			Streak:   8,
			DueDate:  twoWeeksAgo,
		},
	}
	h, _ := tracker.FetchHabit("piano")
	want := 0
	got := h.Streak
	if want != got {
		t.Errorf("want streak to reset to %d, got %d", want, got)
	}
}

func TestAllHabitsReturnsHabitStreaks(t *testing.T) {
	t.Parallel()
	tracker := habit.Tracker{
		"piano": &habit.Habit{
			Name:    "piano",
			Streak:  8,
			DueDate: time.Now().Add(habit.DailyInterval),
		},
	}
	store := habit.FileStore{Tracker: tracker}
	want := "Habits:\nYou're currently on a 8-day streak for 'piano'. Stick to it!\n"
	got := habit.AllHabits(store)
	if want != got {
		t.Errorf("want:\n %s \ngot:\n %s", want, got)
	}
}

func TestAllHabitsReturnsEmptyStringOnNoHabits(t *testing.T) {
	t.Parallel()
	tracker := habit.Tracker{}
	want := ""
	store := habit.FileStore{Tracker: tracker}
	got := habit.AllHabits(store)
	if want != got {
		t.Errorf("want:\n %s \ngot:\n %s", want, got)
	}
}

func TestTracker_CreateHabitCreatesAWeeklyHabit(t *testing.T) {
	t.Parallel()
	tracker := habit.Tracker{}
	newHabit := habit.Habit{
		Name:     "piano",
		Interval: habit.WeeklyInterval,
	}
	tracker.CreateHabit(&newHabit)
	want := time.Now().Add(7 * 24 * time.Hour)
	got := tracker["piano"].DueDate

	if !habit.SameDay(want, got) {
		t.Errorf("For a new habit want %q,\n got %q", want, got)
	}
}

func TestTracker_CreateHabitFailsOnExistingHabit(t *testing.T) {
	t.Parallel()
	tracker := habit.Tracker{
		"piano": &habit.Habit{
			Name: "piano",
		},
	}
	newHabit := habit.Habit{
		Name:     "piano",
		Interval: habit.WeeklyInterval,
	}
	err := tracker.CreateHabit(&newHabit)
	if err == nil {
		t.Errorf("want an error on recreating existing habits, got : %s", err.Error())
	}

}

func TestTracker_CreateHabitFailsOnInvalidInterval(t *testing.T) {
	t.Parallel()
	tracker := habit.Tracker{}
	newHabit := habit.Habit{
		Name:     "piano",
		Interval: time.Hour,
	}
	err := tracker.CreateHabit(&newHabit)
	if err == nil {
		t.Errorf("want an error on creating habit with invalid interval, got : %s", err.Error())
	}

}

func TestTracker_CreateHabitSetsMessageCorrectlyForNewHabit(t *testing.T) {
	t.Parallel()
	ht := habit.Tracker{}
	h := &habit.Habit{Name: "piano",
		Interval: habit.DailyInterval}
	ht.CreateHabit(h)
	want := "Good luck with your new habit 'piano'! Don't forget to do it again tomorrow."
	got := h.String()
	if want != got {
		t.Errorf("For %d day streak: want the Message to be:\n%s,\n got\n%s", h.Streak, want, got)
	}
}

func TestTracker_FetchHabitSetsMessageCorrectlyForStreakBrokenStreak(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		want  string
		habit *habit.Habit
	}{
		{want: "Nice work: you've done the habit 'surf' for 4 days in a row now. Keep it up!", habit: &habit.Habit{Name: "surf", Streak: 3, DueDate: time.Now()}},
		{want: "You last did the habit 'running' 10 days ago, so you're starting a new streak today. Good luck!", habit: &habit.Habit{Name: "running", Streak: 10, DueDate: time.Now().Add(-10 * 24 * time.Hour)}},
	}
	ht := habit.Tracker{}
	for _, tc := range testCases {
		ht[tc.habit.Name] = tc.habit
		h, _ := ht.FetchHabit(tc.habit.Name)
		got := h.String()
		if tc.want != got {
			t.Errorf("For %d day streak: want the Message to be:\n%s,\n got\n%s", tc.habit.Streak, tc.want, got)
		}
	}
}

func TestTracker_FetchHabitSetsMessageCorrectlyForAlreadyIncreasedStreak(t *testing.T) {
	t.Parallel()
	ht := habit.Tracker{
		"piano": &habit.Habit{
			Name:     "piano",
			Interval: habit.DailyInterval,
			Streak:   2,
			DueDate:  time.Now().Add(habit.DailyInterval),
		},
	}
	h, _ := ht.FetchHabit("piano")
	want := "You already logged 'piano' today. Keep it up!"
	got := h.String()
	if want != got {
		t.Errorf("For %d day streak: want the Message to be:\n%s,\n got\n%s", h.Streak, want, got)
	}
}

func TestMessageGenerator(t *testing.T) {
	testCases := []struct {
		h    habit.Habit
		kind habit.MessageKind
		want string
	}{
		{habit.Habit{Name: "piano", Interval: habit.WeeklyInterval}, habit.NewMessage, "Good luck with your new habit 'piano'! Don't forget to do it again in a week."},
		{habit.Habit{Name: "piano", Interval: habit.DailyInterval}, habit.NewMessage, "Good luck with your new habit 'piano'! Don't forget to do it again tomorrow."},
		{habit.Habit{Name: "surfing", Interval: habit.WeeklyInterval}, habit.RepeatMessage, "You already logged 'surfing' today. Keep it up!"},
		{habit.Habit{Name: "meditation", Interval: habit.DailyInterval}, habit.RepeatMessage, "You already logged 'meditation' today. Keep it up!"},
		{habit.Habit{Name: "dancing", Interval: habit.WeeklyInterval, Streak: 2}, habit.StreakMessage, "Nice work: you've done the habit 'dancing' for 2 weeks in a row now. Keep it up!"},
		{habit.Habit{Name: "meditation", Interval: habit.DailyInterval, Streak: 2}, habit.StreakMessage, "Nice work: you've done the habit 'meditation' for 2 days in a row now. Keep it up!"},
		{habit.Habit{Name: "running", Interval: habit.DailyInterval, DueDate: time.Now().Add(-5 * 24 * time.Hour)}, habit.BrokenMessage, "You last did the habit 'running' 5 days ago, so you're starting a new streak today. Good luck!"},
		{habit.Habit{Name: "hiking", Interval: habit.WeeklyInterval, DueDate: time.Now().Add(-3 * 24 * 7 * time.Hour)}, habit.BrokenMessage, "You last did the habit 'hiking' 3 weeks ago, so you're starting a new streak today. Good luck!"},
	}

	for _, tc := range testCases {
		tc.h.GenerateMessage(tc.kind)
		got := tc.h.Message
		if tc.want != got {
			t.Errorf("want Message to be:\n%s\ngot:\n%s", tc.want, got)
		}
	}
}
