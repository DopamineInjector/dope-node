package communication

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

func RunWebsocketListener() {
	http.HandleFunc("/ws", wsHandler)
	log.Printf("Websocket server running")
	http.ListenAndServe(":7312", nil)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	connection := getWebsocketConnection(w, r)

	for {
		_, msg, err := connection.ReadMessage()
		if err != nil {
			log.Println("Error reading message:", err)
			break
		}
		handleMessage(msg)
	}
}

func getWebsocketConnection(w http.ResponseWriter, r *http.Request) *websocket.Conn {
	upgrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Error upgrading to websocket:", err)
		conn.Close()
	}

	return conn
}

func handleMessage(message []byte) {
	log.Printf("Handling message")
	fmt.Println(message)
}

func sendResponse(connection websocket.Conn, response []byte) {
	status := connection.WriteMessage(websocket.TextMessage, response)
	if status != nil {
		log.Println("Error writing message:", status)
	}
}
