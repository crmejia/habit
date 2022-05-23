package habit_test

import (
	"habit"
	"testing"
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

	habit := store.Get("piano")
	if habit == nil {
		t.Fatal()
	}

	want := "piano"
	got := habit.Name
	if want != got {
		t.Errorf("want habit name to be %s, habit %s", want, got)
	}
}

func TestMemoryStore_CreateNewHabit(t *testing.T) {
	t.Parallel()
	store := habit.OpenStore()
	habit := habit.Habit{Name: "piano"}

	store.Create(&habit)

	if _, ok := store.Habits["piano"]; !ok {
		t.Error("want habit to be inserted into store")
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
	habit := habit.Habit{Name: "piano"}
	store.Habits["piano"] = &habit
	err := store.Create(&habit)

	if err == nil {
		t.Error("want Store.Create nil habit to fail with error")
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
