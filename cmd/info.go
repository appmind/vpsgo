package cmd

import (
	"fmt"

	"github.com/appmind/vpsgo/common"
	"github.com/appmind/vpsgo/config"
	"github.com/appmind/vpsgo/ssh"
	"github.com/spf13/cobra"
)

var infoCmd = &cobra.Command{
	Use:   "info HOSTNAME",
	Short: "Get information from the host",
	Long:  `Get information from the host.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		host, err := config.GetHostByName(args[0])
		if err != nil {
			common.Exit(err.Error(), 1)
		}

		msg, err := info(host, "", true)
		if err != nil {
			common.Exit(err.Error(), 1)
		}
		fmt.Print(msg)
	},
}

func init() {
	rootCmd.AddCommand(infoCmd)
}

func info(host config.Host, pwd string, force bool) (string, error) {
	commands := []string{
		"echo 'Kernel Name:          '`uname -s`",
		"echo 'Kernel Release:       '`uname -r`",
		"echo 'Kernel Version:       '`uname -v`",
		"echo 'Network Node Name:    '`uname -n`",
		"echo 'Machine architecture: '`uname -m`",
	}
	return ssh.Exec(commands, host, pwd, force)
}
