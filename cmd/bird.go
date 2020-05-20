package main

import (
	"bird/pkg/bird_server"
	//"bird/pkg/bird_server"
	"errors"
	"github.com/urfave/cli"
)

var (
	start = cli.Command{
		Name:                   "start",
		Usage:					"--conf_path=./config/bird.toml",
		Description:            "start bird",
		Action:                 starts,
		Flags:                  []cli.Flag{
			&cli.StringFlag{
				Name:"conf_path",
				Usage:"bird server config path",
			},
		},
	}
)

func starts(cli *cli.Context)error{
	path := cli.String("conf_path")
	if path == "" {
		return errors.New("please add --conf_path=./config/bird.toml")
	}
	bird_server.Run(path)
	return nil
}

