package cmd

import (
	"os"

	"github.com/codegangsta/cli"
	"github.com/tsaikd/KDGoLib/errutil"
	"github.com/tsaikd/KDGoLib/version"
)

// Flags list all global flags for application
var Flags = []cli.Flag{}

// Commands list all commands for application
var Commands = []cli.Command{}

// Main entry point
func Main() {
	app := cli.NewApp()
	app.Name = "ramlParser"
	app.Usage = "Go RAML Parser"
	app.Version = version.String()
	app.Action = mainAction
	app.Flags = Flags
	app.Commands = Commands

	err := app.Run(os.Args)
	errutil.Trace(err)
}

func mainAction(c *cli.Context) (err error) {
	cli.ShowAppHelp(c)
	return
}
