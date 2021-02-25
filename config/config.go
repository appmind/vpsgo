package config

import (
	"errors"
	"fmt"
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
			common.Exit(err.Error(), 1)
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
	index := -1
	hosts := []Host{}
	viper.UnmarshalKey("hosts", &hosts)
	for key, value := range hosts {
		if value.ID == host.ID {
			index = key
			break
		}
	}

	if index == -1 {
		hosts = append(hosts, host)
	} else {
		host.Name = hosts[index].Name
		hosts[index] = host
	}
	return hosts
}

func SaveHostToConfig(host Host) error {
	return SaveConfig(map[string]interface{}{
		"hosts": AppendHost(host),
		// "active": host.Name,
	})
}

func SetActiveHost(hostname string) error {
	if _, err := GetHostByName(hostname); err != nil {
		common.Exit(err.Error(), 1)
	}
	return SaveConfig(map[string]interface{}{
		"active": hostname,
	})
}

func GetHostByName(name string) (host Host, err error) {
	if name == "." {
		name = viper.GetString("active")
		if name == "" {
			err = errors.New("No default host, 'vps use' first")
			return
		}
	}
	index := -1
	hosts := []Host{}
	viper.UnmarshalKey("hosts", &hosts)
	for key, value := range hosts {
		// Both ID and hostname are supported
		if value.Name == name || value.ID == name {
			index = key
			break
		}
	}

	if index == -1 {
		msg := fmt.Sprintf("'%s' does not exist, 'vps new' or 'vps list' first", name)
		err = errors.New(msg)
	} else {
		host = hosts[index]
	}
	return
}
