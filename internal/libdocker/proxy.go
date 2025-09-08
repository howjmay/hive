package libdocker

import (
	"context"
	"io"
	"log/slog"
	"net"
	"net/http"
	"sync"

	"github.com/iotaledger/wasp-hive/hiveproxy"
	"github.com/iotaledger/wasp-hive/internal/libhive"
)

const hiveproxyTag = "hive/hiveproxy"
const l1nodeTag = "hive/l1-node"

// Build builds the hiveproxy image.
func (cb *ContainerBackend) Build(ctx context.Context, b libhive.Builder) error {
	// build l1 node image first
	if err := b.BuildImageRelative(ctx, l1nodeTag, hiveproxy.Source, "Dockerfile.l1"); err != nil {
		return err
	}

	return b.BuildImage(ctx, hiveproxyTag, hiveproxy.Source)
}

// ServeAPI starts the API server.
func (cb *ContainerBackend) ServeAPI(ctx context.Context, h http.Handler) (libhive.APIServer, error) {
	inR, inW := io.Pipe()
	outR, outW := io.Pipe()

	// Create labels for hiveproxy container.
	proxyLabels := libhive.NewBaseLabels(cb.hiveInstanceID, cb.hiveVersion)
	proxyLabels[libhive.LabelHiveType] = libhive.ContainerTypeProxy

	// Generate container name.
	containerName := libhive.GenerateProxyContainerName()

	opts := libhive.ContainerOptions{Output: outW, Input: inR, Labels: proxyLabels, Name: containerName}
	id, err := cb.CreateContainer(ctx, hiveproxyTag, opts)
	if err != nil {
		return nil, err
	}

	// Launch the proxy server before starting the container.
	var (
		proxy     *hiveproxy.Proxy
		proxyErrC = make(chan error, 1)
	)
	go func() {
		var err error
		proxy, err = hiveproxy.RunBackend(outR, inW, h)
		if err != nil {
			slog.Error("proxy backend startup failed", "err", err)
		}
		proxyErrC <- err
	}()

	// Now start the container.
	info, err := cb.StartContainer(ctx, id, opts)
	if err != nil {
		cb.DeleteContainer(id)
		return nil, err
	}

	// Proxy server should come up.
	select {
	case err := <-proxyErrC:
		if err != nil {
			cb.DeleteContainer(id)
			return nil, err
		}
	}

	// Register proxy in ContainerBackend, so it can be used for CheckLive.
	cb.proxy = proxy

	srv := &proxyContainer{
		cb:              cb,
		containerID:     id,
		containerIP:     net.ParseIP(info.IP),
		containerWait:   info.Wait,
		containerStdin:  inR,
		containerStdout: outW,
		proxy:           proxy,
	}

	// Register proxy in ContainerBackend, so it can be used for CheckLive.
	cb.proxy = proxy
	slog.Info("hiveproxy started", "container", id[:12], "addr", srv.Addr())

	l1ContainerID, err := cb.StartL1Node(ctx)
	if err != nil {
		return nil, err
	}
	srv.l1NodeContainerID = l1ContainerID

	return srv, nil
}

type proxyContainer struct {
	cb *ContainerBackend

	containerID       string
	l1NodeContainerID string
	containerIP       net.IP
	containerStdin    *io.PipeReader
	containerStdout   *io.PipeWriter
	containerWait     func()
	proxy             *hiveproxy.Proxy

	stopping sync.Once
	stopErr  error
}

// Addr returns the listening address of the proxy server.
func (c *proxyContainer) Addr() net.Addr {
	return &net.TCPAddr{IP: c.containerIP, Port: 8081}
}

// Stop terminates the proxy container.
func (c *proxyContainer) Close() error {
	c.stopping.Do(func() {
		// Unregister proxy in backend.
		c.cb.proxy = nil

		// Stop the container.
		c.containerStdin.Close()
		c.containerStdout.Close()
		c.stopErr = c.cb.DeleteContainer(c.containerID)
		c.stopErr = c.cb.DeleteContainer(c.l1NodeContainerID)
		c.containerWait()

		// Stop the local HTTP receiver.
		c.proxy.Close()
	})
	return c.stopErr
}

func (cb *ContainerBackend) StartL1Node(ctx context.Context) (string, error) {
	inR, _ := io.Pipe()
	_, outW := io.Pipe()

	// Create labels for hiveproxy container.
	l1nodeLabels := libhive.NewBaseLabels(cb.hiveInstanceID, cb.hiveVersion)

	// Generate container name.
	containerName := "hive-l1-node"

	opts := libhive.ContainerOptions{Output: outW, Input: inR, Labels: l1nodeLabels, Name: containerName}

	id, err := cb.CreateContainer(ctx, l1nodeTag, opts)
	if err != nil {
		return "", err
	}

	// Now start the container.
	info, err := cb.StartContainer(ctx, id, opts)
	if err != nil {
		cb.DeleteContainer(id)
		return "", err
	}

	return info.ID, nil
}
