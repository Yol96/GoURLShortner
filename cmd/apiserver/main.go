package main

import (
	"log"

	"github.com/Yol96/GoURLShortner/internal/app/apiserver"
)

func main() {
	// Create a new API-server config
	config, err := apiserver.NewConfig()
	if err != nil {
		log.Fatal("Can`t load apiserver config: ", err)
	}

	// Start the server
	if err := apiserver.Start(config); err != nil {
		log.Fatal("Can`t start server: ", err)
	}
}
