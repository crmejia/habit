package habit_test

import (
	"habit"
	"strings"
	"testing"
	"time"
)

func TestNewController(t *testing.T) {
	t.Parallel()
	store := habit.OpenMemoryStore()
	controller, _ := habit.NewController(&store)

	if controller.Store == nil {
		t.Errorf("controller.Store should be initialized by new")
	}
}

func TestNewControllerReturnsErrorOnNilStore(t *testing.T) {
	t.Parallel()
	_, err := habit.NewController(nil)
	if err == nil {
		t.Error("want NewController to return error on nil store")
	}
}

func TestController_HandleReturnsErrorOnNilHabit(t *testing.T) {
	t.Parallel()
	store := habit.OpenMemoryStore()
	controller, _ := habit.NewController(&store)
	_, err := controller.Handle(nil)
	if err == nil {
		t.Error("expected err got nil")
	}
}

func TestController_HandleReturnsErrorOnEmptyHabitName(t *testing.T) {
	t.Parallel()
	store := habit.OpenMemoryStore()
	controller, _ := habit.NewController(&store)
	h := habit.Habit{Name: ""}
	_, err := controller.Handle(&h)
	if err == nil {
		t.Error("expected err got nil")
	}
}

func TestController_HandleUpdatesStreaksDueDateCorrectly(t *testing.T) {
	t.Parallel()
	inputHabit := habit.Habit{
		Name: "piano",
	}
	store := habit.MemoryStore{
		Habits: map[string]*habit.Habit{"piano": &inputHabit},
	}
	controller, _ := habit.NewController(&store)
	testCases := []struct {
		name          string
		streak        int
		dueDate       time.Time
		interval      time.Duration
		wantedStreak  int
		wantedDueDate time.Time
	}{
		{name: "increase streak on daily habit with same day due date", streak: 0, dueDate: time.Now(), interval: habit.DailyInterval, wantedStreak: 1, wantedDueDate: time.Now().Add(habit.DailyInterval)},
		{name: "does not increase streak on already updated daily habit", streak: 1, dueDate: time.Now().Add(habit.DailyInterval), interval: habit.DailyInterval, wantedStreak: 1, wantedDueDate: time.Now().Add(habit.DailyInterval)},
		{name: "resets streak on overdue daily habit", streak: 1, dueDate: time.Now().Add(-1 * habit.DailyInterval), interval: habit.DailyInterval, wantedStreak: 0, wantedDueDate: time.Now().Add(habit.DailyInterval)},
		{name: "increase streak on weekly habit with same day due date", streak: 0, dueDate: time.Now(), interval: habit.WeeklyInterval, wantedStreak: 1, wantedDueDate: time.Now().Add(habit.WeeklyInterval)},
		{name: "does not increase streak on already updated weekly habit", streak: 1, dueDate: time.Now().Add(habit.WeeklyInterval), interval: habit.WeeklyInterval, wantedStreak: 1, wantedDueDate: time.Now().Add(habit.WeeklyInterval)},
		{name: "resets streak on overdue weekly habit", streak: 1, dueDate: time.Now().Add(-1 * habit.WeeklyInterval), interval: habit.WeeklyInterval, wantedStreak: 0, wantedDueDate: time.Now().Add(habit.WeeklyInterval)},
	}

	for _, tc := range testCases {
		inputHabit.Streak = tc.streak
		inputHabit.DueDate = tc.dueDate
		inputHabit.Frequency = tc.interval

		h, _ := controller.Handle(&inputHabit)

		if h.Name != "piano" {
			t.Errorf("wantedStreak piano to be the habit's testName got %s", h.Name)
		}
		if store.Habits["piano"].Streak != tc.wantedStreak {
			t.Errorf("%s. Want habit.Streak to be %d got %d", tc.name, tc.wantedStreak, h.Streak)
		}

		if !habit.SameDay(h.DueDate, tc.wantedDueDate) {
			t.Errorf("%s. Want habit.DueDate to be %s got %s", tc.name, tc.wantedDueDate, h.DueDate)
		}
	}
}

func TestController_HandleCreatesErrorsOnNoInterval(t *testing.T) {
	t.Parallel()
	store := habit.OpenMemoryStore()
	controller, _ := habit.NewController(&store)

	newHabit := habit.Habit{Name: "piano"}
	_, err := controller.Handle(&newHabit)
	if err == nil {
		t.Errorf("expected create new habit with not interval to return error")
	}
}

func TestController_HandleCreatesHabit(t *testing.T) {
	t.Parallel()
	store := habit.OpenMemoryStore()
	controller, _ := habit.NewController(&store)

	newHabit := habit.Habit{Name: "piano", Frequency: habit.DailyInterval}
	_, err := controller.Handle(&newHabit)
	if err != nil {
		t.Errorf("expected handle to return no errors, got: %s", err)
	}

	if store.Habits["piano"] == nil {
		t.Error("want new habit to be inserted into store")
	}
}

func TestController_AllHabits(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		store habit.MemoryStore
		want  string
	}{
		{store: habit.MemoryStore{Habits: map[string]*habit.Habit{}}, want: "no habits have been started"},
		{store: habit.MemoryStore{Habits: map[string]*habit.Habit{"piano": {Name: "piano"}}}, want: "piano"},
	}

	for _, tc := range testCases {
		controller, _ := habit.NewController(&tc.store)
		got := controller.GetAllHabits()
		if !strings.Contains(got, tc.want) {
			t.Errorf("want output to contain %s, got:\n    %s", tc.want, got)
		}
	}
}
