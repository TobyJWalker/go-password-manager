# go-password-manager

This project is a CLI application which can safely store your login credentials. A variety of commands will allow you to do everything you need to do with your passwords.

This project is for me to practice developing my own CLI tool (and hopefully use it plenty myself!). I have chosen Go as I want to develop more confidence with it, and its cross-platform capabilities mean it can be used more widely.

Passwords are hashed for security, and its usage is locked behind a master password. The user will have to login with the master password before using most commands, the session will expire after 5 minutes so that no-one can access your passwords if you leave your device unattended.

# Features

- [x] Create an password for the app
- [x] Verify your identity with the password
- [x] Session expiration after 5 minutes
- [x] Add a new login entry
- [x] List all login entries
- [x] Edit login entry
- [x] Delete login entry
- [x] Encrypted storage
- [ ] Import/Export credentials

# Installing from Source

To build the project from source, you will need to have Go installed on your device.

Build:
```
go build -o build/go-pwm
```

To install the program so that you can run it from anywhere, you can add the Go install directory to your PATH environment variable.

```
# get install directory
go list -f '{{.Target}}'

# add to PATH (Linux/Mac)
export PATH=$PATH:/path/to/install/directory

# add to PATH (Windows)
set PATH=%PATH%;C:\path\to\install\directory

# install
go install

# run
go-pwm [command]
```

# Usage

```
go-pwm [command]

# see help for a list of commands
go-pwm help
```

# Running in development mode

The app determines where to save the database based on an environment variable. Please set this evironment variable before running the app if you want to run it in development mode.

```
export GO_ENV="dev"
go-pwm [command]

# OR

env GO_ENV="dev" go-pwm [command]
```
