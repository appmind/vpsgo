package common

import (
	"crypto/sha1"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
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

func MakeHash(in []string) string {
	s := strings.Join(in, " ")
	v := sha1.Sum([]byte(s))
	return fmt.Sprintf("%x", v)
}

// MakeKeyfile 调用 ssh-keygen 生成密钥保存到应用目录
func MakeKeyfile(name string, force bool) (string, error) {
	keyfile := filepath.Join(GetAppHomeDir(), name)
	if name == "" || keyfile == "" {
		log.Fatal("Name is required.")
	}
	if _, err := os.Stat(keyfile); err == nil {
		if force {
			os.Remove(keyfile)
			os.Remove(keyfile + ".pub")
		} else {
			log.Fatalf("%s already exists.\n", keyfile)
		}
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

func GetRandNumber(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return (rand.Intn(max-min+1) + min)
}
