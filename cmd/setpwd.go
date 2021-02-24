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

var setpwdCmd = &cobra.Command{
	Use:   "setpwd [hostname]",
	Short: "Disable password login or Change password",
	Long:  `Disable password login or Change password.`,
	Run: func(cmd *cobra.Command, args []string) {
		hostname := config.GetHostname(args)
		host, err := config.GetHostByName(hostname)
		if err != nil {
			common.Exit(err.Error(), 1)
		}

		if !Force {
			anwser := common.AskQuestion(
				fmt.Sprintf("Change '%s' password?", hostname),
				[]string{"Y", "n"},
			)
			if strings.ToUpper(anwser) != "Y" {
				os.Exit(1)
			}
		}

		msg, err := setPwd(host, "", true)
		if err != nil {
			common.Exit(err.Error(), 1)
		}

		fmt.Print(msg)
	},
}

func init() {
	rootCmd.AddCommand(setpwdCmd)
	setpwdCmd.Flags().BoolVarP(&Force, "force", "f", false, "no need to confirm")
	setpwdCmd.Flags().StringVarP(&Pwd, "password", "P", "", "\"\" is empty, it is forbidden")
}

func setPwd(host config.Host, pwd string, force bool) (string, error) {
	commands := []string{}
	if Pwd == "" {
		// Confirm that you can login without password
		msg, err := ssh.Exec([]string{"echo 'ok'"}, host, pwd, true)
		if err != nil {
			common.Exit(err.Error(), 1)
		}
		if strings.TrimSpace(msg) != "ok" {
			common.Exit("Maybe need to execute 'vps setkey' first", 1)
		}
		commands = []string{
			"sudo sed -i 's/PermitRootLogin yes/#PermitRootLogin prohibit-password/' /etc/ssh/sshd_config",
			"sudo sed -i 's/#PasswordAuthentication yes/PasswordAuthentication no/' /etc/ssh/sshd_config",
			"sudo service ssh reload 2>&1 >/dev/null",
			"echo 'done.'",
		}
	} else {
		commands = []string{
			"sudo sed -i 's/#PermitRootLogin prohibit-password/PermitRootLogin yes/' /etc/ssh/sshd_config",
			"sudo sed -i 's/PasswordAuthentication no/#PasswordAuthentication yes/' /etc/ssh/sshd_config",
			fmt.Sprintf("echo '%s:%s' | sudo chpasswd", host.User, pwd),
			"sudo service ssh reload 2>&1 >/dev/null",
			"echo 'done.'",
		}
	}

	return ssh.Exec(commands, host, pwd, force)
}
