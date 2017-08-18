package main

import (
	"flag"
	"io"
	"log"
	"net/http"
)

var listen = flag.String("l", "localhost:3000", "hostport")

func main() {
	flag.Parse()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "Hello World\n")
	})
	log.Printf("listening on %s", *listen)
	log.Fatal(http.ListenAndServe(*listen, nil))
}
