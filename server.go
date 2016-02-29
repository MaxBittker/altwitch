package main

import (
  "fmt"
  "net/http"
  "html"
  "log"
)

func main() {
	fmt.Println("Hello, 世界")

  handleRequest := func(w http.ResponseWriter, r *http.Request){
    fmt.Fprintf(w, "<h1>%q</h1>", html.EscapeString(r.URL.Path))
  }

  http.HandleFunc("/test", handleRequest)

  log.Fatal(http.ListenAndServe(":8080", nil))
}
