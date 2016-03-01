package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func postComment(w http.ResponseWriter, req *http.Request) {

	type postJSON struct {
		Sender  string
		Message string
	}

	type responseStruct struct {
		Ok       bool
		ErrorMsg string
	}

	if req.Method != "POST" {
		// User trying to GET the page
		p := responseStruct{false, "Can only POST this endpoint"}
		res, err := json.Marshal(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, string(res), http.StatusBadRequest)
		}

		return
	}

	decoder := json.NewDecoder(req.Body)
	var t postJSON
	err := decoder.Decode(&t)

	if err != nil {
		panic(err)
	}
	var userMessage = t.Message

	if userMessage == "" {
		fmt.Printf("empty")

		p := responseStruct{false, "Message may not be empty"}
		res, err := json.Marshal(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, string(res), http.StatusBadRequest)
		}
		return
	}

	_, err = db.Exec("INSERT INTO messages(id, sender, msg) VALUES (?, ?, ?)", nil, t.Sender, t.Message)

	if err != nil {
		log.Print(err)
		p := responseStruct{false, "DB error"}
		res, err := json.Marshal(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, string(res), http.StatusBadRequest)
		}
		return
	}

	p := responseStruct{true, "TEST"}
	json.NewEncoder(w).Encode(p)
}

func getAllMessages(w http.ResponseWriter, req *http.Request) {
	type messageStruct struct {
		Sender  string
		Message string
	}
	type responseStruct struct {
		Ok       bool
		ErrorMsg string
		Messages []messageStruct
	}
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT * FROM messages ORDER BY id ASC")
	if err != nil {
		log.Print(err)
		p := responseStruct{false, "DB error", make([]messageStruct, 0, 0)}
		res, err := json.Marshal(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, string(res), http.StatusBadRequest)
		}
		return
	}
	defer rows.Close()

	var messages []messageStruct
	for rows.Next() {
		var id int
		var message string
		var sender string

		err = rows.Scan(&id, &sender, &message)
		if err != nil {
			p := responseStruct{false, "DB error", make([]messageStruct, 0, 0)}
			res, err := json.Marshal(p)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				http.Error(w, string(res), http.StatusBadRequest)
			}
			log.Fatal(err)
		}
		messagewrapper := messageStruct{Sender: sender, Message: message}
		messages = append(messages, messagewrapper)

	}
	p := responseStruct{true, "", messages}
	json.NewEncoder(w).Encode(p)
}
