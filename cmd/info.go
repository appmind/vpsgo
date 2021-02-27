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
		"echo 'Local IP:             '`hostname -I`",
		"echo 'Public IP:            '`curl -s ifconfig.me && echo`",
		"echo 'Distribution:         '`cat /etc/os-release | grep NAME | head -n 1 | awk -F '=' '{gsub(/\"/,\"\",$2);print $2}' && cat /etc/os-release | grep VERSION | head -n 1 | awk -F '=' '{gsub(/\"/,\"\",$2);print $2}'`",
		"echo 'Kernel Version:       '`uname -srm`",
		"echo 'OpenSSL Version:      '`openssl version`",
		"echo 'Memory Size:          '`cat /proc/meminfo | grep MemTotal | awk '{print $2,$3}'`",
		"echo 'Disk Size:           '`df -h --total | grep total | awk -F 'total' '{print $2}'`",
		"echo 'CPU Model:           '`cat /proc/cpuinfo | grep 'model name' | uniq | awk -F ':' '{print $2}' && cat /proc/cpuinfo | grep processor | wc -l`",
		"echo 'Time Zone:            '`cat /etc/timezone && date`",
		"echo 'Host Name:            '`hostname`",
	}
	return ssh.Exec(commands, host, pwd, force)
}
