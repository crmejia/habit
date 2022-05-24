package habit

import (
	"errors"
	"time"
)

const (
	//NewMessage MessageKind = iota
	//RepeatMessage
	//StreakMessage
	//BrokenMessage

	DailyInterval  = 24 * time.Hour
	WeeklyInterval = 7 * 24 * time.Hour

	newHabit      = "Good luck with your new habit '%s'! Don't forget to do it again %s."
	streakHabit   = "Nice work: you've done the habit '%s' for %d %s in a row now. Keep it up!"
	repeatedHabit = "You already logged '%s' today. Keep it up!"
	brokeStreak   = "You last did the habit '%s' %.0f %s ago, so you're starting a new streak today. Good luck!"
	habitStatus   = "You're currently on a %d-day streak for '%s'. Stick to it!"
)

type Controller struct {
	Store MemoryStore
}

func NewController(store MemoryStore) Controller {
	return Controller{Store: store}
}

func (c Controller) Handle(input *Habit) (*Habit, error) {
	if input == nil {
		return nil, NilHabitError
	}

	if input.Name == "" {
		return nil, errors.New("input name cannot be empty")
	}

	h := c.Store.Get(input.Name)
	if h != nil {
		h.updateHabit()
		return h, nil
	}

	if !validInterval[input.Interval] {
		return nil, errors.New("invalid interval")
	}
	input.Streak = 0
	input.DueDate = time.Now().Add(input.Interval)
	err := c.Store.Create(input)
	if err != nil {
		return nil, err
	}

	return input, nil
}

func (h *Habit) updateHabit() {
	if SameDay(h.DueDate, time.Now()) {
		//increase streak
		h.Streak++
		h.DueDate = time.Now().Add(h.Interval)
		//h.GenerateMessage(StreakMessage)
	} else if SameDay(h.DueDate, time.Now().Add(h.Interval)) {
		//repeated habit
		//h.GenerateMessage(RepeatMessage)
	} else if !SameDay(h.DueDate, time.Now()) && !SameDay(h.DueDate, time.Now().Add(h.Interval)) {
		//streak lost
		//h.GenerateMessage(BrokenMessage)
		h.Streak = 0
		h.DueDate = time.Now().Add(h.Interval)
	}
}

// SameDay returns true if the days are the same ignoring hours, minutes,etc
func SameDay(d1, d2 time.Time) bool {
	if d1.Year() == d2.Year() && d1.Month() == d2.Month() && d1.Day() == d2.Day() {
		return true
	}
	return false
}

var validInterval = map[time.Duration]bool{
	DailyInterval:  true,
	WeeklyInterval: true,
}
