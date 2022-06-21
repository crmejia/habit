package main

import (
	"github.com/crmejia/habit"
	"os"
)

func main() {
	habit.RunServer(os.Args[1:], os.Stdout)
}
