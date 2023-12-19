package app

import (
	"bufio"
	"fmt"
	"go-pwm/models"
	"os"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var io = bufio.NewScanner(os.Stdin)

// function to hash passwords
func hashPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), 14)
	return string(hash), err
}

// function to compare password to hash
func comparePassword(pwd string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd))
	return err == nil
}

// help function
func Help() {
	fmt.Println(`Usage: go-pwm [command] [args]

Commands:

  help           : display help message
  configure      : setup master password
  login          : login to go-pwm
  add [label]  : add credentials for new service
  list           : list all services
  get [label]  : get credentials for the service
  edit [label] : edit credentials for a service
  rm [label]   : remove a service
    `)
}

func Configure(db *gorm.DB) {

	// check if master password already exists
	if models.CheckMasterExists(db) {
		// check session is valid before continuing
		if !models.CheckSessionValid(db) {
			fmt.Println("Please login before changing master password.")
			os.Exit(0)
		}
	}

	// get desired master password
	fmt.Print("Enter a new password: ")
	io.Scan()
	pwd := io.Text()

	// get confirmation
	fmt.Print("Confirm password: ")
	io.Scan()
	confirm := io.Text()

	// check if passwords match
	if pwd != confirm {
		fmt.Println("Passwords do not match.")
		os.Exit(1)
	}

	// hash password
	hash, err := hashPassword(pwd)
	if err != nil {
		fmt.Println("Error hashing password.")
		os.Exit(1)
	}

	// save the password
	models.SaveMasterPassword(hash, db)
}

func Login(db *gorm.DB) {
	// check if master password already exists
	if !models.CheckMasterExists(db) {
		fmt.Println("Please configure a master password before logging in.")
		os.Exit(0)
	}

	// get master password
	fmt.Print("Enter master password: ")
	io.Scan()
	pwd := io.Text()

	// check if password is correct
	var master models.Master
	db.First(&master)
	if comparePassword(pwd, master.Password) {
		fmt.Println("Login successful.")
		models.UpdateLoginTime(db) // update the login time
	} else {
		fmt.Println("Incorrect password.")
	}
}

func Add(service string, db *gorm.DB) {
	// check if logged in
	if !models.CheckSessionValid(db) {
		fmt.Println("Please login before adding credentials.")
		os.Exit(0)
	}

	// get username
	fmt.Print("Enter username: ")
	io.Scan()
	username := io.Text()

	// get password
	fmt.Print("Enter password: ")
	io.Scan()
	password := io.Text()

	// save credentials
	models.SaveCredentials(service, username, password, db)
}

// list all services available
func List(db *gorm.DB) {
	services := models.GetServices(db)
	for i, service := range services {
		fmt.Printf("%d. %s\n", i+1, service)
	}
}