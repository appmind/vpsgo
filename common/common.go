package common

import (
	"bufio"
	"crypto/sha1"
	"fmt"
	"io/ioutil"
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

func MakeID(in []string) string {
	return MakeHash(in)[0:5]
}

// MakeKeyfile generate the key and save it to the application directory
func MakeKeyfile(name string, force bool) (string, error) {
	timestamp := fmt.Sprintf("%v", time.Now().Unix())
	keyfile := filepath.Join(GetAppHomeDir(), name+timestamp)
	if _, err := os.Stat(keyfile); err == nil {
		// Use timestamp, but still prevent file name conflicts
		Exit(fmt.Sprintf("%s already exists\n", keyfile), 1)
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
	if os.Getenv("GO") == "DEBUG" {
		fmt.Print(string(stdout))
	}
	return keyfile, err
}

func GetKeyString(keyfile string) (string, error) {
	file, err := os.Open(keyfile)
	if err != nil {
		return "", err
	}
	defer file.Close()
	keystr, err := ioutil.ReadAll(file)
	return strings.TrimSpace(string(keystr)), err
}

func GetRandNumber(min, max int) int {
	rand.Seed(time.Now().UnixNano())
	return (rand.Intn(max-min+1) + min)
}

func Exit(msg string, code int) {
	fmt.Println(msg)
	os.Exit(code)
}

// AskQuestion prompt to enter anwser through terminal
func AskQuestion(msg string, in []string) (anwser string) {
	reader := bufio.NewReader(os.Stdin)
	tips := strings.Join(in, "/")

	for {
		fmt.Printf("%s (%s): ", msg, tips)
		anwser, _ = reader.ReadString('\n')
		anwser = strings.TrimSpace(anwser)

		for _, value := range in {
			if anwser == "" {
				if value == strings.ToUpper(value) {
					return value
				}
			} else {
				if strings.EqualFold(value, anwser) {
					return
				}
			}
		}
	}
}
