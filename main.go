package main

import (
	"log"
	"os"

	"github.com/omega1x/crta-box/src/application"
)

func main() {
	application.ListenSigterm()
	if err := application.Cli.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
