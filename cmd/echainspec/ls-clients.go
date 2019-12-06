package main

import (
	"fmt"

	"gopkg.in/urfave/cli.v1"
)

var lsClientsCommand = cli.Command{
	Name:               "ls-clients",
	Usage:              "List supported client configuration formats",
	Action: lsClients,
}

func lsClients(ctx *cli.Context) error {
	for _, name := range chainspecFormats {
		fmt.Println(name)
	}
	return nil
}
