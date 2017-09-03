package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/miku/waste"
	log "github.com/sirupsen/logrus"
)

var (
	listen    = flag.String("listen", "localhost:3000", "hostport")
	imageRef  = flag.String("ref", "docker.io/library/alpine", "image reference")
	imageName = flag.String("image", "alpine", "image name")
)

func main() {
	flag.Parse()
	log.SetLevel(log.DebugLevel)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		var buf bytes.Buffer
		wrapper := waste.WrapDocker{
			ImageRef:  *imageRef,
			ImageName: *imageName,
			Cmd:       []string{"uname", "-a"},
			Writer:    &buf,
			Timeout:   5 * time.Second,
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
