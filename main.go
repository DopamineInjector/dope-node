package main

import (
	"dope-node/config"

	log "github.com/sirupsen/logrus"
	db "github.com/DopamineInjector/go-dope-db"
)

func main() {
	log.Info("Starting node")
	err := config.InitializeConfig();
	if err != nil {
		log.Warn("Could not find read config.toml, resolving to default config values");
	}
	log.Info("Parsed node configuration")
	log.Info("Connecting to node storage")
	dbUrl := config.GetString(config.DbUrlKey);
	checksum, err := db.GetChecksum(dbUrl);
	if err != nil {
		log.Fatalf("Could not connect to db instance, exiting\n%s", err.Error());
	}
	log.Infof("Connected to storage, current state checksum: %s", checksum.Checksum);
}
