package communication

import (
	b "dope-node/blockchain"
	"dope-node/communication/messages"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const (
	ADDRESSES_MESSAGE_TYPE      = "addresses"
	TRANSACTION_MESSAGE_TYPE    = "transaction"
	STRUCTURE_INIT_MESSAGE_TYPE = "structure"
	BLOCK_MESSAGE_TYPE          = "block"
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
		return
	}

	var messType messages.MessageType
	err = json.Unmarshal(mess, &messType)
	if err != nil {
		log.Warnf("Failed to deserialize message. Reason: %s", err)
		return
	}

	log.Infof("Received %s message type", messType.Type)
	handleMessageType(messType.Type, mess)
}

func handleMessageType(messType string, mess []byte) {
	switch messType {
	case ADDRESSES_MESSAGE_TYPE:
		{
			var receivedMessage messages.AvailableNodesAddresses
			err := json.Unmarshal(mess, &receivedMessage)
			if err != nil {
				log.Warnf("Failed to deserialize message. Reason: %s", err)
				break
			}
			knownNodeAddresses = receivedMessage.Addresses
		}
	case TRANSACTION_MESSAGE_TYPE:
		{
			var receivedMessage messages.TransactionRequest
			err := json.Unmarshal(mess, &receivedMessage)
			if err != nil {
				log.Warnf("Failed to deserialize message. Reason: %s", err)
				break
			}

			parsedTrans := receivedMessage.ParseToTransaction()
			err = b.Transact(&parsedTrans, &dbUrl)
			if err != nil {
				log.Infof("Unsuccessful transaction. Reeason: %s", err)
			} else {
				log.Info("Transaction successfull")
			}
		}
	case STRUCTURE_INIT_MESSAGE_TYPE:
		{
			var receivedMessage messages.StructureResponse
			err := json.Unmarshal(mess, &receivedMessage)
			if err != nil {
				log.Warnf("Failed to deserialize message. Reason: %s", err)
				break
			}

			b.InitializeBlockchain(&receivedMessage.Blockchain)
			b.InitalizeTransactions(&receivedMessage.Transactions)
		}
	case BLOCK_MESSAGE_TYPE:
		{
			var receivedMessage messages.BlockMessage
			err := json.Unmarshal(mess, &receivedMessage)
			if err != nil {
				log.Warnf("Failed to deserialize message. Reason: %s", err)
				break
			}

			b.DopeChain = append(b.DopeChain, receivedMessage.Block)
		}
	}

}
