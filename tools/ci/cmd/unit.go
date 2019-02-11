package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/spf13/cobra"
)

const unitHeader = `
################################################################################
# Unit Test Stage ##############################################################
################################################################################
`

var unitTimeoutSec int

var unitCmd = &cobra.Command{
	Use:   "unit",
	Short: "Execute unit tests",
	Long:  "Execute unit tests",
	Run: func(cmd *cobra.Command, args []string) {
		timeoutSec = unitTimeoutSec
		err := withTimeout(func(ctx context.Context, ch chan error) {
			fmt.Println(unitHeader)
			paths := []string{"test", "-v", "-tags", "unit"}
			paths = append(paths, buildPaths()...)
			cmd := exec.CommandContext(ctx, "go", paths...)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				ch <- err
			}
		})
		if err != nil {
			log.Fatal(err)
		}
	},
}

func init() {
	unitCmd.Flags().IntVar(&unitTimeoutSec, "timeout", 5, timeoutFlagDesc)
	unitCmd.Flags().StringVar(&pkg, "pkg", "", fmt.Sprintf("Target package. Options: %s\n", strings.Join(packages, " ")))
	rootCmd.AddCommand(unitCmd)
}
