package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
)

var packages = []string{"db", "rexecd", "util"}
var pkgError = fmt.Sprintf("Invalid target package, options %s\n", strings.Join(packages, " "))
var timeoutFlagDesc = "timeout in seconds"

var pkg string
var timeoutSec int

var rootCmd = &cobra.Command{
	Use:   "ci",
	Short: "Continuous Integration for Metal",
	Long:  "Continuous Integration for Metal",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

// Execute the binary
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func banner(s string) string {
	b := strings.Repeat("#", 80)
	size := len(s)
	m := "#" + " " + s + " " + strings.Repeat("#", 80-size-3)
	return b + "\n" + m + "\n" + b
}

func heading(s string) string {
	return fmt.Sprintf("\n******** %s ********\n", s)
}

func notice(s string) string {
	return fmt.Sprintf("==> %s\n", s)
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

func buildPaths() []string {
	paths := []string{}
	if pkg != "" {
		paths = append(paths, filepath.Join("github.com", "metal-go", "metal", fmt.Sprintf("%s...", pkg)))
	} else {
		for _, p := range packages {
			paths = append(paths, filepath.Join("github.com", "metal-go", "metal", fmt.Sprintf("%s...", p)))
		}
	}
	return paths
}
