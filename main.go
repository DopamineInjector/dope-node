package main

import (
	"dope-node/communication"
	"dope-node/config"

	db "github.com/DopamineInjector/go-dope-db"
	log "github.com/sirupsen/logrus"
)

func main() {
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

	nodeAddress := config.GetString(config.ServerIpAddressKey)
	port := config.GetInt(config.ServerPortKey)
	bootstrapServerAddress := config.GetString(config.BootstrapServerAddressKey)

	go communication.StartConsoleListener()
	err = communication.ConnectToNetwork(&bootstrapServerAddress, &nodeAddress, &port, dbUrl)
	if err != nil {
		log.Errorf("Failed to connect to node network. Reason: %s", err)
	}
}
