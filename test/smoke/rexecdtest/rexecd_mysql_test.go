package rexecdtest

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	docker "docker.io/go-docker"
	"docker.io/go-docker/api/types"
	"google.golang.org/grpc"
	health "google.golang.org/grpc/health/grpc_health_v1"

	proto_rexecd "github.com/metal-go/metal/proto/rexecd"
	"github.com/metal-go/metal/util/user"
)

var sshdHost = "docker_host01_1"
var projectRoot string

func Init(ctx context.Context) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	parts := strings.Split(dir, string(os.PathSeparator))
	projectRoot := strings.Join(parts[:len(parts)-3], string(os.PathSeparator))
	devEnv := exec.Command("docker-compose", "--file", filepath.Join(projectRoot, "docker", "rexecd-mysql-server.yml"), "down")
	devEnv.Dir = projectRoot
	devEnv.Run()
	devEnv = exec.Command("docker-compose", "--file", filepath.Join(projectRoot, "docker", "rexecd-mysql-server.yml"), "up", "--detach")
	devEnv.Dir = projectRoot
	if err = devEnv.Run(); err != nil {
		return err
	}

	for i := 20; i > 0; i-- {
		fmt.Println("Waiting for MySQL to intialize ", i)
		time.Sleep(time.Second)
	}
	return nil
}

func getPublicHostRSA(ctx context.Context, client *docker.Client, containerName string) ([]byte, error) {
	config := types.ExecConfig{
		AttachStdout: true,
		Cmd:          []string{"cat", "/etc/ssh/keys/ssh_host_rsa_key.pub"},
		Tty:          true,
	}
	idResponse, err := client.ContainerExecCreate(ctx, containerName, config)
	if err != nil {
		return []byte{}, err
	}
	hijackedResponse, err := client.ContainerExecAttach(ctx, idResponse.ID, config)
	defer hijackedResponse.Close()
	if err != nil {
		return []byte{}, err
	}
	err = client.ContainerExecStart(ctx, idResponse.ID, types.ExecStartCheck{})
	if err != nil {
		return []byte{}, err
	}
	bytes, err := ioutil.ReadAll(hijackedResponse.Reader)
	return bytes, err
}

func findContainer(ctx context.Context, client *docker.Client, name string) (types.Container, bool, error) {
	containers, err := client.ContainerList(ctx, types.ContainerListOptions{})
	if err != nil {
		return types.Container{}, false, err
	}
	for _, container := range containers {
		for _, n := range container.Names {
			if n == fmt.Sprintf("/%s", name) {
				return container, true, nil
			}
		}
	}
	return types.Container{}, false, nil
}

func getSSHPort(containerJSON types.ContainerJSON) string {
	portMap := containerJSON.NetworkSettings.NetworkSettingsBase.Ports
	return portMap["22/tcp"][0].HostPort
}

func TestRexecd(t *testing.T) {
	ctx := context.Background()
	client, err := docker.NewEnvClient()
	if err != nil {
		t.Fatal(err)
	}
	serviceNet := "0.0.0.0:9000"
	conn, err := grpc.Dial(serviceNet, grpc.WithInsecure())
	if err != nil {
		t.Fatalf("could not create grpcClient at %s", serviceNet)
	}
	defer conn.Close()
	rexecdClient := proto_rexecd.NewRexecdClient(conn)
	healthClient := health.NewHealthClient(conn)
	for count := 0; ; count++ {
		healthCheckResponse, _ := healthClient.Check(ctx, &health.HealthCheckRequest{})
		if healthCheckResponse.GetStatus() == health.HealthCheckResponse_SERVING {
			fmt.Println("health check response:", healthCheckResponse)
			break
		}
		fmt.Printf("waiting for rexecd to start: %d\n", count)
		time.Sleep(1 * time.Second)
	}
	commandRequest := &proto_rexecd.CommandRequest{Cmd: "/bin/true", Username: "dev"}
	container, found, err := findContainer(ctx, client, sshdHost)
	if err != nil {
		t.Fatal(err)
	}
	if !found {
		t.Fatal(fmt.Errorf("could not find container %s", sshdHost))
	}
	containerJSON, err := client.ContainerInspect(ctx, container.ID)
	if err != nil {
		t.Fatal(err)
	}
	port := containerJSON.NetworkSettings.NetworkSettingsBase.Ports["22/tcp"][0].HostPort
	if err != nil {
		t.Fatal(err)
	}
	publicHostRSA, err := getPublicHostRSA(ctx, client, sshdHost)
	if err != nil {
		t.Fatal(err)
	}
	registerHostRequest := &proto_rexecd.RegisterHostRequest{
		Fqdn:      "127.0.0.1",
		Port:      port,
		PublicKey: publicHostRSA,
		KeyType:   proto_rexecd.KeyType_ssh_rsa,
	}
	_, err = rexecdClient.RegisterHost(ctx, registerHostRequest)
	if err != nil {
		t.Fatal(err)
	}

	commandRequest.HostConnect = append(commandRequest.GetHostConnect(), &proto_rexecd.HostConnect{
		Fqdn: "127.0.0.1",
		Port: port,
	})

	privatePath, err := user.Expand("~/.ssh/id_rsa")
	if err != nil {
		t.Fatal(err)
	}
	privateKey, err := ioutil.ReadFile(privatePath)
	if err != nil {
		t.Fatal(err)
	}
	commandRequest.PrivateKey = privateKey
	registerUserRequest := &proto_rexecd.RegisterUserRequest{
		Username: "dev",
	}
	_, err = rexecdClient.RegisterUser(ctx, registerUserRequest)
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}
	stream, err := rexecdClient.Command(ctx, commandRequest)
	if err != nil {
		t.Fatalf("unexpected failure: %s", err.Error())
	}
	for {
		commandResponse, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%s, %v", err.Error(), commandRequest)
		}
		fmt.Println(commandResponse)
		if err != nil {
			t.Errorf("unexpected failure: %s", err.Error())
		}
		exit := commandResponse.GetExitCode()
		if exit != 0 {
			t.Errorf("expected 0 exit status, got %d", exit)
		}
	}
}

func TestMain(t *testing.M) {
	ctx, cancel := context.WithCancel(context.Background())
	Init(ctx)
	exitStatus := t.Run()
	cancel()
	os.Exit(exitStatus)
}
