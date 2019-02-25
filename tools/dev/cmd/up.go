package cmd

import (
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var upCmd = &cobra.Command{
	Use:   "up",
	Short: "Start the development environment",
	Long:  "Start the development environment",
	Run: func(cmd *cobra.Command, args []string) {
		up := exec.Command("docker-compose", "--file", composeFile, "up", "--detach")
		up.Stdout = os.Stdout
		up.Stderr = os.Stderr
		if err := up.Run(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(upCmd)
}
