package habit_test

import (
	"habit"
	"testing"
)

func TestHabitImplementsStringer(t *testing.T) {
	want := "Good luck with your new habit 'piano'! Don't forget to do it again\ntomorrow."
	h := habit.FetchHabit("piano")
	got := h.String()

	if want != got {
		t.Errorf("For day 0: want the message to be %s,\n got %s", want, got)
	}
}

func TestFetchHabitReturnsANewHabitWithZeroDaysStreakOnNewHabit(t *testing.T) {
	h := habit.FetchHabit("piano")
	want := 0
	got := h.Streak

	if want != got {
		t.Errorf("For a new habit want %d,\n got %d", want, got)
	}
}
