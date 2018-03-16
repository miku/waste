package waste

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
)

// WrapDocker wraps information we need to run the isolated process. Reader is
// read and saved inside the container. Writer collects stdout and stderr.
type WrapDocker struct {
	ImageRef  string
	ImageName string
	Cmd       []string
	Reader    io.Reader
	Writer    io.Writer
	Timeout   time.Duration
}

// Run runs docker and executes the command in the container.
func (w WrapDocker) Run() error {
	ctx := context.Background()
	if w.Timeout > 0 {
		log.Debug("running with a timeout of ", w.Timeout)
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, w.Timeout)
		defer cancel()
	}

	log.Debug("creating new docker client")
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	log.Debug("pulling image from ", w.ImageRef)
	_, err = cli.ImagePull(ctx, w.ImageRef, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	imgs, err := cli.ImageList(ctx, types.ImageListOptions{})
	if err != nil {
		return err
	}

	log.Printf("host has %d images", len(imgs))
	for _, summary := range imgs {
		log.Printf("%s: %d", summary.ID, summary.Created)
	}

	log.Debug("creating container from ", w.ImageName)
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image:           w.ImageName,
		Cmd:             w.Cmd,
		NetworkDisabled: true,
		Tty:             true,
	}, nil, nil, "")

	if err != nil {
		return err
	}

	defer func() {
		log.Debug("removing container ", resp.ID)
		err = cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{})
	}()

	cr := &Counter{r: w.Reader}

	log.Debug("copying data into container")
	if err = cli.CopyToContainer(ctx, resp.ID, "/mnt", cr,
		types.CopyToContainerOptions{}); err != nil {
		return err
	}

	log.Debug(cr.N(), " bytes written into container")

	stat, err := cli.ContainerStatPath(ctx, resp.ID, "/mnt/body")
	if err != nil {
		return err
	}

	log.Debug(fmt.Sprintf("stat: %v", stat))

	log.Debug("starting container ", resp.ID)
	if err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	log.Debug("waiting for container ", resp.ID)
	resultC, errC := cli.ContainerWait(ctx, resp.ID, container.WaitConditionNotRunning)
	select {
	case err := <-errC:
		return err
	case <-resultC:
	}

	reader, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return err
	}
	defer reader.Close()
	n, err := io.Copy(w.Writer, reader)
	log.Debug(n, " bytes read from application")
	if bw, ok := w.Writer.(*bufio.Writer); ok {
		bw.Flush()
	}
	return err
}
