package main

import (
	"log"

	"github.com/appmind/vpsgo/cmd"
	"github.com/appmind/vpsgo/config"
)

func main() {
	if err := config.LoadConfig(); err != nil {
		log.Fatal(err)
	}

	cmd.Execute()
}
