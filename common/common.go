package common

import (
	"crypto/sha1"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

// GetHomeDir returns the current user's home directory.
// but the $HOME environment variable is the first.
func GetHomeDir() (home string) {
	if home = os.Getenv("home"); home == "" {
		home, _ = os.UserHomeDir()
	}
	return
}

func GetAppHomeDir(app string) (string, error) {
	path := filepath.Join(GetHomeDir(), app)
	return path, os.MkdirAll(path, 0770)
}

func ConfigFilename() string {
	path, _ := GetAppHomeDir(".vpsgo")
	return filepath.Join(path, "config.yaml")
}

func SafeMode(force bool) bool {
	if force {
		return false
	} else {
		return viper.GetBool("safemode")
	}
}

func MakeHash(in []string) string {
	s := strings.Join(in, " ")
	v := sha1.Sum([]byte(s))
	return fmt.Sprintf("%x", v)
}
