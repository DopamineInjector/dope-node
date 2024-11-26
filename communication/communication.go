package communication

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"dope-node/blockchain"
	"dope-node/communication/messages"

	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var knownNodeAddresses = make([]string, 0)
var dbUrl string
var fullNodeAddress string

const (
	NODE_ENDPOINT          = "/node"
	BOOTSTRAP_ENDPOINT     = "/bootstrap"
	STRUCTURE_INIT_MESSAGE = "structure"
)

func ConnectToNetwork(bootstrapAddr *string, ip *string, port *int, url string) error {
	dbUrl = url
	serverReady := make(chan bool)
	fullNodeAddress = fmt.Sprintf("%s:%d", *ip, *port)

	go func() {
		nodeAddress := fmt.Sprintf("%s:%d", *ip, *port)
		http.HandleFunc(NODE_ENDPOINT, nodeHandler)
		log.Infof("Server running on %s", nodeAddress)

		serverReady <- true
		err := http.ListenAndServe(nodeAddress, nil)
		if err != nil {
			log.Errorf("failed to run server on %s. Reason: %s", nodeAddress, err)
		}
	}()

	<-serverReady
	err := fetchNodeAddresses(bootstrapAddr, ip, port)
	if err != nil {
		return fmt.Errorf("failed to fetch node addresses. Reason: %s", err)
	}

	err = syncBlockchain()
	if err != nil {
		return fmt.Errorf("failed to fetch blockchain structure. Reason: %s", err)
	}

	select {}
}

// For running functions from currently running node
func StartConsoleListener() {
	scanner := bufio.NewScanner(os.Stdin)
	log.Info("Console listener has started")

	for scanner.Scan() {
		input := scanner.Text()
		switch input {
		case "exit":
			os.Exit(0)
		case "transaction":
			fmt.Println("Enter amount: ")
			amount := scanner.Text()
			fmt.Println("Enter receiver: ")
			receiver := scanner.Text()
			beginTransaction(amount, receiver)
		case "block":
			fmt.Println("Block content: ")
			content := scanner.Text()
			digBlock(content)
		default:
			log.Infof("Unknown command: %s\n", input)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Warnf("Error reading console input: %v", err)
	}
}

func digBlock(content string) {
	blockchain.DopeChain.InsertToBlockchain(&content)
	blockchain.DopeTransactions = blockchain.DopeTransactions[:0]
}

func beginTransaction(amount string, receiver string) {
	// TODO: change sender value
	sender := fullNodeAddress
	parsedAmount, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		log.Warnf("cannot parse amount value. Reason: %s", err)
		return
	}

	transToSend := blockchain.Transaction{Amount: parsedAmount, Receiver: receiver, Sender: sender}
	err = blockchain.DopeTransactions.InsertTransaction(&transToSend, &dbUrl)
	if err != nil {
		log.Warnf("cannot make transaction. Reason: %s", err)
		return
	}

	log.Infof("transaction from %s to %s inserted successfully", fullNodeAddress, sender)
}

func syncBlockchain() error {
	if len(knownNodeAddresses) == 0 {
		log.Info("No other nodes. Creating blockchain")
		blockchain.SyncBlockchain(&blockchain.Blockchain{})
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

	return nil
}

func initializeNodeAddresses(addresses []string) {
	knownNodeAddresses = addresses
	knownNodeAddresses = deleteAddress(&fullNodeAddress)
	log.Infof("Bootstrap addresses: %s", knownNodeAddresses)
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
