package habit_test

import (
	"github.com/crmejia/habit"
	"os"
	"testing"
	"time"
)

func TestMemoryStore_GetReturnsNilOnNoHabit(t *testing.T) {
	t.Parallel()
	store := habit.OpenMemoryStore()
	got, _ := store.Get("piano")

	if got != nil {
		t.Error("want Store.Get to return nil")
	}
}

func TestMemoryStore_GetReturnsExistingHabit(t *testing.T) {
	t.Parallel()
	store := habit.OpenMemoryStore()
	store.Habits["piano"] = &habit.Habit{Name: "piano"}

	h, _ := store.Get("piano")
	if h == nil {
		t.Fatal()
	}

	want := "piano"
	got := h.Name
	if want != got {
		t.Errorf("want h testName to be %s, h %s", want, got)
	}
}

func TestMemoryStore_CreateNewHabit(t *testing.T) {
	t.Parallel()
	store := habit.OpenMemoryStore()
	h := habit.Habit{Name: "piano"}

	store.Create(&h)

	if _, ok := store.Habits["piano"]; !ok {
		t.Error("want h to be inserted into store")
	}
}

func TestMemoryStore_CreateUpdateNilHabitFails(t *testing.T) {
	t.Parallel()
	store := habit.OpenMemoryStore()
	err := store.Create(nil)

	if err == nil {
		t.Error("want Store.create nil habit to fail with error")
	}

	err = store.Update(nil)
	if err == nil {
		t.Error("want Store.Update nil habit to fail with error")
	}
}

func TestMemoryStore_CreateExistingHabitFails(t *testing.T) {
	t.Parallel()
	store := habit.OpenMemoryStore()
	h := habit.Habit{Name: "piano"}
	store.Habits["piano"] = &h
	err := store.Create(&h)

	if err == nil {
		t.Error("want Store.Create nil h to fail with error")
	}
}

func TestMemoryStore_UpdateHabit(t *testing.T) {
	t.Parallel()
	store := habit.OpenMemoryStore()
	oldHabit := &habit.Habit{Name: "piano"}
	store.Habits["piano"] = oldHabit

	updateHabit := &habit.Habit{Name: "piano"}
	store.Update(updateHabit)

	if oldHabit == store.Habits["piano"] {
		t.Error("want update to replace habit")
	}
}

func TestMemoryStore_UpdateFailsOnNil(t *testing.T) {
	t.Parallel()
	store := habit.OpenMemoryStore()
	err := store.Update(nil)

	if err == nil {
		t.Error("want Store.Update nil habit to fail with error")
	}
}

func TestMemoryStore_UpdateFailsOnNonExistingHabit(t *testing.T) {
	t.Parallel()
	store := habit.OpenMemoryStore()
	err := store.Update(&habit.Habit{Name: "piano"})

	if err == nil {
		t.Error("want update to fail if habit does not exist")
	}
}

func TestMemoryStore_AllHabitsReturnsSliceOfHabits(t *testing.T) {
	t.Parallel()
	store := habit.OpenMemoryStore()
	store.Habits = map[string]*habit.Habit{
		"piano":   {Name: "piano"},
		"surfing": {Name: "surfing"},
	}

	allHabits := store.GetAllHabits()
	if len(allHabits) != len(store.Habits) {
		t.Error("want GetAllHabits to return a slice of habits")
	}
}

func TestMessageGenerator(t *testing.T) {
	testCases := []struct {
		h    habit.Habit
		kind habit.MessageKind
		want string
	}{
		{habit.Habit{Name: "piano", Frequency: habit.WeeklyInterval}, habit.NewMessage, "Good luck with your new habit 'piano'! Don't forget to do it again in a week."},
		{habit.Habit{Name: "piano", Frequency: habit.DailyInterval}, habit.NewMessage, "Good luck with your new habit 'piano'! Don't forget to do it again tomorrow."},
		{habit.Habit{Name: "surfing", Frequency: habit.WeeklyInterval}, habit.RepeatMessage, "You already logged 'surfing' today. Keep it up!"},
		{habit.Habit{Name: "meditation", Frequency: habit.DailyInterval}, habit.RepeatMessage, "You already logged 'meditation' today. Keep it up!"},
		{habit.Habit{Name: "dancing", Frequency: habit.WeeklyInterval, Streak: 2}, habit.StreakMessage, "Nice work: you've done the habit 'dancing' for 2 weeks in a row now. Keep it up!"},
		{habit.Habit{Name: "meditation", Frequency: habit.DailyInterval, Streak: 2}, habit.StreakMessage, "Nice work: you've done the habit 'meditation' for 2 days in a row now. Keep it up!"},
		{habit.Habit{Name: "running", Frequency: habit.DailyInterval, DueDate: time.Now().Add(-5 * 24 * time.Hour)}, habit.BrokenMessage, "You last did the habit 'running' 5 days ago, so you're starting a new streak today. Good luck!"},
		{habit.Habit{Name: "hiking", Frequency: habit.WeeklyInterval, DueDate: time.Now().Add(-3 * 24 * 7 * time.Hour)}, habit.BrokenMessage, "You last did the habit 'hiking' 3 weeks ago, so you're starting a new streak today. Good luck!"},
	}

	for _, tc := range testCases {
		tc.h.GenerateMessage(tc.kind)
		got := tc.h.Message
		if tc.want != got {
			t.Errorf("want Message to be:\n%s\ngot:\n%s", tc.want, got)
		}
	}
}

