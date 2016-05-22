package parser

// RootDocument The root section of the RAML document describes the basic
// information about an API, such as its title and version. The root section
// also defines assets used elsewhere in the RAML document, such as types and
// traits.
type RootDocument struct {
	LibraryWrap
	RootDocumentExtra

	// directory of RAML file
	WorkingDirectory string
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

// PostProcess for fill some field from RootDocument default config
func (t *RootDocument) PostProcess(parser Parser) (err error) {
	rootdoc := *t
	if err = t.LibraryWrap.PostProcess(rootdoc, parser); err != nil {
		return
	}
	if err = t.RootDocumentExtra.PostProcess(rootdoc); err != nil {
		return
	}
	return
}

// RootDocumentExtra contain fields no in Library
type RootDocumentExtra struct {
	// A short, plain-text label for the API. Its value is a string.
	Title string `yaml:"title"`

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
	MediaType Unimplement `yaml:"mediaType" json:"mediaType,omitempty"`

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

// PostProcess for fill some field from RootDocument default config
func (t *RootDocumentExtra) PostProcess(rootdoc RootDocument) (err error) {
	if err = t.Resources.PostProcess(rootdoc); err != nil {
		return
	}
	return
}
