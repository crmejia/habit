package habit

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"
)

type Habit struct {
	Name    string
	Streak  int
	Period  time.Time
	message string
}

type Tracker map[string]*Habit

func NewTracker() Tracker {
	tracker := Tracker{}
	filename := os.Getenv("HOME") + "/.habitTracker"
	err := tracker.loadFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	return tracker
}

func (ht *Tracker) FetchHabit(name string) *Habit {
	habit, ok := (*ht)[name]
	if !ok { //Create
		habit = &Habit{
			Name:    name,
			Period:  Tomorrow(),
			message: fmt.Sprintf(NewHabit, name),
		}
		(*ht)[habit.Name] = habit
		return habit
	}

	if SameDay(habit.Period, time.Now()) {
		//increase streak
		habit.Streak++
		habit.Period = Tomorrow()
		habit.message = fmt.Sprintf(StreakHabit, habit.Name, habit.Streak)
	} else if SameDay(habit.Period, Tomorrow()) {
		//repeated streak
		habit.message = fmt.Sprintf(StreakHabit, habit.Name, habit.Streak)
	} else if !SameDay(habit.Period, time.Now()) && !SameDay(habit.Period, Tomorrow()) {
		//streak lost
		sinceDuration := time.Since(habit.Period)
		sinceDays := sinceDuration.Hours() / 24.0
		habit.message = fmt.Sprintf(BrokeStreak, habit.Name, sinceDays)
		habit.Streak = 0
		habit.Period = Tomorrow()
	}

	return habit
}
func (h Habit) String() string {
	return h.message
}

const (
	NewHabit    = "Good luck with your new habit '%s'! Don't forget to do it again tomorrow."
	StreakHabit = "Nice work: you've done the habit '%s' for %d days in a row now. Keep it up!"
	BrokeStreak = "You last did the habit '%s' %.0f days ago, so you're starting a new streak today. Good luck!"
)

func Tomorrow() time.Time {
	return time.Now().Add(24 * time.Hour)
}

//func SameDay returns true if the days are the same ignoring hours, mins,etc
func SameDay(d1, d2 time.Time) bool {
	if d1.Year() == d2.Year() && d1.Month() == d2.Month() && d1.Day() == d2.Day() {
		return true
	}
	return false
}

var trackerFile *os.File

//TODO test loading
//why we need to test this?
//loading from a file | db is a fundamental operation but it needs to be tested just to make sure we don't
//ship utterly broken software
//how to test this:
// divide the two behaviors that are bundled together here:
// reading the file | db
// parse the bytes to a map
func (ht *Tracker) loadFile(filename string) error {
	var err error
	trackerFile, err = os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return err
	}

	fileBytes, err := ioutil.ReadAll(trackerFile)
	if err != nil {
		return err
	}
	if len(fileBytes) > 0 {
		err = json.Unmarshal(fileBytes, ht)
		if err != nil {
			return err
		}
	}
	return nil
}

func (ht *Tracker) writeFile() error {
	if trackerFile == nil {
		return errors.New("file is not set")
	}
	fileBytes, err := json.Marshal(ht)
	if err != nil {
		return err
	}
	trackerFile.Truncate(0)
	trackerFile.Seek(0, 0)
	_, err = trackerFile.Write(fileBytes)
	if err != nil {
		return err
	}
	if err != nil {
		return err
	}
	trackerFile.Close()
	return nil
}
