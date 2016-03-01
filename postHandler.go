package main

import (
	"encoding/json"
	"log"
	"net/http"
	"fmt"
)

func postComment(w http.ResponseWriter, req *http.Request) {
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
	req.ParseForm()
	var userMessage = req.Form.Get("message")
	fmt.Printf("%s",userMessage)
	fmt.Printf("myVariable = %#v \n", req.Form)
	if userMessage == "" {
		p := responseStruct{false, "Message may not be empty"}
		res, err := json.Marshal(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, string(res), http.StatusBadRequest)
		}
		return
	}

	_, err := db.Exec("INSERT INTO messages(id, msg) VALUES (?, ?)", nil, userMessage)

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
	type responseStruct struct {
		Ok       bool
		ErrorMsg string
		Messages []string
	}
	w.Header().Set("Content-Type", "application/json")
	rows, err := db.Query("SELECT * FROM messages ORDER BY id ASC")
	defer rows.Close()
	if err != nil {
		log.Print(err)
		p := responseStruct{false, "DB error", make([]string, 0, 0)}
		res, err := json.Marshal(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, string(res), http.StatusBadRequest)
		}
		return
	}
	var messages []string
	for rows.Next() {
		var id int
		var message string
		err = rows.Scan(&id, &message)
		if err != nil {
			p := responseStruct{false, "DB error", make([]string, 0, 0)}
			res, err := json.Marshal(p)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
			} else {
				http.Error(w, string(res), http.StatusBadRequest)
			}
			log.Fatal(err)
		}
		messages = append(messages, message)
	}
	p := responseStruct{true, "", messages}
	json.NewEncoder(w).Encode(p)
}
