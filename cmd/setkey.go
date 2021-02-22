package cmd

import (
	"fmt"
	"log"

	"github.com/appmind/vpsgo/common"
	"github.com/appmind/vpsgo/config"
	"github.com/appmind/vpsgo/ssh"
	"github.com/spf13/cobra"
)

var setkeyCmd = &cobra.Command{
	Use:   "setkey ADDRESS",
	Short: "Generate key file and set password-free login",
	Long:  `Generate key file and set password-free login.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		addr := args[0]
		name := common.MakeHash([]string{addr, User})[0:5]

		// Prepare host login parameters
		host := config.Host{
			ID:      name,
			Name:    name,
			Addr:    addr,
			Port:    Port,
			User:    User,
			Keyfile: Keyfile,
		}

		// Generate new key
		newKey, err := common.MakeKeyfile(host.Name, true)
		if err != nil {
			log.Fatal(err)
		}

		// Set new key
		msg, err := setPubkey(newKey, host, Pwd, true)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(msg)
		// Update key
		host.Keyfile = newKey
		config.SaveHostToConfig(host)
	},
}

func init() {
	rootCmd.AddCommand(setkeyCmd)
	// setkeyCmd.Flags().BoolVarP(&Force, "force", "f", false, "rebuilt id files")
	setkeyCmd.Flags().UintVarP(&Port, "port", "p", 22, "port number of the ssh service")
	setkeyCmd.Flags().StringVarP(&User, "user", "u", "root", "username of the system running ssh")
	setkeyCmd.Flags().StringVarP(&Keyfile, "idfile", "i", "", "id file for ssh login")
	setkeyCmd.Flags().StringVarP(&Pwd, "password", "P", "", "password or passphrase")
}

func setPubkey(file string, host config.Host, pwd string, force bool) (string, error) {
	keystr, err := common.GetKeyString(file + ".pub")
	if err != nil {
		log.Fatal(err)
	}
	if keystr == "" || len(keystr) > 255 {
		log.Fatal("Invalid public key.")
	}

	commands := []string{
		"mkdir -p ~/.ssh",
		fmt.Sprintf("echo '%s' > ~/.ssh/authorized_keys", keystr),
		"chmod -R go= ~/.ssh",
		"echo 'done.'",
	}

	return ssh.Exec(commands, host, pwd, force)
}
