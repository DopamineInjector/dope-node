package main

import (
	"dope-node/communication"
	"flag"
	"fmt"
)

func main() {
	nodeAddress := "127.0.0.1" // Temporary, while running locally
	bootstrapServerAddress := flag.String("bootstrap_address", "127.0.0.1:7312", "An IP address to the bootstrap server")
	port := flag.Int("port", 7313, "Port to run websocket listener on")

	flag.Parse()
	fmt.Println("Current configuration:")
	fmt.Println("\t- Bootstrap server address: ", *bootstrapServerAddress)
	fmt.Println("\t- Server running on port: ", *port)

	// log.Info("Starting node")
	// err := config.InitializeConfig()
	// if err != nil {
	// 	log.Warn("Could not find read config.toml, resolving to default config values")
	// }
	// log.Info("Parsed node configuration")

	// var bc blockchain.Blockchain
	// bc.AddBlock("content1")
	// bc.AddBlock("content2")
	// bc.AddBlock("content3")
	// bc.PrintBlockchain()

	communication.ConnectToNetwork(bootstrapServerAddress, &nodeAddress, port)
}
