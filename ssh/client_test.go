package ssh

import (
	"regexp"
	"testing"

	"github.com/appmind/vpsgo/config"
)

var msg string = `Kernel Name:          Linux
Kernel Release:       4.19.84-microsoft-standard
Kernel Version:       #1 SMP Wed Nov 13 11:44:37 UTC 2019
Network Node Name:    -
Machine architecture: x86_64
`

func TestExec(t *testing.T) {
	host := config.Host{
		Name:    "test",
		Addr:    "127.0.0.1",
		Port:    22,
		User:    "test",
		Keyfile: "",
	}

	commands := []string{
		"echo 'Kernel Name:          '`uname -s`",
		"echo 'Kernel Release:       '`uname -r`",
		"echo 'Kernel Version:       '`uname -v`",
		"echo 'Network Node Name:    '`uname -n`",
		"echo 'Machine architecture: '`uname -m`",
	}

	out, _ := Exec(commands, host, "test", false)
	re1 := regexp.MustCompile(`(Network[ a-zA-Z:]+)(\w+)`)
	out = re1.ReplaceAllString(out, "${1}-")

	if out != msg {
		t.Error("result is not equal")
	}
}
