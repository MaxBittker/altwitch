package main

import (
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
	client.conn.WriteMessage(websocket.TextMessage, []byte("Welcome to the chat room"))
	// Get this client on the list
	theLobby.register <- client
	// Parallelize writing messages
	go client.writeMessages()
	// May as well read in this thread/goroutine, we're done with it
	client.readMessages()
}
