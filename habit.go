package habit

import (
	"fmt"
	"time"
)

type Habit struct {
	Name   string
	Streak int
	Period time.Time
}
type Tracker map[string]Habit

func (ht *Tracker) FetchHabit(name string) Habit {
	habit, ok := (*ht)[name]
	if !ok { //Create
		habit = Habit{
			Name:   name,
			Period: Tomorrow(),
		}
		(*ht)[habit.Name] = habit
		return habit
	}

	if SameDay(habit.Period, time.Now()) {
		habit.Streak++
		habit.Period = Tomorrow()
	}
	return habit
}
func (h Habit) String() string {
	var message string
	if h.Streak == 0 {
		message = fmt.Sprintf(NewHabit, h.Name)
	} else if h.Streak > 0 {
		message = fmt.Sprintf(StreakHabit, h.Name, h.Streak)
	}
	return message
}

const (
	NewHabit    = "Good luck with your new habit '%s'! Don't forget to do it again tomorrow."
	StreakHabit = "Nice work: you've done the habit '%s' for %d days in a row now. Keep it up!"
)

//var messages = map[string]string{
//	"newHabit": "Good luck with your new habit '%s'! Don't forget to do it again\ntomorrow.",
//	1:          "Nice work: you've done the habit '%s' for %s days in a row now.\nKeep it up!",
//}

func Tomorrow() time.Time {
	return time.Now().Add(24 * time.Hour)
}

//func Same Day returns true if the days are the same ignoring hours, mins,etc
func SameDay(d1, d2 time.Time) bool {
	if d1.Year() == d2.Year() && d1.Month() == d2.Month() && d1.Day() == d2.Day() {
		return true
	}
	return false
}
