package main

import (
	"dope-node/communication"
	"flag"
	"fmt"
)

func main() {
	isBoostrapServer := flag.Bool("bootstrap", false, "Is this node running as a bootstrap")
	bootstrapServerAddress := flag.String("bootstrap_address", "127.0.0.1:7312", "An IP address to the bootstrap server")
	port := flag.Int("port", 7312, "Port to run websocket listener on")

	flag.Parse()
	fmt.Println("Current configuration:")
	fmt.Println("\t- Is running as bootstrap: ", *isBoostrapServer)
	fmt.Println("\t- Bootstrap server address: ", *bootstrapServerAddress)
	fmt.Println("\t- Server running on port: ", *port)

	communication.RegisterEndpoints(*bootstrapServerAddress, *isBoostrapServer, *port)
}
