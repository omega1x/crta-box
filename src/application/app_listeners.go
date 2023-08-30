package application

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/omega1x/crta-box/src/conf"
)

func ListenSigterm() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		const EXIT_CODE = 0
		log.Printf(conf.LOG_WARN+"SIGTERM detected. Exit with code [%d]", EXIT_CODE)
		os.Exit(EXIT_CODE)
	}()
}
