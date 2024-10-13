package main

import (
	"dope-node/communication"
	"flag"
	"fmt"
)

func main() {
	bootstrapServer := flag.String("bootstrap", "127.0.0.1", "An IP address to the bootstrap server")
	fmt.Println(bootstrapServer)

	communication.RunWebsocketListener()
}
