package rexecd

import (
	proto_rexecd "github.com/metal-go/metal/proto/rexecd"
)

// Server implements proto_rexecd.RexecdServer
type Server interface {
	proto_rexecd.RexecdServer
	Run() error
}
