package communication

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"dope-node/blockchain"
	"dope-node/communication/messages"

	db "github.com/DopamineInjector/go-dope-db"
	"github.com/gorilla/websocket"
	log "github.com/sirupsen/logrus"
)

var knownNodeAddresses = make([]string, 0)
var dbUrl string
var fullNodeAddress string
var AddressesFetched = make(chan bool)

const (
	BASE_URL                       = "/api"
	ACCOUNTS_ENDPOINT              = BASE_URL + "/account"
	ACCOUNTS_INFO_ENDPOINT         = ACCOUNTS_ENDPOINT + "/info"
	TRANSFER_ENDPOINT              = BASE_URL + "/transfer"
	SMARTCONTRACT_ENDPOINT         = BASE_URL + "/smartContract"
	NODE_ENDPOINT                  = "/node"
	BOOTSTRAP_ENDPOINT             = "/bootstrap"
	STRUCTURE_SYNC_REQUEST_MESSAGE = "sync"
)

func ConnectToNetwork(bootstrapAddr *string, ip *string, port *int, url string) error {
	dbUrl = url
	serverReady := make(chan bool)
	fullNodeAddress = fmt.Sprintf("%s:%d", *ip, *port)
	insertNamespaces()

	go func() {
		nodeAddress := fmt.Sprintf("%s:%d", *ip, *port)
		http.HandleFunc(NODE_ENDPOINT, nodeHandler)
		http.HandleFunc(ACCOUNTS_ENDPOINT, handleAccounts)
		http.HandleFunc(ACCOUNTS_INFO_ENDPOINT, handleAccountsInfo)
		http.HandleFunc(TRANSFER_ENDPOINT, handleTransfer)
		http.HandleFunc(SMARTCONTRACT_ENDPOINT, handleSmartContract)

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

	<-AddressesFetched
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
		case "block":
			fmt.Println("Block content: ")
			scanner.Scan()
			content := scanner.Text()
			digBlock(content)
		case "status":
			fmt.Println("Blockchain: ")
			blockchain.DopeChain.Print()
			fmt.Println("Transactions: ")
			blockchain.DopeTransactions.Print()
		default:
			log.Infof("Unknown command: %s\n", input)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Warnf("Error reading console input: %v", err)
	}
}

func digBlock(content string) {
	b := blockchain.DopeChain.InsertToBlockchain(&content)

	connectMessage := messages.BlockMessage{Type: "block", Block: *b}
	serializedMess, err := json.Marshal((connectMessage))
	if err != nil {
		log.Warnf("Cannot serialize. Reason: %s", err)
		return
	}

	log.Infof("Block with %s content initialized", content)
	resolveTransactions()
	sendWsMessageToAllNodes(serializedMess)
}

func resolveTransactions() {
	for _, tr := range blockchain.DopeTransactions {
		err := blockchain.DopeTransactions.InsertTransaction(&tr, &dbUrl)
		if err != nil {
			log.Infof("cannot resolve transaction: %s", err)
		}
	}

	blockchain.DopeTransactions = blockchain.DopeTransactions[:0]
}

func syncBlockchain() error {
	if len(knownNodeAddresses) == 0 {
		log.Info("No other nodes. Creating blockchain")
		blockchain.SyncBlockchain(&blockchain.Blockchain{})
		return nil
	}

	initMess := messages.StructureRequest{Type: STRUCTURE_SYNC_REQUEST_MESSAGE, Requester: fullNodeAddress}
	serializedMess, err := json.Marshal(initMess)
	if err != nil {
		return err
	}

	// assuming that all the nodes have the same blockchain - so sending request to only one
	err = sendWsMessage(&knownNodeAddresses[0], serializedMess, NODE_ENDPOINT)
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

	err = sendWsMessage(bootstrapAddress, serializedMess, BOOTSTRAP_ENDPOINT)
	if err != nil {
		return err
	}

	return nil
}

func sendWsMessageToAllNodes(message []byte) {
	for _, addr := range knownNodeAddresses {
		if addr != fullNodeAddress {
			err := sendWsMessage(&addr, message, NODE_ENDPOINT)
			if err != nil {
				log.Warnf("Cannot send digged block to %s. Reason: %s", addr, err)
			}
		}
	}
}

func sendWsMessage(targetAddress *string, message []byte, ep string) error {
	u := url.URL{Scheme: "ws", Host: *targetAddress, Path: ep}
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		return err
	}
	defer c.Close()

	err = c.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return err
	}
	log.Debugf("Message sent to %s", *targetAddress)

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

func insertNamespaces() {
	db.CreateNamespace(dbUrl, db.SelectNamespaceRequest{Namespace: "balance"})
}
