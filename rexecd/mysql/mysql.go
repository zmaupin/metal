package mysql

import (
	"context"
	"fmt"
	"net"

	_ "github.com/go-sql-driver/mysql" // driver

	"github.com/metal-go/metal/config"
	"github.com/metal-go/metal/db/mysql/migration"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	proto_rexecd "github.com/metal-go/metal/proto/rexecd"
)

// MySQL driver
type MySQL struct{}

// New pointer to MySQL
func New() *MySQL {
	return &MySQL{}
}

// Command executes a command
func (m *MySQL) Command(c *proto_rexecd.CommandRequest, s proto_rexecd.Rexecd_CommandServer) error {
	return nil
}

// RegisterHost registers a Host
func (m *MySQL) RegisterHost(ctx context.Context, r *proto_rexecd.RegisterHostRequest) (
	*proto_rexecd.RegisterHostResponse, error,
) {
	return &proto_rexecd.RegisterHostResponse{}, nil
}

// RegisterUser registers a User
func (m *MySQL) RegisterUser(ctx context.Context, r *proto_rexecd.RegisterUserRequest) (
	*proto_rexecd.RegisterUserResponse, error,
) {
	return &proto_rexecd.RegisterUserResponse{}, nil
}

// Run starts the server
func (m *MySQL) Run() error {
	migrate := migration.New()
	if err := migrate.Run(); err != nil {
		return err
	}
	network := fmt.Sprintf("%s:%s", config.RexecdGlobal.Address, config.RexecdGlobal.Port)
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
