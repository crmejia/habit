package main

import (
	"fmt"
	"github.com/mitchellh/go-homedir"
	"habit"
	"os"
)

func main() {
	dir, err := homedir.Dir()
	if err != nil {
		fmt.Println(err)
	}
	store, _ := habit.OpenDBStore(dir + "/.habitTracker.db")
	habit.RunCLI(os.Args[1:], os.Stdout, &store)
}
