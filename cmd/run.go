package cmd

import (
	"fmt"

	"github.com/appmind/vpsgo/common"
	"github.com/appmind/vpsgo/config"
	"github.com/appmind/vpsgo/ssh"
	"github.com/spf13/cobra"
)

var command string

var runCmd = &cobra.Command{
	Use:   "run COMMAND [HOSTNAME]",
	Short: "Execute a command on the remote host",
	Long:  `Execute a command on the remote host.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		command := args[0]
		hostname := "."
		if len(args) > 1 {
			hostname = args[1]
		}

		host, err := config.GetHostByName(hostname)
		if err != nil {
			common.Exit(err.Error(), 1)
		}

		msg, err := ssh.Exec([]string{command}, host, "", true)
		if err != nil {
			common.Exit(err.Error(), 1)
		}
		fmt.Print(msg)
	},
}

func init() {
	rootCmd.AddCommand(runCmd)
}
