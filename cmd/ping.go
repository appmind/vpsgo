package cmd

import (
	"fmt"
	"log"
	"runtime"

	"github.com/appmind/vpsgo/ssh"
	"github.com/spf13/cobra"
)

var (
	addr string = "127.0.0.1"
	port uint   = 22
	user string = "root"
	pwd  string = ""
)

var pingCmd = &cobra.Command{
	Use:   "ping IP_ADDR",
	Short: "Connect and get information from the server",
	Long:  `Connect and get information from the server.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addr = args[0]
		vps := ssh.Vps{
			Name: "ping",
			Addr: addr,
			Port: uint(port),
			User: user,
			Pwd:  pwd,
			Key:  "",
		}
		msg, err := ssh.Ping(vps, false)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(msg)
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
	pingCmd.Flags().UintVarP(&port, "port", "p", 22, "IP port number of the SSH service")
	pingCmd.Flags().StringVarP(&user, "user", "u", "root", "Username of the OS running the SSH service")
	if runtime.GOOS == "windows" {
		pingCmd.Flags().StringVarP(&pwd, "password", "P", "", "Password of the user (required)")
		// In windows, the password flag must be set, because terminal don't support input
		pingCmd.MarkFlagRequired("password")
	} else {
		pingCmd.Flags().StringVarP(&pwd, "password", "P", "", "Password of the user")
	}
}
