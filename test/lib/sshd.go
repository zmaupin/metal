package lib

import (
	"context"
	"fmt"
	"io/ioutil"
	"time"

	docker "docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"docker.io/go-docker/api/types/container"

	"github.com/docker/go-connections/nat"
	"github.com/metal-go/metal/util/user"
	"github.com/metal-go/metal/util/worker"
)

// SSHDWorkerOpt is an option for an SSHDWorker
type SSHDWorkerOpt func(*SSHDWorker)

// SSHDWorker runs the SSHD container for integration tests
type SSHDWorker struct {
	Func          worker.Func
	PublicHostKey []byte
	docker        *docker.Client
	id            string
}

// SSHDWorkerWithFunc adds a worker func
func SSHDWorkerWithFunc(workerFn worker.Func) SSHDWorkerOpt {
	return func(s *SSHDWorker) {
		s.Func = workerFn
	}
}

// NewSSHDWorker returns a new SSHDWorker
func NewSSHDWorker(opts ...SSHDWorkerOpt) (*SSHDWorker, error) {
	docker, err := docker.NewEnvClient()
	if err != nil {
		return nil, err
	}
	sshd := &SSHDWorker{docker: docker}
	for _, fn := range opts {
		fn(sshd)
	}
	if sshd.Func == nil {
		sshd.Func = func(ctx context.Context, ch chan error) {}
	}
	return sshd, nil
}

// Work sets up an SSHD container and does work
func (s *SSHDWorker) Work(ctx context.Context, ch chan error) {
	if err := s.Setup(ctx); err != nil {
		ch <- err
		return
	}
	s.Func(ctx, ch)
}

// Teardown destroys the SSHD container
func (s *SSHDWorker) Teardown(ctx context.Context) error {
	return s.docker.ContainerRemove(ctx, s.id, types.ContainerRemoveOptions{Force: true})
}

// Setup launches an SSHD container
func (s *SSHDWorker) Setup(ctx context.Context) error {
	if err := s.pullImage(ctx); err != nil {
		return err
	}

	if err := s.runContainer(ctx); err != nil {
		return err
	}

	return s.setPublicHostKey(ctx)
}

func (s *SSHDWorker) pullImage(ctx context.Context) error {
	_, err := s.docker.ImagePull(ctx, "panubo/sshd", types.ImagePullOptions{})
	return err
}

func (s *SSHDWorker) runContainer(ctx context.Context) error {
	p, err := user.Expand("~/.ssh/id_rsa.pub")
	if err != nil {
		return err
	}

	volume := fmt.Sprintf("%s:/root/.ssh/authorized_keys", p)

	response, err := s.docker.ContainerCreate(ctx,
		&container.Config{
			Image:     "panubo/sshd",
			Tty:       true,
			OpenStdin: true,
		},
		&container.HostConfig{
			Binds: []string{volume},
			PortBindings: nat.PortMap(map[nat.Port][]nat.PortBinding{
				"22/tcp": []nat.PortBinding{nat.PortBinding{HostIP: "0.0.0.0", HostPort: "2222"}},
			}),
		}, nil, "")

	s.id = response.ID

	if err != nil {
		return err
	}

	if err = s.docker.ContainerStart(ctx, response.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}
	for {
		containerJSON, err := s.docker.ContainerInspect(ctx, response.ID)
		if err != nil {
			return err
		}
		if containerJSON.ContainerJSONBase.State.Status == "running" {
			time.Sleep(time.Second)
			return nil
		}
	}
}

func (s *SSHDWorker) setPublicHostKey(ctx context.Context) error {
	containerJSON, err := s.docker.ContainerInspect(ctx, s.id)
	if err != nil {
		return err
	}
	execConfig := types.ExecConfig{
		User:         "root",
		AttachStdout: true,
		AttachStderr: true,
		Cmd:          []string{"cat", "/etc/ssh/keys/ssh_host_rsa_key.pub"},
		Tty:          true,
	}
	execResponse, err := s.docker.ContainerExecCreate(ctx, containerJSON.ContainerJSONBase.Name, execConfig)
	if err != nil {
		return err
	}
	hijackedResponse, err := s.docker.ContainerExecAttach(ctx, execResponse.ID, execConfig)
	defer hijackedResponse.Close()
	if err != nil {
		return err
	}
	if err = s.docker.ContainerExecStart(ctx, execResponse.ID, types.ExecStartCheck{}); err != nil {
		return err
	}

	bytes, err := ioutil.ReadAll(hijackedResponse.Reader)
	if err != nil {
		return err
	}
	s.PublicHostKey = bytes
	return nil
}
