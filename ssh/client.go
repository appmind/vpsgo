package ssh

import (
	"fmt"
	"os"
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

	return fmt.Sprintf("%s/.ssh/known_hosts", home), err
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
func Ping(vps Vps, checkKnownhosts bool) (string, error) {

	var auth goph.Auth
	var callback ssh.HostKeyCallback

	if checkKnownhosts {
		var err error
		callback, err = DefaultKnownHosts()
		if err != nil {
			return "", err
		}
	} else {
		// if there a "man in the middle proxy", this can harm you!
		callback = ssh.InsecureIgnoreHostKey()
	}

	if vps.Key != "" {
		// Start new ssh connection with private key.
		var err error
		auth, err = goph.Key(vps.Key, "")
		if err != nil {
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
	out, err := client.Run("lsb_release -a && echo 'Hostname: '`hostname` && echo 'Username: '`whoami`")
	if err != nil {
		return "", err
	}

	// Get your output as []byte.
	return string(out), nil
}
