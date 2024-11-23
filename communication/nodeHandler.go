package communication

import (
	"dope-node/blockchain"
	"dope-node/communication/messages"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

func nodeHandler(w http.ResponseWriter, r *http.Request) {
	connection, err := getWebsocketConnection(w, r)
	if err != nil {
		log.Warnf("Failed to establish websocket connection. Reason: %s", err)
		return
	}

	_, mess, err := connection.ReadMessage()
	if err != nil {
		log.Warnf("Failed to read message. Reason: %s", err)
	}

	var messType messages.MessageType
	err = json.Unmarshal(mess, &messType)
	if err != nil {
		log.Warnf("Failed to deserialize message. Reason: %s", err)
	}

	switch messType.Type {
	case "addresses":
		{
			log.Info("received addresses message")
			var receivedMessage messages.AvailableNodesAddresses
			err = json.Unmarshal(mess, &receivedMessage)
			if err != nil {
				log.Warnf("Failed to deserialize message. Reason: %s", err)
			}
			knownNodeAddresses = receivedMessage.Addresses
		}
	case "transaction":
		{
			log.Info("received transaction message")
			var receivedMessage messages.Transaction
			err = json.Unmarshal(mess, &receivedMessage)
			if err != nil {
				log.Warnf("Failed to deserialize message. Reason: %s", err)
			}
			_, err = blockchain.Transact(receivedMessage)
			if err == nil {
				log.Info("Unsuccessful transaction")
			} else {
				log.Info("Transaction successfull")
			}
		}
	}
}
