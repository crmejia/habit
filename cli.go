package habit

import (
	"flag"
	"fmt"
	"io"
	"log"
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

	var serverMode bool
	flagSet.BoolVar(&serverMode, "server", false, serverMode_usage)
	flagSet.BoolVar(&serverMode, "s", false, serverMode_usage+shorthand)

	flagSet.Parse(args)

	if len(flagSet.Args()) > 1 {
		fmt.Fprintln(output, "too many args")
		fmt.Fprintln(output, help_intro)
		flagSet.Usage()
		return
	}
	if len(flagSet.Args()) == 0 && !serverMode {
		// no habit specified
		fmt.Fprintln(output, help_intro)
		flagSet.Usage()
		return
	}
	if serverMode {
		address := defaultTCPAddress
		if len(flagSet.Args()) > 0 {
			address = flagSet.Args()[0]
		}
		runHTTPServer(filename, address)
	} else {
		habitName := flagSet.Args()[0]
		runCLI(filename, habitName, frequency, &ht)
	}
	return
}

const (
	frequency_usage  = "Set the frecuency of the habit: daily, weekly."
	serverMode_usage = "Runs habit as a HTTP Server"
	shorthand        = " (shorthand)"
	help_intro       = `habit is an application to assist you in building habits
     habit <options> <HABIT_NAME> -- to create/update a new habit
habit  -- to list all habits`
)

func runCLI(filename, habitName, frequency string, ht *Tracker) {

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

func runHTTPServer(filename, address string) {
	server := NewServer(filename, address)
	server.Run()
}
