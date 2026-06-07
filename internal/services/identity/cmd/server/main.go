package main

import (
	"log"
	"os"

	"github.com/lxbachit03/ygz-microservices-golang/internal/services/identity/internal/app"
)

func main() {

	app, err := app.NewApplication()
	if err != nil {
		// log

		log.Printf("Failed to create application: %v", err)
		os.Exit(1)
	}

	if err := app.Run(); err != nil {
		// log

		log.Printf("Failed to run application: %v", err)
		os.Exit(1)
	}
}
