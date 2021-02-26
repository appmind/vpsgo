package cmd

import (
	"fmt"

	"github.com/appmind/vpsgo/common"
	"github.com/appmind/vpsgo/config"
	"github.com/appmind/vpsgo/ssh"
	"github.com/spf13/cobra"
)

var newName string

var setnameCmd = &cobra.Command{
	Use:   "setname HOSTNAME",
	Short: "Change host name (perhaps need to restart the host)",
	Long:  `Change the host name (perhaps need to restart the host).`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		host := getHostConfirm(args[0], "Change name '%s' to '"+newName+"' ?")

		msg, err := setName(newName, host, Pwd, true)
		if err != nil {
			common.Exit(err.Error(), 1)
		}

		fmt.Print(msg)
		host.Name = newName
		config.SaveHostToConfig(host)
	},
}

func init() {
	rootCmd.AddCommand(setnameCmd)
	setnameCmd.Flags().BoolVarP(&Force, "force", "f", false, "no confirmation")
	setnameCmd.Flags().StringVar(&newName, "to", "", "a new name")
	setnameCmd.MarkFlagRequired("to")
}

func setName(name string, host config.Host, pwd string, force bool) (string, error) {
	commands := []string{
		fmt.Sprintf("sudo sed -i 's/`hostname`/%s/' /etc/hosts", name),
		fmt.Sprintf("sudo echo '%s' > /etc/hostname", name),
		fmt.Sprintf("sudo hostname %s", name),
		"echo 'done.'",
	}

	return ssh.Exec(commands, host, pwd, force)
}
