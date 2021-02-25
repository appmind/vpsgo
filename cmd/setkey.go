package cmd

import (
	"fmt"
	"os"

	"github.com/appmind/vpsgo/common"
	"github.com/appmind/vpsgo/config"
	"github.com/appmind/vpsgo/ssh"
	"github.com/spf13/cobra"
)

var setkeyCmd = &cobra.Command{
	Use:   "setkey HOSTNAME",
	Short: "Generate key file and set password-free login",
	Long:  `Generate key file and set password-free login.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		host := getHostConfirm(args[0], "Change '%s' Key file?")

		// Generate new key
		newKey, err := common.MakeKeyfile(host.ID, Force)
		if err != nil {
			common.Exit(err.Error(), 1)
		}

		// Set new key
		msg, err := setPubkey(newKey, host, Pwd, true)
		if err != nil {
			common.Exit(err.Error(), 1)
		}

		fmt.Print(msg)
		os.Remove(host.Keyfile)
		os.Remove(host.Keyfile + ".pub")

		// Update key
		host.Keyfile = newKey
		config.SaveHostToConfig(host)
	},
}

func init() {
	rootCmd.AddCommand(setkeyCmd)
	setkeyCmd.Flags().BoolVarP(&Force, "force", "f", false, "no confirmation")
}

func setPubkey(file string, host config.Host, pwd string, force bool) (string, error) {
	keystr, err := common.GetKeyString(file + ".pub")
	if err != nil {
		common.Exit(err.Error(), 1)
	}
	if keystr == "" || len(keystr) > 255 {
		common.Exit("Invalid public key", 1)
	}

	commands := []string{
		"mkdir -p ~/.ssh",
		fmt.Sprintf("echo '%s' > ~/.ssh/authorized_keys", keystr),
		"chmod -R go= ~/.ssh",
		"echo 'done.'",
	}

	return ssh.Exec(commands, host, pwd, force)
}
