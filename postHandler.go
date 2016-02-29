package main

import (
	"net/http"
	"log"
	"encoding/json"
)

func postComment(w http.ResponseWriter, req *http.Request) {
	type successOrFail struct {
		ok bool
		errormsg string
	}
	if req.Method != "POST" {
		// User trying to GET the page
		p := successOrFail{false, "Can only POST this endpoint"}
		json.NewEncoder(w).Encode(p)
		return
	}
	req.ParseForm()
	var userMessage = req.Form.Get("message")
	_, err := db.Exec("INSERT INTO messages(id, url) VALUES (?, ?)", nil, userMessage)
	if err != nil {
		log.Fatal(err)
		p := successOrFail{false, "DB error"}
		json.NewEncoder(w).Encode(p)
	}
	
	p := successOrFail{true, ""}
	json.NewEncoder(w).Encode(p)
}
