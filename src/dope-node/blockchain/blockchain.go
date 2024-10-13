package blockchain

import (
	"fmt"
	"log"
)

var blockchain = []Block{}

func AddBlock(content string) {
	var newBlock Block

	if len(blockchain) == 0 {
		newBlock = createGenesisBlock(content)
		log.Println("Created genesis block: " + newBlock.ToString())
	} else {
		newBlock = createBlock(blockchain[len(blockchain)-1], content)
		log.Println("Created new block: " + newBlock.ToString())
	}

	blockchain = append(blockchain, newBlock)
}

func PrintBlockchain() {
	for i := 0; i < len(blockchain); i++ {
		fmt.Println(blockchain[i])
	}
}
