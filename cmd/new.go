package cmd

import (
	"strings"

	"github.com/appmind/vpsgo/common"
	"github.com/appmind/vpsgo/config"
	"github.com/appmind/vpsgo/ssh"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new ADDRESS",
	Short: "Add a host to VPS list",
	Long:  `Add a host to the VPS list.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addr := args[0]
		id := common.MakeID([]string{addr, User})

		// Prepare host parameters
		host := config.Host{
			ID:      id,
			Name:    id,
			Addr:    addr,
			Port:    Port,
			User:    User,
			Keyfile: Keyfile,
		}

		// Confirm you can login
		msg, err := ssh.Exec([]string{"echo 'ok'"}, host, Pwd, true)
		if err != nil {
			common.Exit(err.Error(), 1)
		}
		if strings.TrimSpace(msg) != "ok" {
			common.Exit("unknown error", 1)
		}

		// If ok, save it
		config.SaveHostToConfig(host)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().UintVarP(&Port, "port", "p", 22, "port number of the ssh service")
	newCmd.Flags().StringVarP(&User, "user", "u", "root", "username of the system running ssh")
	newCmd.Flags().StringVarP(&Keyfile, "idfile", "i", "", "identity file (private key)")
	newCmd.Flags().StringVarP(&Pwd, "password", "P", "", "password or passphrase")
}
