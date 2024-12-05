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
	Type     string `json:"type"`
	Sender   string `json:"sender"`
	Receiver string `json:"receiver"`
	Amount   int    `json:"amount"`
}

type StructureRequest struct {
	Type      string `json:"type"`
	Requester string `json:"requester"`
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

var SmartContractRequest struct {
	Payload struct {
		Sender     []byte `json:"sender"`
		Contract   []byte `json:"contract"`
		Entrypoint string `json:"entrypoint"`
		Args       string `json:"args"`
	} `json:"payload"`
	Signature []byte `json:"signature"`
	View      bool   `json:"view"`
}

func (req *TransactionRequest) ParseToTransaction() blockchain.Transaction {
	return blockchain.Transaction{
		Sender:   req.Sender,
		Receiver: req.Receiver,
		Amount:   req.Amount,
	}
}

func (req *SmartContractRequest) ParseToSmartContract() blockchain.SmartContract {
	return blockchain.SmartContract{
		Sender:     req.Payload.Sender,
		Contract:   req.Payload.Contract,
		Entrypoint: req.Payload.Entrypoint,
		Args:       req.Payload.Args,
	}
}
