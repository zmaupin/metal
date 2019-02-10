package cmd

import (
	"context"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var unitCmd = &cobra.Command{
	Use:   "unit",
	Short: "Execute unit tests",
	Long:  "Execute unit tests",
	RunE: func(cmd *cobra.Command, args []string) error {
		return withTimeout(func(ctx context.Context, ch chan error) {
			paths := []string{"test", "-v", "-tags", "unit"}
			paths = append(paths, buildPaths()...)
			cmd := exec.CommandContext(ctx, "go", paths...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				ch <- err
			}
		})
	},
}

func init() {
	unitCmd.Flags().IntVar(&timeoutSec, "timeout", 5, timeoutFlagDesc)
	rootCmd.AddCommand(unitCmd)
}
