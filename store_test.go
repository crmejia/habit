package habit_test

import (
	"habit"
	"testing"
	"time"
)

func TestMemoryStore_GetReturnsNilOnNoHabit(t *testing.T) {
	t.Parallel()
	store := habit.OpenStore()
	got := store.Get("piano")

	if got != nil {
		t.Error("want Store.Get to return nil")
	}
}

func TestMemoryStore_GetReturnsExistingHabit(t *testing.T) {
	t.Parallel()
	store := habit.OpenStore()
	store.Habits["piano"] = &habit.Habit{Name: "piano"}

	h := store.Get("piano")
	if h == nil {
		t.Fatal()
	}

	want := "piano"
	got := h.Name
	if want != got {
		t.Errorf("want h testName to be %s, h %s", want, got)
	}
}

func TestMemoryStore_CreateNewHabit(t *testing.T) {
	t.Parallel()
	store := habit.OpenStore()
	h := habit.Habit{Name: "piano"}

	store.Create(&h)

	if _, ok := store.Habits["piano"]; !ok {
		t.Error("want h to be inserted into store")
	}
}

func TestMemoryStore_CreateNilHabitFails(t *testing.T) {
	t.Parallel()
	store := habit.OpenStore()
	err := store.Create(nil)

	if err == nil {
		t.Error("want Store.create nil habit to fail with error")
	}
}

func TestMemoryStore_CreateExistingHabitFails(t *testing.T) {
	t.Parallel()
	store := habit.OpenStore()
	h := habit.Habit{Name: "piano"}
	store.Habits["piano"] = &h
	err := store.Create(&h)

	if err == nil {
		t.Error("want Store.Create nil h to fail with error")
	}
}

func TestMemoryStore_UpdateHabit(t *testing.T) {
	t.Parallel()
	store := habit.OpenStore()
	oldHabit := &habit.Habit{Name: "piano"}
	store.Habits["piano"] = oldHabit

	updateHabit := &habit.Habit{Name: "piano"}
	store.Update(updateHabit)

	if oldHabit == store.Habits["piano"] {
		t.Error("want update to replace habit")
	}
}

func TestMemoryStore_UpdateFailsOnNil(t *testing.T) {
	t.Parallel()
	store := habit.OpenStore()
	err := store.Update(nil)

	if err == nil {
		t.Error("want Store.Update nil habit to fail with error")
	}
}

func TestMemoryStore_UpdateFailsOnNonExistingHabit(t *testing.T) {
	t.Parallel()
	store := habit.OpenStore()
	err := store.Update(&habit.Habit{Name: "piano"})

	if err == nil {
		t.Error("want update to fail if habit does not exist")
	}
}

func TestMemoryStore_AllHabitsReturnsSliceOfHabits(t *testing.T) {
	t.Parallel()
	store := habit.OpenStore()
	store.Habits = map[string]*habit.Habit{
		"piano":   &habit.Habit{Name: "piano"},
		"surfing": &habit.Habit{Name: "surfing"},
	}

	allHabits := store.AllHabits()
	if len(allHabits) != len(store.Habits) {
		t.Error("want AllHabits to return a slice of habits")
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

func TestController_HandleSetsMessageCorrectlyForNewHabit(t *testing.T) {
	t.Parallel()
	store := habit.OpenStore()
	controller := habit.NewController(store)
	h := &habit.Habit{Name: "piano",
		Interval: habit.DailyInterval}
	controller.Handle(h)

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
	store := habit.OpenStore()
	controller := habit.NewController(store)
	for _, tc := range testCases {
		controller.Store.Habits[tc.habit.Name] = tc.habit
		controller.Handle(tc.habit)

		got := tc.habit.String()
		if tc.want != got {
			t.Errorf("For %d day streak: want the Message to be:\n%s,\n got\n%s", tc.habit.Streak, tc.want, got)
		}
	}
}
