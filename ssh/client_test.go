package ssh

import (
	"regexp"
	"testing"
)

var msg string = `Kernel Name:                   Linux
Kernel Release:                4.19.84-microsoft-standard
Kernel Version:                #1 SMP Wed Nov 13 11:44:37 UTC 2019
Network Node Name:             -
Machine architecture:          x86_64
Processor architecture:        x86_64
HD Platform (OS architecture): x86_64
Operating System:              GNU/Linux
Hostname:                      -
Username:                      test
`

func TestPing(t *testing.T) {
	vps := Vps{
		Name: "test",
		Addr: "127.0.0.1",
		Port: 22,
		User: "test",
		Pwd:  "test",
		Key:  "",
	}

	out, _ := Ping(vps, false)
	re1 := regexp.MustCompile(`(Network[ a-zA-Z:]+)(\w+)`)
	re2 := regexp.MustCompile(`(Host[ a-zA-Z:]+)(\w+)`)
	out = re1.ReplaceAllString(out, "${1}-")
	out = re2.ReplaceAllString(out, "${1}-")

	if out != msg {
		t.Error("result is not equal")
	}
}
