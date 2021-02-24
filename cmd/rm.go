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
		idname := args[0]
		hosts := []config.Host{}
		newHosts := []config.Host{}
		viper.UnmarshalKey("hosts", &hosts)

		index := -1
		for key, value := range hosts {
			if value.ID != idname && value.Name != idname {
				newHosts = append(newHosts, value)
			} else {
				index = key
			}
		}

		active := viper.GetString("active")
		if hosts[index].ID == active || hosts[index].Name == active {
			active = ""
		}

		if index >= 0 {
			config.SaveConfig(map[string]interface{}{
				"hosts":  newHosts,
				"active": active,
			})
			fmt.Printf("host '%s' is removed", idname)
		} else {
			fmt.Printf("host '%s' not found", idname)
		}
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
}
