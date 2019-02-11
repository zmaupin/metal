package cmd

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/spf13/cobra"
)

var protoCmd = &cobra.Command{
	Use:   "proto",
	Short: "Build proto",
	Long:  "Build proto",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(heading("starting proto file compilation for language target: Go"))
		infos, err := ioutil.ReadDir("proto")
		if err != nil {
			log.Fatal(err)
		}
		for _, info := range infos {
			svc := info.Name()
			svcDir, err := ioutil.ReadDir(filepath.Join("proto", svc))
			if err != nil {
				log.Fatal(err)
			}
			for _, i := range svcDir {
				if filepath.Ext(i.Name()) == ".proto" {
					base := filepath.Join("proto", svc)
					name := filepath.Join(base, i.Name())
					fmt.Println(notice(fmt.Sprintf("compiling %s", name)))
					cmd := exec.Command("protoc", "-I", base, "--go_out", base, name)
					cmd.Stdout = os.Stdout
					cmd.Stderr = os.Stderr
					if err := cmd.Run(); err != nil {
						log.Fatal(err)
					}
				}
			}
		}
		fmt.Println(notice("done"))
	},
}

func init() {
	rootCmd.AddCommand(protoCmd)
}
