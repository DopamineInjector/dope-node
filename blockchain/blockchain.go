package blockchain

import (
	"fmt"
	"log"
)

type Blockchain []Block

func (blockchain *Blockchain) AddBlock(content string) {
	var newBlock Block

	if len(*blockchain) == 0 {
		newBlock = createGenesisBlock(content)
		log.Println("Created genesis block: " + newBlock.ToString())
	} else {
		prevBlock := (*blockchain)[len(*blockchain)-1]
		newBlock = createBlock(&prevBlock, content)
		log.Println("Created new block: " + newBlock.ToString())
	}

	*blockchain = append(*blockchain, newBlock)
}

func (blockchain *Blockchain) PrintBlockchain() {
	for i := 0; i < len(*blockchain); i++ {
		fmt.Println((*blockchain)[i])
	}
}
