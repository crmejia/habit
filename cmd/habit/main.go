package main

import (
	"habit"
	"os"
)

func main() {
	filename := os.Getenv("HOME") + "/.habitTracker"
	habit.RunCLI(filename, os.Args[1:], os.Stdout)
}
