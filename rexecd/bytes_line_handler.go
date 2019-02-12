package rexecd

import (
	"context"
)

// BytesLineHandler handles a []byte
type BytesLineHandler interface {
	Handle(ctx context.Context, b []byte) error
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
