package parser

import "strconv"

// Unimplement For extra clarity
type Unimplement struct {
	Value
}

// HTTPCode For extra clarity
type HTTPCode int // e.g. 200

func (t HTTPCode) String() string {
	return strconv.Itoa(int(t))
}

// PostProcessConfig used for PostProcess()
type PostProcessConfig interface {
	RootDocument() RootDocument
	Library() Library
	Parser() Parser
}

func newPostProcessConfig(
	rootdoc RootDocument,
	library Library,
	parser Parser,
) PostProcessConfig {
	return postProcessConfigImpl{
		dataRootDocument: rootdoc,
		dataLibrary:      library,
		dataParser:       parser,
	}
}

type postProcessConfigImpl struct {
	dataRootDocument RootDocument
	dataLibrary      Library
	dataParser       Parser
}

func (t postProcessConfigImpl) RootDocument() RootDocument {
	return t.dataRootDocument
}

func (t postProcessConfigImpl) Library() Library {
	return t.dataLibrary
}

func (t postProcessConfigImpl) Parser() Parser {
	return t.dataParser
}

const typeInteger = "integer"
const typeNumber = "number"
const typeString = "string"
const typeObject = "object"
const typeArray = "array"
