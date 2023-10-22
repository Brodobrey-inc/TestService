package main

import (
	"fmt"

	"github.com/Brodobrey-inc/TestService/config"
	"github.com/Brodobrey-inc/TestService/database"
	"github.com/Brodobrey-inc/TestService/logging"
	"github.com/Brodobrey-inc/TestService/webserver"
)

func main() {
	fmt.Println("1. Initialize config from .env file")
	config.Initialize()

	fmt.Println("2. Initialize logger")
	logging.Initialize()

	fmt.Println("3. Initialize webserver")
	router := webserver.Initialize()

	fmt.Println("4. Connect to database")
	database.Initialize()

	fmt.Println("5. Start webserver")
	webserver.StartServer(router)
}
