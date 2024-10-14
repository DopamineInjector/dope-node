package communication

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func RegisterEndpoints(boostrapServerAddress string, isBoostrapServer bool) {
	http.HandleFunc("/ws", wsHandler)
	if isBoostrapServer {
		http.HandleFunc("/addresses", httpHandler)
	} else {
		nodeAddresses := fetchNodeAddresses
		log.Println(nodeAddresses)
	}

	http.ListenAndServe(":7312", nil)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	connection := getWebsocketConnection(w, r)

	for {
		_, msg, err := connection.ReadMessage()
		if err != nil {
			log.Fatal("Error reading message:", err)
			break
		}

		log.Println("Handling websocket message from: " + r.Host)
		var receivedMessage Message
		json.Unmarshal(msg, &receivedMessage)
		fmt.Println(receivedMessage.Content)
	}
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling http request from: " + r.Host)
}

func getWebsocketConnection(w http.ResponseWriter, r *http.Request) *websocket.Conn {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal("Error upgrading to websocket:", err)
		conn.Close()
	} else {
		log.Println("Connection handled successfully")
	}

	return conn
}

func sendMessage(connection websocket.Conn, response []byte) {
	status := connection.WriteMessage(websocket.TextMessage, response)
	if status != nil {
		log.Println("Error writing message:", status)
	}
}

func fetchNodeAddresses(bootstrapSeverAddress string) []string {
	resp, err := http.Get("http://" + bootstrapSeverAddress + "/addresses")
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
