package cmd

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/spf13/cobra"
)

var unitCmd = &cobra.Command{
	Use:   "unit",
	Short: "Execute unit tests",
	Long:  "Execute unit tests",
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(timeoutSec))
		defer cancel()
		run := func() chan error {
			ch := make(chan error)
			go func() {
				if err := validatePkgArg(); err != nil {
					ch <- err
				}
				paths := []string{"test", "-v", "-tags", "unit"}
				if pkg != "" {
					paths = append(paths, filepath.Join("github.com", "metal-go", "metal", fmt.Sprintf("%s...", pkg)))
				} else {
					for _, p := range packages {
						paths = append(paths, filepath.Join("github.com", "metal-go", "metal", fmt.Sprintf("%s...", p)))
					}
				}
				cmd := exec.CommandContext(ctx, "go", paths...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				if err := cmd.Run(); err != nil {
					ch <- err
				}
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
	},
}

func init() {
	unitCmd.Flags().IntVar(&timeoutSec, "timeout", 5, "timeout in seconds")
	rootCmd.AddCommand(unitCmd)
}
