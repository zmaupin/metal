package mysql

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql" // driver
	log "github.com/sirupsen/logrus"

	"github.com/metal-go/metal/config"
	"github.com/metal-go/metal/db/mysql/migration"
	"github.com/metal-go/metal/rexecd"
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
func (m *Server) Command(ctx context.Context, c *proto_rexecd.CommandRequest, s proto_rexecd.Rexecd_CommandServer) error {
	log.WithFields(log.Fields{
		"username": c.GetUsername(),
		"cmd":      c.GetCmd(),
	}).Info("command request submitted")

	errorChan := make(chan error)
	select {
	case <-m.commandScheduler(ctx, c, s, errorChan):
		return nil
	case err := <-errorChan:
		return err
	case <-ctx.Done():
		return errors.New("context timeout exceeded")
	}
}

func (m *Server) commandScheduler(ctx context.Context, c *proto_rexecd.CommandRequest, s proto_rexecd.Rexced_CommandServer, errorChan chan error) {
	doneChan := make(chan bool)
	go func() {
		defer close(errorChan)
		defer close(doneChan)

		t := time.Now()
		wg := &sync.WaitGroup{}
		for _, hostConnect := range c.GetHostConnect() {
			wg.Add(1)
			go m.command(ctx, hostConnect, c, s, wg, t, errorChan)
		}
		wg.Wait()
		doneChan <- true
	}()
	return doneChan
}

func (m *Server) command(ctx context.Context, hostConnect *proto_rexecd.HostConnect,
	c *proto_rexecd.CommandRequest, server proto_rexecd.Rexecd_CommandServer,
	wg *sync.WaitGroup, t time.Time, ch chan error) {

	host := NewHost(m.db)
	host.Read(ctx, hostConnect.GetFqdn())
	sshConfig, err := rexecd.NewSSHClientConfig(c.GetUsername(), c.GetPrivateKey(), host.PublicKey, host.KeyType)
	if err != nil {
		m.exitStatus(host.ID, err, 1)
		wg.Done()
		return
	}

	sshSession, err := rexecd.NewSSHSessionBuilder(hostConnect.GetFqdn(), sshConfig,
		rexecd.WithSSHSessionBuilderPort(host.Port),
		rexecd.WithSSHSessionBuilderEnv(c.GetEnv())).Build()

	defer sshSession.Close()

	if err != nil {
		m.exitStatus(host.ID, err, 1)
		wg.Done()
		return
	}

	rexecd.NewExecRunner(c.GetCmd(), sshSession)

}

func (m *Server) exitStatus(id int64, err error, exitCode int64) *proto_rexecd.CommandResponse {
	return &proto_rexecd.CommandResponse{
		Id:        id,
		ErrorMsg:  error.Error(),
		ExitCode:  exitCode,
		Timestamp: time.Now().Unix(),
	}
}

// RegisterHost registers a Host
func (m *Server) RegisterHost(ctx context.Context, r *proto_rexecd.RegisterHostRequest) (
	*proto_rexecd.RegisterHostResponse, error,
) {
	host := NewHost(m.db, WithHostFQDN(r.GetFqdn()), WithHostPort(r.GetPort()),
		WithHostPublicKey(r.GetPublicKey()), WithHostKeyType(r.GetKeyType()))
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
