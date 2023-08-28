package application

import (
	"log"

	"github.com/muonsoft/validation/validate"
	"github.com/omega1x/crta-box/src/com"
	"github.com/omega1x/crta-box/src/conf"
	"github.com/urfave/cli/v2"
)

func DoBox(cCtx *cli.Context) error {

	var (
		host     string = cCtx.String("ipv4")
		port     uint16 = uint16(cCtx.Uint("port"))
		database string = cCtx.String("database")
		username string = cCtx.String("username")
		password string = cCtx.String("password")
	)

	if err := validate.IPv4(host); err != nil {
		log.Fatalf(conf.LOG_ERROR+"[%s] is invalid IPv4-address", host)
	}

	if port != conf.POSTGRES_DEFAULT_PORT {
		log.Printf(conf.LOG_WARN+"[%d] is a non-standard port for *PostgreSQL*", port)
	}

	log.Printf(conf.LOG_INFO+"Check connection to *CRTA-BOX* on host [%s:%d]...", host, port)

	if err := com.BoxPing(host, port, database, username, password); err != nil {
		log.Fatalf(conf.LOG_ERROR + "Unreachable!")
	}

	log.Printf(conf.LOG_INFO + "Ok")
	return nil
}
