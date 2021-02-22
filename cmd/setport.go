package cmd

import (
	"fmt"
	"log"

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
			log.Fatal(err)
		}

		if !Force && (Port != 22 && Port != 0) {
			log.Fatal("Only 22 or 0 is allowed")
		}
		if Port == 0 {
			Port = uint(common.GetRandNumber(32768, 61000))
		}
		if Port == host.Port {
			log.Fatalf("Port %v in use.", Port)
		}

		commands := []string{
			fmt.Sprintf("sudo sed -i 's/^.*#*.*Port.*$/Port %v/' /etc/ssh/sshd_config", Port),
			fmt.Sprintf("sudo sed -i 's/^Port.*$/Port %v/' /etc/ssh/sshd_config", Port),
			"sudo service ssh reload 2>&1 >/dev/null",
			fmt.Sprintf("echo 'Port %v is ready.'", Port),
		}

		msg, err := ssh.Exec(commands, host, Pwd, true)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(msg)
		host.Port = Port
		config.SaveHostToConfig(host)
	},
}

func init() {
	rootCmd.AddCommand(setportCmd)
	setportCmd.Flags().BoolVarP(&Force, "force", "f", false, "be careful")
	setportCmd.Flags().UintVarP(&Port, "number", "N", 22, "0 is a random number")
}
