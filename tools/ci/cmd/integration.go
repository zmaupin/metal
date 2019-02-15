package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/metal-go/metal/util/with"
	"github.com/metal-go/metal/util/worker"
	"github.com/spf13/cobra"
)

var integrationTimeoutSec time.Duration
var baseIntegrationArgs = []string{"test", "-v"}

type integrationConfig struct {
	name   string
	worker worker.Interface
}

var integrationSuite = []integrationConfig{
	integrationConfig{name: "mysql", worker: worker.Func(func(ctx context.Context, ch chan error) {
		testArgs := append(baseIntegrationArgs, "-tags", "mysql", "-p", "1")
		testArgs = append(testArgs, buildPaths()...)
		cmd := exec.CommandContext(ctx, "go", testArgs...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			ch <- err
			return
		}
	})},
}

var integrationCmd = &cobra.Command{
	Use:   "integration",
	Short: "Execute integration tests",
	Long:  "Execute integration tests",
	Run: func(cmd *cobra.Command, args []string) {
		err := with.Timeout(worker.Func(func(ctx context.Context, ch chan error) {
			fmt.Println(banner("Test Integration Stage"))
			for _, suite := range integrationSuite {
				fmt.Printf(heading(suite.name))
				suite.worker.Work(ctx, ch)
				select {
				case err := <-ch:
					if err != nil {
						log.Fatal(err)
					}
				default:
					ch <- nil
				}
			}
		}), time.Duration(time.Second*integrationTimeoutSec))
		if err != nil {
			log.Fatal(err.Error())
		}
	},
}

func init() {
	integrationCmd.Flags().DurationVar(&integrationTimeoutSec, "timeout", time.Duration(time.Second*300), timeoutFlagDesc)
	integrationCmd.Flags().StringVar(&pkg, "pkg", "", fmt.Sprintf("Target package. Options: %s\n", strings.Join(packages, " ")))
	rootCmd.AddCommand(integrationCmd)
}
