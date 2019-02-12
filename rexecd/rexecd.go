package rexecd

import (
	"bufio"
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	proto_rexecd "github.com/metal-go/metal/proto/rexecd"
	"golang.org/x/crypto/ssh"
)

// Server implements proto_rexecd.RexecdServer
type Server interface {
	proto_rexecd.RexecdServer
	Run() error
}

// SSHEnv is a map of environment variables and their corresponding values
type SSHEnv map[string]string

// SSHSessionBuilder is holds the information required to build a pointer to a
// sss.Session
type SSHSessionBuilder struct {
	fqdn         string
	port         string
	env          SSHEnv
	clientConfig *ssh.ClientConfig
}

// SSHSessionBuilderOpt is an option for an SSHSessionBuilder
type SSHSessionBuilderOpt func(*SSHSessionBuilder)

// WithSSHSessionBuilderEnv adds an SSHEnv to the SSHSessionBuilder
func WithSSHSessionBuilderEnv(env SSHEnv) SSHSessionBuilderOpt {
	return func(s *SSHSessionBuilder) {
		s.env = env
	}
}

// WithSSHSessionBuilderPort adds a port to the SSHSessionBuilder
func WithSSHSessionBuilderPort(port string) SSHSessionBuilderOpt {
	return func(s *SSHSessionBuilder) {
		s.port = port
	}
}

// NewSSHSessionBuilder returns a pointer to an ssh.Session
func NewSSHSessionBuilder(fqdn string, sshConfig *ssh.ClientConfig, opts ...SSHSessionBuilderOpt) *SSHSessionBuilder {
	builder := &SSHSessionBuilder{fqdn: fqdn, port: "22", clientConfig: sshConfig}
	for _, fn := range opts {
		fn(builder)
	}
	return builder
}

// Build returns a pointer to an ssh.Session
func (s SSHSessionBuilder) Build() (*ssh.Session, error) {
	network := fmt.Sprintf("%s:%s", s.fqdn, s.port)
	client, err := ssh.Dial("tcp", network, s.clientConfig)
	if err != nil {
		return &ssh.Session{}, err
	}
	sshSession, err := client.NewSession()
	if err != nil {
		return sshSession, err
	}
	if s.env != nil {
		for key, val := range s.env {
			if err = sshSession.Setenv(key, val); err != nil {
				return sshSession, err
			}
		}
	}
	return sshSession, nil
}

// ExecRunner runs a command on a remote host
type ExecRunner struct {
	cmd           string
	sshSession    *ssh.Session
	stdoutHandler BytesLineHandler
	stderrHandler BytesLineHandler
}

// ExecRunnerOpt is an option for an ExecRunner
type ExecRunnerOpt func(*ExecRunner)

// NewExecRunner returns a pointer to an ExecRunner
func NewExecRunner(cmd string, sshSession *ssh.Session, stdoutHandler BytesLineHandler, stderrHandler BytesLineHandler, opts ...ExecRunnerOpt) *ExecRunner {
	runner := &ExecRunner{
		cmd:           cmd,
		sshSession:    sshSession,
		stdoutHandler: stdoutHandler,
		stderrHandler: stderrHandler,
	}

	for _, fn := range opts {
		fn(runner)
	}

	return runner
}

// Run executes the command, feeding stdout into the stdout pipeline and stderr
// into the stderr pipeline
func (e *ExecRunner) Run(ctx context.Context) (statusCode int64, err error) {
	defer e.sshSession.Close()
	// Setup stdout and stderr readers and scanners
	outReader, err := e.sshSession.StdoutPipe()
	if err != nil {
		return 1, err
	}
	errReader, err := e.sshSession.StderrPipe()
	if err != nil {
		return 1, err
	}
	outScanner := bufio.NewScanner(outReader)
	errScanner := bufio.NewScanner(errReader)

	// Feed bytes of lines to the given pipeline
	feeder := func(scanner *bufio.Scanner, handler BytesLineHandler) {
		for scanner.Scan() {
			line := append(scanner.Bytes(), byte('\n'))
			handler.Handle(ctx, line)
		}
	}
	go func() { feeder(outScanner, e.stdoutHandler) }()
	go func() { feeder(errScanner, e.stderrHandler) }()

	// Run it
	err = e.sshSession.Run(e.cmd)

	// Check for errors
	if err == nil {
		return int64(0), nil
	}
	exitErr, ok := err.(*ssh.ExitError)
	if ok {
		return int64(exitErr.Waitmsg.ExitStatus()), nil
	}
	return 1, err
}

// BuildAuthMethod returns an ssh.AuthMethod from the given private key presented
// as a byte array
func BuildAuthMethod(userPrivateKey []byte) (ssh.AuthMethod, error) {
	var sshAuthMethod ssh.AuthMethod
	signer, err := ssh.ParsePrivateKey(userPrivateKey)
	if err != nil {
		return sshAuthMethod, err
	}
	sshAuthMethod = ssh.PublicKeys(signer)
	return sshAuthMethod, nil
}

// NewSSHClientConfig builds an ssh.ClientConfig based on the given
// proto_rexecd.RegisterUserRequest and proto_rexecd.RegisterHostRequest. This
// will enforce FixedHostKey checking.
func NewSSHClientConfig(username string, privateUserKey, publicHostKey []byte, hostKeyType proto_rexecd.KeyType) (*ssh.ClientConfig, error) {
	key, _, _, _, err := ssh.ParseAuthorizedKey(publicHostKey)
	if err != nil {
		return &ssh.ClientConfig{}, err
	}
	hostKeyCallback := ssh.FixedHostKey(key)
	authMethod, err := BuildAuthMethod(privateUserKey)
	if err != nil {
		return &ssh.ClientConfig{}, err
	}
	keyType := proto_rexecd.KeyType_name[int32(hostKeyType)]
	return &ssh.ClientConfig{
		User:              username,
		Auth:              []ssh.AuthMethod{authMethod},
		HostKeyAlgorithms: []string{keyType},
		HostKeyCallback:   hostKeyCallback,
	}, nil
}

// GeneratePrivateKey creates a RSA Private Key of specified byte size
func GeneratePrivateKey(bitSize int) (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, err
	}
	err = privateKey.Validate()
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

// EncodePrivateKeyToPEM encodes Private Key from RSA to PEM format
func EncodePrivateKeyToPEM(privateKey *rsa.PrivateKey) []byte {
	privDER := x509.MarshalPKCS1PrivateKey(privateKey)
	privBlock := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDER,
	}
	privatePEM := pem.EncodeToMemory(&privBlock)
	return privatePEM
}

// GeneratePublicKey take a rsa.PublicKey and return bytes suitable for writing
// to .pub file returns in the format "ssh-rsa ..."
func GeneratePublicKey(publicKey *rsa.PublicKey) ([]byte, error) {
	publicRsaKey, err := ssh.NewPublicKey(publicKey)
	if err != nil {
		return nil, err
	}
	pubKeyBytes := ssh.MarshalAuthorizedKey(publicRsaKey)
	return pubKeyBytes, nil
}
