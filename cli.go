package habit

import (
	"flag"
	"fmt"
	"io"
	"log"
)

func RunCLI(filename string, args []string, output io.Writer) {
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
	store := NewFileStore(filename)
	tracker, err := store.Load()
	if err != nil {
		log.Fatal(err)
	}
	//TODO move this to runCLI or at least move after parsing
	//this might mean making changing to AllHabits(store) instead of a method on Tracker
	if len(args) == 0 {
		if len(tracker) > 0 {
			fmt.Fprintln(output, AllHabits(store))
			return
		} else {
			// no previous habits
			fmt.Fprintln(output, help_intro)
			flagSet.Usage()
			return
		}
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
		runHTTPServer(store, address)
	} else {
		habitName := flagSet.Args()[0]
		runCLI(store, habitName, frequency)
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

func runCLI(store Storable, habitName, frequency string) {
	ht := NewTracker(store)
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

	err := store.Write(&ht)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(habit)
}

func runHTTPServer(store Storable, address string) {
	server := NewServer(store, address)
	server.Run()
}