func TestController_HandleSetsMessageCorrectlyForNewHabit(t *testing.T) {
	t.Parallel()
	store := habit.OpenMemoryStore()
	controller, _ := habit.NewController(&store)
	h := &habit.Habit{Name: "piano",
		Frequency: habit.DailyInterval}
	controller.Handle(h)

	want := "Good luck with your new habit 'piano'! Don't forget to do it again tomorrow."
	got := h.String()
	if want != got {
		t.Errorf("For %d day streak: want the Message to be:\n%s,\n got\n%s", h.Streak, want, got)
	}
}

func TestTracker_FetchHabitSetsMessageCorrectlyForStreakBrokenStreak(t *testing.T) {
	t.Parallel()
	testCases := []struct {
		want  string
		habit *habit.Habit
	}{
		{want: "Nice work: you've done the habit 'surf' for 4 days in a row now. Keep it up!", habit: &habit.Habit{Name: "surf", Streak: 3, DueDate: time.Now()}},
		{want: "You last did the habit 'running' 10 days ago, so you're starting a new streak today. Good luck!", habit: &habit.Habit{Name: "running", Streak: 10, DueDate: time.Now().Add(-10 * 24 * time.Hour)}},
	}
	store := habit.OpenMemoryStore()
	controller, _ := habit.NewController(&store)
	for _, tc := range testCases {

		store.Habits[tc.habit.Name] = tc.habit
		controller.Handle(tc.habit)

		got := tc.habit.String()
		if tc.want != got {
			t.Errorf("For %d day streak: want the Message to be:\n%s,\n got\n%s", tc.habit.Streak, tc.want, got)
		}
	}
}

func TestOpenDBStoreErrorsOnEmptyDBSource(t *testing.T) {
	t.Parallel()
	_, err := habit.OpenDBStore("")
	if err == nil {
		t.Error("Want error on empty string db source")
	}
}

func TestDBStore_CreateUpdateNilHabitFails(t *testing.T) {
	t.Parallel()
	dbSource := t.TempDir() + "test.db"
	store, _ := habit.OpenDBStore(dbSource)

	err := store.Create(nil)

	if err == nil {
		t.Error("want Store.create nil habit to fail with error")
	}

	err = store.Update(nil)
	if err == nil {
		t.Error("want Store.Update nil habit to fail with error")
	}
}

func TestDBStore_GetCreateRoundTrip(t *testing.T) {
	t.Parallel()
	dbSource := t.TempDir() + "test.db"
	dbStore, err := habit.OpenDBStore(dbSource)
	if err != nil {
		t.Fatal(err)
	}

	h, err := dbStore.Get("piano")
	if err != nil {
		t.Fatal(err)
	}

	if h != nil {
		t.Error("expected get to return nil on empty db")
	}
	h = &habit.Habit{
		Name: "piano",
	}
	err = dbStore.Create(h)
	if err != nil {
		t.Fatal(err)
	}

	got, err := dbStore.Get("piano")
	if err != nil {
		t.Fatal(err)
	}
	if got == nil {
		t.Error("wanted habit piano. Got nil.")
	}
}

func TestDBStore_CreateUpdateRoundTrip(t *testing.T) {
	t.Parallel()
	dbSource := t.TempDir() + "test.db"
	dbStore, err := habit.OpenDBStore(dbSource)
	if err != nil {
		t.Fatal(err)
	}

	h := &habit.Habit{
		Name: "piano",
	}
	err = dbStore.Create(h)
	if err != nil {
		t.Fatal()
	}

	intermediateHabit, err := dbStore.Get("piano")
	if err != nil {
		t.Fatal(err)
	}
	if intermediateHabit == nil {
		t.Error("wanted habit piano. Got nil.")
	}

	intermediateHabit.Streak = 5
	intermediateHabit.Frequency = habit.DailyInterval
	now := time.Now().Truncate(time.Second)
	intermediateHabit.DueDate = now

	err = dbStore.Update(intermediateHabit)
	if err != nil {
		t.Fatal(err)
	}
	got, err := dbStore.Get("piano")
	if err != nil {
		t.Fatal(err)
	}
	if got.Streak != 5 || got.Frequency != habit.DailyInterval || !habit.SameDay(got.DueDate, now) {
		t.Error("wanted habit piano. To be updated.")
	}
}

