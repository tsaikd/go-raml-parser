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

// SupportToCheckEmpty implement IsEmpty() instance
type SupportToCheckEmpty interface {
	IsEmpty() bool
}

// PostProcessConfig used for PostProcess()
type PostProcessConfig interface {
	RootDocument() RootDocument
	Library() Library
	Parser() Parser
	AnnotationUsage() map[string]bool
	TraitUsage() map[string]bool
}

func newPostProcessConfig(
	rootdoc RootDocument,
	library Library,
	parser Parser,
) PostProcessConfig {
	return postProcessConfigImpl{
		dataRootDocument:    rootdoc,
		dataLibrary:         library,
		dataParser:          parser,
		dataAnnotationUsage: map[string]bool{},
		dataTraitUsage:      map[string]bool{},
	}
}

type postProcessConfigImpl struct {
	dataRootDocument    RootDocument
	dataLibrary         Library
	dataParser          Parser
	dataAnnotationUsage map[string]bool
	dataTraitUsage      map[string]bool
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

func (t postProcessConfigImpl) AnnotationUsage() map[string]bool {
	return t.dataAnnotationUsage
}

func (t postProcessConfigImpl) TraitUsage() map[string]bool {
	return t.dataTraitUsage
}

// RAML built-in types
const (
	TypeNull    = "null"
	TypeBoolean = "boolean"
	TypeInteger = "integer"
	TypeNumber  = "number"
	TypeString  = "string"
	TypeObject  = "object"
	TypeArray   = "array"
	TypeFile    = "file"
	TypeBinary  = "binary"
)
