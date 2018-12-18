package cmd

import "fmt"
import "os"

import "github.com/spf13/cobra"

var rootCmd = &cobra.Command{
	Use:   "metal",
	Short: "Metal is a server lifecycle management tool",
	Long: `
Metal is a server lifecycle management tool

See: https://github.com/metal-go/metal/blob/master/README.md`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute the CLI
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
