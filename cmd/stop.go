/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"fmt"
	"os/exec"

	"github.com/spf13/cobra"
)

// stopCmd represents the stop command
var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stops the fly container and shuts down the IDE",
	Run: func(cmd *cobra.Command, args []string) {

		command := exec.Command("docker", "container", "stop", "fly-container")
		err := command.Run()
		if err != nil {
			fmt.Println("FLY Container is not running")
		} else {
			fmt.Println("FLY Container stopped")
		}

	},
}

func init() {
	rootCmd.AddCommand(stopCmd)
}
