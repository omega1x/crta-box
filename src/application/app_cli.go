package application

import (
	"github.com/omega1x/crta-box/src/conf"
	"github.com/urfave/cli/v2"
)

var Cli = &cli.App{
	Name:    "crta-box",
	Usage:   "stream data from *culvert rupture telltale aggregation* boxes (*CRTA-BOX*es).",
	Version: "0.0.1",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:   "log",
			Usage:  "force logging to `FILE`",
			Action: DoLog,
		},
	},
	Commands: []*cli.Command{
		{
			Name:    "box",
			Aliases: []string{"b"},
			Usage:   "check connection with *CRTA-BOX*",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "ipv4",
					Value: conf.BOX_DEFAULT_HOST,
					Usage: "IPv4-address of *CRTA-BOX* (PostgreSQL)",
				},
				&cli.UintFlag{
					Name:  "port",
					Value: conf.POSTGRES_DEFAULT_PORT,
					Usage: "IP-port of *CRTA-BOX*'s PostgreSQL",
				},
				&cli.StringFlag{
					Name:    "database",
					Aliases: []string{"d", "n"},
					Value:   conf.BOX_DEFAULT_DATABASE,
					Usage:   "name of *CRTA-BOX*'s PostgreSQL to write data to",
				},
				&cli.StringFlag{
					Name:    "username",
					Aliases: []string{"u"},
					Value:   conf.BOX_DEFAULT_USERNAME,
					Usage:   "name of *CRTA-BOX*'s PostgreSQL user with read/write privileges",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"a", "s"},
					Value:   conf.BOX_DEFAULT_PASSWORD,
					Usage:   "*CRTA-BOX*'s PostgreSQL user password",
				},
			},
			Action: DoBox,
		},
		{
			Name:    "house",
			Aliases: []string{"h"},
			Usage:   "check connection with *ClickHouse*-database",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "ipv4",
					Aliases: []string{"i", "4"},
					Value:   conf.HOUSE_DEFAULT_HOST,
					Usage:   "IPv4-address of target *ClickHouse*-database",
				},
				&cli.UintFlag{
					Name:    "port",
					Aliases: []string{"p", "t"},
					Value:   conf.HOUSE_DEFAULT_PORT,
					Usage:   "target *ClickHouse* IP-port",
				},
				&cli.StringFlag{
					Name:    "database",
					Aliases: []string{"d", "n"},
					Value:   conf.HOUSE_DEFAULT_DATABASE,
					Usage:   "name of *ClickHouse* database to write data to",
				},
				&cli.StringFlag{
					Name:    "username",
					Aliases: []string{"u"},
					Value:   conf.HOUSE_DEFAULT_USERNAME,
					Usage:   "name of *ClickHouse* user with read/write privileges",
				},
				&cli.StringFlag{
					Name:    "password",
					Aliases: []string{"a", "s"},
					Value:   conf.HOUSE_DEFAULT_PASSWORD,
					Usage:   "*ClickHouse* user password",
				},
			},
			Action: DoHouse,
		},
		{
			Name:    "stream",
			Aliases: []string{"s"},
			Usage:   "stream data from *CRTA-BOX* to target *ClickHouse*-database",
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:  "box_ipv4",
					Value: conf.BOX_DEFAULT_HOST,
					Usage: "IPv4-address of *CRTA-BOX* (PostgreSQL)",
				},
				&cli.UintFlag{
					Name:  "box_port",
					Value: conf.POSTGRES_DEFAULT_PORT,
					Usage: "IP-port of *CRTA-BOX*'s PostgreSQL",
				},
				&cli.StringFlag{
					Name:  "box_database",
					Value: conf.BOX_DEFAULT_DATABASE,
					Usage: "name of *CRTA-BOX*'s PostgreSQL to write data to",
				},
				&cli.StringFlag{
					Name:  "box_username",
					Value: conf.BOX_DEFAULT_USERNAME,
					Usage: "name of *CRTA-BOX*'s PostgreSQL user with read/write privileges",
				},
				&cli.StringFlag{
					Name:  "box_password",
					Value: conf.BOX_DEFAULT_PASSWORD,
					Usage: "*CRTA-BOX*'s PostgreSQL user password",
				},
				&cli.StringFlag{
					Name:  "house_ipv4",
					Value: conf.HOUSE_DEFAULT_HOST,
					Usage: "IPv4-address of target *ClickHouse*-database",
				},
				&cli.UintFlag{
					Name:  "house_port",
					Value: conf.HOUSE_DEFAULT_PORT,
					Usage: "target *ClickHouse* IP-port",
				},
				&cli.StringFlag{
					Name:  "house_database",
					Value: conf.HOUSE_DEFAULT_DATABASE,
					Usage: "name of *ClickHouse* database to write data to",
				},
				&cli.StringFlag{
					Name:  "house_username",
					Value: conf.HOUSE_DEFAULT_USERNAME,
					Usage: "name of *ClickHouse* user with read/write privileges",
				},
				&cli.StringFlag{
					Name:  "house_password",
					Value: conf.HOUSE_DEFAULT_PASSWORD,
					Usage: "*ClickHouse* user password",
				},
			},
			Action: DoStream,
		},
	},
}
