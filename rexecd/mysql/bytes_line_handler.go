package mysql

import "context"

// BytesLineHandlerType type
type BytesLineHandlerType int

// // BytesLineQueue is a queue of bytes
// type BytesLineQueue struct {
// 	data  []byte
// 	input chan []byte
// 	done  chan bool
// }
//
// // NewBytesLineQueue returns a new BytesLineQueue
// func NewBytesLineQueue(max uint64) *BytesLineQueue {
// 	return &BytesLineQueue{
// 		data:  make([]byte, 0, max),
// 		input: make(chan []byte),
// 		done:  make(chan bool),
// 	}
// }
//
// func (b *BytesLineQueue) enqueue() {
// 	go func() {
// 		for {
// 			select {
// 			case d := <-b.input:
// 				b.data = append(b.data, d...)
// 			case <-b.done:
// 				break
// 			}
// 		}
// 	}()
// }
//
// func (b *BytesLineQueue) Dequeue() []byte {
//
// }
//
// // Enqueue pushes bytes onto the queue
// func (b *BytesLineQueue) Enqueue(d []byte) {
// 	b.input <- d
// }

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
