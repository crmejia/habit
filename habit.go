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
	message  string
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
		habit.message = fmt.Sprintf(streakHabit, habit.Name, habit.Streak)
	} else if SameDay(habit.DueDate, time.Now().Add(habit.Interval)) {
		//repeated habit
		habit.message = fmt.Sprintf(repeatedHabit, habit.Name)
	} else if !SameDay(habit.DueDate, time.Now()) && !SameDay(habit.DueDate, time.Now().Add(habit.Interval)) {
		//streak lost
		sinceDuration := time.Since(habit.DueDate)
		sinceDays := sinceDuration.Hours() / 24.0
		habit.message = fmt.Sprintf(brokeStreak, habit.Name, sinceDays)
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
	habit.message = fmt.Sprintf(newHabit, habit.Name)
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
	return h.message
}

//TODO improve messaging for weekly intervals
// could also add variation to messages~
const (
	newHabit      = "Good luck with your new habit '%s'! Don't forget to do it again tomorrow."
	repeatedHabit = "You already logged '%s' today. Keep it up!"
	streakHabit   = "Nice work: you've done the habit '%s' for %d days in a row now. Keep it up!"
	brokeStreak   = "You last did the habit '%s' %.0f days ago, so you're starting a new streak today. Good luck!"
	habitStatus   = "You're currently on a %d-day streak for '%s'. Stick to it!"
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
