package parser

// RootDocument The root section of the RAML document describes the basic
// information about an API, such as its title and version. The root section
// also defines assets used elsewhere in the RAML document, such as types and
// traits.
type RootDocument struct {
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

	// An alias for the equivalent "types" node for compatibility with
	// RAML 0.8. Deprecated - API definitions should use the "types" node
	// because a future RAML version might remove the "schemas" alias with
	// that node. The "types" node supports XML and JSON schemas.
	Schemas Unimplement `yaml:"schemas" json:"schemas,omitempty"`

	// Declarations of (data) types for use within the API.
	Types APITypes `yaml:"types" json:"types,omitempty"`

	// Declarations of traits for use within the API.
	Traits Unimplement `yaml:"traits" json:"traits,omitempty"`

	// Declarations of resource types for use within the API.
	ResourceTypes Unimplement `yaml:"resourceTypes" json:"resourceTypes,omitempty"`

	// Declarations of annotation types for use by annotations.
	AnnotationTypes Unimplement `yaml:"annotationTypes" json:"annotationTypes,omitempty"`

	// Annotations to be applied to this API. An annotation is a map having
	// a key that begins with "(" and ends with ")" where the text enclosed
	// in parentheses is the annotation name, and the value is an instance of
	// that annotation.
	Annotations map[string]Unimplement `yaml:",regexp:\\(.*\\)" json:"annotations,omitempty"`

	// Declarations of security schemes for use within the API.
	SecuritySchemes Unimplement `yaml:"securitySchemes" json:"securitySchemes,omitempty"`

	// The security schemes that apply to every resource and method in the API.
	SecuredBy Unimplement `yaml:"securedBy" json:"securedBy,omitempty"`

	// Imported external libraries for use within the API.
	Uses Unimplement `yaml:"uses" json:"uses,omitempty"`

	// The resources of the API, identified as relative URIs that begin with
	// a slash (/). A resource node is one that begins with the slash and is
	// either at the root of the API definition or a child of a resource node.
	// For example, /users and /{groupId}.
	Resources Resources `yaml:",regexp:/.*" json:"resources,omitempty"`
}
