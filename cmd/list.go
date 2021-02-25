package cmd

import (
	"fmt"

	"github.com/appmind/vpsgo/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "Print VPS host list",
	Long:  `Print VPS list. A flag points to the default host.`,
	Run: func(cmd *cobra.Command, args []string) {
		hosts := []config.Host{}
		active := viper.GetString("active")
		viper.UnmarshalKey("hosts", &hosts)
		for _, v := range hosts {
			if v.ID == active || v.Name == active {
				fmt.Printf("* %s@%s:%d %s:%s %s\n", v.User, v.Addr, v.Port, v.ID, v.Name, v.Keyfile)
			} else {
				fmt.Printf("  %s@%s:%d %s:%s %s\n", v.User, v.Addr, v.Port, v.ID, v.Name, v.Keyfile)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
