package main

import (
	"database/sql"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"html"
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

	sqlStmt := "CREATE TABLE IF NOT EXISTS messages ( id INTEGER NOT NULL PRIMARY KEY, msg TEXT NOT NULL);"
	_, err = db.Exec(sqlStmt)
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/newMessage", postComment)
	http.HandleFunc("/getAllMessages", getAllMessages)

	handleRequest := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<h1>%q</h1>", html.EscapeString(r.URL.Path))
	}

	http.HandleFunc("/test", handleRequest)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
