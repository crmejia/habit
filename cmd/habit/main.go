package main

import (
	"github.com/crmejia/habit"
	"os"
)

func main() {
	habit.RunCLI(os.Args[1:], os.Stdout)
}
