package app

import (
	"io"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func actionLog(ctx *cli.Context, file_name string) error {
	logFile, err := os.OpenFile(file_name, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		return err
	}
	log.SetOutput(io.MultiWriter(os.Stdout, logFile))
	log.Printf(LOG_INFO+"Forced logging to [%s]", file_name)
	return nil
}

var Cli = &cli.App{
	Name:    "crta-box",
	Usage:   "stream data from acoustic-based *culvert rupture telltale aggregation boxes* (*CRTA-BOX*es) to the dedicated ClickHouse-database",
	Version: "0.1.0",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:   "log",
			Usage:  "force logging to `FILE`",
			Action: actionLog,
		},
	},
	Commands: []*cli.Command{
		{
			Name:    "box",
			Aliases: []string{"b"},
			Usage:   "check connection with *CRTA-BOX*",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "ipv4",
					Aliases: []string{"4"},
					Value:   BOX_DEFAULT_HOST,
					Usage:   "IPv4-address",
				},
				&cli.UintFlag{
					Name:    "port",
					Aliases: []string{"p"},
					Value:   BOX_DEFAULT_PORT,
					Usage:   "TCP-port",
				},
				&cli.StringFlag{
					Name:    "database",
					Aliases: []string{"d"},
					Value:   BOX_DEFAULT_DATABASE,
					Usage:   "database name",
				},
				&cli.StringFlag{
					Name:    "username",
					Aliases: []string{"u"},
					Value:   BOX_DEFAULT_USERNAME,
					Usage:   "user name with read privileges",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"a"},
					Value:   BOX_DEFAULT_PASSWORD,
					Usage:   "user's password",
				},
			},
			Action: actionPingBox,
		},
		{
			Name:    "house",
			Aliases: []string{"c"},
			Usage:   "check connection with *ClickHouse*-database",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "ipv4",
					Aliases: []string{"4"},
					Value:   HOUSE_DEFAULT_HOST,
					Usage:   "IPv4-address",
				},
				&cli.UintFlag{
					Name:    "port",
					Aliases: []string{"p"},
					Value:   HOUSE_DEFAULT_PORT,
					Usage:   "TCP-port",
				},
				&cli.StringFlag{
					Name:    "database",
					Aliases: []string{"d"},
					Value:   HOUSE_DEFAULT_DATABASE,
					Usage:   "database name",
				},
				&cli.StringFlag{
					Name:    "table",
					Aliases: []string{"t"},
					Value:   "",
					Usage:   "table name",
				},
				&cli.StringFlag{
					Name:    "username",
					Aliases: []string{"u"},
					Value:   HOUSE_DEFAULT_USERNAME,
					Usage:   "user name with read and write privileges",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"a"},
					Value:   HOUSE_DEFAULT_PASSWORD,
					Usage:   "user's password",
				},
			},
			Action: actionPingHouse,
		},
		{
			Name:    "stream",
			Aliases: []string{"s"},
			Usage:   "stream data from *CRTA-BOX* server to target *ClickHouse*-database",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "box_ipv4",
					Value: BOX_DEFAULT_HOST,
					Usage: "IPv4-address for *CRTA-BOX* server",
				},
				&cli.UintFlag{
					Name:  "box_port",
					Value: BOX_DEFAULT_PORT,
					Usage: "TCP-port for *CRTA-BOX* server",
				},
				&cli.StringFlag{
					Name:  "box_database",
					Value: BOX_DEFAULT_DATABASE,
					Usage: "database name for *CRTA-BOX* server",
				},
				&cli.StringFlag{
					Name:  "box_username",
					Value: BOX_DEFAULT_USERNAME,
					Usage: "user name for *CRTA-BOX* server",
				},
				&cli.StringFlag{
					Name:  "box_password",
					Value: BOX_DEFAULT_PASSWORD,
					Usage: "user's password for *CRTA-BOX* server",
				},
				&cli.StringFlag{
					Name:  "house_ipv4",
					Value: HOUSE_DEFAULT_HOST,
					Usage: "IPv4-address for *ClickHouse*-database",
				},
				&cli.UintFlag{
					Name:  "house_port",
					Value: HOUSE_DEFAULT_PORT,
					Usage: "TCP-port for *ClickHouse*-database",
				},
				&cli.StringFlag{
					Name:  "house_database",
					Value: HOUSE_DEFAULT_DATABASE,
					Usage: "database name for *ClickHouse*-database",
				},

				&cli.StringFlag{
					Name:  "house_table",
					Value: HOUSE_DEFAULT_TABLE,
					Usage: "table name for *ClickHouse*-database",
				},

				&cli.StringFlag{
					Name:  "house_username",
					Value: HOUSE_DEFAULT_USERNAME,
					Usage: "username name for *ClickHouse*-database",
				},
				&cli.StringFlag{
					Name:  "house_password",
					Value: HOUSE_DEFAULT_PASSWORD,
					Usage: "user's password for *ClickHouse*-database",
				},
			},
			Action: actionStream,
		},
	},
}
