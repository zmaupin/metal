package cmd

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/Unknwon/log"
	"github.com/spf13/cobra"
)

var baseIntegrationArgs = []string{"test", "-v"}

type integrationConfig struct {
	name   string
	worker worker
}

var integrationSuite = []integrationConfig{
	integrationConfig{name: "mysql", worker: MySQLWorker},
}

// MySQLWorker execute tests with mysql tags
func MySQLWorker(ctx context.Context, ch chan error) {
	buf := &bytes.Buffer{}
	args := strings.Split("run --interactive --publish 3306:3306 --env MYSQL_ROOT_PASSWORD=password --detach mysql", " ")
	cmd := exec.CommandContext(ctx, "docker", args...)
	cmd.Stdout = buf
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		ch <- err
		return
	}
	teardown := func() {
		id := strings.TrimSpace(buf.String())
		fmt.Printf(notice(fmt.Sprintf("tearing down %s", id)))
		cmd = exec.CommandContext(ctx, "docker", "rm", "--force", id)
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Error(err.Error())
		}
	}

	for {
		conn, err := net.Dial("tcp", "127.0.0.1:3306")
		if err != nil {
			continue
		}
		fmt.Println(notice("Connection established. Waiting for server initialization to complete..."))
		for i := 20; i > 0; i-- {
			fmt.Println(i)
			time.Sleep(time.Second)
		}
		conn.Close()
		break
	}

	os.Setenv("METAL_REXECD_DATA_SOURCE_NAME", "root:password@tcp(127.0.0.1:3306)/")
	os.Setenv("METAL_REXECD_SERVER_TYPE", "mysql")

	testArgs := append(baseIntegrationArgs, "-tags", "mysql")
	testArgs = append(testArgs, buildPaths()...)
	cmd = exec.CommandContext(ctx, "go", testArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		teardown()
		ch <- err
		return
	}
	teardown()
}

const integrationHeader = `
################################################################################
# Integration Test Stage #######################################################
################################################################################
`

var integrationCmd = &cobra.Command{
	Use:   "integration",
	Short: "Execute integration tests",
	Long:  "Execute integration tests",
	Run: func(cmd *cobra.Command, args []string) {
		err := withTimeout(func(ctx context.Context, ch chan error) {
			fmt.Println(integrationHeader)
			for _, suite := range integrationSuite {
				fmt.Printf(heading(suite.name))
				suite.worker(ctx, ch)
			}
		})
		if err != nil {
			log.Fatal(err.Error())
		}
	},
}

func init() {
	rootCmd.AddCommand(integrationCmd)
}
