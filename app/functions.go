package app

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
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

// get a services details
func Get(service string, db *gorm.DB) {
	// check if logged in
	if !models.CheckSessionValid(db) {
		fmt.Println("Please login before getting credentials.")
		os.Exit(0)
	}

	// get credentials
	uname, pwd := models.GetCredentials(service, db)
	fmt.Printf("Username: %s\n", uname)
	fmt.Printf("Password: %s\n", pwd)
}

// remove a service
func Remove(service string, db *gorm.DB) {
	// check if logged in
	if !models.CheckSessionValid(db) {
		fmt.Println("Please login before removing credentials.")
		os.Exit(0)
	}

	// check service exists
	if !models.CheckServiceExists(service, db) {
		fmt.Println("Service not found.")
		os.Exit(0)
	}

	// delete credentials
	models.DeleteCredentials(service, db)
}

// edit a service
func Edit(service string, db *gorm.DB) {
	// check if logged in
	if !models.CheckSessionValid(db) {
		fmt.Println("Please login before editing credentials.")
		os.Exit(0)
	}

	// check service exists
	if !models.CheckServiceExists(service, db) {
		fmt.Println("Service not found.")
		os.Exit(0)
	}

	// get username
	fmt.Print("Enter new username: ")
	io.Scan()
	username := io.Text()

	// get password
	fmt.Print("Enter new password: ")
	io.Scan()
	password := io.Text()

	// save credentials
	models.EditCredentials(service, username, password, db)
}

// export a service creds
func Export(service string, db *gorm.DB) {
	// check if logged in
	if !models.CheckSessionValid(db) {
		fmt.Println("Please login before exporting credentials.")
		os.Exit(0)
	}

	// check service exists
	if !models.CheckServiceExists(service, db) {
		fmt.Println("Service not found.")
		os.Exit(0)
	}

	// get credentials
	cred := models.GetEntry(service, db)

	// encode to json
	data, err := json.Marshal(cred)
	if err != nil {
		fmt.Println("Error exporting credentials.")
		os.Exit(1)
	}

	// encode to base64
	enc_data := base64.StdEncoding.EncodeToString(data)

	newFile := service + ".data"

	// write to file
	os.WriteFile(newFile, []byte(enc_data), os.ModePerm)

	fmt.Printf("Exported credentials for %s to %s.\n", service, newFile)
}

// import a service credential from file
func Import(file string, db *gorm.DB) {
	// check if logged in
	if !models.CheckSessionValid(db) {
		fmt.Println("Please login before importing credentials.")
		os.Exit(0)
	}

	// check file exists
	if _, err := os.Stat(file); os.IsNotExist(err) {
		fmt.Println("File not found.")
		os.Exit(0)
	}

	// read file
	data, err := os.ReadFile(file)
	if err != nil {
		fmt.Println("Error reading file.")
		os.Exit(1)
	}

	// decode from base64
	dec_data, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		fmt.Println("Error decoding file.")
		os.Exit(1)
	}

	// decode from json
	var cred models.Credential
	err = json.Unmarshal(dec_data, &cred)
	if err != nil {
		fmt.Println("Error decoding file.")
		os.Exit(1)
	}

	// save credentials
	models.ImportCredentials(cred.Service, cred.Username, cred.Password, cred.Key, db)

	fmt.Printf("Imported credentials for %s.\n", cred.Service)
}