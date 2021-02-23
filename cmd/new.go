package cmd

import (
	"log"
	"strings"

	"github.com/appmind/vpsgo/common"
	"github.com/appmind/vpsgo/config"
	"github.com/appmind/vpsgo/ssh"
	"github.com/spf13/cobra"
)

var newCmd = &cobra.Command{
	Use:   "new ADDRESS",
	Short: "Add a VPS host to list",
	Long:  `Add a VPS host to the VPS list.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addr := args[0]
		code := []string{addr, User}
		name := common.MakeHostID(code)

		// Prepare host parameters
		host := config.Host{
			ID:      name,
			Name:    name,
			Addr:    addr,
			Port:    Port,
			User:    User,
			Keyfile: Keyfile,
		}

		// Confirm you can login
		commands := []string{"echo 'done.'"}
		msg, err := ssh.Exec(commands, host, Pwd, true)
		if err != nil {
			log.Fatal(err)
		}
		if strings.TrimSpace(msg) != "done." {
			log.Fatal("Please check the parameters.")
		}

		// If ok, save it
		config.SaveHostToConfig(host)
	},
}

func init() {
	rootCmd.AddCommand(newCmd)
	newCmd.Flags().UintVarP(&Port, "port", "p", 22, "port number of the ssh service")
	newCmd.Flags().StringVarP(&User, "user", "u", "root", "username of the system running ssh")
	newCmd.Flags().StringVarP(&Keyfile, "idfile", "i", "", "id file for ssh login")
	newCmd.Flags().StringVarP(&Pwd, "password", "P", "", "password or passphrase")
}
