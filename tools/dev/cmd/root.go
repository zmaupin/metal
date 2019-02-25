package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

const composeFile = "docker/rexecd-mysql-server.yml"

var rootCmd = &cobra.Command{
	Use:   "dev",
	Short: "Development environment utility",
	Long:  "Development environment utility",
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
