package cmd

import (
	"fmt"
	"log"
	"strings"

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
			log.Fatal(err)
		}

		commands := []string{}
		if Pwd == "" {
			// Confirm that you can login without password
			commands = []string{
				"echo 'done.'",
			}
			msg, err := ssh.Exec(commands, host, "", true)
			if err != nil {
				log.Fatal(err)
			}
			if strings.TrimSpace(msg) != "done." {
				log.Fatal("Maybe need to execute 'vps setkey' first")
			}
			commands = []string{
				// #PasswordAuthentication yes
				"sudo sed -i 's/PermitRootLogin yes/PermitRootLogin prohibit-password/' /etc/ssh/sshd_config",
				"sudo service ssh reload 2>&1 >/dev/null",
				"echo 'done.'",
			}
		} else {
			commands = []string{
				"sudo sed -i 's/.*PermitRootLogin.*/PermitRootLogin yes/' /etc/ssh/sshd_config",
				fmt.Sprintf("echo '%s:%s' | sudo chpasswd", host.User, Pwd),
				"sudo service ssh reload 2>&1 >/dev/null",
				"echo 'done.'",
			}
		}

		msg, err := ssh.Exec(commands, host, "", true)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(msg)
	},
}

func init() {
	rootCmd.AddCommand(setpwdCmd)
	setpwdCmd.Flags().StringVarP(&Pwd, "password", "P", "", "default is empty, it is forbidden")
}
