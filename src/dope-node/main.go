package main

import (
	"dope-node/communication"
	"flag"
	"fmt"
)

func main() {
	isBoostrapServer := flag.Bool("bootstrap", false, "Is this node running as a bootstrap")
	bootstrapServerAddress := flag.String("bootstrap_address", "127.0.0.1:7312", "An IP address to the bootstrap server")

	flag.Parse()
	fmt.Println(*isBoostrapServer)
	fmt.Println(*bootstrapServerAddress)

	communication.RegisterEndpoints(*bootstrapServerAddress, *isBoostrapServer)
}
