package main

import (
	"archive/tar"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/miku/waste"
	log "github.com/sirupsen/logrus"
)

var (
	listen    = flag.String("listen", "localhost:3000", "hostport")
	imageRef  = flag.String("ref", "docker.io/library/alpine", "image reference")
	imageName = flag.String("image", "alpine", "image name")
	timeout   = flag.Duration("timeout", 10*time.Second, "timeout")
)

func main() {
	flag.Parse()
	log.SetLevel(log.DebugLevel)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		b, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("could not read request body: %s", err),
				http.StatusInternalServerError)
			return
		}

		log.Debug("request body contains newlines: ", bytes.Contains(b, []byte("\n")))

		// Tar the body, see: https://git.io/v50x4.
		hdr := &tar.Header{
			Name:       "body",
			Mode:       0644,
			Size:       int64(len(b)),
			AccessTime: time.Now(),
			ChangeTime: time.Now(),
			ModTime:    time.Now(),
		}
		var buf bytes.Buffer
		tw := tar.NewWriter(&buf)
		if err := tw.WriteHeader(hdr); err != nil {
			http.Error(w, fmt.Sprintf("could write tar header: %s", err),
				http.StatusInternalServerError)
			return
		}
		n, err := io.Copy(tw, bytes.NewReader(b))
		if err != nil {
			http.Error(w, fmt.Sprintf("could not tar content: %s", err),
				http.StatusInternalServerError)
			return
		}
		log.Debug("archived ", n, " bytes from request body")

		// Collection container output here.
		var bufOut bytes.Buffer

		wrapper := waste.WrapDocker{
			ImageRef:  *imageRef,
			ImageName: *imageName,
			Cmd:       []string{"cat", "-u", "/mnt/body"},
			Reader:    &buf,
			Writer:    &bufOut,
			Timeout:   *timeout,
		}

		defer r.Body.Close()

		if err := wrapper.Run(); err != nil {
			http.Error(w,
				fmt.Sprintf("failed to run container: %s", err),
				http.StatusInternalServerError)
			return
		}
		if _, err := io.Copy(w, &bufOut); err != nil {
			http.Error(w,
				fmt.Sprintf("failed copy stream: %s", err),
				http.StatusInternalServerError)
			return
		}
		log.Debug("operation finished successfully")
	})
	log.Printf("listening on %s", *listen)
	log.Fatal(http.ListenAndServe(*listen, nil))
}
