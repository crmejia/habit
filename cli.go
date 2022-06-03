package habit

import (
	"flag"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"io"
)

//RunCLI parses arguments and runs habit tracker
func RunCLI(args []string, output io.Writer) {
	flagSet := flag.NewFlagSet("habit", flag.ContinueOnError)
	flagSet.SetOutput(output)
	flagSet.Usage = func() {
		fmt.Fprintln(output,
			`habit is an application to assist you in building habits
Usage: habit <Option Flags> <HABIT_NAME> -- to create/update a new habit
       habit all   --   to list all habits
Option Flags:`)
		flagSet.PrintDefaults()
	}

	frequency := flagSet.String("f", "daily", "Set the frequency of the habit: daily(default), weekly.")
	storeType := flagSet.String("s", "db", "Set the store backend for habit tracker: db(default), file")

	err := flagSet.Parse(args)
	if err != nil {
		fmt.Fprintln(output, err)
		return
	}

	if len(flagSet.Args()) == 0 {
		flagSet.Usage()
		return
	}

	if len(flagSet.Args()) > 1 {
		fmt.Fprintln(output, "too many args")
		flagSet.Usage()
		return
	}

	store, err := storeFactory(*storeType)
	if err != nil {
		fmt.Fprintln(output, err)
		flagSet.Usage()
		return
	}
	controller, err := NewController(*store)

	if flagSet.Args()[0] == "all" {
		fmt.Fprintln(output, controller.AllHabits())
		return
	}

	h, err := ParseHabit(flagSet.Args()[0], *frequency)
	if err != nil {
		fmt.Fprintln(output, err)
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

func storeFactory(storeType string) (*Store, error) {
	dir, err := homedir.Dir()
	if err != nil {
		return nil, err
	}
	var opener func(string) (Store, error)
	var source string
	switch storeType {
	case "db":
		opener = OpenDBStore
		source = dir + "/.habitTracker.db"
	case "file":
		opener = OpenFileStore
		source = dir + "/.habitTracker"
	default:
		return nil, fmt.Errorf("unknown store type %s\n", storeType)
	}
	store, err := opener(source)
	if err != nil {
		return nil, err
	}
	return &store, nil
}
