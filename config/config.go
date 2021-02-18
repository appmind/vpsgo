package config

import (
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/appmind/vpsgo/common"
	"github.com/appmind/vpsgo/ssh"
	"github.com/spf13/viper"
)

type Config struct {
	safemode bool      `mapstructure:"safemode"`
	phost    uint      `mapstructure:"phost"`
	hosts    []Host    `mapstructure:"hosts"`
	services []Service `mapstructure:"services"`
}

type Host struct {
	Name    string `mapstructure:"name"`
	Addr    string `mapstructure:"addr"`
	Port    string `mapstructure:"port"`
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
	viper.AutomaticEnv()
	if _, err := os.Stat(filename); err != nil {
		viper.Set("safemode", false)
		viper.Set("phost", 0)
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

func AddHostToConfig(vps ssh.Vps) {
	port := strconv.Itoa(int(vps.Port))
	name := common.MakeHash([]string{vps.Addr, port, vps.User})
	host := Host{
		Name:    name[0:5],
		Addr:    vps.Addr,
		Port:    port,
		User:    vps.User,
		Keyfile: vps.Key,
	}
	SaveConfig(map[string]interface{}{
		"hosts": AppendHost(host),
	})
}
