package cmd

import (
	"fmt"

	"github.com/appmind/vpsgo/common"
	"github.com/appmind/vpsgo/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var useCmd = &cobra.Command{
	Use:   "use HOSTNAME",
	Short: "Set the default VPS host",
	Long:  `Set the default host by the name in the 'vps list'.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		hostname := args[0] // cannot call config.GetHostname(args)
		host, err := config.GetHostByName(hostname)
		if err != nil {
			common.Exit(err.Error(), 1)
		}

		active := viper.GetString("active")
		if host.ID != active && host.Name != active {
			config.SetActiveHost(hostname)
			fmt.Printf("The default host is '%s', represented by '.'\n", hostname)
		}
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
}
