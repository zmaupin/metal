package cmd

import (
	"github.com/spf13/cobra"

	"github.com/metal-go/metal/apploader"
	"github.com/metal-go/metal/config"
	log "github.com/sirupsen/logrus"
)

var appCommand = &cobra.Command{
	Use:   "app",
	Short: "The metal webapp",
	Long:  "The metal webapp",
	Run: func(cmd *cobra.Command, args []string) {
		config.AppInit()
		log.Fatal(apploader.Run())
	},
}

func init() {
	rootCmd.AddCommand(appCommand)
}
