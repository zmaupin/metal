package mysql

import (
	"context"
)

// BytesLineHandlerType type
type BytesLineHandlerType int

const (
	// MySQLStdout type
	MySQLStdout BytesLineHandlerType = iota
	// MySQLStderr type
	MySQLStderr
)

// BytesLineHandler handles
type BytesLineHandler struct {
	command     *Command
	handlerType BytesLineHandlerType
	stdout      []byte
	stderr      []byte
}

// NewBytesLineHandler returns a new BytesLineHandler
func NewBytesLineHandler(command *Command, handlerType BytesLineHandlerType) *BytesLineHandler {
	return &BytesLineHandler{
		command:     command,
		handlerType: handlerType,
		stdout:      []byte{},
		stderr:      []byte{},
	}
}

// Handle satisfies rexecd.BytesLineHandler
func (m *BytesLineHandler) Handle(ctx context.Context, b []byte) error {
	switch m.handlerType {
	case MySQLStdout:
		m.stdout = append(m.stdout, b...)
	case MySQLStderr:
		m.stderr = append(m.stderr, b...)
	}
	return nil
}

// Finish wraps up the handling of bytes
func (m *BytesLineHandler) Finish(ctx context.Context) error {
	switch m.handlerType {
	case MySQLStdout:
		if len(m.stdout) > 0 {
			if err := m.command.AddStdout(ctx, m.stdout); err != nil {
				return err
			}
		}
	case MySQLStderr:
		if len(m.stderr) > 0 {
			if err := m.command.AddStderr(ctx, m.stderr); err != nil {
				return err
			}
		}
	}
	return nil
}
