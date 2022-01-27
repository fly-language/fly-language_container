/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/spf13/cobra"
)

// runCmd represents the run command
var runCmd = &cobra.Command{
	Use:   "run",
	Short: "Compiles and run a FLY script",
	Args:  cobra.ExactArgs(1),

	Long: `Once you have written your script you can execute it with this command.
For example:
	fly run foo.fly`,
	Run: func(cmd *cobra.Command, args []string) {
		if info, err := os.Stat(args[0]); err == nil {

			if !strings.HasSuffix(info.Name(), ".fly") {
				fmt.Printf("%sCompilation failed:%s\nNot a Fly script.\n", BoldLine+ColorRed, ColorReset)
				os.Exit(0)
			}

			baseName := info.Name()[:len(info.Name())-4]

			ctx := context.Background()

			err = flyRun(baseName, ctx)
			if err != nil {
				fmt.Printf("%sError on running FLY script:%s\n%v\n", BoldLine, ColorReset+ColorRed, err)
			}

		} else if errors.Is(err, os.ErrNotExist) {
			fmt.Printf("%sCompilation failed:%s\nFile `%s` does not exists\n", BoldLine+ColorRed, ColorReset, args[0])
			os.Exit(0)

		} else {
			panic(err)
		}
	},
}

func init() {
	rootCmd.AddCommand(runCmd)

	// TODO add flags
}

func flyRun(fname string, ctx context.Context) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	runExec, err := cli.ContainerExecCreate(ctx, "fly-container", types.ExecConfig{
		WorkingDir:   "/home/project",
		AttachStdout: true,
		AttachStderr: true,
		Detach:       false,
		Cmd:          []string{"java", "-cp", "/home/fly/target/fly-project-0.0.1-jar-with-dependencies.jar", fmt.Sprintf("src-gen/%s.java", fname)},
	})
	if err != nil {
		return fmt.Errorf("container not started (you need to run `fly start` first")
	}

	hijResp, err := cli.ContainerExecAttach(ctx, runExec.ID, types.ExecStartCheck{})
	if err != nil {
		return err
	}
	defer hijResp.Close()

	data, _ := ioutil.ReadAll(hijResp.Reader)
	fmt.Print(string(data))

	return nil
}
