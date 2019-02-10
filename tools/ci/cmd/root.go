package cmd

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

var packages = []string{"db", "rexecd", "util"}
var pkgError = fmt.Sprintf("Invalid target package, options %s\n", strings.Join(packages, " "))

var pkg string
var timeoutSec int

var rootCmd = &cobra.Command{
	Use:   "ci",
	Short: "Continuous Integration for Metal",
	Long:  "Continuous Integration for Metal",
	RunE: func(cmd *cobra.Command, args []string) error {
		return nil
	},
}

// Execute the binary
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func validatePkgArg() error {
	if pkg == "" {
		return nil
	}
	for _, p := range packages {
		if pkg == p {
			return nil
		}
	}
	return errors.New(pkgError)
}

func init() {
	rootCmd.PersistentFlags().StringVar(&pkg, "pkg", "", fmt.Sprintf("Target package. Options: %s\n", strings.Join(packages, " ")))
}
