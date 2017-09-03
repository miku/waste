package main

import (
	"context"
	"flag"
	"io"
	"log"
	"net/http"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

var listen = flag.String("listen", "localhost:3000", "hostport")

func main() {
	flag.Parse()
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		log.Println("creating new docker client")
		cli, err := client.NewEnvClient()
		if err != nil {
			http.Error(w, "cannot create docker client", http.StatusInternalServerError)
			return
		}

		log.Println("pulling image")
		_, err = cli.ImagePull(ctx, "docker.io/library/alpine", types.ImagePullOptions{})
		if err != nil {
			http.Error(w, "cannot pull image", http.StatusInternalServerError)
			return
		}

		log.Println("creating container")
		resp, err := cli.ContainerCreate(ctx, &container.Config{
			Image: "alpine",
			Cmd:   []string{"uname", "-a"},
		}, nil, nil, "")

		if err != nil {
			http.Error(w, "cannot create container", http.StatusInternalServerError)
			return
		}

		log.Printf("starting container: %s", resp.ID)
		if err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
			http.Error(w, "cannot start container", http.StatusInternalServerError)
			return
		}

		log.Printf("waiting for container: %s", resp.ID)
		if _, err = cli.ContainerWait(ctx, resp.ID); err != nil {
			http.Error(w, "container wait failed", http.StatusInternalServerError)
			return
		}

		out, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
		if err != nil {
			http.Error(w, "cannot access logs", http.StatusInternalServerError)
			return
		}
		if _, err := io.Copy(w, out); err != nil {
			http.Error(w, "cannot write response", http.StatusInternalServerError)
			return
		}

		log.Printf("removing container: %s", resp.ID)
		if err := cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{}); err != nil {
			http.Error(w, "cannot remove container", http.StatusInternalServerError)
			return
		}
	})
	log.Printf("listening on %s", *listen)
	log.Fatal(http.ListenAndServe(*listen, nil))
}
