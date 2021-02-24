package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/appmind/vpsgo/common"
	"github.com/appmind/vpsgo/config"
	"github.com/appmind/vpsgo/ssh"
	"github.com/spf13/cobra"
)

var newName string

var setnameCmd = &cobra.Command{
	Use:   "setname [hostname]",
	Short: "Change hostname (perhaps need to restart VPS host)",
	Long:  `Change the host name (perhaps need to restart VPS host).`,
	Run: func(cmd *cobra.Command, args []string) {
		hostname := config.GetHostname(args)
		host, err := config.GetHostByName(hostname)
		if err != nil {
			common.Exit(err.Error(), 1)
		}

		if !Force {
			anwser := common.AskQuestion(
				fmt.Sprintf("Change '%s' hostname?", hostname),
				[]string{"Y", "n"},
			)
			if strings.ToUpper(anwser) != "Y" {
				os.Exit(1)
			}
		}

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
	setnameCmd.Flags().StringVar(&newName, "to", "", "rename host to a new name")
	setnameCmd.Flags().BoolVarP(&Force, "force", "f", false, "no need to confirm")
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
