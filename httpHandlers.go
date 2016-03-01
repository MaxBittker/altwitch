package main

import (
	"encoding/json"
	"log"
	"net/http"
)

func newMessage(w http.ResponseWriter, req *http.Request) {

	// Input from the user/browser
	type postJSON struct {
		Sender  string
		Message string
	}

	// What our response is going to look like
	type responseStruct struct {
		Ok       bool
		ErrorMsg string
	}

	// We are always going to response with JSON, so set the Content-Type now
	w.Header().Set("Content-Type", "application/json")

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

	// Decode the user's input to us
	decoder := json.NewDecoder(req.Body)
	var t postJSON
	// Convert the user's data into our struct
	err := decoder.Decode(&t)

	if err != nil {
		// Don't crash if we get malformed input
		log.Println(err)
		p := responseStruct{false, "Badly formed JSON input"}
		res, err := json.Marshal(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, string(res), http.StatusBadRequest)
		}
		return
	}
	// Get the actual message from the properly constructed struct
	userMessage := t.Message

	// Sanity check: don't allow blank messages
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

	// Insert this new message into the database.
	// We don't care about the result return value because we don't need
	// the ID of this new message
	_, err = db.Exec("INSERT INTO messages(id, sender, msg) VALUES (?, ?, ?)", nil, t.Sender, t.Message)

	if err != nil {
		log.Print(err)
		p := responseStruct{false, "Internal database error"}
		res, err := json.Marshal(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, string(res), http.StatusBadRequest)
		}
		return
	}

	// Everything is good! Tell the user there was no error, and of course there is no
	// corresponding error string
	p := responseStruct{true, ""}
	json.NewEncoder(w).Encode(p)
}

func getAllMessages(w http.ResponseWriter, req *http.Request) {
	// One particular message by one particular author
	type messageStruct struct {
		Sender  string
		Message string
	}
	// What we are going to return to the user
	type responseStruct struct {
		Ok       bool
		ErrorMsg string
		Messages []messageStruct
	}
	// We are always going to response with JSON, so set the Content-Type now
	w.Header().Set("Content-Type", "application/json")
	// Get all the messages, but preserve the temporal order (by id)
	rows, err := db.Query("SELECT * FROM messages ORDER BY id ASC")
	if err != nil {
		log.Print(err)
		// Database error, report it
		p := responseStruct{false, "DB error", make([]messageStruct, 0, 0)}
		// Try to turn our database error into a valid JSON string
		res, err := json.Marshal(p)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
		} else {
			http.Error(w, string(res), http.StatusBadRequest)
		}
		return
	}
	// Remember to close at the end of this function
	defer rows.Close()

	// Slice to hold all our messages
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
		// Construct one method struct from one database row
		messagewrapper := messageStruct{Sender: sender, Message: message}
		// Append to our struct
		messages = append(messages, messagewrapper)

	}
	// Return all our messages, with unbounded(!) size
	p := responseStruct{true, "", messages}
	json.NewEncoder(w).Encode(p)
}
