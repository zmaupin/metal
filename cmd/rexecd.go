package cmd

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/metal-go/metal/config"
	"github.com/metal-go/metal/rexecd/mysql"
)

var rexecdCommand = &cobra.Command{
	Use:   "rexecd",
	Short: "The remote execution service",
	Long: strings.TrimSpace(`
rexecd allows gRPC clients and CLI clients to execute remote
commands at global scale.`),
	Run: func(cmd *cobra.Command, args []string) {
		config.RexecdInit()
		switch config.RexecdGlobal.ServerType {
		case "mysql":
			err := mysql.New().Run()
			if err != nil {
				log.Fatal(err)
			}
		default:
			log.Fatal(fmt.Errorf("unknown server type"))
		}
	},
}

func init() {
	rootCmd.AddCommand(rexecdCommand)
}
