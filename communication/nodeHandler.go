package communication

import (
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

	var receivedMessage AvailableNodesAddresses
	err = json.Unmarshal(mess, &receivedMessage)
	if err != nil {
		log.Warnf("Failed to deserialize message. Reason: %s", err)
	}

	knownNodeAddresses = receivedMessage.Addresses
}
