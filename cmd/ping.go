package cmd

import (
	"fmt"
	"log"

	"github.com/appmind/vpsgo/ssh"
	"github.com/spf13/cobra"
)

var (
	addr    string = "127.0.0.1"
	port    uint   = 22
	user    string = "root"
	pwd     string = ""
	force   bool   = false
	keyfile string = ""
)

var pingCmd = &cobra.Command{
	Use:   "ping ADDRESS",
	Short: "Connect and get information from the server",
	Long:  `Connect and get information from the remote server.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addr = args[0]
		vps := ssh.Vps{
			Name: "ping",
			Addr: addr,
			Port: port,
			User: user,
			Pwd:  pwd,
			Key:  keyfile,
		}
		msg, err := ssh.Ping(vps, !force)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(msg)
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
	pingCmd.Flags().UintVarP(&port, "port", "p", 22, "IP port number of the SSH service")
	pingCmd.Flags().StringVarP(&user, "user", "u", "root", "Username of the remote system")
	pingCmd.Flags().StringVarP(&keyfile, "keyfile", "k", "", "Private key (Identity file)")
	pingCmd.Flags().BoolVarP(&force, "force", "f", false, "Ignore known hosts or Passphrase")
	pingCmd.Flags().StringVarP(&pwd, "password", "P", "", "Password or Passphrase")
}
