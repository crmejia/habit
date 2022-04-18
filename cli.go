package habit

import (
	"fmt"
	"log"
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

	ht := NewTracker()
	if len(os.Args) == 1 {
		fmt.Print(ht.AllHabits())
		return
	}
	//TODO parse weekly habit
	habit := ht.FetchHabit(os.Args[1])
	err := ht.writeFile()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(habit)
}
