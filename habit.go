package habit

import "fmt"

type Habit struct {
	name   string
	Streak int
}

func FetchHabit(name string) Habit {
	return Habit{name: name}
}
func (h Habit) String() string {
	return fmt.Sprintf(messages[0], h.name)
}

var messages = map[int]string{
	0: "Good luck with your new habit '%s'! Don't forget to do it again\ntomorrow.",
}
