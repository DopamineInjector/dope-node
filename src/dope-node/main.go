package main

import (
	"dope-node/blockchain"
)

func main() {
	blockchain.AddBlock("content")
	blockchain.AddBlock("content2")
	blockchain.AddBlock("content3")
	blockchain.AddBlock("content4")

	blockchain.PrintBlockchain()
}
