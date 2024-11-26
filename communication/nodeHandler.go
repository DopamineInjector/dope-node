package communication

import (
	b "dope-node/blockchain"
	"dope-node/communication/messages"
	"encoding/json"
	"net/http"

	log "github.com/sirupsen/logrus"
)

const (
	ADDRESSES_MESSAGE_TYPE              = "addresses"
	TRANSACTION_MESSAGE_TYPE            = "transaction"
	SYNC_STRUCTURE_MESSAGE_TYPE         = "structure"
	BLOCK_MESSAGE_TYPE                  = "block"
	STRUCTURE_SYNC_REQUEST_MESSAGE_TYPE = "sync"
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
			initializeNodeAddresses(receivedMessage.Addresses)
			AddressesFetched <- true
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
			err = b.DopeTransactions.InsertTransaction(&parsedTrans, &dbUrl)
			if err != nil {
				log.Infof("Unsuccessful transaction. Reeason: %s", err)
			} else {
				log.Info("Transaction successfull")
			}
		}
	case SYNC_STRUCTURE_MESSAGE_TYPE:
		{
			var receivedMessage messages.StructureResponse
			err := json.Unmarshal(mess, &receivedMessage)
			if err != nil {
				log.Warnf("Failed to deserialize message. Reason: %s", err)
				break
			}

			b.SyncBlockchain(&receivedMessage.Blockchain)
			b.SyncTransactions(&receivedMessage.Transactions)
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
			b.DopeTransactions = b.DopeTransactions[:0]
		}
	case STRUCTURE_SYNC_REQUEST_MESSAGE_TYPE:
		{
			var receivedMessage messages.StructureRequest
			err := json.Unmarshal(mess, &receivedMessage)
			if err != nil {
				log.Warnf("Failed to deserialize message. Reason: %s", err)
				break
			}

			mess := messages.StructureResponse{Type: SYNC_STRUCTURE_MESSAGE_TYPE, Blockchain: b.DopeChain, Transactions: b.DopeTransactions}
			serializedMess, err := json.Marshal(mess)
			if err != nil {
				log.Warnf("Cannot serialize. Reason: %s", err)
				return
			}

			log.Info("Sending: ")
			b.DopeChain.Print()
			sendWsMessage(&receivedMessage.Requester, serializedMess, "/node")
		}
	}

}
