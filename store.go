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
	//Message  string
}

// the point of the store interface is to wrap data around the basic go map struct.
//that means that for testing a memory still will do see store_test for more comments.
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

//Get should get return a habit always? No, because that would limit the user of the package.It would mean get does
//get+create changing the signature to Get(habit) because name is only one attribute not the whole habit and the data
// needs to be UPDATE in that case, not created... See where this rabbit hole goes?
//in case the get fails the function has to go ahead and do some other two behaviors(create, update)? The answer is
//no, this is business logic, not a CRUD logic. So the user( habit manager or habit consumer, whatever other package
//leverages this one) is free to implement to meet their needs. The user gets to decide what to do with it. If the user
//wants to create a habit that doesn't exist(after trying a get) they can do so at their own discretion not my library
//limiting their use case by short vision.
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
