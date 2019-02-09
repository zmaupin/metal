package cmd

import (
	"github.com/spf13/cobra"

	"github.com/metal-go/metal/config"
)

var appCommand = &cobra.Command{
	Use:   "app",
	Short: "The metal webapp",
	Long:  "The metal webapp",
	Run: func(cmd *cobra.Command, args []string) {
		config.AppInit()
	},
}

func init() {
	rootCmd.AddCommand(appCommand)
}
