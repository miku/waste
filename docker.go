package waste

import (
	"bufio"
	"context"
	"io"
	"time"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	log "github.com/sirupsen/logrus"
)

// WrapDocker wraps information we need to run the isolated process.
type WrapDocker struct {
	ImageLocation string
	ImageName     string
	Cmd           []string
	Writer        io.Writer
	Timeout       time.Duration
}

// Run runs docker and executes the command in the container.
func (w WrapDocker) Run() error {
	ctx := context.Background()
	if w.Timeout > 0 {
		log.Debug("running with a timeout of %v", w.Timeout)
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, w.Timeout)
		defer func() {
			log.Debug("operation timed out")
			cancel()
		}()
	}

	log.Debug("creating new docker client")
	cli, err := client.NewEnvClient()
	if err != nil {
		return err
	}

	log.Debug("pulling image from %s", w.ImageLocation)
	_, err = cli.ImagePull(ctx, w.ImageLocation, types.ImagePullOptions{})
	if err != nil {
		return err
	}

	log.Debug("creating container from %s", w.ImageName)
	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: w.ImageName,
		Cmd:   w.Cmd,
	}, nil, nil, "")

	if err != nil {
		return err
	}

	log.Debug("starting container %s", resp.ID)
	if err = cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	defer func() {
		log.Debug("removing container %s", resp.ID)
		err = cli.ContainerRemove(ctx, resp.ID, types.ContainerRemoveOptions{})
	}()

	log.Debug("waiting for container %s", resp.ID)
	if _, err = cli.ContainerWait(ctx, resp.ID); err != nil {
		return err
	}

	reader, err := cli.ContainerLogs(ctx, resp.ID, types.ContainerLogsOptions{ShowStdout: true})
	if err != nil {
		return err
	}
	defer reader.Close()
	n, err := io.Copy(w.Writer, reader)
	log.Debug("%d bytes read from application", n)
	if bw, ok := w.Writer.(*bufio.Writer); ok {
		bw.Flush()
	}
	return err
}
