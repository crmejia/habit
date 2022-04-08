package habit

import (
	"fmt"
	"os"
)

const (
	help_intro = "habit is an application to assist you in building habits"
)

func RunCLI() {
	if len(os.Args) > 2 {
		fmt.Println("too many args")
		return
	}

	if len(os.Args) == 1 {
		fmt.Println("display all habits")
		return
	}

	ht := NewTracker()
	habit := ht.FetchHabit(os.Args[1])
	ht.writeFile()
	fmt.Println(habit)
}
