package cmd

import (
	"fmt"
	"log"

	"github.com/appmind/vpsgo/common"
	"github.com/appmind/vpsgo/config"
	"github.com/appmind/vpsgo/ssh"
	"github.com/spf13/cobra"
)

var pingCmd = &cobra.Command{
	Use:   "ping ADDRESS",
	Short: "Connect and get information from the host",
	Long:  `Connect and get information from the remote server.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addr := args[0]
		host := config.Host{
			Name:    "ping",
			Addr:    addr,
			Port:    port,
			User:    user,
			Keyfile: keyfile,
		}
		msg, err := ping(host, pwd, common.SafeMode(force))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(msg)
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
	pingCmd.Flags().BoolVarP(&force, "force", "f", false, "Ignore known_hosts check")
	pingCmd.Flags().UintVarP(&port, "port", "p", 22, "IP port number of the SSH service")
	pingCmd.Flags().StringVarP(&user, "user", "u", "root", "Username of the remote system")
	pingCmd.Flags().StringVarP(&keyfile, "idfile", "i", "", "Identity file (Private key)")
	pingCmd.Flags().StringVarP(&pwd, "password", "P", "", "Password or Passphrase")
}

func ping(host config.Host, pwd string, issafe bool) (string, error) {
	commands := []string{
		"echo 'Kernel Name:          '`uname -s`",
		"echo 'Kernel Release:       '`uname -r`",
		"echo 'Kernel Version:       '`uname -v`",
		"echo 'Network Node Name:    '`uname -n`",
		"echo 'Machine architecture: '`uname -m`",
	}
	return ssh.Exec(commands, host, pwd, issafe)
}
