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

			err := flyCompile(baseName, ctx)
			if err != nil {
				fmt.Printf("%sError on running FLY script:%s\n%v\n", BoldLine, ColorReset+ColorRed, err)
			} else {
				err = flyRun(baseName, ctx)
				if err != nil {
					fmt.Printf("%sError on running FLY script:%s\n%v\n", BoldLine, ColorReset+ColorRed, err)
				}
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

func flyCompile(fname string, ctx context.Context) error {

	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		panic(err)
	}
	defer cli.Close()

	compileExec, err := cli.ContainerExecCreate(ctx, "fly-container", types.ExecConfig{
		WorkingDir: "/home/project/src-gen",
		Cmd:        []string{"javac", "-classpath", "/home/fly/lib/*", fmt.Sprintf("%s.java", fname)},
	})
	if err != nil {
		return fmt.Errorf("container not started (you need to run `fly start` first)")
	}

	_, err = cli.ContainerExecAttach(ctx, compileExec.ID, types.ExecStartCheck{})
	if err != nil {
		return err
	}

	status, _ := cli.ContainerExecInspect(ctx, compileExec.ID)
	for status.Running {
		status, _ = cli.ContainerExecInspect(ctx, compileExec.ID)
	}
	return nil
}

func flyRun(fname string, ctx context.Context) error {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	runExec, err := cli.ContainerExecCreate(ctx, "fly-container", types.ExecConfig{
		WorkingDir:   "/home/project/src-gen",
		AttachStdout: true,
		Tty:          true,
		Cmd:          []string{"java", fname},
	})
	if err != nil {
		return err
	}

	hijResp, err := cli.ContainerExecAttach(ctx, runExec.ID, types.ExecStartCheck{
		Tty: true,
	})
	if err != nil {
		return err
	}

	buf, err := ioutil.ReadAll(hijResp.Reader)
	if err != nil {
		return err
	}
	fmt.Print(string(buf))
	return nil
}
