package cmd

import (
	"context"

	"github.com/spf13/cobra"
)

type integrationConfig struct {
	name string
}

var integrationSuite = []integrationConfig{
	integrationConfig{name: "mysql"},
}

var integrationCmd = &cobra.Command{
	Use:   "integration",
	Short: "Execute integration tests",
	Long:  "Execute integration tests",
	RunE: func(cmd *cobra.Command, args []string) error {
		return withTimeout(func(ctx context.Context, ch chan error) {})
	},
}

func init() {
	integrationCmd.Flags().IntVar(&timeoutSec, "timeout", 60, timeoutFlagDesc)
	rootCmd.AddCommand(integrationCmd)
}
