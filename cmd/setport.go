package cmd

import (
	"fmt"
	"strconv"

	"github.com/appmind/vpsgo/common"
	"github.com/appmind/vpsgo/config"
	"github.com/appmind/vpsgo/ssh"
	"github.com/spf13/cobra"
)

var pno uint = 0

var setportCmd = &cobra.Command{
	Use:   "setport HOSTNAME",
	Short: "Change port (perhaps need to configure firewall)",
	Long:  `Change the ssh port number (perhaps need to configure firewall).`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if !Force && (pno != 0 && pno != 22 && pno < 1024) {
			common.Exit("Port number: 22 | >=1024 | 0(random)", 1)
		}
		if pno == 0 {
			pno = uint(common.GetRandNumber(32768, 61000))
		}

		host := getHostConfirm(args[0],
			"Change '%s' port to "+strconv.Itoa(int(pno))+" ?",
		)
		if pno == host.Port {
			common.Exit(fmt.Sprintf("Port %v in use.", pno), 1)
		}

		msg, err := setPort(pno, host, Pwd, true)
		if err != nil {
			common.Exit(err.Error(), 1)
		}

		fmt.Print(msg)
		host.Port = pno
		config.SaveHostToConfig(host)
	},
}

func init() {
	rootCmd.AddCommand(setportCmd)
	setportCmd.Flags().BoolVarP(&Force, "force", "f", false, "no limit, no confirmation")
	setportCmd.Flags().UintVar(&pno, "to", 0, "22 or >=1024 or 0(random number)")
	setportCmd.MarkFlagRequired("to")
}

func setPort(port uint, host config.Host, pwd string, force bool) (string, error) {
	commands := []string{
		fmt.Sprintf("sudo sed -i 's/^#Port .*/Port %v/' /etc/ssh/sshd_config", Port),
		fmt.Sprintf("sudo sed -i 's/^Port .*/Port %v/' /etc/ssh/sshd_config", Port),
		"sudo service ssh reload 2>&1 >/dev/null",
		fmt.Sprintf("echo 'Port %v is ready.'", Port),
	}

	return ssh.Exec(commands, host, pwd, force)
}
