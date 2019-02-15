package rexecd

import (
	"fmt"

	"golang.org/x/crypto/ssh"
)

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
