package common

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
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

func GetAppHomeDir() (path string) {
	path = filepath.Join(GetHomeDir(), ".vps")
	os.MkdirAll(path, 0770)
	return
}

func ConfigFilename() string {
	return filepath.Join(GetAppHomeDir(), "config.yaml")
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

// MakeKeyfile 调用 ssh-keygen 生成密钥保存到应用目录
func MakeKeyfile(name string) (string, error) {
	keyfile := filepath.Join(GetAppHomeDir(), name)
	if name == "" || keyfile == "" {
		log.Fatal("Message: name is required.")
	}
	if _, err := os.Stat(keyfile); err == nil {
		log.Fatalf("%s already exists.\n", keyfile)
	}

	cmd := exec.Command(
		"ssh-keygen",
		"-t", "ed25519",
		"-f", keyfile,
		"-C", name,
		"-q",
		"-N", "''",
	)

	stdout, err := cmd.Output()
	fmt.Print(string(stdout))
	return keyfile, err
}

func GetKeyString(keyfile string) (string, error) {
	file, err := os.Open(keyfile)
	if err != nil {
		return "", err
	}
	defer file.Close()
	keystr, err := ioutil.ReadAll(file)
	return string(keystr), err
}
