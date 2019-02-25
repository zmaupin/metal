package cmd

import (
	"log"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var downCmd = &cobra.Command{
	Use:   "down",
	Short: "Stop the development environment",
	Long:  "Stop the development environment",
	Run: func(cmd *cobra.Command, args []string) {
		down := exec.Command("docker-compose", "--file", composeFile, "down")
		down.Stdout = os.Stdout
		down.Stderr = os.Stderr
		if err := down.Run(); err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(downCmd)
}
