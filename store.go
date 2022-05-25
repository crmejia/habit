package habit

import (
	"errors"
	"time"
)

type Habit struct {
	Name     string
	Streak   int
	DueDate  time.Time
	Interval time.Duration
	Message  string
}

type MemoryStore struct {
	Habits map[string]*Habit
}

var NilHabitError = errors.New("habit cannot be nil")

func OpenStore() MemoryStore {
	//here a file store or a db store would get the data from persistence.
	memoryStore := MemoryStore{
		Habits: map[string]*Habit{},
	}
	return memoryStore
}

func (s *MemoryStore) Get(name string) *Habit {
	habit, ok := s.Habits[name]
	if ok {
		return habit
	}
	return nil
}

func (s *MemoryStore) Create(habit *Habit) error {
	if habit == nil {
		return NilHabitError
	}

	if _, ok := s.Habits[habit.Name]; ok {
		return errors.New("habit already exists")
	}
	s.Habits[habit.Name] = habit
	return nil
}

func (s *MemoryStore) Update(habit *Habit) error {
	if habit == nil {
		return NilHabitError
	}

	if _, ok := s.Habits[habit.Name]; !ok {
		return errors.New("cannot update habit does not exists")
	}

	s.Habits[habit.Name] = habit
	return nil
}

func (s MemoryStore) AllHabits() []*Habit {
	allHabits := make([]*Habit, 0, len(s.Habits))
	for _, h := range s.Habits {
		allHabits = append(allHabits, h)
	}
	return allHabits
}
