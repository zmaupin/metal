package cmd

import "github.com/spf13/cobra"

var integrationCmd = &cobra.Command{
	Use:   "integration",
	Short: "Execute integration tests",
	Long:  "Execute integration tests",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

func init() {
	rootCmd.AddCommand(integrationCmd)
}
