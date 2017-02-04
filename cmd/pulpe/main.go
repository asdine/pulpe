package main

import (
	"log"

	"github.com/blankrobot/pulpe/cli"
)

func main() {
	cmd := cli.New()

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
