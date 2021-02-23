package cmd

import (
	"fmt"

	"github.com/appmind/vpsgo/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rmCmd = &cobra.Command{
	Use:   "remove HOSTNAME",
	Short: "Remove a VPS host from the VPS list",
	Long:  `Remove a VPS host from the VPS list.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		hostname := args[0]
		hosts := []config.Host{}
		newHosts := []config.Host{}
		viper.UnmarshalKey("hosts", &hosts)

		hostid := ""
		for _, value := range hosts {
			if value.Name != hostname {
				newHosts = append(newHosts, value)
			} else {
				hostid = value.ID
			}
		}

		active := viper.GetString("active")
		if hostid == active {
			active = ""
		}

		if hostid != "" {
			fmt.Printf("host '%s' removed.", hostname)
			config.SaveConfig(map[string]interface{}{
				"hosts":  newHosts,
				"active": active,
			})
		} else {
			fmt.Printf("host '%s' not found.", hostname)
		}
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
