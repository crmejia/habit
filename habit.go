package habit

import (
	"fmt"
	"log"
	"time"
)

import (
	"errors"
)

type Habit struct {
	Name     string
	Streak   int
	DueDate  time.Time
	Interval time.Duration
	Message  string
}

type Tracker map[string]*Habit

func NewTracker(filename string) Tracker {
	tracker := Tracker{}

	err := tracker.LoadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return tracker
}

func (ht *Tracker) FetchHabit(name string) (*Habit, bool) {
	habit, ok := (*ht)[name]
	if !ok {
		return nil, false
	}
	if SameDay(habit.DueDate, time.Now()) {
		//increase streak
		habit.Streak++
		habit.DueDate = time.Now().Add(habit.Interval)
		habit.GenerateMessage(StreakMessage)
	} else if SameDay(habit.DueDate, time.Now().Add(habit.Interval)) {
		//repeated habit
		habit.GenerateMessage(RepeatMessage)
	} else if !SameDay(habit.DueDate, time.Now()) && !SameDay(habit.DueDate, time.Now().Add(habit.Interval)) {
		//streak lost
		habit.GenerateMessage(BrokenMessage)
		habit.Streak = 0
		habit.DueDate = time.Now().Add(habit.Interval)
	}

	return habit, true
}

func (ht *Tracker) CreateHabit(habit *Habit) error {
	_, ok := (*ht)[habit.Name]
	if ok {
		return errors.New("habit already exists")
	}
	if !validInterval[habit.Interval] {
		return errors.New("not a valid interval")
	}
	habit.DueDate = time.Now().Add(habit.Interval)
	habit.GenerateMessage(NewMessage)
	(*ht)[habit.Name] = habit
	return nil
}
func (ht *Tracker) AllHabits() string {
	message := "Habits:\n"
	for _, habit := range *ht {
		message += fmt.Sprintf(habitStatus+"\n", habit.Streak, habit.Name)
	}
	return message
}

func (h Habit) String() string {
	return h.Message
}

//TODO add variation to messages
func (h *Habit) GenerateMessage(kind MessageKind) {
	var intervalString string
	switch kind {
	case NewMessage:
		if h.Interval == WeeklyInterval {
			intervalString = "in a week"
		} else {
			intervalString = "tomorrow"
		}
		h.Message = fmt.Sprintf(newHabit, h.Name, intervalString)
	case RepeatMessage:
		h.Message = fmt.Sprintf(repeatedHabit, h.Name)
	case StreakMessage:
		if h.Interval == WeeklyInterval {
			intervalString = "weeks"
		} else {
			intervalString = "days"
		}
		h.Message = fmt.Sprintf(streakHabit, h.Name, h.Streak, intervalString)
	case BrokenMessage:
		sinceDuration := time.Since(h.DueDate)
		sinceDays := sinceDuration.Hours() / 24.0
		intervalString = "days"
		if h.Interval == WeeklyInterval {
			intervalString = "weeks"
			sinceDays = (sinceDuration.Hours() / 24.0) / 7.0
		}
		h.Message = fmt.Sprintf(brokeStreak, h.Name, sinceDays, intervalString)
	}
}

const (
	newHabit      = "Good luck with your new habit '%s'! Don't forget to do it again %s."
	streakHabit   = "Nice work: you've done the habit '%s' for %d %s in a row now. Keep it up!"
	repeatedHabit = "You already logged '%s' today. Keep it up!"
	brokeStreak   = "You last did the habit '%s' %.0f %s ago, so you're starting a new streak today. Good luck!"
	habitStatus   = "You're currently on a %d-day streak for '%s'. Stick to it!"
)

type MessageKind int

const (
	NewMessage MessageKind = iota
	RepeatMessage
	StreakMessage
	BrokenMessage
)

//func SameDay returns true if the days are the same ignoring hours, mins,etc
func SameDay(d1, d2 time.Time) bool {
	if d1.Year() == d2.Year() && d1.Month() == d2.Month() && d1.Day() == d2.Day() {
		return true
	}
	return false
}

const (
	DailyInterval  time.Duration = 24 * time.Hour
	WeeklyInterval               = 7 * 24 * time.Hour
)

var validInterval = map[time.Duration]bool{
	DailyInterval:  true,
	WeeklyInterval: true,
}
