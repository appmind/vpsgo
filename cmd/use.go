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
		idname := args[0]
		if idname == "." {
			common.Exit("Need a hostname in the 'vps list'", 1)
		}
		host, err := config.GetHostByName(idname)
		if err != nil {
			common.Exit(err.Error(), 1)
		}

		active := viper.GetString("active")
		if host.ID != active && host.Name != active {
			config.SetActiveHost(host.Name)
			fmt.Printf("The default host is '%s', represented by '.'\n", host.Name)
		}
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
}
