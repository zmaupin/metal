package memory

import (
	"context"
	"fmt"
	"net"
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/metal-go/metal/config"
	proto_rexecd "github.com/metal-go/metal/proto/rexecd"
	"github.com/metal-go/metal/rexecd"
)

// Server is a single-instance, in-memory RexecdServer implementation.
// This implementation should not be used for production but is available for
// experimentation and demonstration.
type Server struct {
	registerHostRequestTable *registerHostRequestTable
	registerUserRequestTable *registerUserRequestTable
	commandStore             *commandStore
}

// NewServer returns a pointer to a new MemoryServer
func NewServer() *Server {
	return &Server{
		registerHostRequestTable: newRegisterHostRequestTable(),
		registerUserRequestTable: newRegisterUserRequestTable(),
		commandStore:             newCommandStore(),
	}
}

// Run runs Rexecd
func (m *Server) Run() error {
	network := fmt.Sprintf("%s:%s", config.RexecdGlobal.GetAddress(), config.RexecdGlobal.GetPort())
	lis, err := net.Listen("tcp", network)
	if err != nil {
		return fmt.Errorf("failed to listen: %v", err)
	}
	server := grpc.NewServer()
	proto_rexecd.RegisterRexecdServer(server, m)
	grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	server.Serve(lis)
	return nil
}

// Command sends a CommandRequest to a proto_rexecd.Rexecd_CommandServer and
// streams a response to the client
//
// TODO: Use real pipeline to capture stdout and stderr for each host in
//       UI
func (m *Server) Command(commandRequest *proto_rexecd.CommandRequest, server proto_rexecd.Rexecd_CommandServer) error {
	log.WithFields(log.Fields{
		"username": commandRequest.GetUsername(),
		"cmd":      commandRequest.GetCmd(),
	}).Info("command request submitted")
	t := time.Now()
	m.commandStore.addCommand(*commandRequest, t)
	wg := new(sync.WaitGroup)
	for _, hostConfig := range commandRequest.HostConfig {
		wg.Add(1)
		go m.command(hostConfig, commandRequest, wg, server, t)
	}
	wg.Wait()
	return nil
}

func (m *Server) command(hostConfig *proto_rexecd.HostConfig, commandRequest *proto_rexecd.CommandRequest, wg *sync.WaitGroup, server proto_rexecd.Rexecd_CommandServer, t time.Time) {
	defer wg.Done()
	registerHostRequest, found := m.registerHostRequestTable.get(hostConfig.GetHostId())
	if !found {
		server.Send(m.exitStatus(fmt.Sprintf("PublicKey for %s not found", hostConfig.GetHostId()), rexecd.ExitUnknown, hostConfig))
	}
	registerUserRequest, found := m.registerUserRequestTable.get(commandRequest.GetUsername())
	if !found {
		server.Send(m.exitStatus(fmt.Sprintf("PublicKey for %s not found", commandRequest.GetUsername()), rexecd.ExitUnknown, hostConfig))
	}
	clientConfig, err := rexecd.BuildClientConfig(registerUserRequest, registerHostRequest)
	if err != nil {
		server.Send(m.exitStatus(err.Error(), rexecd.ExitUnknown, hostConfig))
		return
	}
	session := rexecd.NewSSHSessionBuilder(hostConfig.GetAddress(), clientConfig)
	session.AddEnv(commandRequest.GetEnv())
	port := hostConfig.GetPort()
	if port != "" {
		session.AddPort(port)
	}
	sshSession, err := session.Build()
	if err != nil {
		server.Send(m.exitStatus(err.Error(), rexecd.ExitUnknown, hostConfig))
		return
	}

	runner := rexecd.NewExecRunner(commandRequest.GetCmd(), sshSession)
	hostID := hostConfig.GetHostId()
	execData := m.commandStore.hostCommandStore.data[hostID][t]
	runner.AddStdoutPipeline(execData.addStdoutLine)
	runner.AddStdoutPipeline(execData.addStderrLine)

	status, err := runner.Run()
	execData.setDone()

	if err != nil {
		server.Send(m.exitStatus(err.Error(), rexecd.ExitUnknown, hostConfig))
		return
	}
	server.Send(m.exitStatus("", status, hostConfig))
}

func (m *Server) exitStatus(msg string, status int32, hostConfig *proto_rexecd.HostConfig) *proto_rexecd.CommandResponse {
	return &proto_rexecd.CommandResponse{
		HostConfig: hostConfig,
		Exit: &proto_rexecd.Exit{
			Msg:    msg,
			Status: status,
		},
	}
}

// RegisterHost registers a host's public key in memory
func (m *Server) RegisterHost(ctx context.Context, registerHostRequest *proto_rexecd.RegisterHostRequest) (*proto_rexecd.RegisterHostResponse, error) {
	m.registerHostRequestTable.set(*registerHostRequest)
	return &proto_rexecd.RegisterHostResponse{}, nil
}

// RegisterUser registers a rexecd user. It generates and stores a private and
// public keypair if KeyPair is not passed in the
// proto_rexecd.RegisterUserRequest
func (m *Server) RegisterUser(ctx context.Context, registerUserRequest *proto_rexecd.RegisterUserRequest) (*proto_rexecd.RegisterUserResponse, error) {
	_, found := m.registerUserRequestTable.get(registerUserRequest.GetUsername())
	if found {
		return &proto_rexecd.RegisterUserResponse{}, fmt.Errorf("user %s already exists", registerUserRequest.GetUsername())
	}
	privateKey := registerUserRequest.GetPrivateKey()
	if privateKey == nil {
		key, err := rexecd.GeneratePrivateKey(4096)
		if err != nil {
			return &proto_rexecd.RegisterUserResponse{}, fmt.Errorf("failed to generate private key")
		}
		publicKey, err := rexecd.GeneratePublicKey(&key.PublicKey)
		if err != nil {
			return &proto_rexecd.RegisterUserResponse{}, fmt.Errorf("failed to generate public key")
		}
		privateKey := rexecd.EncodePrivateKeyToPEM(key)
		registerUserRequest.PublicKey = publicKey
		registerUserRequest.PrivateKey = privateKey
	}
	m.registerUserRequestTable.set(*registerUserRequest)
	return &proto_rexecd.RegisterUserResponse{}, nil
}
