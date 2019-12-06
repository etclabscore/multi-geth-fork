package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"strings"

	paramtypes "github.com/ethereum/go-ethereum/params/types"
	"github.com/ethereum/go-ethereum/params/types/common"
	"github.com/ethereum/go-ethereum/params/types/goethereum"
	"gopkg.in/urfave/cli.v1"
)

func readInputData(ctx *cli.Context) ([]byte, error) {
	if !ctx.GlobalIsSet(fileInFlag.Name) {
		if fi, err := os.Stdin.Stat(); err != nil {
			return nil, err
		} else if fi.Size() == 0 {
			return nil, errEmptyChainspecValue
		}
		return ioutil.ReadAll(os.Stdin)
	}
	return ioutil.ReadFile(ctx.GlobalString(fileInFlag.Name))
}

func unmarshalChainSpec(format string, data []byte) (conf common.Configurator, err error) {
	conf, ok := chainspecFormatTypes[format]
	if !ok {
		return nil, errInvalidChainspecValue
	}
	err = json.Unmarshal(data, conf)
	if err != nil {
		return conf, err
	}
	if !strings.Contains(format, "geth") {
		return
	}
	// Logic in params/types/gen_genesis.go already "auto-magically"
	// handles genesis Config unmarshaling, and IT PREFERS MULTIGETH,
	// and the two data types are not mutually exclusive (are overlapping).
	// So we need to redo custom unmarshaling logic to enforce data type
	// preference based on passed format value.
	type dec struct {
		Config common.ChainConfigurator `json:"config"`
	}
	var d dec
	if format == "geth" {
		d.Config = &goethereum.ChainConfig{}
	} else if format == "multigeth" {
		d.Config = &paramtypes.MultiGethChainConfig{}
	} else {
		panic("impossible")
	}
	t := chainspecFormatTypes[format].(*paramtypes.Genesis)
	err = json.Unmarshal(data, &d)
	if err != nil {
		return conf, err
	}
	t.Config = d.Config
	conf = t
	return
}

func jsonMarshalPretty(i interface{}) ([]byte, error) {
	return json.MarshalIndent(i, "", "    ")
}
