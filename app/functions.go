package app

import (
	"fmt"
)

func Help() {
	fmt.Println(`Usage: go-pwm [command] [args]

Commands:

  help           : display help message
  init           : setup credentials for the first time
  add [service]  : add credentials for new service
  list           : list all services
  get [service]  : get credentials for the service
  edit [service] : edit credentials for a service
  rm [service]   : remove a service
    `)
}