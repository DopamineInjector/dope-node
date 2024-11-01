package main

import (
	"dope-node/communication"
	"dope-node/config"
	"flag"
	"fmt"

	log "github.com/sirupsen/logrus"
)

func main() {
	nodeAddress := "127.0.0.1" // Temporary, while running locally
	bootstrapServerAddress := flag.String("bootstrap_address", "127.0.0.1:7312", "An IP address to the bootstrap server")
	port := flag.Int("port", 7313, "Port to run websocket listener on")

	flag.Parse()
	fmt.Println("Current configuration:")
	fmt.Println("\t- Bootstrap server address: ", *bootstrapServerAddress)
	fmt.Println("\t- Server running on port: ", *port)

	log.Info("Starting node")
	err := config.InitializeConfig()
	if err != nil {
		log.Warn("Could not find read config.toml, resolving to default config values")
	}
	log.Info("Parsed node configuration")

	communication.ConnectToNetwork(bootstrapServerAddress, &nodeAddress, port)
}
