package cmd

import (
	"fmt"

	"github.com/codegangsta/cli"
	"github.com/tsaikd/KDGoLib/version"
)

func init() {
	version.VERSION = "1.0.0"

	Commands = append(Commands, cli.Command{
		Name:   "version",
		Usage:  "Show version detail",
		Action: versionAction,
	})
}

func versionAction(c *cli.Context) (err error) {
	verjson, err := version.Json()
	if err != nil {
		return
	}
	fmt.Println(verjson)
	return
}
