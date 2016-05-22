package parser

// Method RESTful API methods are operations that are performed on a resource.
// The OPTIONAL properties get, patch, put, post, delete, head, and options of
// a resource define its methods; these correspond to the HTTP methods defined
// in the HTTP version 1.1 specification RFC2616 and its extension, RFC5789.
// The value of these method properties is a map that has the following
// key-value pairs:
type Method struct {
	// An alternate, human-friendly method name in the context of the resource.
	// If the displayName node is not defined for a method, documentation tools
	// SHOULD refer to the resource by its key, which acts as the method name.
	DisplayName string `yaml:"displayName" json:"displayName,omitempty"`

	// A longer, human-friendly description of the method in the context of the
	// resource. Its value is a string and MAY be formatted using markdown.
	Description string `yaml:"description" json:"description,omitempty"`

	// Annotations to be applied to this API. An annotation is a map having
	// a key that begins with "(" and ends with ")" where the text enclosed in
	// parentheses is the annotation name, and the value is an instance of
	// that annotation.
	Annotations map[string]Unimplement `yaml:",regexp:\\(.*\\)" json:"annotations,omitempty"`

	// Detailed information about any query parameters needed by this method.
	// Mutually exclusive with queryString.
	QueryParameters QueryParameters `yaml:"queryParameters" json:"queryParameters,omitempty"`

	// Detailed information about any request headers needed by this method.
	Headers Headers `yaml:"headers" json:"headers,omitempty"`

	// The query string needed by this method. Mutually exclusive with queryParameters.
	QueryString Unimplement `yaml:"queryString" json:"queryString,omitempty"`

	// Information about the expected responses to a request.
	Responses Responses `yaml:"responses"`

	// A request body that the method admits.
	Bodies Bodies `yaml:"body" json:"body,omitempty"`

	// Explicitly specify the protocol(s) used to invoke a method, thereby
	// overriding the protocols set elsewhere, for example in the baseUri
	// or the root-level protocols node.
	Protocols Unimplement `yaml:"protocols" json:"protocols,omitempty"`

	// A list of the traits to apply to this method.
	Is []*Trait `yaml:"is" json:"is,omitempty"`

	// The security schemes that apply to this method.
	SecuredBy Unimplement `yaml:"securedBy" json:"securedBy,omitempty"`
}

// PostProcess for fill some field from RootDocument default config
func (t *Method) PostProcess(conf PostProcessConfig) (err error) {
	if err = t.Headers.PostProcess(conf); err != nil {
		return
	}
	if err = t.Responses.PostProcess(conf); err != nil {
		return
	}
	return
}

// IsEmpty return true if Method is empty
func (t *Method) IsEmpty() bool {
	return t.Headers.IsEmpty() &&
		t.Responses.IsEmpty() &&
		t.Bodies.IsEmpty()
}
