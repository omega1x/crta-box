/*
`crta-box` - industrial monitoring systems for power plants (command line utility).

Stream data from acoustic-based *culvert rupture telltale aggregation* boxes
(*CRTA-BOX*es) to the dedicated
[ClickHouse](https://clickhouse.com/)-database).

Usage and detailed manual: github.com/omega1x/crta-box
*/

package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/omega1x/crta-box/internal/app"
)

func listenSigterm() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		const L_EXIT_CODE_OK = 0
		log.Printf(app.LOG_WARN+"SIGTERM detected. Exit with code [%d]", L_EXIT_CODE_OK)
		os.Exit(L_EXIT_CODE_OK)
	}()
}

func main() {
	listenSigterm()
	if err := app.Cli.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
