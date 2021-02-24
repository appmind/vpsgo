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

var setportCmd = &cobra.Command{
	Use:   "setport [hostname]",
	Short: "Change port (perhaps need to configure firewall)",
	Long:  `Change the ssh port number (perhaps need to configure firewall).`,
	Run: func(cmd *cobra.Command, args []string) {
		hostname := config.GetHostname(args)
		host, err := config.GetHostByName(hostname)
		if err != nil {
			common.Exit(err.Error(), 1)
		}

		if !Force && (Port != 22 && Port != 0) {
			common.Exit("Only 22 or 0 is allowed", 1)
		}
		if Port == 0 {
			Port = uint(common.GetRandNumber(32768, 61000))
		}
		if Port == host.Port {
			common.Exit(fmt.Sprintf("Port %v in use.", Port), 1)
		}

		if !Force {
			anwser := common.AskQuestion(
				fmt.Sprintf("Change '%s' port?", hostname),
				[]string{"Y", "n"},
			)
			if strings.ToUpper(anwser) != "Y" {
				os.Exit(1)
			}
		}

		msg, err := setPort(Port, host, Pwd, true)
		if err != nil {
			common.Exit(err.Error(), 1)
		}

		fmt.Print(msg)
		host.Port = Port
		config.SaveHostToConfig(host)
	},
}

func init() {
	rootCmd.AddCommand(setportCmd)
	setportCmd.Flags().UintVarP(&Port, "number", "N", 22, "0 is a random number")
	setportCmd.Flags().BoolVarP(&Force, "force", "f", false, "no need to confirm")
}

func setPort(port uint, host config.Host, pwd string, force bool) (string, error) {
	commands := []string{
		fmt.Sprintf("sudo sed -i 's/^.*#*.*Port.*$/Port %v/' /etc/ssh/sshd_config", Port),
		fmt.Sprintf("sudo sed -i 's/^Port.*$/Port %v/' /etc/ssh/sshd_config", Port),
		"sudo service ssh reload 2>&1 >/dev/null",
		fmt.Sprintf("echo 'Port %v is ready.'", Port),
	}

	return ssh.Exec(commands, host, pwd, force)
}
