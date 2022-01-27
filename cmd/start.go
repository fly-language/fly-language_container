/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/docker/go-connections/nat"
	"github.com/go-git/go-git/v5"
	"github.com/spf13/cobra"
)

// startCmd represents the start command
var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the FLY IDE",
	Long: `Docker is required to run this command. It starts a Theia IDE on localhost:3000 on which you can use 
all FLY IDE integration like Compilation and Validation.`,

	Run: func(cmd *cobra.Command, args []string) {
		var err error
		var workspace string
		avoidBuild, _ := cmd.Flags().GetBool("avoid-build")
		if len(args) == 0 {
			workspace, _ = os.Getwd()
		} else {
			workspace, err = filepath.Abs(args[0])
			if err != nil {
				log.Fatalf("Invalid path: %s\n", args[0])
			}
		}
		err = startFly(avoidBuild, workspace)
		if err != nil {
			fmt.Printf("%sError on starting FLY server:%s\n%v\n", BoldLine, ColorReset+ColorRed, err)
		} else {
			fmt.Printf("%sFLY server succesfully started.%s\nActive workspace: %s\n", BoldLine, ColorReset, BoldLine+ColorCyan+workspace)
		}
	},
}

func init() {
	rootCmd.AddCommand(startCmd)

	startCmd.Flags().BoolP("avoid-build", "a", false, "Doesn't create the Docker image")
}

func startFly(avoidBuild bool, workspace string) error {
	ctx := context.Background()
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return err
	}
	defer cli.Close()

	images, err := cli.ImageList(ctx, types.ImageListOptions{
		All: true,
	})
	if err != nil {
		return err
	}

	toBuild := true
	for _, image := range images {
		if image.RepoTags[0] == fmt.Sprintf("fly:%s", VERSION) {
			toBuild = false
			break
		}
	}

	if toBuild && !avoidBuild {
		fmt.Printf("%sBuilding fly image%s\n", BoldLine, ColorReset)
		err := buildImage()
		if err != nil {
			return err
		}
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	resp, err := cli.ContainerCreate(ctx, &container.Config{
		Image: fmt.Sprintf("fly:%s", VERSION),
		Tty:   false,
	}, &container.HostConfig{
		Binds: []string{
			fmt.Sprintf("%s:/home/project:cached", workspace),
			fmt.Sprintf("%s/.aws:/home/theia/.aws", home),
		},
		AutoRemove: true,
		PortBindings: nat.PortMap{
			"3000/tcp": []nat.PortBinding{
				{
					HostIP:   "0.0.0.0",
					HostPort: "3000",
				},
			},
		},
	}, nil, nil, "fly-container")
	if err != nil {
		return err
	}

	if err := cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return err
	}

	return nil
}

func buildImage() error {
	dir, err := ioutil.TempDir(".", ".fly_tmp_*")
	if err != nil {
		log.Fatal(nil)
	}
	// defer os.RemoveAll(dir)

	git.PlainClone(dir, false, &git.CloneOptions{
		URL:      "https://github.com/fly-language/fly-language_container.git",
		Progress: os.Stdout,
	})

	out, err := exec.Command("docker", "build", "-t", fmt.Sprintf("fly:%s", VERSION), fmt.Sprintf("%s/installation", dir)).Output()
	if err != nil {
		fmt.Println(dir)
		return err
	}

	if len(out) > 0 {
		fmt.Println(string(out))
	} else {
		fmt.Println("FLY Image succesfully created")
	}

	return nil
}
