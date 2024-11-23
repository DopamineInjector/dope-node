package messages

import (
	"dope-node/blockchain"
)

type NewConnectionMessage struct {
	Ip   string `json:"ip"`
	Port int    `json:"port"`
}

type MessageType struct {
	Type string `json:"type"`
}

type AvailableNodesAddresses struct {
	Type      string   `json:"type"`
	Addresses []string `json:"addresses"`
}

type TransactionRequest struct {
	Type     string  `json:"type"`
	Sender   string  `json:"sender"`
	Receiver string  `json:"receiver"`
	Amount   float64 `json:"amount"`
}

type StructureRequest struct {
	Type string `json:"type"`
}

type StructureResponse struct {
	Type         string                  `json:"type"`
	Blockchain   blockchain.Blockchain   `json:"blockchain"`
	Transactions blockchain.Transactions `json:"transactions"`
}

type BlockMessage struct {
	Type  string           `json:"type"`
	Block blockchain.Block `json:"block"`
}

func (req *TransactionRequest) ParseToTransaction() blockchain.Transaction {
	return blockchain.Transaction{
		Sender:   req.Sender,
		Receiver: req.Receiver,
		Amount:   req.Amount,
	}
}
