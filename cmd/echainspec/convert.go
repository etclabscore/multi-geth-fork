package main

import (
	"errors"
	"fmt"
	"strings"

	"github.com/ethereum/go-ethereum/params/convert"
	"gopkg.in/urfave/cli.v1"
)

var convertCommand = cli.Command{
	Name:        "convert",
	ShortName:   "c",
	Usage:       "",
	UsageText:   "",
	Description: "Convert a chain configuration from one format to another",
	ArgsUsage:   fmt.Sprintf("Requires flag --%s"),
	Flags: []cli.Flag{
		convertOutputFormatFlag,
	},
	Action: convertCmd,
}

var errInvalidOutputFlag = errors.New("invalid output format type")
var convertOutputFormatFlag = cli.StringFlag{
	Name:  "outputf",
	Usage: fmt.Sprintf("Output format type for converted configuration file [%s]", strings.Join(chainspecFormats, "|")),
}

func convertCmd(ctx *cli.Context) error {
	c, ok := chainspecFormatTypes[ctx.String(convertOutputFormatFlag.Name)]
	if !ok {
		return errInvalidOutputFlag
	}
	err := convert.Convert(globalChainspecValue, c)
	if err != nil {
		return err
	}
	b, err := jsonMarshalPretty(c)
	if err != nil {
		return err
	}
	fmt.Println(string(b))
	return nil
}

