package mysql

import "context"

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
}

// NewBytesLineHandler returns a new BytesLineHandler
func NewBytesLineHandler(command *Command, handlerType BytesLineHandlerType) *BytesLineHandler {
	return &BytesLineHandler{
		command:     command,
		handlerType: handlerType,
	}
}

// Handle satisfies rexecd.BytesLineHandler
func (m *BytesLineHandler) Handle(ctx context.Context, b []byte) error {
	switch m.handlerType {
	case MySQLStdout:
		return m.command.AddStdoutLine(ctx, b)
	case MySQLStderr:
		return m.command.AddStderrLine(ctx, b)
	default:
		return nil
	}
}
