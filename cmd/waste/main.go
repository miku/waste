package main

import (
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello World\n")
	})
	log.Fatal(http.ListenAndServe("localhost:3000", nil))
}
