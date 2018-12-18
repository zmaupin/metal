package user

import (
	"os/user"
	"strings"
)

// Home returns the path to the home directory for user of the current process
func Home() (string, error) {
	current, err := user.Current()
	if err != nil {
		return "", err
	}
	return current.HomeDir, nil
}

// Expand expands the tilde path
func Expand(path string) (string, error) {
	if !strings.HasPrefix(path, "~") {
		return path, nil
	}

	home, err := Home()
	if err != nil {
		return "", err
	}

	return home + path[1:], nil
}
