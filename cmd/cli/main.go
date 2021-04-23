package main

import (
	"log"
	"stylight/internal/cli"
)

func main() {
	err := cli.RunCLI()
	if err != nil {
		log.Fatalf(err.Error())
	}
}