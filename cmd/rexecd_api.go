package cmd

import (
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/metal-go/metal/config"
	"github.com/metal-go/metal/rexecd/api"
)

var rexecdAPICommand = &cobra.Command{
	Use:   "rexecd",
	Short: "The remote execution service",
	Long: strings.TrimSpace(`
rexecd allows gRPC clients and CLI clients to execute remote
commands at global scale.`),
	Run: func(cmd *cobra.Command, args []string) {
		config.RexecdInit()
		doneCh := make(chan struct{})
		err := api.Run(doneCh)
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(rexecdAPICommand)
}
