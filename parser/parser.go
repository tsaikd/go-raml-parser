package parser

import (
	"bytes"
	"io/ioutil"

	"github.com/advance512/yaml"
	"github.com/tsaikd/go-raml-parser/parser/parserConfig"
)

// NewParser create Parser instance
func NewParser() Parser {
	parser := &parserImpl{}
	return parser
}

// Parser used to parse raml file
type Parser interface {
	// Config Parser to change the behavior of parsing
	Config(config parserConfig.Enum, value interface{}) (err error)

	// ParseFile Parse a RAML file.
	// Return RootDocument or an error if something went wrong.
	ParseFile(filePath string) (rootdoc RootDocument, err error)

	// ParseFile Parse RAML from bynary data.
	// Return RootDocument or an error if something went wrong.
	ParseData(data []byte) (rootdoc RootDocument, err error)
}

type parserImpl struct {
	checkRAMLVersion bool
}

func (t *parserImpl) Config(config parserConfig.Enum, value interface{}) (err error) {
	switch config {
	case parserConfig.CheckRAMLVersion:
		return t.ConfigCheckRAMLVersion(value)
	default:
		return ErrorUnsupportedParserConfig1.New(nil, config)
	}
}

func (t *parserImpl) ConfigCheckRAMLVersion(value interface{}) (err error) {
	switch value.(type) {
	case bool:
		t.checkRAMLVersion = value.(bool)
		return nil
	default:
		return ErrorInvaludParserConfigValueType3.New(nil, parserConfig.CheckRAMLVersion, true, value)
	}
}

func (t parserImpl) CheckRAMLVersion(data []byte) (err error) {
	buffer := bytes.NewBuffer(data)
	firstLine, err := buffer.ReadString('\n')
	if err != nil {
		return
	}
	if firstLine[:10] != "#%RAML 1.0" {
		return ErrorUnexpectedRAMLVersion2.New(nil, "#%RAML 1.0", firstLine[:10])
	}
	return nil
}

func (t parserImpl) ParseFile(filePath string) (rootdoc RootDocument, err error) {
	filedata, err := ioutil.ReadFile(filePath)
	if err != nil {
		return
	}

	return t.ParseData(filedata)
}

func (t parserImpl) ParseData(data []byte) (rootdoc RootDocument, err error) {
	if t.checkRAMLVersion {
		if err = t.CheckRAMLVersion(data); err != nil {
			return
		}
	}

	err = yaml.Unmarshal(data, &rootdoc)
	return
}
