package ssh

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/appmind/vpsgo/config"
	"github.com/melbahja/goph"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

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
	fd := int(os.Stdin.Fd())

	if terminal.IsTerminal(fd) {
		pass, err := terminal.ReadPassword(fd)
		if err != nil {
			panic(err)
		}

		fmt.Println("")
		return strings.TrimSpace(string(pass))
	}

	return ""
}

// Exec send and execute host commands via ssh
func Exec(cmds []string, host config.Host, pwd string, force bool) (string, error) {
	var err error
	var auth goph.Auth
	var callback ssh.HostKeyCallback

	if force {
		callback = ssh.InsecureIgnoreHostKey()
	} else {
		if callback, err = DefaultKnownHosts(); err != nil {
			return "", err
		}
	}

	if host.Keyfile != "" {
		// Start new ssh connection with private key.
		if auth, err = goph.Key(host.Keyfile, pwd); err != nil {
			// ssh: this private key is passphrase protected
			pwd = askPass("Private key passphrase: ")
			if auth, err = goph.Key(host.Keyfile, pwd); err != nil {
				return "", err
			}
		}
	} else {
		if pwd == "" {
			pwd = askPass(fmt.Sprintf("%s@%s's password: ", host.User, host.Addr))
		}
		auth = goph.Password(pwd)
	}

	if os.Getenv("GO") == "DEBUG" {
		fmt.Println(host, pwd, force)
	}

	client, err := goph.NewConn(&goph.Config{
		User:     host.User,
		Addr:     host.Addr,
		Port:     host.Port,
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
	out, err := client.Run(strings.Join(cmds, " && "))
	if err != nil {
		return "", err
	}

	// Get your output as []byte.
	return string(out), nil
}
