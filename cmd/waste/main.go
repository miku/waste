package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/miku/waste"
	log "github.com/sirupsen/logrus"
)

var listen = flag.String("listen", "localhost:3000", "hostport")

func main() {
	flag.Parse()

	log.SetLevel(log.DebugLevel)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		var buf bytes.Buffer
		wrapper := waste.WrapDocker{
			ImageLocation: "docker.io/library/alpine",
			ImageName:     "alpine",
			Cmd:           []string{"uname", "-a"},
			Writer:        &buf,
		}

		if err := wrapper.Run(); err != nil {
			http.Error(w,
				fmt.Sprintf("failed to run container: %s", err),
				http.StatusInternalServerError)
			return
		}
		if _, err := io.Copy(w, &buf); err != nil {
			http.Error(w,
				fmt.Sprintf("failed copy stream: %s", err),
				http.StatusInternalServerError)
		}
	})
	log.Printf("listening on %s", *listen)
	log.Fatal(http.ListenAndServe(*listen, nil))
}
