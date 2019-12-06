package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/params/convert"
	paramtypes "github.com/ethereum/go-ethereum/params/types"
	"github.com/ethereum/go-ethereum/params/types/common"
	"github.com/ethereum/go-ethereum/params/types/goethereum"
	"github.com/ethereum/go-ethereum/params/types/parity"
	"gopkg.in/urfave/cli.v1"
)

/*

formats: [parity|multigeth|geth|~~aleth(TODO)~~]

? If -[i|in] is not passed, then GUESS the proper config by trial and error. Exit 1 if not found.

> echainspec -[i|in] <format> convert -[o|out] multigeth [<stdin>|<my/file/path.json]
#@1> <JSON>

> echainspec -[i|in] <format> validate [<stdin>|<my/file/path.json]
#> <exitcode=(0|1)>

> echainspec -[i|in] <format> forks [<stdin>|<my/file/path.json]
#> 1150000
#> 1920000
#> 2250000
#> ...

> echainspec -[i|in] <format> ips [<stdin>|<my/file/path.json]
#> eip2 1150000
#> eip7 1150000
#> eip150 2250000
#> eip155 2650000
#> eip161abc 3000000
#> eip161d 3000000
#> eip170 3000000

*/

var gitCommit = "" // Git SHA1 commit hash of the release (set via linker flags)
var gitDate = ""

var (
	chainspecFormatTypes = map[string]common.Configurator{
		"parity": &parity.ParityChainSpec{},
		"multigeth": &paramtypes.Genesis{
			Config: &paramtypes.MultiGethChainConfig{},
		},
		"geth": &paramtypes.Genesis{
			Config: &goethereum.ChainConfig{},
		},
		// TODO
		// "aleth"
		// "retesteth"
	}
)

var chainspecFormats = func() []string {
	names := []string{}
	for k := range chainspecFormatTypes {
		names = append(names, k)
	}
	return names
}()

var defaultChainspecValues = map[string]common.Configurator{
	"foundation": params.DefaultGenesisBlock(),
	"classic":    params.DefaultClassicGenesisBlock(),
	// TODO
}

var defaultChainspecNames = func() []string {
	names := []string{}
	for k := range defaultChainspecValues {
		names = append(names, k)
	}
	return names
}()

var (
	app = utils.NewApp(gitCommit, gitDate, "the evm command line interface")

	FormatInFlag = cli.StringFlag{
		Name:  "inputf",
		Usage: fmt.Sprintf("Input format type [%s]", strings.Join(chainspecFormats, "|")),
		Value: "",
	}
	FileInFlag = cli.StringFlag{
		Name:        "file",
		Usage:       "Path to chain configuration file (JSON)",
	}
	DefaultValueFlag = cli.StringFlag{
		Name:  "default",
		Usage: fmt.Sprintf("Default chainspec values [%s]", strings.Join(defaultChainspecNames, "|")),
	}
	convertOutputFormatFlag = cli.StringFlag{
		Name:  "outputf",
		Usage: fmt.Sprintf("Output format type for converted configuration file [%s]", strings.Join(chainspecFormats, "|")),
	}
)

var globalChainspecValue common.Configurator

var errInvalidOutputFlag = errors.New("invalid output format type")
var errNoChainspecValue = errors.New("undetermined chainspec value")
var errInvalidDefaultValue = errors.New("no default chainspec found for name given")
var errInvalidChainspecValue = errors.New("could not read given chainspec")
var errEmptyChainspecValue = errors.New("missing chainspec data")

func mustGetChainspecValue(ctx *cli.Context) error {
	if ctx.GlobalIsSet(DefaultValueFlag.Name) {
		if ctx.GlobalString(DefaultValueFlag.Name) == "" {
			return errNoChainspecValue
		}
		v, ok := defaultChainspecValues[ctx.GlobalString(DefaultValueFlag.Name)]
		if !ok {
			return fmt.Errorf("error: %v, name: %s", errInvalidDefaultValue, ctx.GlobalString(DefaultValueFlag.Name))
		}
		globalChainspecValue = v
		return nil
	}
	data, err := readInputData(ctx)
	if err != nil {
		return err
	}
	configurator, err := unmarshalChainSpec(ctx.GlobalString(FormatInFlag.Name), data)
	if err != nil {
		return err
	}
	globalChainspecValue = configurator
	return nil
}

func convertf(ctx *cli.Context) error {
	c, ok := chainspecFormatTypes[ctx.String(convertOutputFormatFlag.Name)]
	if !ok {
		b, err := jsonMarshalPretty(globalChainspecValue)
		if err != nil {
			return err
		}
		fmt.Println(string(b))
		return nil
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

func init() {
	log.SetFlags(0)
	app.Flags = []cli.Flag{
		FormatInFlag,
		FileInFlag,
		DefaultValueFlag,
		convertOutputFormatFlag,
	}
	app.Commands = []cli.Command{
		validateCommand,
		forksCommand,
		ipsCommand,
	}
	app.Before = mustGetChainspecValue
	app.Action = convertf
}

func main() {
	if err := app.Run(os.Args); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
