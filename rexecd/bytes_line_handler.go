package rexecd

import (
	"database/sql"

	proto_rexecd "github.com/metal-go/metal/proto/rexecd"
	"github.com/metal-go/metal/rexecd/mysql"
)

// BytesLineHandler handles a []byte
type BytesLineHandler interface {
	Handle(b []byte, ch chan error)
}

// BytesLineHandlerFunc enables functions to implement CommandResponseHandler
type BytesLineHandlerFunc func(b []byte, ch chan error)

// Handle satisfies CommandResponseHandler
func (blhf BytesLineHandlerFunc) Handle(c *proto_rexecd.CommandResponse, ch chan error) {
	blhf(c, ch)
}

// BytesLineMiddleware enables response wrapping
type BytesLineMiddleware interface {
	// Wrap must take a CommandResponseHandler and return a new CommandResponseHandler
	// that calls next. This wraps next and builds a frame in the handler call
	// stack.
	Wrap(next BytesLineHandler) BytesLineHandler
}

// NewWrappedBytesLineHandler composes a wrapped BytesLineHandler
func NewWrappedBytesLineHandler(h BytesLineHandler, stages ...BytesLineMiddleware) BytesLineHandler {
	last := h
	for i := len(stages) - 1; i >= 0; i-- {
		last = stages[i].Wrap(last)
	}
	return last
}

// DefaultBytesLineHandler handles
type DefaultBytesLineHandler struct {
	command *mysql.Command
	db      *sql.DB
}
