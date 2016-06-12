package main

import (
	"github.com/tsaikd/KDGoLib/cliutil/cmder"
	"github.com/tsaikd/go-raml-parser/cmd"
	"github.com/tsaikd/go-raml-parser/cmd/parse"
)

func main() {
	cmder.Main(
		*cmd.Module,
		*parse.Module,
	)
}
