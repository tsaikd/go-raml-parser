package main

import (
	"github.com/tsaikd/KDGoLib/cliutil/cmder"

	// load cmd modules
	_ "github.com/tsaikd/go-raml-parser/cmd"
	_ "github.com/tsaikd/go-raml-parser/cmd/parse"
)

func main() {
	cmder.Main()
}
