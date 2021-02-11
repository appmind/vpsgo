package ssh

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

// Vps define vps server config
type Vps struct {
	Name string
	Addr string
	Port uint
	User string
	Pwd  string
	Key  string
}

// DefaultKnownHostsPath returns default user knows hosts file.
func DefaultKnownHostsPath() (string, error) {

	var err error = nil
	var home string = ""

	// Support environment variable $env:home in windows
	if runtime.GOOS == "windows" {
		home = os.Getenv("home")
	}

	if home == "" {
		home, err = os.UserHomeDir()
		if err != nil {
			return "", err
		}
	}

	path := filepath.Join(home, ".ssh", "known_hosts")
	return fmt.Sprintf("%s", path), err
}

// DefaultKnownHosts returns host key callback from default known hosts path, and error if any.
func DefaultKnownHosts() (ssh.HostKeyCallback, error) {

	path, err := DefaultKnownHostsPath()
	if err != nil {
		return nil, err
	}

	// fmt.Println(path)
	return goph.KnownHosts(path)
}

// prompt to enter password through terminal
func askPass(msg string) string {

	fmt.Print(msg)

	pass, err := terminal.ReadPassword(0)

	if err != nil {
		panic(err)
	}

	fmt.Println("")

	return strings.TrimSpace(string(pass))
}

// Ping connect and get information from the vps
func Ping(vps Vps, issafe bool) (string, error) {

	var err error
	var auth goph.Auth
	var callback ssh.HostKeyCallback

	if issafe {
		if callback, err = DefaultKnownHosts(); err != nil {
			return "", err
		}
	} else {
		// if there a "man in the middle proxy", this can harm you!
		callback = ssh.InsecureIgnoreHostKey()
	}

	if vps.Key != "" {
		if vps.Pwd == "" && issafe {
			msg := fmt.Sprintf("Private key passphrase: ")
			vps.Pwd = askPass(msg)
		}
		// Start new ssh connection with private key.
		if auth, err = goph.Key(vps.Key, vps.Pwd); err != nil {
			return "", err
		}
	} else {
		if vps.Pwd == "" {
			msg := fmt.Sprintf("%s@%s's password: ", vps.User, vps.Addr)
			vps.Pwd = askPass(msg)
		}
		auth = goph.Password(vps.Pwd)
	}

	client, err := goph.NewConn(&goph.Config{
		User:     vps.User,
		Addr:     vps.Addr,
		Port:     vps.Port,
		Auth:     auth,
		Timeout:  5 * time.Second,
		Callback: callback,
	})

	if err != nil {
		return "", err
	}

	// Defer closing the network connection.
	defer client.Close()

	// Execute your command.
	commands := []string{
		"echo 'Kernel Name:                   '`uname -s`",
		"echo 'Kernel Release:                '`uname -r`",
		"echo 'Kernel Version:                '`uname -v`",
		"echo 'Network Node Name:             '`uname -n`",
		"echo 'Machine architecture:          '`uname -m`",
		"echo 'Processor architecture:        '`uname -p`",
		"echo 'HD Platform (OS architecture): '`uname -i`",
		"echo 'Operating System:              '`uname -o`",
		"echo 'Hostname:                      '`hostname`",
		"echo 'Username:                      '`whoami`",
	}
	out, err := client.Run(strings.Join(commands, " && "))
	if err != nil {
		return "", err
	}

	// Get your output as []byte.
	return string(out), nil
}
