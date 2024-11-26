package blockchain

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type Blockchain []Block

var DopeChain Blockchain

func SyncBlockchain(chain *Blockchain) {
	DopeChain = *chain
	log.Info("Synchronized blockchain: ")
	DopeChain.Print()
}

func (bchain *Blockchain) InsertToBlockchain(content *string) *Block {
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
	return newBlock
}

func (bchain *Blockchain) Print() {
	for _, b := range *bchain {
		fmt.Println(b)
	}
}
