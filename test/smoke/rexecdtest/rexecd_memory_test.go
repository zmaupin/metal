package rexecdtest

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
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

var sshdHosts = []string{"docker_host01_1", "docker_host02_1"}

func Init(ctx context.Context) error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	parts := strings.Split(dir, string(os.PathSeparator))
	projectRoot := strings.Join(parts[:len(parts)-3], string(os.PathSeparator))

	devEnv := exec.Command("make", "rexecd-memory-server-restart")
	devEnv.Stdout = os.Stdout
	devEnv.Stderr = os.Stderr
	devEnv.Dir = projectRoot
	if err = devEnv.Run(); err != nil {
		return err
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
	var count uint
	for {
		healthCheckResponse, err := healthClient.Check(ctx, &health.HealthCheckRequest{})
		if healthCheckResponse.GetStatus() == health.HealthCheckResponse_SERVING {
			fmt.Println("health check response:", healthCheckResponse)
			break
		}
		count++
		fmt.Printf("waiting for rexecd to start: %d: %s\n", count, err.Error())
		time.Sleep(1 * time.Second)
	}
	commandRequest := &proto_rexecd.CommandRequest{Cmd: "/bin/true", Username: "dev"}
	for _, name := range sshdHosts {
		container, found, err := findContainer(ctx, client, name)
		if err != nil {
			t.Fatal(err)
		}
		if !found {
			t.Fatal(fmt.Errorf("could not find container %s", name))
		}
		containerJSON, err := client.ContainerInspect(ctx, container.ID)
		if err != nil {
			t.Fatal(err)
		}
		var ipaddress string
		for _, data := range containerJSON.NetworkSettings.Networks {
			ipaddress = data.IPAddress
			break
		}
		if err != nil {
			t.Fatal(err)
		}
		publicHostRSA, err := getPublicHostRSA(ctx, client, name)
		if err != nil {
			t.Fatal(err)
		}
		registerHostRequest := &proto_rexecd.RegisterHostRequest{
			HostId:    name,
			PublicKey: publicHostRSA,
			KeyType:   "rsa-sha2-512",
		}
		_, err = rexecdClient.RegisterHost(ctx, registerHostRequest)
		if err != nil {
			t.Fatal(err)
		}

		commandRequest.HostConfig = append(commandRequest.GetHostConfig(), &proto_rexecd.HostConfig{
			HostId:  name,
			Address: ipaddress,
			Port:    "22",
		})
	}
	privatePath, err := user.Expand("~/.ssh/id_rsa")
	if err != nil {
		t.Fatal(err)
	}
	publicPath, err := user.Expand("~/.ssh/id_rsa.pub")
	if err != nil {
		t.Fatal(err)
	}
	privateKey, err := ioutil.ReadFile(privatePath)
	if err != nil {
		t.Fatal(err)
	}
	publicKey, err := ioutil.ReadFile(publicPath)
	if err != nil {
		t.Fatal(err)
	}
	registerUserRequest := &proto_rexecd.RegisterUserRequest{
		Username:   "dev",
		PrivateKey: privateKey,
		PublicKey:  publicKey,
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
			t.Errorf("unexpected failure: %s", err.Error())
		}
		exit := commandResponse.GetExit()
		if exit.GetStatus() != 0 {
			t.Errorf("expected 0 exit status, got %d", exit.GetStatus())
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
