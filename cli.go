package habit

import (
	"flag"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"io"
)

const (
	frequencyUsage = "Set the frequency of the habit: daily(default), weekly."
	storeTypeUsage = "Set the store backend for habit tracker: db(default), file"
	helpIntro      = `habit is an application to assist you in building habits
     habit <options> <HABIT_NAME> -- to create/update a new habit
habit  -- to list all habits`
)

//RunCLI parses arguments and runs habit tracker
func RunCLI(args []string, output io.Writer) {
	flagSet := flag.NewFlagSet("habit", flag.ContinueOnError)
	flagSet.SetOutput(output)

	frequency := flagSet.String("f", "daily", frequencyUsage)
	storeType := flagSet.String("s", "db", storeTypeUsage)

	err := flagSet.Parse(args)
	if err != nil {
		fmt.Fprintln(output, err)
		return
	}

	if len(flagSet.Args()) > 1 {
		fmt.Fprintln(output, "too many args")
		fmt.Fprintln(output, helpIntro)
		flagSet.Usage()
		return
	}
	dir, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
	}

	var store Store
	var opener func(string) (Store, error)
	var source string
	switch *storeType {
	case "db":
		opener = OpenDBStore
		source = dir + "/.habitTracker.db"
	case "file":
		opener = OpenFileStore
		source = dir + "/.habitTracker"
	default:
		fmt.Fprintf(output, "unknown store type %s\n", *storeType)
		fmt.Fprintln(output, helpIntro)
		flagSet.Usage()
		return
	}
	store, err = opener(source)
	if err != nil {
		fmt.Fprintln(output, err)
	}

	controller, err := NewController(store)
	if err != nil {
		fmt.Fprintln(output, err)
	}
	if len(args) == 0 {
		allHabits := controller.AllHabits()
		if allHabits == "" {
			fmt.Fprintln(output, helpIntro)
			flagSet.Usage()
			return
		}
		fmt.Fprintln(output, allHabits)
		return
	}
	if len(flagSet.Args()) == 0 {
		fmt.Fprintln(output, helpIntro)
		flagSet.Usage()
		return
	}

	h, err := ParseHabit(flagSet.Args()[0], *frequency)
	if err != nil {
		fmt.Fprintln(output, err)
		fmt.Fprintln(output, helpIntro)
		flagSet.Usage()
		return
	}

	h, err = controller.Handle(h)
	if err != nil {
		fmt.Fprintln(output, err)
		return
	}
	fmt.Fprintln(output, h)

	return
}
