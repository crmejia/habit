package main

import (
	"habit"
	"os"
)

func main() {
	habit.RunCLI(os.Args[1:], os.Stdout)
}
