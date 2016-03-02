package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

func upgradeToWebsockets(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		// Allow connections from any origin (for testing)
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Websocket upgrader: ", err)
		return
	}

	client := &wsClient{conn: conn, outboundMsgs: make(chan []byte), userId: getNextId()}
	// Welcome message, removable
	outgoingMsg := externalWebsocketMessageStruct{Message: "Welcome to the chat room", Sender: "The Admins"}
	marshalled, err := json.Marshal(outgoingMsg)
	if err != nil {
		log.Println("Error marshalling welcome message")
	}
	client.conn.WriteMessage(websocket.TextMessage, marshalled)
	// Get this client on the list
	theLobby.register <- client
	// Parallelize writing messages
	go client.writeMessages()
	// May as well read in this thread/goroutine, we're done with it
	client.readMessages()
}
