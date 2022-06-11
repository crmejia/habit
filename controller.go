package habit

import (
	"errors"
	"fmt"
	"time"
)

const (
	//DailyInterval a constant representing a time.Duration of a day
	DailyInterval = 24 * time.Hour
	//WeeklyInterval a constant representing a time.Duration of a week
	WeeklyInterval = 7 * 24 * time.Hour

	newHabit      = "Good luck with your new habit '%s'! Don't forget to do it again %s."
	streakHabit   = "Nice work: you've done the habit '%s' for %d %s in a row now. Keep it up!"
	repeatedHabit = "You already logged '%s' today. Keep it up!"
	brokeStreak   = "You last did the habit '%s' %.0f %s ago, so you're starting a new streak today. Good luck!"
	habitStatus   = "You're currently on a %d-day streak for '%s'. Stick to it!"

	NewMessage MessageKind = iota
	RepeatMessage
	StreakMessage
	BrokenMessage
)

//Controller enforces business logic on Habits
type Controller struct {
	Store Store
}

//NewController returns a new Controller which uses the given store
func NewController(store Store) (Controller, error) {
	if store == nil {
		return Controller{}, errors.New("store cannot be nil")
	}
	return Controller{Store: store}, nil
}

//Handle Creates, Delete, or Updates the provided habit based on the status
func (c Controller) Handle(input *Habit) (*Habit, error) {
	if input == nil {
		return nil, ErrNilHabit
	}

	if input.Name == "" {
		return nil, errors.New("inputHabit name cannot be empty")
	}

	h, err := c.Store.Get(input.Name)
	if err != nil {
		return nil, err
	}
	if h != nil {
		h.updateHabit()
		err = c.Store.Update(h)
		if err != nil {
			return nil, err
		}
		return h, nil
	}

	if input.Frequency != DailyInterval && input.Frequency != WeeklyInterval {
		return nil, errors.New("invalid interval")
	}
	input.Streak = 0
	input.DueDate = time.Now().Add(input.Frequency)
	input.GenerateMessage(NewMessage)
	err = c.Store.Create(input)
	if err != nil {
		return nil, err
	}

	return input, nil
}

//GetAllHabits wraps Store.GetAllHabits and returns a string representation of the existing habits
func (c Controller) GetAllHabits() string {
	allHabits := c.Store.GetAllHabits()
	if len(allHabits) == 0 {
		return "no habits have been started"
	}
	message := "Habits:\n"
	for _, h := range allHabits {
		message += fmt.Sprintf(habitStatus+"\n", h.Streak, h.Name)
	}
	return message
}

//MessageKind represents the message to be displayed
type MessageKind int

// SameDay returns true if the days are the same ignoring hours, minutes,etc
func SameDay(d1, d2 time.Time) bool {
	if d1.Year() == d2.Year() && d1.Month() == d2.Month() && d1.Day() == d2.Day() {
		return true
	}
	return false
}

func parseHabit(name, frequency string) (*Habit, error) {
	if name == "" {
		return nil, errors.New("habit testName cannot be empty")
	}

	if frequency == "" {
		return nil, errors.New("habit frequency cannot be empty")
	}

	h := Habit{Name: name}
	switch frequency {
	case "daily":
		h.Frequency = DailyInterval
	case "weekly":
		h.Frequency = WeeklyInterval
	default:
		return nil, fmt.Errorf("unknown frequency: %s", frequency)
	}
	return &h, nil
}
func (h Habit) String() string {
	return h.Message
}

func (h *Habit) updateHabit() {
	if SameDay(h.DueDate, time.Now()) {
		//increase streak
		h.Streak++
		h.DueDate = time.Now().Add(h.Frequency)
		h.GenerateMessage(StreakMessage)
	} else if SameDay(h.DueDate, time.Now().Add(h.Frequency)) {
		//repeated habit
		h.GenerateMessage(RepeatMessage)
	} else if !SameDay(h.DueDate, time.Now()) && !SameDay(h.DueDate, time.Now().Add(h.Frequency)) {
		//streak lost
		h.GenerateMessage(BrokenMessage)
		h.Streak = 0
		h.DueDate = time.Now().Add(h.Frequency)
	}
}

//GenerateMessage creates the appropriate message for a given habit.
func (h *Habit) GenerateMessage(kind MessageKind) {
	var intervalString string
	switch kind {
	case NewMessage:
		if h.Frequency == WeeklyInterval {
			intervalString = "in a week"
		} else {
			intervalString = "tomorrow"
		}
		h.Message = fmt.Sprintf(newHabit, h.Name, intervalString)
	case RepeatMessage:
		h.Message = fmt.Sprintf(repeatedHabit, h.Name)
	case StreakMessage:
		if h.Frequency == WeeklyInterval {
			intervalString = "weeks"
		} else {
			intervalString = "days"
		}
		h.Message = fmt.Sprintf(streakHabit, h.Name, h.Streak, intervalString)
	case BrokenMessage:
		sinceDuration := time.Since(h.DueDate)
		sinceDays := sinceDuration.Hours() / 24.0
		intervalString = "days"
		if h.Frequency == WeeklyInterval {
			intervalString = "weeks"
			sinceDays = (sinceDuration.Hours() / 24.0) / 7.0
		}
		h.Message = fmt.Sprintf(brokeStreak, h.Name, sinceDays, intervalString)
	}
}
