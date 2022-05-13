package habit_test

import (
	"habit"
	"os"
	"testing"
	"time"
)

func TestFileStoreRoundtripWriteRead(t *testing.T) {
	t.Parallel()
	writeTracker := habit.Tracker{
		"piano": &habit.Habit{
			Name:     "piano",
			Interval: habit.WeeklyInterval,
			Streak:   1,
			DueDate:  time.Now().Add(habit.WeeklyInterval),
		},
	}
	tmpFile := CreateTmpFile(t)

	writeFileStore := habit.NewFileStore(tmpFile.Name())
	writeFileStore.Write(&writeTracker)

	loadFileStore := habit.NewFileStore(tmpFile.Name())
	loadTracker, _ := loadFileStore.Load()

	_, ok := loadTracker["piano"]
	if !ok {
		t.Errorf("want loaded file to contain the same habit that was written")
	}
}

func TestDBStoreCreatesHabitTableIfItDoesNotExist(t *testing.T) {
	t.Parallel()
	dbsource := os.TempDir() + "test.db"
	dbStore := habit.NewDBStore(dbsource)

	tracker, err := dbStore.Load()
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(dbsource)

	got := len(tracker)
	if len(tracker) > 0 {
		t.Errorf("want no habits on new db, got %d", got)
	}
}

func TestDBStoreRoundtripWriteRead(t *testing.T) {
	t.Parallel()
	writeTracker := habit.Tracker{
		"piano": &habit.Habit{
			Name:     "piano",
			Interval: habit.WeeklyInterval,
			Streak:   1,
			DueDate:  time.Now().Add(habit.WeeklyInterval),
		},
	}
	dbsource := os.TempDir() + "roundtrip.db"
	dbStore := habit.NewDBStore(dbsource)
	defer os.Remove(dbsource)

	dbStore.Write(&writeTracker)

	loadTracker := habit.Tracker{}
	loadTracker, err := dbStore.Load()
	if err != nil {
		t.Fatal(err)
	}

	_, ok := loadTracker["piano"]
	if !ok {
		t.Errorf("want loaded file to contain the same habit that was written")
	}
}

// CreateTmpFile creates a temporary file with a random name in a temporary directory. This temporary directory gets
// removed by t.Cleanup as part of the test so there no need to defer os.Remove the temp file
func CreateTmpFile(t *testing.T) *os.File {
	tmpDir := t.TempDir()
	tmpFile, err := os.CreateTemp(tmpDir, "")
	if err != nil {
		t.Fatal("couldn't create tmp file")
	}
	return tmpFile
}
