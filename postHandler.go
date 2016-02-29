package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func postComment(w http.ResponseWriter, req *http.Request) {
	type successOrFail struct {
		Ok       bool
		ErrorMsg string
	}
	if req.Method != "POST" {
		// User trying to GET the page
		p := successOrFail{false, "Can only POST this endpoint"}
		json.NewEncoder(w).Encode(p)
		return
	}
	req.ParseForm()
	var userMessage = req.Form.Get("message")
	_, err := db.Exec("INSERT INTO messages(id, msg) VALUES (?, ?)", nil, userMessage)

	if err != nil {
		log.Print(err)
		p := successOrFail{false, "DB error"}
		json.NewEncoder(w).Encode(p)
		return
	}

	p := successOrFail{true, ""}
	json.NewEncoder(w).Encode(p)
}
