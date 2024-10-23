package main

import (
	"dope-node/config"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Starting node")
	err := config.InitializeConfig();
	if err != nil {
		log.Warn("Could not find read config.toml, resolving to default config values");
	}
	log.Info("Parsed node configuration")
}
