package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
	"sync"
	"time"
)

type wsClient struct {
	// Pointer to a websocket connection
	conn *websocket.Conn

	// Messages that we need to send to this client
	outboundMsgs chan []byte

	// Unique atomic ID
	userId *uint32
}

var firstId uint32 = 0
var mutex = &sync.Mutex{}

func getNextId() *uint32 {
	mutex.Lock()
	returnInt := firstId
	firstId += 1
	mutex.Unlock()
	return &returnInt
}

type internalWebsocketMessageStruct struct {
	Message []byte
	Sender  []byte
	UserId  *uint32
}

type externalWebsocketMessageStruct struct {
	Message string
	Sender  string
}

// Reads input messages from this client in an infinite loop
func (client *wsClient) readMessages() {
	defer func() {
		// Unregister with lobby
		theLobby.unregister <- client
		// Close connection with client
		client.conn.Close()
		log.Println("readMessages exiting")
	}()
	// Infinite loop
	for {
		_, msg, err := client.conn.ReadMessage()
		if err != nil {
			log.Println("readMessages: ", err)
			closeMessage := websocket.FormatCloseMessage(websocket.CloseAbnormalClosure, "Unknown error!")
			client.conn.WriteControl(websocket.CloseMessage, closeMessage, time.Now().Add(5*time.Second))
			return
		}
		// Send it to the lobby
		incomingStruct := externalWebsocketMessageStruct{}
		err = json.Unmarshal(msg, &incomingStruct)
		if err != nil {
			log.Println("Error unmarshalling input", err)
			continue
		}
		theLobby.broadcast <- internalWebsocketMessageStruct{Message: []byte(incomingStruct.Message), Sender: []byte(incomingStruct.Sender), UserId: client.userId}
	}
}

func (client *wsClient) writeMessages() {
	defer client.conn.Close()
	defer log.Println("writeMessages exiting")
	for {
		// Read a message from the outbound channel
		msg := <-client.outboundMsgs
		// Send message to the browser
		err := client.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			log.Println("writeMessages: ", err)
			return
		}
	}
}

type lobby struct {
	// Map of clients, where the bool is pointless (we just want the map)
	clients map[*wsClient]bool

	// Channel on which to receive messages
	broadcast chan internalWebsocketMessageStruct

	// Make a new connection
	register chan *wsClient

	// Remove an old connection
	unregister chan *wsClient
}

var theLobby = lobby{
	clients:    make(map[*wsClient]bool),
	broadcast:  make(chan internalWebsocketMessageStruct),
	register:   make(chan *wsClient),
	unregister: make(chan *wsClient),
}

func (l *lobby) serveLobby() {
	for {
		select {
		case conn := <-l.register:
			// Get the map set up for this input
			l.clients[conn] = true
		case conn := <-l.unregister:
			// Delete the map entry
			delete(l.clients, conn)
			// Close the channel - prevent a resource leak
			close(conn.outboundMsgs)
		case msgStruct := <-l.broadcast:
			// We have a new inbound message!
			for conn := range l.clients {
				if msgStruct.UserId != conn.userId {
					externalStruct := externalWebsocketMessageStruct{Message: string(msgStruct.Message), Sender: string(msgStruct.Sender)}
					marshalled, err := json.Marshal(externalStruct)
					if err == nil {
						select {
						case conn.outboundMsgs <- marshalled:

						// do nothing, we just sent the message!
						default:
							// message did not send successfully
							close(conn.outboundMsgs)
							delete(l.clients, conn)
						}
					} else {
						log.Println("Error marshalling: ", err)
					}
				}

			}
		}
	}
}
