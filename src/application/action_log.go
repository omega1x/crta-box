package application

import (
	"io"
	"log"
	"os"

	"github.com/omega1x/crta-box/src/conf"
	"github.com/urfave/cli/v2"
)

func DoLog(ctx *cli.Context, file_name string) error {
	logFile, err := os.OpenFile(file_name, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	log.Printf(conf.LOG_INFO+"Forced logging to [%s]", file_name)
	return nil
}
