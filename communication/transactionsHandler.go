package communication

import (
	"dope-node/blockchain"
	"dope-node/communication/messages"
	"dope-node/utils"
	"encoding/json"
	"fmt"
	"net/http"

	log "github.com/sirupsen/logrus"
)

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

		result, err := utils.VerifySignature(input.Payload.Sender, fmt.Sprintf("%v", input.Payload), input.Signature)
		if err != nil || !result {
			log.Infof("Invalid signature. Reason: %s", err)
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
	err := blockchain.DopeTransactions.InsertTransaction(&transToSend, &dbUrl)
	if err != nil {
		log.Warnf("cannot make transaction. Reason: %s", err)
		return
	}

	transMess := messages.TransactionRequest{Type: "transaction", Amount: amount, Receiver: receiver, Sender: sender}
	serializedMess, err := json.Marshal(transMess)
	if err != nil {
		log.Warnf("Cannot serialize transaction. Reason: %s", err)
		return
	}

	log.Infof("transaction from %s to %s inserted successfully", fullNodeAddress, sender)
	sendWsMessageToAllNodes(serializedMess)
}
