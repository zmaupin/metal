package rexecd

import (
	"strings"

	proto_rexecd "github.com/metal-go/metal/proto/rexecd"
)

// Server implements proto_rexecd.RexecdServer
type Server interface {
	proto_rexecd.RexecdServer
	Run() error
}

// KeyTypeKey takes a KeyType and returns the corresponding string
func KeyTypeKey(k proto_rexecd.KeyType) string {
	return strings.Replace(proto_rexecd.KeyType_name[int32(k)], "_", "-", -1)
}

// KeyTypeValue takes a KeyType value and returns the corresponding KeyType
func KeyTypeValue(s string) proto_rexecd.KeyType {
	return proto_rexecd.KeyType(proto_rexecd.KeyType_value[strings.Replace(s, "-", "_", -1)])
}
