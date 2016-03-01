package main

import (
	"database/sql"
	"github.com/mattn/go-sqlite3"
	"log"
	"net/http"
)

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite3", ":memory:")
	if err != nil {
		// Cast the error interface to a sqlite error
		// May allow us to get some more information
		trueErr, ok := err.(sqlite3.Error)
		// Check if cast failed
		if !ok {
			log.Fatal(err)
		}
		log.Fatal(trueErr)
	}
	// Remember to close the database
	defer db.Close()
	// Actually make a connection, will write the database to disk if necessary
	db.Ping()

	sqlStmt := "CREATE TABLE IF NOT EXISTS messages ( id INTEGER NOT NULL PRIMARY KEY, sender TEXT NOT NULL, msg TEXT NOT NULL);"
	// We don't care about the result, only whether or not it failed
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/newMessage", newMessage)
	http.HandleFunc("/getAllMessages", getAllMessages)
	http.HandleFunc("/websocket", upgradeToWebsockets)
	go theLobby.serveLobby()
	http.Handle("/", http.FileServer(http.Dir("./client")))
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatal(err)
	}

}
