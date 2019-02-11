package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"

	_ "github.com/go-sql-driver/mysql" // driver

	"github.com/metal-go/metal/config"
	"github.com/metal-go/metal/db/mysql/migration"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	proto_rexecd "github.com/metal-go/metal/proto/rexecd"
)

// Server driver
type Server struct {
	db *sql.DB
}

// New pointer to MySQL
func New() *Server {
	return &Server{}
}

// Command executes a command
func (m *Server) Command(c *proto_rexecd.CommandRequest, s proto_rexecd.Rexecd_CommandServer) error {
	return nil
}

// RegisterHost registers a Host
func (m *Server) RegisterHost(ctx context.Context, r *proto_rexecd.RegisterHostRequest) (
	*proto_rexecd.RegisterHostResponse, error,
) {
	host := NewHost(m.db, r.GetFqdn(), r.GetPort(), r.GetPublicKey(), r.GetKeyType())
	id, err := host.Create(ctx)
	return &proto_rexecd.RegisterHostResponse{Id: id}, err
}

// RegisterUser registers a User
func (m *Server) RegisterUser(ctx context.Context, r *proto_rexecd.RegisterUserRequest) (
	*proto_rexecd.RegisterUserResponse, error,
) {
	u := NewUser(m.db, r.GetUsername(), WithUserFirstName(r.GetFirstName()),
		WithUserLastName(r.GetLastName()), WithUserAdmin(r.GetAdmin()))
	err := u.Create(ctx)
	return &proto_rexecd.RegisterUserResponse{}, err
}

// Run starts the server
func (m *Server) Run() error {
	migrate := migration.New()
	if err := migrate.Run(); err != nil {
		return err
	}
	dsn := config.RexecdGlobal.DataSourceName + "rexecd"
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal(err)
	}
	m.db = db
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
