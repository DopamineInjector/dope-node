package communication

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/websocket"
)

var knownAddresses = []string{}

func RegisterEndpoints(boostrapServerAddress string, isBoostrapServer bool, port int) {
	http.HandleFunc("/ws", wsHandler)
	log.Println("Running server")
	if isBoostrapServer {
		http.HandleFunc("/addresses", httpHandler)
	} else {
		// This sends 'hello' message to all known nodes
		addresses := fetchNodeAddresses(boostrapServerAddress, "127.0.0.1:"+strconv.Itoa(port))
		for _, addr := range addresses {
			u := url.URL{Scheme: "ws", Host: addr, Path: "/ws"}
			conn, _, err := websocket.DefaultDialer.Dial(u.String(), nil)
			if err != nil {
				log.Printf("Failed to connect to %s: %v", addr, err)
				continue
			}
			defer conn.Close()

			messageToSend, _ := json.Marshal(Message{Content: "hello"})
			sendMessage(*conn, messageToSend)
		}
	}

	http.ListenAndServe(":"+strconv.Itoa(port), nil)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	connection := getWebsocketConnection(w, r)

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

		log.Println("Handling WebSocket message from:", r.Host)
		var receivedMessage Message
		err = json.Unmarshal(msg, &receivedMessage)
		if err != nil {
			log.Println("Error unmarshalling message:", err)
			continue
		}
		fmt.Println("Received message content:", receivedMessage.Content)
	}
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling http request")
	addressesResponse := FetchBootstrapAddressesResponse{
		Addresses: knownAddresses,
	}

	ip := r.URL.Query().Get("sender")
	log.Println("Added address: " + ip)
	knownAddresses = append(knownAddresses, ip)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(addressesResponse)
}

func getWebsocketConnection(w http.ResponseWriter, r *http.Request) *websocket.Conn {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Error upgrading to WebSocket:", err)
		return nil
	}

	return conn
}

func sendMessage(connection websocket.Conn, mess []byte) {
	err := connection.WriteMessage(websocket.TextMessage, mess)
	if err != nil {
		log.Println("Error writing message:", err)
		return
	}
}

func fetchNodeAddresses(bootstrapSeverAddress string, addressToRegister string) []string {
	resp, err := http.Get("http://" + bootstrapSeverAddress + "/addresses?sender=" + addressToRegister)
	if err != nil {
		log.Fatal("Error making the GET request:", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal("Error reading the response body:", err)
	}

	var parsedResponse FetchBootstrapAddressesResponse
	err = json.Unmarshal(body, &parsedResponse)
	if err != nil {
		log.Fatal("Error unmarshaling the JSON:", err)
	}

	log.Println("Received node addresses:")
	for _, nodeAddress := range parsedResponse.Addresses {
		log.Println("\t- " + nodeAddress)
	}

	return parsedResponse.Addresses
}
