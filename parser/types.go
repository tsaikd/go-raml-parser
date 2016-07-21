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
	Parser() Parser
	RootDocument() *RootDocument
	Library() *Library
	AnnotationUsage() map[string]bool
	TraitUsage() map[string]bool
}

func newPostProcessConfig(
	parser Parser,
	rootdoc *RootDocument,
	library *Library,
	annotationUsage map[string]bool,
	traitUsage map[string]bool,
) PostProcessConfig {
	if parser == nil {
		parser = NewParser()
	}
	if rootdoc == nil {
		rootdoc = &RootDocument{}
	}
	if library == nil {
		library = &rootdoc.Library
	}
	if annotationUsage == nil {
		annotationUsage = map[string]bool{}
	}
	if traitUsage == nil {
		traitUsage = map[string]bool{}
	}
	return postProcessConfigImpl{
		dataParser:          parser,
		dataRootDocument:    rootdoc,
		dataLibrary:         library,
		dataAnnotationUsage: annotationUsage,
		dataTraitUsage:      traitUsage,
	}
}

type postProcessConfigImpl struct {
	dataParser          Parser
	dataRootDocument    *RootDocument
	dataLibrary         *Library
	dataAnnotationUsage map[string]bool
	dataTraitUsage      map[string]bool
}

func (t postProcessConfigImpl) Parser() Parser {
	return t.dataParser
}

func (t postProcessConfigImpl) RootDocument() *RootDocument {
	return t.dataRootDocument
}

func (t postProcessConfigImpl) Library() *Library {
	return t.dataLibrary
}

func (t postProcessConfigImpl) AnnotationUsage() map[string]bool {
	return t.dataAnnotationUsage
}

func (t postProcessConfigImpl) TraitUsage() map[string]bool {
	return t.dataTraitUsage
}

type typoCheck map[string]*Value

func (t typoCheck) IsEmpty() bool {
	return len(t) == 0
}

func (t typoCheck) Names() []string {
	names := []string{}
	for name := range t {
		names = append(names, name)
	}
	return names
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
