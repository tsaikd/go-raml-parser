package parse

import (
	"fmt"

	"github.com/tsaikd/KDGoLib/cliutil/cmder"
	"github.com/tsaikd/KDGoLib/jsonex"
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
		&cli.BoolFlag{
			Name:        "allowIntToBeNum",
			Usage:       "Allow integer type to be number type when checking",
			Destination: &allowIntToBeNum,
		},
	).
	SetAction(action)

var ramlFile string
var checkRAMLVersion bool
var allowIntToBeNum bool

var checkOptions = []parser.CheckValueOption{}

func action(c *cli.Context) (err error) {
	ramlParser := parser.NewParser()

	if allowIntToBeNum {
		checkOptions = append(checkOptions, parser.CheckValueOptionAllowIntegerToBeNumber(true))
	}

	if err = ramlParser.Config(parserConfig.CheckRAMLVersion, checkRAMLVersion); err != nil {
		return
	}

	if err = ramlParser.Config(parserConfig.CheckValueOptions, checkOptions); err != nil {
		return
	}

	rootdoc, err := ramlParser.ParseFile(ramlFile)
	if err != nil {
		return
	}

	jsondata, err := jsonex.MarshalIndent(rootdoc, "", "  ")
	if err != nil {
		return
	}
	fmt.Println(string(jsondata))

	return
}
