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
	tmpFile := CreateTmpFile()
	defer os.Remove(tmpFile.Name())

	writeFileStore := habit.NewFileStore(tmpFile.Name())
	writeFileStore.Write(writeTracker)

	loadFileStore := habit.NewFileStore(tmpFile.Name())
	loadTracker, _ := loadFileStore.Load()

	_, ok := loadTracker["piano"]
	if !ok {
		t.Errorf("want loaded file to contain the same habit that was written")
	}
}

//func TestRoundtripDBWriteRead(t *testing.T) {
//	t.Parallel()
//	writeTracker := habit.Tracker{
//		"piano": &habit.Habit{
//			Name:     "piano",
//			Interval: habit.WeeklyInterval,
//			Streak:   1,
//			DueDate:  time.Now().Add(habit.WeeklyInterval),
//		},
//	}
//
//	writeTracker.WriteDB()
//
//	loadTracker := habit.Tracker{}
//	loadTracker.LoadDB()
//
//	_, ok := loadTracker["piano"]
//	if !ok {
//		t.Errorf("want loaded file to contain the same habit that was written")
//	}
//}

func CreateTmpFile() *os.File {
	tmpFile, err := os.CreateTemp("", "")
	if err != nil {
		log.Fatal("couldn't create tmp file")
	}
	defer os.Remove(tmpFile.Name())
	return tmpFile
}
