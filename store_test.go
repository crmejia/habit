package habit_test

import (
	"habit"
	"log"
	"os"
	"testing"
	"time"
)

func TestRoundtripWriteRead(t *testing.T) {
	t.Parallel()
	writeTracker := habit.Tracker{
		"piano": &habit.Habit{
			Name:     "piano",
			Interval: habit.WeeklyInterval,
			Streak:   1,
			DueDate:  time.Now().Add(habit.WeeklyInterval),
		},
	}
	tmpFile := tmpFile()
	defer os.Remove(tmpFile.Name())
	writeTracker.WriteFile(tmpFile.Name())

	loadTracker := habit.Tracker{}
	loadTracker.LoadFile(tmpFile.Name())

	_, ok := loadTracker["piano"]
	if !ok {
		t.Errorf("want loaded file to contain the same habit that was written")
	}
}
func tmpFile() *os.File {
	tmpFile, err := os.CreateTemp("", "")
	if err != nil {
		log.Fatal("couldn't create tmp file")
	}
	defer os.Remove(tmpFile.Name())
	return tmpFile
}
