package communication

import (
	"dope-node/blockchain"
	"dope-node/communication/messages"
	"dope-node/utils"
	"encoding/base64"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const MAX_TRANSACTIONS_PER_BLOCK = 5

func handleTransfer(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var input struct {
			Payload struct {
				Sender    string `json:"sender"`
				Recipient string `json:"recipient"`
				Amount    int    `json:"amount"`
			} `json:"payload"`
			Signature string `json:"signature"`
		}
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		marshalledPayload, err := json.Marshal(input.Payload)
		if err != nil {
			log.Warnf("error marshalling payload")
		}
		sndr, err := base64.StdEncoding.DecodeString(input.Payload.Sender)
		if err != nil {
			http.Error(w, "Invalid sender encoding", http.StatusBadRequest)
			return
		}
		sig, err := base64.StdEncoding.DecodeString(input.Signature)
		if err != nil {
			http.Error(w, "Invalid signature encoding", http.StatusBadRequest)
			return
		}

		result, err := utils.VerifySignature(sndr, marshalledPayload, sig)
		if err != nil || !result {
			log.Infof("Invalid signature. Reason: %s", err)
			w.WriteHeader(http.StatusForbidden)
			return
		}

		beginTransaction(input.Payload.Sender, input.Payload.Amount, input.Payload.Recipient)
		w.WriteHeader(http.StatusCreated)
	} else {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func beginTransaction(sender string, amount int, receiver string) {
	transToSend := blockchain.Transaction{Amount: amount, Receiver: receiver, Sender: sender}
	if sender == receiver {
		log.Info("cannot send $ to yourself")
		return
	}
	blockchain.DopeTransactables.InsertTransactable(transToSend)

	transMess := messages.TransactionRequest{Type: "transaction", Amount: amount, Receiver: receiver, Sender: sender}
	serializedMess, err := json.Marshal(transMess)
	if err != nil {
		log.Warnf("Cannot serialize transaction. Reason: %s", err)
		return
	}

	log.Infof("transaction from %s to %s inserted successfully", fullNodeAddress, sender)
	sendWsMessageToAllNodes(serializedMess)

	if len(blockchain.DopeTransactables) >= MAX_TRANSACTIONS_PER_BLOCK {
		digBlock("bloczek")
	}
}
