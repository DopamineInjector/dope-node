package blockchain

import (
	"crypto/sha256"
	"encoding/hex"
)

type TrieNode struct {
	Hash     string
	Children map[rune]*TrieNode
	Value    string
}

type PMTrie struct {
	Root *TrieNode
}

func (trie PMTrie) InsertValue(key string, value string) {
	currentNode := trie.Root

	for _, chr := range key {
		if currentNode.Children[chr] == nil {
			newNode := TrieNode{
				Children: map[rune]*TrieNode{},
			}
			currentNode.Children[chr] = &newNode
		}
		currentNode = currentNode.Children[chr]
	}

	currentNode.Value = value
	// tmp solution - later will recalculate only modified parts
	trie.Root.recalculateHashes()
}

func (trie *PMTrie) Get(key string) string {
	currentNode := trie.Root
	for _, val := range key {
		if currentNode.Children[val] == nil {
			return ""
		}
		currentNode = currentNode.Children[val]
	}

	return currentNode.Value
}

func (node *TrieNode) recalculateHashes() {
	for _, v := range node.Children {
		if v != nil {
			v.recalculateHashes()
		}
	}
	node.Hash = node.calculateHash()
}

func (node *TrieNode) ToString() string {
	return node.Hash + ":" + node.Value
}

func (node *TrieNode) calculateHash() string {
	contentToHash := ""
	for _, val := range node.Children {
		if val != nil {
			contentToHash += val.ToString()
		}
	}
	contentToHash += node.Value

	hash := sha256.New()
	hash.Write([]byte(contentToHash))
	hashValue := hash.Sum(nil)

	return hex.EncodeToString(hashValue)
}
