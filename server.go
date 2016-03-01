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
		trueErr, ok := err.(sqlite3.Error)
		if !ok {
			log.Fatal(err)
		}
		log.Fatal(trueErr)
	}
	defer db.Close()
	db.Ping()

	sqlStmt := "CREATE TABLE IF NOT EXISTS messages ( id INTEGER NOT NULL PRIMARY KEY, sender TEXT NOT NULL, msg TEXT NOT NULL);"
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
