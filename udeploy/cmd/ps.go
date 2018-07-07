package cmd

import (
	"github.com/spf13/cobra"
	"net"
	"fmt"
	"os"
)

func init() {
	rootCmd.AddCommand(psCmd)
}

var psCmd = &cobra.Command{
	Use:   "ps",
	Short: "Print the current deployments",
	Long:  "Print the current deployments",
	Run: func(cmd *cobra.Command, args []string) {
		conn, err := net.Dial("unix", "/tmp/udeploy.sock")

		if err != nil {
			fmt.Printf("Error connecting %v\n", err)
			os.Exit(1)
		}

		conn.Write([]byte("ps"))

		buff := make([]byte, 1024)
		n, err := conn.Read(buff)

		if err != nil {
			panic(err)
		}

		fmt.Printf(string(buff[:n]))

	},
}
