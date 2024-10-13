package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
	"strconv"
	"time"
)

type Block struct {
	Index        int
	Timestamp    string
	Content      string
	PreviousHash string
	Hash         string
}

func createBlock(previousBlock Block, content string) Block {
	newBlock := Block{
		Index:        previousBlock.Index + 1,
		Timestamp:    strconv.FormatInt(time.Now().Unix(), 10),
		Content:      content,
		PreviousHash: previousBlock.Hash,
		Hash:         "",
	}
	newBlock.Hash = calculateHash(newBlock)
	return newBlock
}

func createGenesisBlock(content string) Block {
	genesis := Block{
		Index:        0,
		Timestamp:    strconv.FormatInt(time.Now().Unix(), 10),
		Content:      content,
		PreviousHash: "0",
		Hash:         "",
	}
	genesis.Hash = calculateHash(genesis)
	return genesis
}

func calculateHash(block Block) string {
	dataToHash := strconv.Itoa(block.Index) + block.Timestamp + block.Content + block.PreviousHash
	hash := sha256.New()
	hash.Write([]byte(dataToHash))
	hashValue := hash.Sum(nil)

	return hex.EncodeToString(hashValue)
}
