package rexecd

import (
	"strings"

	"golang.org/x/crypto/ssh"

	proto_rexecd "github.com/metal-go/metal/proto/rexecd"
)

// NewSSHClientConfig builds an ssh.ClientConfig based on the given
// proto_rexecd.RegisterUserRequest and proto_rexecd.RegisterHostRequest. This
// will enforce FixedHostKey checking.
func NewSSHClientConfig(username string, privateUserKey, publicHostKey []byte, hostKeyType proto_rexecd.KeyType) (*ssh.ClientConfig, error) {
	key, _, _, _, err := ssh.ParseAuthorizedKey(publicHostKey)
	if err != nil {
		return &ssh.ClientConfig{}, err
	}
	hostKeyCallback := ssh.FixedHostKey(key)
	authMethod, err := BuildAuthMethod(privateUserKey)
	if err != nil {
		return &ssh.ClientConfig{}, err
	}
	keyType := strings.Replace(proto_rexecd.KeyType_name[int32(hostKeyType)], "_", "-", -1)
	return &ssh.ClientConfig{
		User:              username,
		Auth:              []ssh.AuthMethod{authMethod},
		HostKeyAlgorithms: []string{keyType},
		HostKeyCallback:   hostKeyCallback,
	}, nil
}

// BuildAuthMethod returns an ssh.AuthMethod from the given private key presented
// as a byte array
func BuildAuthMethod(userPrivateKey []byte) (ssh.AuthMethod, error) {
	var sshAuthMethod ssh.AuthMethod
	signer, err := ssh.ParsePrivateKey(userPrivateKey)
	if err != nil {
		return sshAuthMethod, err
	}
	sshAuthMethod = ssh.PublicKeys(signer)
	return sshAuthMethod, nil
}
