package blockchain

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

type Blockchain []Block
type Transactables []Transactable

var TransactionsNumber int32 = 0
var DopeChain Blockchain
var DopeTransactables Transactables

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
	executeTransactables()

	*bchain = append(*bchain, *newBlock)
	log.Infof("blockchain: ")
	bchain.Print()
	return newBlock
}

func (dTransactable *Transactables) InsertTransactable(t Transactable) {
	log.Infof("Transactables: ")
	DopeTransactables.Print()
	*dTransactable = append(*dTransactable, t)
}

func (bchain *Blockchain) Print() {
	for _, b := range *bchain {
		fmt.Println(b)
	}
}

func (trans *Transactables) Print() {
	for _, t := range *trans {
		t.print()
	}
}

func executeTransactables() {
	for _, t := range DopeTransactables {
		out, err := t.run()
		if err != nil {
			log.Warnf("error while running transactable. Reason: %s", err)
			t.print()
		} else {
			log.Infof("transactable result: %s", *out)
		}
	}

	DopeTransactables = DopeTransactables[:0]
}
