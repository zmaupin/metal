package rexecd

import (
	"bufio"
	"context"
	"sync"

	log "github.com/sirupsen/logrus"
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

// NewExecRunner returns a pointer to an ExecRunner
func NewExecRunner(cmd string, sshSession *ssh.Session, stdoutHandler BytesLineHandler, stderrHandler BytesLineHandler) *ExecRunner {
	runner := &ExecRunner{
		cmd:           cmd,
		sshSession:    sshSession,
		stdoutHandler: stdoutHandler,
		stderrHandler: stderrHandler,
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
		return -1, err
	}
	errReader, err := e.sshSession.StderrPipe()
	if err != nil {
		return -1, err
	}

	// Configure Scanners for stdout and stdin
	outScanner := bufio.NewScanner(outReader)
	outScanner.Buffer([]byte{}, 1e+9)
	errScanner := bufio.NewScanner(errReader)
	errScanner.Buffer([]byte{}, 1e+9)

	// Create a waitgroup stdout and stderr processing
	wg := &sync.WaitGroup{}
	wg.Add(2)

	// Feed bytes of lines to the given handler
	feeder := func(scanner *bufio.Scanner, handler BytesLineHandler, w *sync.WaitGroup) {
		for scanner.Scan() {
			b := scanner.Bytes()
			b = append(b, []byte("\n")...)
			e := handler.Handle(ctx, b)
			for {
				if err != nil {
					log.Error(e)
					e = handler.Handle(ctx, b)
				} else {
					break
				}
			}
		}
		w.Done()
	}

	// Injest bytes
	go func() { feeder(outScanner, e.stdoutHandler, wg) }()
	go func() { feeder(errScanner, e.stderrHandler, wg) }()

	// Run it
	err = e.sshSession.Run(e.cmd)

	// Check for errors
	if err != nil {
		return int64(-1), err
	}

	// Wait for the ingestors to finish
	wg.Wait()

	if err = e.stdoutHandler.Finish(ctx); err != nil {
		return -1, err
	}

	if err = e.stderrHandler.Finish(ctx); err != nil {
		return -1, err
	}

	exitErr, ok := err.(*ssh.ExitError)
	if ok {
		return int64(exitErr.Waitmsg.ExitStatus()), nil
	}

	return 0, nil
}
