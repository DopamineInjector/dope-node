package communication

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/gorilla/websocket"
)

var knownAddresses = []string{}

func RegisterEndpoints(boostrapServerAddress string, isBoostrapServer bool, port int) {
	const WS_ENDPOINT = "/ws"
	const ADDRESSES_ENDPOINT = "/addresses"
	senderAddr := fmt.Sprintf("127.0.0.1:%d", port) // For now 127.0.0.1 - later real IP

	http.HandleFunc(WS_ENDPOINT, wsHandler)
	if isBoostrapServer {
		http.HandleFunc(ADDRESSES_ENDPOINT, addressesRequestHandler)
	} else {
		addresses, err := fetchNodeAddresses(boostrapServerAddress, senderAddr)
		if err != nil {
			log.Println("Failed to fetch addresses from the boostrap server - " + boostrapServerAddress)
			addresses = []string{}
		}
		for _, addr := range addresses {
			u := url.URL{Scheme: "ws", Host: addr, Path: WS_ENDPOINT}
			conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
			if err != nil {
				log.Printf("Failed to connect to %s: %v", addr, err)
				continue
			}
			defer conn.Close()

			// for now it sends hello message.
			messageToSend, err := json.Marshal(Message{Content: "hello"})
			if err != nil {
				log.Println("Parsing error", err)
			}

			sendMessage(*conn, messageToSend)
		}
	}

	log.Println("Running server")
	declaredPort := fmt.Sprintf(":%d", port)
	log.Println(http.ListenAndServe(declaredPort, nil))
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	connection, err := getWebsocketConnection(w, r)
	if err != nil {
		log.Println("Failed to initialize websocket connection")
	}

	for {
		_, msg, err := connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("WebSocket closed unexpectedly:", err)
			} else {
				log.Println("Error reading message:", err)
			}
			break
		}

		var receivedMessage Message
		err = json.Unmarshal(msg, &receivedMessage)
		if err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}
		fmt.Println("Received message content:", receivedMessage.Content)
	}
}

func addressesRequestHandler(w http.ResponseWriter, r *http.Request) {
	addressesResponse := FetchBootstrapAddressesResponse{
		Addresses: knownAddresses,
	}

	ip := r.URL.Query().Get("sender")
	w.Header().Set("Content-Type", "application/json")

	if ip != "" {
		log.Println("Added address: " + ip)
		knownAddresses = append(knownAddresses, ip)
		json.NewEncoder(w).Encode(addressesResponse)
	} else {
		log.Println("Received request with no sender specified in a http parameter")
		json.NewEncoder(w).Encode([]string{})
	}
}

func getWebsocketConnection(w http.ResponseWriter, r *http.Request) (*websocket.Conn, error) {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to WebSocket:", err)
	}

	return conn, err
}

func sendMessage(connection websocket.Conn, mess []byte) error {
	err := connection.WriteMessage(websocket.TextMessage, mess)
	if err != nil {
		log.Println("Error writing message:", err)
		return err
	}

	return nil
}

func fetchNodeAddresses(bootstrapSeverAddress string, addressToRegister string) ([]string, error) {
	resp, err := http.Get("http://" + bootstrapSeverAddress + "/addresses?sender=" + addressToRegister)
	if err != nil {
		log.Println("Error making the GET request:", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("Error reading the response body:", err)
	}

	var parsedResponse FetchBootstrapAddressesResponse
	err = json.Unmarshal(body, &parsedResponse)
	if err != nil {
		log.Println("Error unmarshaling the JSON:", err)
	}

	log.Println("Received node addresses:")
	for _, nodeAddress := range parsedResponse.Addresses {
		log.Println("\t- " + nodeAddress)
	}

	return parsedResponse.Addresses, err
}
