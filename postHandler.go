package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func postComment(w http.ResponseWriter, req *http.Request) {
	type responseStruct struct {
		Ok       bool
		ErrorMsg string
	}
	if req.Method != "POST" {
		// User trying to GET the page
		p := responseStruct{false, "Can only POST this endpoint"}
		json.NewEncoder(w).Encode(p)
		return
	}
	req.ParseForm()
	var userMessage = req.Form.Get("message")
	if userMessage == "" {
		p := responseStruct{false, "Message may not be empty"}
		json.NewEncoder(w).Encode(p)
		return
	}

	_, err := db.Exec("INSERT INTO messages(id, msg) VALUES (?, ?)", nil, userMessage)

	if err != nil {
		log.Print(err)
		p := responseStruct{false, "DB error"}
		json.NewEncoder(w).Encode(p)
		return
	}

	p := responseStruct{true, ""}
	json.NewEncoder(w).Encode(p)
}

func getAllMessages(w http.ResponseWriter, req *http.Request) {
	type resultStruct struct {
		Ok       bool
		ErrorMsg string
		Messages []string
	}
	rows, err := db.Query("SELECT * FROM messages")
	defer rows.Close()
	if err != nil {
		log.Print(err)
		p := resultStruct{false, "DB error", make([]string, 0, 0)}
		json.NewEncoder(w).Encode(p)
		return
	}
	for rows.Next() {
		var id int
		var message string
		err = rows.Scan(&id, &message)
	}

}
