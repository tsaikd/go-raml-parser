package module

import (
	"encoding/json"
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/tsaikd/KDGoLib/cliutil/cmder"
	"github.com/tsaikd/go-raml-parser/parser"
	"github.com/tsaikd/go-raml-parser/parser/parserConfig"
)

func init() {
	cmder.Commands = append(cmder.Commands, cli.Command{
		Name:   "parse",
		Usage:  "Parse RAML file and show API in json format",
		Action: action,
		Flags: []cli.Flag{
			cli.StringFlag{
				Name:        "f, ramlfile",
				Value:       "api.raml",
				Usage:       "Source RAML file",
				Destination: &ramlFile,
			},
			cli.BoolFlag{
				Name:        "checkRAMLVersion",
				Usage:       "Check RAML Version",
				Destination: &checkRAMLVersion,
			},
		},
	})
}

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
