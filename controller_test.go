package habit_test

import (
	"habit"
	"testing"
)

func TestNewController(t *testing.T) {
	store := habit.OpenStore()
	controller := habit.NewController(store)

	if controller.Store.Habits != nil {
		t.Errorf("controller.Store should be initialized by new")
	}
}

func TestController_HandleReturnsErrorOnNilHabit(t *testing.T) {
	store := habit.OpenStore()
	controller := habit.NewController(store)
	_, err := controller.Handle(nil)
	if err == nil {
		t.Error("expected err got nil")
	}
}

func TestController_HandleReturnsErrorOnEmptyHabitName(t *testing.T) {
	store := habit.OpenStore()
	controller := habit.NewController(store)
	h := habit.Habit{Name: ""}
	_, err := controller.Handle(&h)
	if err == nil {
		t.Error("expected err got nil")
	}
}

func TestController_GetHabitReturnsHabitOnNoHabits(t *testing.T) {
	store := habit.OpenStore()
	controller := habit.NewController(store)

	newHabit := habit.Habit{Name: "piano"}
	got, err := controller.Handle(&newHabit)
	if err != nil {
		t.Fatal(err)
	}
	if got == nil {
		t.Error("want controller.Get to return an existing habit or create a new one. Cannot be nil")
	}
}
