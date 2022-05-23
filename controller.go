package habit

import "errors"

type Controller struct {
	Store MemoryStore
}

func NewController(store MemoryStore) Controller {
	return Controller{Store: store}
}

func (c Controller) Handle(input *Habit) (*Habit, error) {
	if input == nil {
		return nil, NilHabitError
	}

	if input.Name == "" {
		return nil, errors.New("input name cannot be empty")
	}

	h := c.Store.Get(input.Name)
	if h != nil {
		return input, nil
	}

	err := c.Store.Create(input)

	if err != nil {
		return nil, err
	}

	return input, nil
}
