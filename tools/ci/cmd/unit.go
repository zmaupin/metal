package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/spf13/cobra"

	"github.com/metal-go/metal/util/with"
	"github.com/metal-go/metal/util/worker"
)

var unitTimeoutSec time.Duration

var unitCmd = &cobra.Command{
	Use:   "unit",
	Short: "Execute unit tests",
	Long:  "Execute unit tests",
	Run: func(cmd *cobra.Command, args []string) {
		err := with.Timeout(worker.Func(func(ctx context.Context, ch chan error) {
			fmt.Println(heading("Unit Test Stage"))
			paths := []string{"test", "-v", "-tags", "unit"}
			paths = append(paths, buildPaths()...)
			cmd := exec.CommandContext(ctx, "go", paths...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				ch <- err
			}
			close(ch)
		}), time.Duration(time.Second*unitTimeoutSec))
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	unitCmd.Flags().DurationVar(&unitTimeoutSec, "timeout", 5, timeoutFlagDesc)
	unitCmd.Flags().StringVar(&pkg, "pkg", "", fmt.Sprintf("Target package. Options: %s\n", strings.Join(packages, " ")))
	rootCmd.AddCommand(unitCmd)
}
