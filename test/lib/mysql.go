package lib

import (
	"context"
	"fmt"
	"io"
	"os"
	"testing"
	"time"

	docker "docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"
	"github.com/docker/go-connections/nat"
	log "github.com/sirupsen/logrus"

	"github.com/metal-go/metal/util/worker"
)

// MySQLWorker is response for setting up and tearing down mysql during
// integration and smoke testing. It's also responsible for running the actual
// tests.
type MySQLWorker struct {
	docker *docker.Client
	Func   worker.Func
	id     string
}

// NewMySQLWorker returns a new MySQLWorker
func NewMySQLWorker(workerFunc worker.Func) (*MySQLWorker, error) {
	client, err := docker.NewEnvClient()
	if err != nil {
		return nil, err
	}
	return &MySQLWorker{
		docker: client,
		Func:   workerFunc,
	}, nil
}

// Work does the actual work
func (m *MySQLWorker) Work(ctx context.Context, ch chan error) {
	defer m.teardown(ctx, ch)
	if err := m.setup(ctx); err != nil {
		fmt.Println(err)
		ch <- err
		return
	}
	m.Func(ctx, ch)
}

// Net returns the address:port for the MySQL container
func (m *MySQLWorker) net(ctx context.Context) (string, error) {
	var data []nat.PortBinding
	var ok bool
	for {
		containerJSON, err := m.docker.ContainerInspect(ctx, m.id)
		if err != nil {
			return "", err
		}

		data, ok = containerJSON.NetworkSettings.NetworkSettingsBase.Ports["3306/tcp"]
		if !ok {
			fmt.Println("Waiting for networking to initialize")
			fmt.Println(containerJSON.NetworkSettings.NetworkSettingsBase.Ports)
			time.Sleep(time.Second)
			continue
		} else {
			return fmt.Sprintf("%s:%s", data[0].HostIP, data[0].HostPort), nil
		}
	}
}

func (m *MySQLWorker) runContainer(ctx context.Context) error {
	reader, err := m.docker.ImagePull(ctx, "mysql", types.ImagePullOptions{})
	if err != nil {
		return err
	}
	io.Copy(os.Stdout, reader)

	container, err := m.docker.ContainerCreate(
		ctx,
		&container.Config{Image: "mysql", Env: []string{"MYSQL_ROOT_PASSWORD=password"}, OpenStdin: true, Tty: true},
		&container.HostConfig{PublishAllPorts: true},
		nil,
		"")

	if err != nil {
		return err
	}

	m.id = container.ID
	return m.docker.ContainerStart(ctx, container.ID, types.ContainerStartOptions{})
}

// Setup MySQL
func (m *MySQLWorker) setup(ctx context.Context) error {
	time.Sleep(time.Second)
	if err := m.runContainer(ctx); err != nil {
		return err
	}
	for {
		for i := 20; i > 0; i-- {
			fmt.Println("Waiting for mysql to initialize", i)
			time.Sleep(time.Second)
		}
		break
	}
	net, err := m.net(ctx)
	if err != nil {
		return err
	}
	dsn := fmt.Sprintf("root:password@tcp(%s)/", net)
	os.Setenv("METAL_REXECD_DATA_SOURCE_NAME", dsn)
	os.Setenv("METAL_REXECD_SERVER_TYPE", "mysql")
	return nil
}

// Teardown destroys the mysql container
func (m *MySQLWorker) teardown(ctx context.Context, ch chan error) {
	if m.id == "" {
		return
	}
	defer m.docker.Close()
	timeout := time.Second
	if err := m.docker.ContainerStop(ctx, m.id, &timeout); err != nil {
		ch <- err
		return
	}
}

// MySQLTestMain to be called in integration test suite
func MySQLTestMain(t *testing.M) {
	var exitCode int
	ch := make(chan error)

	// Configure a Worker that has setup and teardown
	worker, err := NewMySQLWorker(worker.Func(func(ctx context.Context, ch chan error) {
		exitCode = t.Run()
	}))

	if err != nil {
		log.Fatal(err)
	}

	worker.Work(context.Background(), ch)

	select {
	case err := <-ch:
		log.Fatal(err)
	default:
		os.Exit(exitCode)
	}
}
