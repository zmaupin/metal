package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

// worker is a blocking function. Do not close the channel within a worker!
type worker func(ctx context.Context, ch chan error)

var packages = []string{"db", "rexecd", "util"}
var pkgError = fmt.Sprintf("Invalid target package, options %s\n", strings.Join(packages, " "))
var timeoutFlagDesc = "timeout in seconds"

var pkg string
var timeoutSec int

var rootCmd = &cobra.Command{
	Use:   "ci",
	Short: "Continuous Integration for Metal",
	Long:  "Continuous Integration for Metal",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute the binary
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func heading(s string) string {
	return fmt.Sprintf("\n******** %s ********\n", s)
}

func notice(s string) string {
	return fmt.Sprintf("==> %s\n", s)
}

func validatePkgArg() error {
	if pkg == "" {
		return nil
	}
	for _, p := range packages {
		if pkg == p {
			return nil
		}
	}
	return errors.New(pkgError)
}

func buildPaths() []string {
	paths := []string{}
	if pkg != "" {
		paths = append(paths, filepath.Join("github.com", "metal-go", "metal", fmt.Sprintf("%s...", pkg)))
	} else {
		for _, p := range packages {
			paths = append(paths, filepath.Join("github.com", "metal-go", "metal", fmt.Sprintf("%s...", p)))
		}
	}
	return paths
}

func withTimeout(worker worker) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeoutSec))
	defer cancel()
	run := func() chan error {
		ch := make(chan error)
		go func() {
			if err := validatePkgArg(); err != nil {
				ch <- err
				close(ch)
				return
			}
			worker(ctx, ch)
			close(ch)
		}()
		return ch
	}
	select {
	case err := <-run():
		return err
	case <-ctx.Done():
		return errors.New("timeout exceeded")
	}
}
