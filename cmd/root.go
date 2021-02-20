package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/appmind/vpsgo/common"
	"github.com/appmind/vpsgo/config"
	"github.com/appmind/vpsgo/ssh"
	"github.com/spf13/cobra"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

var (
	port    uint   = 22
	user    string = "root"
	pwd     string = ""
	force   bool   = false
	keyfile string = ""
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "vps",
	Short: "vpsgo is a CLI tool for VPS services management",
	Long: `vpsgo is a command line tool developed in golang that helps you manage
VPS services more simply and easily. For more information,
please visit https://github.com/appmind/vpsgo.`,
	// Run:   func(cmd *cobra.Command, args []string) {},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// cobra.OnInitialize(initConfig)
	if err := config.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.vpsgo.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".vpsgo" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".vpsgo")
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}

func setPubkey(host config.Host, pwd string, issafe bool) (string, error) {
	keystr, err := common.GetKeyString(keyfile + ".pub")
	if err != nil {
		log.Fatal(err)
	}
	if keystr == "" {
		log.Fatal("Invalid public key.")
	}

	commands := []string{
		"mkdir -p ~/.ssh",
		// "touch ~/.ssh/authorized_keys",
		fmt.Sprintf("echo '%s' > ~/.ssh/authorized_keys", keystr),
		"chmod -R go= ~/.ssh",
		"echo 'done.'",
	}

	return ssh.Exec(commands, host, pwd, issafe)
}
