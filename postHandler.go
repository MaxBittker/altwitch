package main

import (
	"net/http"
)

func postComment(w http.ResponseWriter, req *http.Request) {
	if req.Method != "POST" {
		// User trying to GET the page
		http.Error(w, "Only POST", 400);
		return
	}
	req.ParseForm()
	var userMessage = req.Form.Get("message")
	res, err := db.Exec("INSERT INTO messages(id, url) VALUES (?, ?)", nil, userMessage)
	if err != nil {
		
	}
}
