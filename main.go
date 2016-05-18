package main

import (
	"github.com/tsaikd/go-raml-parser/cmd"

	// load cmd modules
	_ "github.com/tsaikd/go-raml-parser/cmd/module"
)

func main() {
	cmd.Main()
}
