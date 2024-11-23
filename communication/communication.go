package communication

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"

	"dope-node/communication/messages"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var knownNodeAddresses = make([]string, 0)

const (
	NODE_ENDPOINT      = "/node"
	BOOTSTRAP_ENDPOINT = "/bootstrap"
)

func ConnectToNetwork(bootstrapAddr *string, ip *string, port *int) {
	go func() {
		nodeAddress := fmt.Sprintf("%s:%d", *ip, *port)
		http.HandleFunc(NODE_ENDPOINT, nodeHandler)
		err := http.ListenAndServe(nodeAddress, nil)
		if err != nil {
			log.Errorf("Failed to run server on %s. Reason: %s", nodeAddress, err)
		}
	}()

	for {
		err := fetchNodeAddressesFromBootstrap(*bootstrapAddr, *ip, *port)
		if err != nil {
			log.Warnf("Failed to fetch node addresses. Reason: %s", err)
		}

		log.Infof("fetched bootstrap addresses: %s", knownNodeAddresses)
		knownNodeAddresses = deleteAddress(fmt.Sprintf("%s:%d", *ip, *port))
		log.Infof("Cleaned bootstrap addresses: %s", knownNodeAddresses)

		// Update addresses every minute
		time.Sleep(1 * time.Minute)
	}

}

func fetchNodeAddressesFromBootstrap(bootstrapAddress string, ip string, port int) error {
	connectMessage := messages.NewConnectionMessage{Ip: ip, Port: port}
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

func sendWsMessage(targetAddress string, message []byte) error {
	u := url.URL{Scheme: "ws", Host: targetAddress, Path: BOOTSTRAP_ENDPOINT}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	defer c.Close()

	err = c.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return err
	}
	log.Infof("Message sent to %s", targetAddress)

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

func deleteAddress(address string) []string {
	for i, v := range knownNodeAddresses {
		if v == address {
			return append(knownNodeAddresses[:i], knownNodeAddresses[i+1:]...)
		}
	}
	return knownNodeAddresses
}
