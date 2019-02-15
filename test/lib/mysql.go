package lib

import (
	"bytes"
	"context"
	"net"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/metal-go/metal/util/worker"
	log "github.com/sirupsen/logrus"
)

// MySQLWorker is response for setting up and tearing down mysql during
// integration and smoke testing. It's also responsible for running the actual
// tests.
type MySQLWorker struct {
	Func worker.Func
	id   string
}

// NewMySQLWorker returns a new MySQLWorker
func NewMySQLWorker(workerFunc worker.Func) *MySQLWorker {
	return &MySQLWorker{Func: workerFunc}
}

// Setup MySQL
func (m *MySQLWorker) Setup(ctx context.Context, ch chan error) {
	buf := &bytes.Buffer{}
	args := strings.Split("run --interactive --publish 3306:3306 --env MYSQL_ROOT_PASSWORD=password --detach mysql", " ")
	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		ch <- err
		return
	}
	m.id = strings.TrimSpace(buf.String())

	for {
		conn, err := net.Dial("tcp", "127.0.0.1:3306")
		if err != nil {
			continue
		}
		conn.Close()
		log.Info("Connection established. Waiting for server initialization to complete...")
		for i := 20; i > 0; i-- {
			log.Info(i)
			time.Sleep(time.Second)
		}
		break
	}
	os.Setenv("METAL_REXECD_DATA_SOURCE_NAME", "root:password@tcp(127.0.0.1:3306)/")
	os.Setenv("METAL_REXECD_SERVER_TYPE", "mysql")
}

// Work does the actual work
func (m *MySQLWorker) Work(ctx context.Context, ch chan error) {
	defer m.Teardown(ctx, ch)
	m.Setup(ctx, ch)
	m.Func(ctx, ch)
}

// Teardown destroys the mysql container
func (m *MySQLWorker) Teardown(ctx context.Context, ch chan error) {
	cmd := exec.CommandContext(ctx, "docker", "rm", "-f", m.id)
	if err := cmd.Run(); err != nil {
		ch <- err
	}
}

// MySQLTestMain to be called in integration test suite
func MySQLTestMain(t *testing.M) {
	var exitCode int
	ch := make(chan error)

	// Configure a Worker that has setup and teardown
	NewMySQLWorker(worker.Func(func(ctx context.Context, ch chan error) {
		exitCode = t.Run()
		close(ch)
	})).Work(context.Background(), ch)

	// If an error comes along, exit
	if err := <-ch; err != nil {
		log.Fatal(err)
	}

	os.Exit(exitCode)
}
