package cmd

import (
	"fmt"
	"log"

	"github.com/appmind/vpsgo/common"
	"github.com/appmind/vpsgo/config"
	"github.com/spf13/cobra"
)

var setkeyCmd = &cobra.Command{
	Use:   "setkey ADDRESS",
	Short: "Generate key file and set password-free login",
	Long:  `Generate key file and set password-free login.`,
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

		var err error // Module-level keyfile must have a value
		keyfile, err = common.MakeKeyfile(host.Name)
		if err != nil {
			log.Fatal(err)
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
	rootCmd.AddCommand(setkeyCmd)
	setkeyCmd.Flags().BoolVarP(&force, "force", "f", false, "Ignore known_hosts check")
	setkeyCmd.Flags().UintVarP(&port, "port", "p", 22, "IP port number of the SSH service")
	setkeyCmd.Flags().StringVarP(&user, "user", "u", "root", "Username of the remote system")
	setkeyCmd.Flags().StringVarP(&pwd, "password", "P", "", "Password of the user")
}
