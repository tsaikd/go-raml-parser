package parser

import "github.com/tsaikd/go-raml-parser/parser/parserConfig"

// RootDocument The root section of the RAML document describes the basic
// information about an API, such as its title and version. The root section
// also defines assets used elsewhere in the RAML document, such as types and
// traits.
type RootDocument struct {
	LibraryWrap
	RootDocumentExtra

	// directory of RAML file
	WorkingDirectory string `json:",omitempty"`
}

// UnmarshalYAML unmarshal RootDocument from YAML
func (t *RootDocument) UnmarshalYAML(unmarshaler func(interface{}) error) (err error) {
	if err = unmarshaler(&t.LibraryWrap); err != nil {
		return
	}
	if err = unmarshaler(&t.RootDocumentExtra); err != nil {
		return
	}
	return
}

// MarshalJSON marshal to json
func (t RootDocument) MarshalJSON() ([]byte, error) {
	return MarshalJSONWithoutEmptyStruct(t)
}

// IsEmpty return true if it is empty
func (t RootDocument) IsEmpty() bool {
	return t.LibraryWrap.IsEmpty() &&
		t.RootDocumentExtra.IsEmpty() &&
		t.WorkingDirectory == ""
}

var _ afterCheckUnusedTrait = RootDocument{}

func (t RootDocument) afterCheckUnusedTrait(conf PostProcessConfig) (err error) {
	ignore, err := conf.Parser().Get(parserConfig.IgnoreUnusedTrait)
	if err != nil {
		return
	}
	if ignore.(bool) {
		return
	}
	for name := range conf.TraitUsage() {
		return ErrorUnusedTrait1.New(nil, name)
	}
	return
}

// RootDocumentExtra contain fields no in Library
type RootDocumentExtra struct {
	// A short, plain-text label for the API. Its value is a string.
	Title string `yaml:"title" json:"title,omitempty"`

	// A substantial, human-friendly description of the API. Its value is a
	// string and MAY be formatted using markdown.
	Description string `yaml:"description" json:"description,omitempty"`

	// The version of the API, for example "v1". Its value is a string.
	Version string `yaml:"version" json:"version,omitempty"`

	// A URI that serves as the base for URIs of all resources. Often used as
	// the base of the URL of each resource containing the location of the API.
	// Can be a template URI.
	BaseURI string `yaml:"baseUri" json:"baseUri,omitempty"`

	// Named parameters used in the baseUri (template).
	BaseURIParameters Unimplement `yaml:"baseUriParameters" json:"baseUriParameters,omitempty"`

	// The protocols supported by the API.
	Protocols Unimplement `yaml:"protocols" json:"protocols,omitempty"`

	// The default media types to use for request and response bodies
	// (payloads), for example "application/json".
	MediaType string `yaml:"mediaType" json:"mediaType,omitempty"`

	// Additional overall documentation for the API.
	Documentation Unimplement `yaml:"documentation" json:"documentation,omitempty"`

	// The security schemes that apply to every resource and method in the API.
	SecuredBy Unimplement `yaml:"securedBy" json:"securedBy,omitempty"`

	// The resources of the API, identified as relative URIs that begin with
	// a slash (/). A resource node is one that begins with the slash and is
	// either at the root of the API definition or a child of a resource node.
	// For example, /users and /{groupId}.
	Resources Resources `yaml:",regexp:/.*" json:"resources,omitempty"`
}

// IsEmpty return true if it is empty
func (t RootDocumentExtra) IsEmpty() bool {
	return t.Title == "" &&
		t.Description == "" &&
		t.Version == "" &&
		t.BaseURI == "" &&
		t.BaseURIParameters.IsEmpty() &&
		t.Protocols.IsEmpty() &&
		t.MediaType == "" &&
		t.Documentation.IsEmpty() &&
		t.SecuredBy.IsEmpty() &&
		t.Resources.IsEmpty()
}
