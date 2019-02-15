package rexecd

import (
	"bufio"
	"context"

	"golang.org/x/crypto/ssh"
)

// SSHEnv is a map of environment variables and their corresponding values
type SSHEnv map[string]string

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
