package cmd

import (
	"github.com/spf13/cobra"
	"os"
	"fmt"
)

var rootCmd = &cobra.Command{
	Use:   "udeploy",
	Short: "Micro Deploy is a small deployment tool",
	Long:  "Micro Deploy takes care of deploying very small applications",
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Help()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
