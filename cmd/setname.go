package cmd

import (
	"fmt"
	"log"

	"github.com/appmind/vpsgo/config"
	"github.com/appmind/vpsgo/ssh"
	"github.com/spf13/cobra"
)

var newName string

var setnameCmd = &cobra.Command{
	Use:   "setname HOSTNAME",
	Short: "Change hostname (perhaps need to restart VPS host)",
	Long:  `Change the host name (perhaps need to restart VPS host).`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		hostname := config.GetHostname(args)
		host, err := config.GetHostByName(hostname)
		if err != nil {
			log.Fatal(err)
		}

		commands := []string{
			fmt.Sprintf("sudo sed -i 's/`hostname`/%s/' /etc/hosts", newName),
			fmt.Sprintf("sudo echo '%s' > /etc/hostname", newName),
			fmt.Sprintf("sudo hostname %s", newName),
			"echo 'done.'",
		}

		msg, err := ssh.Exec(commands, host, "", true)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(msg)
		host.Name = newName
		config.SaveHostToConfig(host)
	},
}

func init() {
	rootCmd.AddCommand(setnameCmd)
	setnameCmd.Flags().StringVar(&newName, "to", "", "rename host to a new name")
}
