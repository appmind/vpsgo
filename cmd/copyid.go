package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/appmind/vpsgo/common"
	"github.com/appmind/vpsgo/config"
	"github.com/spf13/cobra"
)

var copyidCmd = &cobra.Command{
	Use:   "copyid ADDRESS",
	Short: "Upload identity file and set password-free login",
	Long:  `Upload identity file and set password-free login.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addr := args[0]
		name := common.MakeHash([]string{addr, user})[0:5]
		host := config.Host{
			ID:      name,
			Name:    name,
			Addr:    addr,
			Port:    port,
			User:    user,
			Keyfile: "",
		}

		if keyfile == "" {
			log.Fatal("Identity file is required.")
		}
		if _, err := os.Stat(keyfile); err != nil {
			log.Fatal("Identity file is not exists.")
		}

		msg, err := setPubkey(host, pwd, common.SafeMode(force))
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(msg)
		host.Keyfile = keyfile
		config.SaveHostToConfig(host)
	},
}

func init() {
	rootCmd.AddCommand(copyidCmd)
	copyidCmd.Flags().BoolVarP(&force, "force", "f", false, "Ignore known_hosts check")
	copyidCmd.Flags().UintVarP(&port, "port", "p", 22, "IP port number of the SSH service")
	copyidCmd.Flags().StringVarP(&user, "user", "u", "root", "Username of the remote system")
	copyidCmd.Flags().StringVarP(&keyfile, "idfile", "i", "", "Identity file (Private key)")
	copyidCmd.Flags().StringVarP(&pwd, "password", "P", "", "Password of the user")
}
