// Almighty Father, You are the source of all wisdom, talents, and skills.
// Thank You for these beautiful gifts, opportunities to serve you, and the work that You have entrusted to me.
// Please open my heart and enlighten my mind, so that I may be fully in tune with Your divine purpose in calling me to the engineering profession.
// Lord God, You are the greatest engineer. Please infuse me even with just the tiniest spark of Your Divine Wisdom so that as I  do my work, it is really your work that is done.
// Loving God, make my heart Your Heart, make my mind Your Mind, and make my hands Your Hands.
// Make me your instrument, so that I may be always conscious and mindful of the fact that on my work depend the lives and properties of my fellow human. Bless me also with the gift of love and sensitivity to respect the people who build and use the products of my work.
// This I ask in Jesus’ name.
// AMEN.
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
