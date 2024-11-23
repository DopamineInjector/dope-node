package communication

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"dope-node/blockchain"
	"dope-node/communication/messages"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var knownNodeAddresses = make([]string, 0)
var dbUrl string

const (
	NODE_ENDPOINT          = "/node"
	BOOTSTRAP_ENDPOINT     = "/bootstrap"
	STRUCTURE_INIT_MESSAGE = "structure"
)

func ConnectToNetwork(bootstrapAddr *string, ip *string, port *int, url string) error {
	dbUrl = url
	go func() {
		nodeAddress := fmt.Sprintf("%s:%d", *ip, *port)
		http.HandleFunc(NODE_ENDPOINT, nodeHandler)
		err := http.ListenAndServe(nodeAddress, nil)
		if err != nil {
			log.Errorf("failed to run server on %s. Reason: %s", nodeAddress, err)
		}
	}()

	err := fetchNodeAddresses(bootstrapAddr, ip, port)
	if err != nil {
		return fmt.Errorf("failed to fetch node addresses. Reason: %s", err)
	}

	err = initializeBlockchain()
	if err != nil {
		return fmt.Errorf("failed to fetch blockchain structure. Reason: %s", err)
	}

	return nil
}

func initializeBlockchain() error {
	if len(knownNodeAddresses) == 0 {
		log.Info("No other nodes. Creating blockchain")
		blockchain.InitializeBlockchain(&blockchain.Blockchain{})
		return nil
	}

	initMess := messages.StructureRequest{Type: STRUCTURE_INIT_MESSAGE}
	serializedMess, err := json.Marshal((initMess))
	if err != nil {
		return err
	}

	// assuming that all the nodes have the same blockchain - so sending request to only one
	err = sendWsMessage(&knownNodeAddresses[0], serializedMess)
	if err != nil {
		return err
	}

	return nil
}

func fetchNodeAddresses(bootstrapAddr *string, ip *string, port *int) error {
	err := fetchNodeAddressesFromBootstrap(bootstrapAddr, ip, port)
	if err != nil {
		return err
	}

	ownAddress := fmt.Sprintf("%s:%d", *ip, *port)
	knownNodeAddresses = deleteAddress(&ownAddress)
	log.Infof("Bootstrap addresses: %s", knownNodeAddresses)

	return nil
}

func fetchNodeAddressesFromBootstrap(bootstrapAddress *string, ip *string, port *int) error {
	connectMessage := messages.NewConnectionMessage{Ip: *ip, Port: *port}
	serializedMess, err := json.Marshal((connectMessage))
	if err != nil {
		return err
	}

	err = sendWsMessage(bootstrapAddress, serializedMess)
	if err != nil {
		return err
	}

	return nil
}

func sendWsMessage(targetAddress *string, message []byte) error {
	u := url.URL{Scheme: "ws", Host: *targetAddress, Path: BOOTSTRAP_ENDPOINT}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	defer c.Close()

	err = c.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return err
	}
	log.Infof("Message sent to %s", *targetAddress)

	return nil
}

func getWebsocketConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Warn("Error upgrading to WebSocket:", err)
	}

	return conn, err
}

func deleteAddress(address *string) []string {
	for i, v := range knownNodeAddresses {
		if v == *address {
			return append(knownNodeAddresses[:i], knownNodeAddresses[i+1:]...)
		}
	}

	return knownNodeAddresses
}
