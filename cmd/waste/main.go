package main

import (
	"archive/tar"
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/miku/waste"
	"github.com/moby/moby/client"
)

const version = "0.1.0"

var (
	listen    = flag.String("listen", "localhost:3000", "hostport")
	imageRef  = flag.String("ref", "docker.io/library/alpine", "image reference")
	imageName = flag.String("image", "alpine", "image name")
	timeout   = flag.Duration("timeout", 10*time.Second, "timeout")
)

var banner = fmt.Sprintf(`
██╗    ██╗ █████╗ ███████╗████████╗███████╗
██║    ██║██╔══██╗██╔════╝╚══██╔══╝██╔════╝
██║ █╗ ██║███████║███████╗   ██║   █████╗  
██║███╗██║██╔══██║╚════██║   ██║   ██╔══╝  
╚███╔███╔╝██║  ██║███████║   ██║   ███████╗
 ╚══╝╚══╝ ╚═╝  ╚═╝╚══════╝   ╚═╝   ╚══════╝
                                           

Welcome to waste - your premium "cat as a service" provider.

This server accepts HTTP requests and will copy the request body into a
container, run the "cat" command on the input and stream the output back to
stdout.

Example, inspect a local file:

    $ curl http://%s --data-binary @README.md

Or run the docker webpage to a docker container first:

    $ curl http://%s --data-binary @<(curl -sL http://www.docker.io)

Version: %s
Startup: %s
`, *listen, *listen, version, time.Now())

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
		_, err = io.Copy(tw, bytes.NewReader(b))
		if err != nil {
			http.Error(w, fmt.Sprintf("could not tar content: %s", err),
				http.StatusInternalServerError)
			return
		}
		log.Debug("archived ", len(b), " bytes from request body")

		// Collect container output here.
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
		io.WriteString(w, "\n")
		log.Debug("operation finished successfully")
	})

	fmt.Println(banner)

	// Make sure docker runs, or fail early.
	cli, err := client.NewEnvClient()
	if err != nil {
		log.Fatal(err)
	}
	ping, err := cli.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	log.Debug("docker is up: ", ping.APIVersion)

	log.Fatal(http.ListenAndServe(*listen, nil))
}
