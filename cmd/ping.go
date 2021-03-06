package cmd

import (
	"fmt"

	"github.com/appmind/vpsgo/common"
	"github.com/appmind/vpsgo/config"
	"github.com/appmind/vpsgo/ssh"
	"github.com/spf13/cobra"
)

var pingCmd = &cobra.Command{
	Use:   "ping ADDRESS",
	Short: "Get response from the host",
	Long:  `Get response from the host.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addr := args[0]
		host := config.Host{
			Name:    "ping",
			Addr:    addr,
			Port:    Port,
			User:    User,
			Keyfile: Keyfile,
		}

		msg, err := ping(host, Pwd, true)
		if err != nil {
			common.Exit(err.Error(), 1)
		}
		fmt.Print(msg)
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
	pingCmd.Flags().UintVarP(&Port, "port", "p", 22, "port number of the ssh service")
	pingCmd.Flags().StringVarP(&User, "user", "u", "root", "username of the system running ssh")
	pingCmd.Flags().StringVarP(&Keyfile, "idfile", "i", "", "identity file (private key)")
	pingCmd.Flags().StringVarP(&Pwd, "password", "P", "", "password or passphrase")
}

func ping(host config.Host, pwd string, force bool) (string, error) {
	commands := []string{
		"echo 'Network Node Name:    '`uname -n`",
		"echo 'OpenSSL Version:      '`openssl version`",
		"echo 'Distribution:         '`cat /etc/os-release | grep NAME | head -n 1 | awk -F '=' '{gsub(/\"/,\"\",$2);print $2}' && cat /etc/os-release | grep VERSION | head -n 1 | awk -F '=' '{gsub(/\"/,\"\",$2);print $2}'`",
	}
	return ssh.Exec(commands, host, pwd, force)
}
