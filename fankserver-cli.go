package main

import (
	"os"
	"time"

	"github.com/fankserver/fankserver-cli/api"

	cli "gopkg.in/urfave/cli.v2"
)

func main() {
	app := &cli.App{
		Description: "Fankserver cli",
		Compiled:    time.Now(),
		Copyright:   "(c) 2017 Fankserver Gaming Community",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Aliases: []string{"c"},
				Name:    "config",
				Usage:   "config file",
				Value:   "config.toml",
			},
		},
		Commands: []*cli.Command{
			{
				Action: api.Listen,
				Name:   "api",
				Usage:  "Starting the api",
				Flags: []cli.Flag{
					&cli.UintFlag{
						Aliases: []string{"p"},
						Name:    "port",
						Usage:   "http port",
						Value:   8080,
					},
					&cli.StringFlag{
						Aliases: []string{"i"},
						Name:    "interface",
						Usage:   "http interface",
						Value:   "",
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
