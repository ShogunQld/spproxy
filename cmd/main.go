package main

import (
	"fmt"
	"log"
	"os"
	"spproxy/internal/configs"
	"spproxy/internal/server"
)

func main() {
	args := os.Args
	configFile := "config.json"
	if len(args) > 1 {
		configFile = args[1]
	}

	// Load configurations from config file
	config, err := configs.NewConfiguration(configFile)
	if err != nil {
		log.Fatalf("could not load configuration: %v", err)
	}

	fmt.Printf("\n*********************\n")
	fmt.Printf("* Sticky Port Proxy *\n")
	fmt.Printf("*********************\n\n")

	log.Fatal(server.Run(config))
}
