package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strconv"
	"time"
)

type Block struct {
	Index        int
	Timestamp    int64
	Content      string
	PreviousHash string
	Hash         string
	Transactions []Transactable
	MPTHash      string
	dbChecksum   string
}

func (block *Block) ToString() string {
	dateOfCreation := time.Unix(block.Timestamp, 0)
	return fmt.Sprintf("Block: {index: %d, created: %s, hash: %s}", block.Index, dateOfCreation, block.Hash)
}

func (block *Block) calculateHash() string {
	dataToHash := strconv.Itoa(block.Index) + strconv.FormatInt(block.Timestamp, 10) + block.Content + block.PreviousHash + block.MPTHash + block.dbChecksum
	hash := sha256.New()
	hash.Write([]byte(dataToHash))
	hashValue := hash.Sum(nil)

	return hex.EncodeToString(hashValue)
}

func createBlock(previousBlock *Block, content *string) *Block {
	newBlock := Block{
		Index:        previousBlock.Index + 1,
		Timestamp:    time.Now().Unix(),
		Content:      *content,
		PreviousHash: previousBlock.Hash,
		Transactions: DopeTransactables,
		MPTHash:      "hash",
		dbChecksum:   "checksum",
	}
	newBlock.Hash = newBlock.calculateHash()

	return &newBlock
}

func createGenesisBlock(content *string) *Block {
	genesis := Block{
		Index:        0,
		Timestamp:    time.Now().Unix(),
		Content:      *content,
		PreviousHash: "0",
		Transactions: DopeTransactables,
		MPTHash:      "hash",
		dbChecksum:   "checksum",
	}
	genesis.Hash = genesis.calculateHash()

	return &genesis
}
