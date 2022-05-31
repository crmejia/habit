package habit

import (
	"flag"
	"fmt"
	"io"
)

const (
	frequencyUsage = "Set the frequency of the habit: daily, weekly."
	helpIntro      = `habit is an application to assist you in building habits
     habit <options> <HABIT_NAME> -- to create/update a new habit
habit  -- to list all habits`
)

//RunCLI parses arguments and runs habit tracker
func RunCLI(args []string, output io.Writer) error {
	flagSet := flag.NewFlagSet("habit", flag.ContinueOnError)
	flagSet.SetOutput(output)

	frequency := flagSet.String("f", "daily", frequencyUsage)

	err := flagSet.Parse(args)
	if err != nil {
		return err
	}

	if len(flagSet.Args()) > 1 {
		fmt.Fprintln(output, "too many args")
		fmt.Fprintln(output, helpIntro)
		flagSet.Usage()
		return nil
	}

	store := OpenMemoryStore()
	controller := NewController(store)
	if len(args) == 0 {
		allHabits := controller.AllHabits()
		if allHabits == "" {
			fmt.Fprintln(output, helpIntro)
			flagSet.Usage()
			return nil
		}
		fmt.Fprintln(output, allHabits)
		return nil
	}

	if len(flagSet.Args()) == 0 {
		// no habit specified
		fmt.Fprintln(output, helpIntro)
		flagSet.Usage()
		return nil
	}

	h, err := ParseHabit(flagSet.Args()[0], *frequency)
	if err != nil {
		return err
	}

	h, err = controller.Handle(h)
	if err != nil {
		return err
	}
	fmt.Fprint(output, h)

	return nil
}
