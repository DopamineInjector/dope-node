package blockchain

import (
	"fmt"
	"log"
)

type Blockchain []Block

var DopeChain Blockchain

func InitializeBlockchain(chain *Blockchain) {
	DopeChain = *chain
}

func (bchain *Blockchain) CreateBlock(content *string) {
	var newBlock *Block

	if len(*bchain) == 0 {
		newBlock = createGenesisBlock(content)
		log.Println("Created genesis block: " + newBlock.ToString())
	} else {
		prevBlock := (*bchain)[len(*bchain)-1]
		newBlock = createBlock(&prevBlock, content)
		log.Println("Created new block: " + newBlock.ToString())
	}

	*bchain = append(*bchain, *newBlock)
}

func (bchain *Blockchain) ToString() {
	for i := 0; i < len(*bchain); i++ {
		fmt.Println((*bchain)[i])
	}
}
