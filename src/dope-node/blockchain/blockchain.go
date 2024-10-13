package blockchain

import "fmt"

var blockchain = []Block{}

func AddBlock(content string) {
	var newBlock Block

	if len(blockchain) == 0 {
		newBlock = createGenesisBlock(content)
	} else {
		newBlock = createBlock(blockchain[len(blockchain)-1], content)
	}

	blockchain = append(blockchain, newBlock)
}

func PrintBlockchain() {
	for i := 0; i < len(blockchain); i++ {
		fmt.Println(blockchain[i])
	}
}