func TestDBStore_AllHabits(t *testing.T) {
	t.Parallel()
	dbSource := t.TempDir() + "test.db"
	dbStore, err := habit.OpenDBStore(dbSource)
	if err != nil {
		t.Fatal(err)
	}
	habits := []*habit.Habit{
		{Name: "piano"},
		{Name: "surfing"},
	}

	for _, h := range habits {
		dbStore.Create(h)
	}

	got := dbStore.GetAllHabits()
	if len(got) != len(habits) {
		t.Errorf("want GetAllHabits to return %d habits, got %d", len(habits), len(got))
	}
}

func TestOpenFileStoreCreatesFile(t *testing.T) {
	t.Parallel()
	filename := t.TempDir() + ".habitTracker"
	_, err := habit.OpenFileStore(filename)
	if err != nil {
		t.Fatal(err)
	}
	_, err = os.Stat(filename)
	if err != nil {
		t.Error("wanted file to be created by OpenFileStore")
	}
}

func TestOpenFileStoreErrorsOnEmptyFilename(t *testing.T) {
	t.Parallel()
	_, err := habit.OpenFileStore("")
	if err == nil {
		t.Error("Want error on empty string db source")
	}
}
func TestFileStore_CreateUpdateNilHabitFails(t *testing.T) {
	t.Parallel()
	filename := t.TempDir() + ".habitTracker"
	store, _ := habit.OpenFileStore(filename)
	err := store.Create(nil)
	if err == nil {
		t.Error("want Store.create nil habit to fail with error")
	}

	err = store.Update(nil)
	if err == nil {
		t.Error("want Store.Update nil habit to fail with error")
	}
}

func TestFileStore_GetReturnsNilOnNoHabit(t *testing.T) {
	t.Parallel()
	filename := t.TempDir() + ".habitTracker"
	fileStore, err := habit.OpenFileStore(filename)
	if err != nil {
		t.Fatal(err)
	}
	h, err := fileStore.Get("piano")
	if err != nil {
		t.Fatal(err)
	}
	if h != nil {
		t.Error("want Store.Get to return nil")
	}
}

func TestFileStore_GetCreateRoundTrip(t *testing.T) {
	t.Parallel()
	filename := t.TempDir() + ".habitTracker"
	fileStore, err := habit.OpenFileStore(filename)
	if err != nil {
		t.Fatal(err)
	}

	h, err := fileStore.Get("piano")
	if err != nil {
		t.Fatal(err)
	}

	if h != nil {
		t.Error("expected get to return nil on empty db")
	}
	h = &habit.Habit{
		Name: "piano",
	}
	err = fileStore.Create(h)
	if err != nil {
		t.Fatal(err)
	}

	got, err := fileStore.Get("piano")
	if err != nil {
		t.Fatal(err)
	}
	if got == nil {
		t.Error("wanted habit piano. Got nil.")
	}
}

func TestFileStore_CreateUpdateRoundTrip(t *testing.T) {
	t.Parallel()
	filename := t.TempDir() + ".habitTracker"
	fileStore, err := habit.OpenFileStore(filename)
	if err != nil {
		t.Fatal(err)
	}

	h := &habit.Habit{
		Name: "piano",
	}
	err = fileStore.Create(h)
	if err != nil {
		t.Fatal()
	}

	intermediateHabit, err := fileStore.Get("piano")
	if err != nil {
		t.Fatal(err)
	}
	if intermediateHabit == nil {
		t.Error("wanted habit piano. Got nil.")
	}

	intermediateHabit.Streak = 5
	intermediateHabit.Frequency = habit.DailyInterval
	now := time.Now().Truncate(time.Second)
	intermediateHabit.DueDate = now

	err = fileStore.Update(intermediateHabit)
	if err != nil {
		t.Fatal(err)
	}
	got, err := fileStore.Get("piano")
	if err != nil {
		t.Fatal(err)
	}
	if got.Streak != 5 || got.Frequency != habit.DailyInterval || !habit.SameDay(got.DueDate, now) {
		t.Error("wanted habit piano. To be updated.")
	}
}

func TestFileStore_AllHabitsReturnsSliceOfHabits(t *testing.T) {
	t.Parallel()
	filename := t.TempDir() + ".habitTracker"
	fileStore, err := habit.OpenFileStore(filename)
	if err != nil {
		t.Fatal(err)
	}

	habits := []*habit.Habit{
		{Name: "piano"},
		{Name: "surfing"},
	}

	for _, h := range habits {
		fileStore.Create(h)
	}

	got := fileStore.GetAllHabits()
	if len(got) != len(habits) {
		t.Errorf("want GetAllHabits to return %d habits, got %d", len(habits), len(got))
	}
}
