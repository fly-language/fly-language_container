/*
Copyright © 2022 Antonio De Caro antonio.decaro99@outlook.it

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

const VERSION = "0.0.1"

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "fly",
	Short: "A FLY CLI for starting and running fly scripts",
	Long: `A command line tool for compiling and executing FLY scripts.
Visit the FLY documentation to know more about this language.
(https://fly-language.github.io)`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().BoolP("help", "h", false, "Help message for toggle")
}
