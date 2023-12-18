package main

import (
	"fmt"
	"go-pwm/app"
	"os"
)

func main() {

	// check amount of args
	if len(os.Args) < 2 {
		app.Help()
		os.Exit(0)
	}

	// get args
	args := os.Args[1:]

	switch args[0] {

	case "help":
		app.Help()
	
	default:
		fmt.Printf("'%s' is an unrecognised command. See 'go-pwm help' for a list of commands.", args[0])
		os.Exit(1)
	}
}