package config

import (
	"errors"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/appmind/vpsgo/common"
	"github.com/spf13/viper"
)

type Config struct {
	active   string    `mapstructure:"active"`
	hosts    []Host    `mapstructure:"hosts"`
	services []Service `mapstructure:"services"`
}

type Host struct {
	ID      string `mapstructure:"id"`
	Name    string `mapstructure:"name"`
	Addr    string `mapstructure:"addr"`
	Port    uint   `mapstructure:"port"`
	User    string `mapstructure:"user"`
	Keyfile string `mapstructure:"keyfile"`
}

type Service struct {
	Name string `mapstructure:"name"`
	Url  string `mapstructure:"url"`
}

func setConfig() {
	filename := common.ConfigFilename()
	dir, file := filepath.Split(filename)
	name := strings.Split(file, ".")
	viper.SetConfigName(name[0])
	viper.SetConfigType(name[1])
	viper.AddConfigPath(dir)
	// viper.AutomaticEnv()
	if _, err := os.Stat(filename); err != nil {
		viper.Set("active", "")
		viper.Set("hosts", []Host{})
		viper.Set("services", []Service{})
		if err = viper.SafeWriteConfig(); err != nil {
			log.Fatal(err)
		}
	}
}

func SaveConfig(values map[string]interface{}) error {
	for key, value := range values {
		viper.Set(key, value)
	}
	setConfig()
	return viper.WriteConfig()
}

func LoadConfig() error {
	setConfig()
	return viper.ReadInConfig()
}

func AppendHost(host Host) []Host {
	var found = false
	var index int
	hosts := []Host{}
	viper.UnmarshalKey("hosts", &hosts)
	for key, value := range hosts {
		if value.Name == host.Name {
			found = true
			index = key
			break
		}
	}
	if found {
		hosts[index] = host
	} else {
		hosts = append(hosts, host)
	}
	return hosts
}

func SaveHostToConfig(host Host) error {
	return SaveConfig(map[string]interface{}{
		"hosts":  AppendHost(host),
		"active": host.Name,
	})
}

func SetActiveHost(hostname string) error {
	if _, err := GetHostByName(hostname); err != nil {
		log.Fatal(err)
	}
	return SaveConfig(map[string]interface{}{
		"active": hostname,
	})
}

// GetHostname by the first parameter or configration
func GetHostname(args []string) (hostname string) {
	if len(args) > 0 {
		hostname = args[0]
	} else {
		hostname = viper.GetString("active")
	}
	if hostname == "" {
		log.Fatal("No default host.")
	}

	return
}

func GetHostByName(name string) (host Host, err error) {
	index := -1
	hosts := []Host{}
	err = errors.New("The host does not exist")
	viper.UnmarshalKey("hosts", &hosts)
	for key, value := range hosts {
		if value.Name == name {
			index = key
			break
		}
	}

	if index >= 0 {
		host = hosts[index]
		err = nil
	}
	return
}
