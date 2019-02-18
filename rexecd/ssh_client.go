package rexecd

import (
	"bytes"

	"golang.org/x/crypto/ssh"
)

// NewSSHClientConfig builds an ssh.ClientConfig based on the given
// proto_rexecd.RegisterUserRequest and proto_rexecd.RegisterHostRequest. This
// will enforce FixedHostKey checking.
func NewSSHClientConfig(username string, privateUserKey, publicHostKey []byte) (*ssh.ClientConfig, error) {
	publicKey, _, _, _, err := ssh.ParseAuthorizedKey(bytes.TrimSpace(publicHostKey))
	if err != nil {
		return nil, err
	}
	signer, err := ssh.ParsePrivateKey(bytes.TrimSpace(privateUserKey))
	if err != nil {
		return nil, err
	}
	return &ssh.ClientConfig{
		User:              username,
		Auth:              []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback:   ssh.InsecureIgnoreHostKey(),
		HostKeyAlgorithms: []string{publicKey.Type()},
	}, nil
}
