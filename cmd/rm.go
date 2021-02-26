package cmd

import (
	"fmt"
	"os"
	"strings"

	"github.com/appmind/vpsgo/common"
	"github.com/appmind/vpsgo/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rmCmd = &cobra.Command{
	Use:   "rm HOSTNAME",
	Short: "Remove a host from the VPS list",
	Long:  `Remove a VPS host from the VPS list.`,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		idname := args[0]
		hosts := []config.Host{}
		newHosts := []config.Host{}
		viper.UnmarshalKey("hosts", &hosts)

		active := viper.GetString("active")
		if idname == "." && active != "" {
			idname = active
		}

		index := -1
		for key, value := range hosts {
			if value.ID != idname && value.Name != idname {
				newHosts = append(newHosts, value)
			} else {
				index = key
			}
		}

		if index >= 0 {
			if !Force {
				anwser := common.AskQuestion(
					fmt.Sprintf("Remove '%s' host?", idname),
					[]string{"Y", "n"},
				)
				if strings.ToUpper(anwser) != "Y" {
					os.Exit(1)
				}
			}

			if hosts[index].ID == active || hosts[index].Name == active {
				active = ""
			}

			config.SaveConfig(map[string]interface{}{
				"hosts":  newHosts,
				"active": active,
			})
		} else {
			fmt.Printf("host '%s' not found\n", idname)
		}
	},
}

func init() {
	rootCmd.AddCommand(rmCmd)
	rmCmd.Flags().BoolVarP(&Force, "force", "f", false, "no confirmation")
}
