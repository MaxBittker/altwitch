package main

import (
	"database/sql"
	"fmt"
	"html"
	"log"
	"net/http"
	"github.com/mattn/go-sqlite3"
)



func main() {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		trueErr, ok := err.(sqlite3.Error)
		if !ok {
			log.Fatal(err)
		}
		log.Fatal(trueErr)
	}
	defer db.Close()
	db.Ping()


	handleRequest := func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<h1>%q</h1>", html.EscapeString(r.URL.Path))
	}

	http.HandleFunc("/test", handleRequest)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
