package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net/http"

	"github.com/moby/moby/api/types"
	"github.com/moby/moby/api/types/container"
	"github.com/moby/moby/client"
)

var listen = flag.String("l", "localhost:3000", "hostport")

func main() {
	flag.Parse()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		c, err := client.NewEnvClient()
		if err != nil {
			http.Error(w, "cannot create docker client", http.StatusInternalServerError)
			return
		}
		_, err := c.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{})
		if err != nil {
			http.Error(w, "cannot pull image", http.StatusInternalServerError)
			return
		}
		resp, err := c.ContainerCreate(ctx, &container.Config{
			Image: "alpine",
			Cmd: []string{"uname", "-a"}
		}, nil, nil, "")

		if err != nil {
			http.Error(w, "cannot create container", http.StatusInternalServerError)
			return
		}

		if err := c.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			http.Error(w, "cannot start container", http.StatusInternalServerError)
			return
		}
		if err := c.ContainerWait(ctx, resp.ID); err != nil {
			http.Error(w, "container wait failed", http.StatusInternalServerError)
			return
		}
		out, err := c.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
		if err != nil {
			http.Error(w, "cannot access logs", http.StatusInternalServerError)
			return
		}
		if _, err := io.Copy(w, out); err != nil {
			http.Error(w, "cannot write response", http.StatusInternalServerError)
			return
		}
	})
	log.Printf("listening on %s", *listen)
	log.Fatal(http.ListenAndServe(*listen, nil))
}
