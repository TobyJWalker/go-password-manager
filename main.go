package main

import (
	"fmt"
	"go-pwm/app"
	"go-pwm/models"
	"os"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {

	// check amount of args
	if len(os.Args) < 2 {
		app.Help()
		os.Exit(0)
	}

	// check if data folder exists
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		os.Mkdir("data", 0755)
	}

	// attempt connection to database
	db, err := gorm.Open(sqlite.Open("data/go-pwm.sqlite"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&models.Master{}, &models.Credential{})

	// get args
	args := os.Args[1:]

	switch args[0] {

	case "help":
		app.Help()

	case "configure":
		app.Configure(db)

	case "login":
		app.Login(db)

	case "add":
		if len(args) < 2 {
			fmt.Println("Please specify a service.")
		} else {
			app.Add(args[1], db)
		}
	
	case "list":
		app.List(db)

	case "get":
		if len(args) < 2 {
			fmt.Println("Please specify a service.")
		} else {
			app.Get(args[1], db)
		}

	case "rm":
		if len(args) < 2 {
			fmt.Println("Please specify a service.")
		} else {
			app.Remove(args[1], db)
		}
	
	default:
		fmt.Printf("'%s' is an unrecognised command. See 'go-pwm help' for a list of commands.", args[0])
	}
}