package habit

import (
	"flag"
	"fmt"
	"io"
	"log"
)

const (
	frequency_usage = "Set the frecuency of the habit: daily(default), weekly."
	shorthand       = " (shorthand)"
	help_intro      = `habit is an application to assist you in building habits
     habit <options> <HABIT_NAME> -- to create/update a new habit
habit  -- to list all habits`
)

func RunCLI(filename string, args []string, output io.Writer) {

	ht := NewTracker(filename)
	if len(args) == 0 {
		fmt.Fprintln(output, ht.AllHabits())
		return
	}
	flagSet := flag.NewFlagSet("habit", flag.ExitOnError)
	flagSet.SetOutput(output)

	var frequency string
	flagSet.StringVar(&frequency, "frequency", "daily", frequency_usage)
	flagSet.StringVar(&frequency, "f", "daily", frequency_usage+shorthand)

	flagSet.Parse(args)

	if len(flagSet.Args()) > 1 {
		fmt.Fprintln(output, "too many args")
		fmt.Fprintln(output, help_intro)
		flagSet.Usage()
		return
	}
	if len(flagSet.Args()) == 0 {
		// no habit specified
		fmt.Fprintln(output, help_intro)
		flagSet.Usage()
		return
	}

	habitName := flagSet.Args()[0]
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

	err := ht.WriteFile(filename)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(habit)
}
