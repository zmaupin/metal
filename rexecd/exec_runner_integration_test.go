// +build integration

package rexecd

import (
	"context"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/metal-go/metal/test/lib"
	"golang.org/x/crypto/ssh"
)

// TODO setup ssh.ClientConfig for testing remote exec with PublicKeys authentication type
func TestExecRunner(t *testing.T) {
	sshd, err := lib.NewSSHDWorker()
	if err != nil {
		t.Fatal(err)
	}
	ctx := context.Background()
	if err = sshd.Setup(ctx); err != nil {
		t.Fatal(err)
	}
	defer sshd.Teardown(ctx)
	privateKey, err := ioutil.ReadFile(filepath.Join(os.Getenv("HOME"), ".ssh", "id_rsa"))
	if err != nil {
		t.Fatal(err)
	}

	signer, err := ssh.ParsePrivateKey(privateKey)
	if err != nil {
		t.Fatal(err)
	}

	publicKey, _, _, _, err := ssh.ParseAuthorizedKey(sshd.PublicHostKey)
	if err != nil {
		t.Fatal(err)
	}

	client, err := ssh.Dial("tcp", "127.0.0.1:2222", &ssh.ClientConfig{
		User:              "root",
		Auth:              []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback:   ssh.FixedHostKey(publicKey),
		HostKeyAlgorithms: []string{publicKey.Type()},
	})
	if err != nil {
		t.Fatal(err)
	}

	session, err := client.NewSession()
	if err != nil {
		t.Fatal(err)
	}
	defer session.Close()

	stdout, err := session.StdoutPipe()
	if err != nil {
		t.Fatal(err)
	}

	if err = session.Run(`bash -c "echo {1..10} | sed 's/1/swapped/g'"`); err != nil {
		t.Fatal(err)
	}

	stdoutBytes, err := ioutil.ReadAll(stdout)
	if err != nil {
		t.Fatal(err)
	}

	got := string(stdoutBytes)
	expected := "swapped 2 3 4 5 6 7 8 9 swapped0\n"

	if got != expected {
		t.Errorf("expected %s got %s", expected, got)
	}
}
