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
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"

	"github.com/metal-go/metal/config"
	"github.com/metal-go/metal/db/mysql/migration"
	"github.com/metal-go/metal/rexecd"

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
	log.WithFields(log.Fields{
		"username":   c.GetUsername(),
		"cmd":        c.GetCmd(),
		"host_count": len(c.GetHostConnect()),
	}).Info("command request submitted")

	// TODO: Make this a context with a configurable timeout
	ctx := context.Background()

	select {
	case <-m.commandScheduler(ctx, c, s):
		return nil
	case <-ctx.Done():
		return errors.New("context timeout exceeded")
	}
}

func (m *Server) commandScheduler(ctx context.Context, c *proto_rexecd.CommandRequest, s proto_rexecd.Rexecd_CommandServer) chan bool {
	doneChan := make(chan bool)
	go func() {
		t := time.Now()
		wg := &sync.WaitGroup{}
		for _, hostConnect := range c.GetHostConnect() {
			wg.Add(1)
			go m.command(ctx, hostConnect, c, s, wg, t)
		}
		wg.Wait()
		doneChan <- true
		close(doneChan)
	}()
	return doneChan
}

func (m *Server) command(ctx context.Context, hostConnect *proto_rexecd.HostConnect,
	c *proto_rexecd.CommandRequest, s proto_rexecd.Rexecd_CommandServer,
	wg *sync.WaitGroup, t time.Time) {

	// Load Host
	host := NewHost(m.db)
	if err := host.Read(ctx, hostConnect.GetFqdn()); err != nil {
		s.Send(m.exitStatus(ctx, nil, host, err, 1, wg, c))
		return
	}

	// Load Command
	command := NewCommand(m.db)
	if err := command.Create(ctx, c.GetCmd(), c.GetUsername(), hostConnect.GetFqdn()); err != nil {
		s.Send(m.exitStatus(ctx, command, host, err, 1, wg, c))
		return
	}

	// Build sshConfig
	sshConfig, err := rexecd.NewSSHClientConfig(c.GetUsername(), c.GetPrivateKey(), host.PublicKey, host.KeyType)
	if err != nil {
		s.Send(m.exitStatus(ctx, command, host, err, 1, wg, c))
		return
	}
	fmt.Println("I MADE IT HERE")

	// Build SSH Session
	sshSession, err := rexecd.NewSSHSessionBuilder(hostConnect.GetFqdn(), sshConfig,
		rexecd.WithSSHSessionBuilderPort(host.Port),
		rexecd.WithSSHSessionBuilderEnv(c.GetEnv())).Build()

	if err != nil {
		s.Send(m.exitStatus(ctx, command, host, err, 1, wg, c))
		return
	}

	defer sshSession.Close()

	// Build ExecRunner
	execRunner := rexecd.NewExecRunner(c.GetCmd(), sshSession, NewBytesLineHandler(command, MySQLStdout),
		NewBytesLineHandler(command, MySQLStderr))

	// Run it
	exitCode, err := execRunner.Run(ctx)
	s.Send(m.exitStatus(ctx, command, host, err, exitCode, wg, c))
}

// Update the command table with the appropriate exit code and return a
// CommandResponse to the client
func (m *Server) exitStatus(ctx context.Context, command *Command, host *Host, err error, exitCode int64, wg *sync.WaitGroup, commandRequest *proto_rexecd.CommandRequest) *proto_rexecd.CommandResponse {
	var id int64
	if command == nil {
		id = 0
	} else {
		id = command.ID
	}

	var hostFQDN string
	if host == nil {
		hostFQDN = ""
	} else {
		hostFQDN = host.FQDN
	}

	logFields := log.Fields{
		"cmd_id":    id,
		"cmd":       commandRequest.GetCmd(),
		"username":  commandRequest.GetUsername(),
		"host":      hostFQDN,
		"exit_code": exitCode,
	}

	if exitCode == 0 {
		log.WithFields(logFields).Info("command execution succeeded")
	} else {
		log.WithFields(logFields).Info("command execution failed")
	}

	if command != nil {
		if e := command.SetExitCode(ctx, exitCode); err != nil {
			log.WithFields(logFields).Error(e)
		}
	}

	wg.Done()

	return &proto_rexecd.CommandResponse{
		Id:        id,
		ErrorMsg:  err.Error(),
		ExitCode:  exitCode,
		Timestamp: time.Now().Unix(),
	}
}

// RegisterHost registers a Host
func (m *Server) RegisterHost(ctx context.Context, r *proto_rexecd.RegisterHostRequest) (
	*proto_rexecd.RegisterHostResponse, error,
) {
	host := NewHost(m.db)
	id, err := host.Create(ctx, r.GetFqdn(), WithHostPort(r.GetPort()),
		WithHostPublicKey(r.GetPublicKey()), WithHostKeyType(r.GetKeyType()))
	return &proto_rexecd.RegisterHostResponse{Id: id}, err
}

// RegisterUser registers a User
func (m *Server) RegisterUser(ctx context.Context, r *proto_rexecd.RegisterUserRequest) (
	*proto_rexecd.RegisterUserResponse, error,
) {
	u := NewUser(m.db)
	err := u.Create(ctx, r.GetUsername(), WithUserFirstName(r.GetFirstName()),
		WithUserLastName(r.GetLastName()), WithUserAdmin(r.GetAdmin()))

	return &proto_rexecd.RegisterUserResponse{}, err
}

// Run starts the server
func (m *Server) Run(done chan bool) error {
	migrate := migration.New()
	if err := migrate.Run(); err != nil {
		return err
	}
	dsn := config.RexecdGlobal.DataSourceName + "rexecd"
	db, err := sql.Open("mysql", dsn)
	defer db.Close()
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
	run := func() chan error {
		ch := make(chan error)
		go func() {
			ch <- server.Serve(lis)
		}()
		return ch
	}
	select {
	case err := <-run():
		return err
	case <-done:
		return nil
	}
}
