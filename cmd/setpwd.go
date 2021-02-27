package cmd

import (
	"fmt"
	"strings"

	"github.com/appmind/vpsgo/common"
	"github.com/appmind/vpsgo/config"
	"github.com/appmind/vpsgo/ssh"
	"github.com/spf13/cobra"
)

var setpwdCmd = &cobra.Command{
	Use:   "setpwd HOSTNAME",
	Short: "Disable password login or Reset password",
	Long:  `Disable password login or Reset password.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		host, err := config.GetHostByName(args[0])
		if err != nil {
			common.Exit(err.Error(), 1)
		}

		pwd := ""
		if Force {
			pwd = common.AskPass(
				fmt.Sprintf("%s@%s's password: ", host.User, host.Addr),
			)
		}

		msg, err := setPwd(host, pwd, true)
		if err != nil {
			common.Exit(err.Error(), 1)
		}

		fmt.Print(msg)
	},
}

func init() {
	rootCmd.AddCommand(setpwdCmd)
	setpwdCmd.Flags().BoolVar(&Force, "reset", false, "reset password")
}

func setPwd(host config.Host, pwd string, force bool) (string, error) {
	// Confirm that you can login without password
	msg, err := ssh.Exec([]string{"echo 'ok'"}, host, "", true)
	if err != nil {
		common.Exit(err.Error(), 1)
	}
	if strings.TrimSpace(msg) != "ok" {
		common.Exit("unknown error", 1)
	}

	commands := []string{}
	if pwd == "" {
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
