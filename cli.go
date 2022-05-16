package habit

import (
	"flag"
	"fmt"
	"io"
	"log"
)

const (
	frequencyUsage  = "Set the frequency of the habit: daily, weekly."
	serverModeUsage = "Runs habit as a HTTP Server"
	helpIntro       = `habit is an application to assist you in building habits
     habit <options> <HABIT_NAME> -- to create/update a new habit
habit  -- to list all habits`
)

//RunCLI parses arguments and runs habit tracker
func RunCLI(filename string, args []string, output io.Writer) {
	flagSet := flag.NewFlagSet("habit", flag.ContinueOnError)
	flagSet.SetOutput(output)

	frequency := flagSet.String("f", "daily", frequencyUsage)
	serverMode := flagSet.Bool("s", false, serverModeUsage)

	err := flagSet.Parse(args)
	if err != nil {
		log.Fatal(err)
	}

	if len(flagSet.Args()) > 1 {
		fmt.Fprintln(output, "too many args")
		fmt.Fprintln(output, helpIntro)
		flagSet.Usage()
		return
	}
	store := NewFileStore(filename)
	tracker, err := store.Load()
	if err != nil {
		log.Fatal(err)
	}

	if len(args) == 0 {
		if len(tracker) > 0 {
			fmt.Fprintln(output, AllHabits(store))
			return
		}
		// no previous habits
		fmt.Fprintln(output, helpIntro)
		flagSet.Usage()
		return
	}

	if len(flagSet.Args()) == 0 && !(*serverMode) {
		// no habit specified
		fmt.Fprintln(output, helpIntro)
		flagSet.Usage()
		return
	}

	if *serverMode {
		address := defaultTCPAddress
		if len(flagSet.Args()) > 0 {
			address = flagSet.Args()[0]
		}
		runHTTPServer(store, address)
	} else {
		habitName := flagSet.Args()[0]
		runCLI(store, habitName, *frequency, output)
	}
}

func runCLI(store Store, habitName, frequency string, output io.Writer) {
	ht := NewTracker(store)
	habit, ok := ht.FetchHabit(habitName)

	if !ok {
		habit = &Habit{
			Name: habitName,
		}
		switch frequency {
		case "daily":
			habit.Interval = DailyInterval
		case "weekly":
			habit.Interval = WeeklyInterval
		default:
			fmt.Fprintf(output, "unknown frequency %s", frequency)
			return
		}
		err := ht.CreateHabit(habit)
		if err != nil {
			log.Fatal(err)
		}
	}

	err := store.Write(&ht)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprint(output, habit)
}

func runHTTPServer(store Store, address string) {
	server := NewServer(store, address)
	server.Run()
}
