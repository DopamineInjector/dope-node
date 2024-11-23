package main

import (
	"dope-node/communication"
	"dope-node/config"
	"flag"
	"fmt"

	db "github.com/DopamineInjector/go-dope-db"
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

	log.Info("Connecting to node storage")
	dbUrl := config.GetString(config.DbUrlKey)
	checksum, err := db.GetChecksum(dbUrl)
	if err != nil {
		log.Fatalf("Could not connect to db instance, exiting\n%s", err.Error())
	}
	log.Infof("Connected to storage, current state checksum: %s", checksum.Checksum)

	err = communication.ConnectToNetwork(bootstrapServerAddress, &nodeAddress, port, dbUrl)
	if err != nil {
		log.Errorf("Failed to connect to node network. Reason: %s", err)
	}
}
