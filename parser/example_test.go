package parser_test

import (
	"fmt"
	"os"
	"strings"

	"github.com/tsaikd/go-raml-parser/parser"
	"github.com/tsaikd/go-raml-parser/parser/parserConfig"
)

func ExampleParser() {
	ramlParser := parser.NewParser()
	chkopts := []parser.CheckValueOption{
		parser.CheckValueOptionAllowArrayToBeNull(true),
		parser.CheckValueOptionAllowIntegerToBeNumber(true),
	}
	if err := ramlParser.Config(parserConfig.CheckValueOptions, chkopts); err != nil {
		fmt.Println(err)
	}
	data := []byte(strings.TrimSpace(`
#%RAML 1.0
types:
    User:
        type: object
        properties:
            name: string
        examples:
            example1:
                name: Alice
            example2:
                name: Bob
/user:
    get:
        responses:
            200:
                body:
                    application/json:
                        type: User[]
	`))
	workdir := "."

	rootdoc, err := ramlParser.ParseData(data, workdir)
	if err != nil {
		fmt.Println(err)
	}

	if userType := rootdoc.Types["User"]; userType != nil {
		fmt.Println("User type:", userType.Type)
		for _, property := range userType.Properties.Slice() {
			fmt.Println("Property:", property.Name)
		}
	}

	// Output:
	// User type: object
	// Property: name
}

func ExampleParser_cache() {
	cacheDirectory := ".cache"
	ramlFilePath := "./raml-examples/others/mobile-order-api/api.raml"

	if err := os.RemoveAll(cacheDirectory); err != nil {
		fmt.Println(err)
	}

	ramlParser := parser.NewParser()
	if err := ramlParser.Config(parserConfig.CacheDirectory, cacheDirectory); err != nil {
		fmt.Println(err)
	}

	rootdoc, err := ramlParser.ParseFile(ramlFilePath)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rootdoc.Title)

	if _, err := os.Stat(cacheDirectory); os.IsNotExist(err) {
		fmt.Println("cache directory should exist")
	}

	rootdoc, err = ramlParser.ParseFile(ramlFilePath)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(rootdoc.Title)

	if err := os.RemoveAll(cacheDirectory); err != nil {
		fmt.Println(err)
	}

	if _, err := os.Stat(cacheDirectory); err == nil {
		fmt.Println("cache directory should not exist")
	}

	// Output:
	// Mobile Order API
	// Mobile Order API
}
