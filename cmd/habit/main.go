package main

import (
	"habit"
	"os"
)

func main() {
	filename := os.Getenv("HOME") + "/.habitTracker"
	habit.RunCLI(filename, os.Args[1:], os.Stdout)
}

//TODO fix installation go install https://github.com/crmejia/habit/cmd/habit@latest
//TODO add HTTP server documentation
//TODO add DB store
