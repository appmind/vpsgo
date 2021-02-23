package cmd

import (
	"log"

	"github.com/appmind/vpsgo/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var useCmd = &cobra.Command{
	Use:   "use HOSTNAME",
	Short: "Set the default VPS host",
	Long:  `Set the default VPS host by the host name.`,
	Run: func(cmd *cobra.Command, args []string) {
		hostname := ""
		if len(args) > 0 {
			hostname = args[0]
		}
		if hostname == "" {
			log.Fatal("A host name in the 'vps list' is required.")
		}

		if hostname != viper.GetString("active") {
			config.SetActiveHost(hostname)
		}
	},
}

func init() {
	rootCmd.AddCommand(useCmd)
}
