package habit

import (
	"flag"
	"fmt"
	"log"
	"os"
)

const (
	frequency_usage = "Set the frecuency of the habit: daily(default), weekly."
	shorthand       = " (shorthand)"
	help_intro      = `habit is an application to assist you in building habits
     habit <options> <HABIT_NAME> -- to create/update a new habit
habit  -- to list all habits`
)

func RunCLI() {
	var frequency string
	flag.StringVar(&frequency, "frequency", "daily", frequency_usage)
	flag.StringVar(&frequency, "f", "daily", frequency_usage+shorthand)

	ht := NewTracker()
	if len(os.Args) == 1 {
		fmt.Print(ht.AllHabits())
		return
	}
	flag.Parse()

	if len(flag.Args()) > 1 {
		fmt.Println("too many args")
		fmt.Println(help_intro)
		flag.Usage()
		return
	}

	habitName := flag.Args()[0]
	habit, ok := ht.FetchHabit(habitName)

	if !ok {
		habit = &Habit{
			Name: habitName,
		}
		if frequency == "daily" {
			habit.Interval = DailyInterval
		} else if frequency == "weekly" {
			habit.Interval = WeeklyInterval
		} else {
			fmt.Printf("unknown frecuency %s", frequency)
			return
		}
		ht.CreateHabit(habit)
	}

	err := ht.writeFile()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(habit)
}
