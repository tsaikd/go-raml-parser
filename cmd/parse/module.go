package parse

import (
	"encoding/json"
	"fmt"

	"github.com/tsaikd/KDGoLib/cliutil/cmder"
	"github.com/tsaikd/go-raml-parser/parser"
	"github.com/tsaikd/go-raml-parser/parser/parserConfig"
	"gopkg.in/urfave/cli.v2"
)

// Module info
var Module = cmder.NewModule("parse").
	SetUsage("Parse RAML file and show API in json format").
	AddFlag(
		&cli.StringFlag{
			Name:        "f",
			Aliases:     []string{"ramlfile"},
			Value:       "api.raml",
			Usage:       "Source RAML file",
			Destination: &ramlFile,
		},
		&cli.BoolFlag{
			Name:        "checkRAMLVersion",
			Usage:       "Check RAML Version",
			Destination: &checkRAMLVersion,
		},
	).
	SetAction(action)

var ramlFile string
var checkRAMLVersion bool

func action(c *cli.Context) (err error) {
	parser := parser.NewParser()

	if err = parser.Config(parserConfig.CheckRAMLVersion, checkRAMLVersion); err != nil {
		return
	}

	rootdoc, err := parser.ParseFile(ramlFile)
	if err != nil {
		return
	}

	jsondata, err := json.MarshalIndent(rootdoc, "", "  ")
	if err != nil {
		return
	}
	fmt.Println(string(jsondata))

	return
}
