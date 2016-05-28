package parser

// Resources map of Resource
type Resources map[string]*Resource

// PostProcess for fill some field from RootDocument default config
func (t *Resources) PostProcess(conf PostProcessConfig) (err error) {
	if t == nil {
		return
	}
	for _, resource := range *t {
		if err = resource.PostProcess(conf); err != nil {
			return
		}
	}
	return
}

// IsEmpty return true if it is empty
func (t Resources) IsEmpty() bool {
	for _, elem := range t {
		if elem != nil {
			if !elem.IsEmpty() {
				return false
			}
		}
	}
	return true
}

// Resource is identified by its relative URI, which MUST begin with a slash
// ("/"). Every node whose key begins with a slash, and is either at the root
// of the API definition or is the child node of a resource node, is such
// a resource node.
type Resource struct {
	// An alternate, human-friendly name for the resource. If the displayName
	// node is not defined for a resource, documentation tools SHOULD refer to
	// the resource by its key, which acts as the resource name. For example,
	// tools should refer to the relative URI /jobs.
	DisplayName string `yaml:"displayName" json:"displayName,omitempty"`

	// A substantial, human-friendly description of a resource. Its value is a
	// string and MAY be formatted using markdown.
	Description string `yaml:"description" json:"description,omitempty"`

	// Annotations to be applied to this API. An annotation is a map having
	// a key that begins with "(" and ends with ")" where the text enclosed in
	// parentheses is the annotation name, and the value is an instance of that
	// annotation.
	Annotations Annotations `yaml:",regexp:\\(.*\\)" json:"annotations,omitempty"`

	// The object describing the method.
	Get     Method `yaml:"get" json:"get,omitempty"`
	Patch   Method `yaml:"patch" json:"patch,omitempty"`
	Put     Method `yaml:"put" json:"put,omitempty"`
	Post    Method `yaml:"post" json:"post,omitempty"`
	Delete  Method `yaml:"delete" json:"delete,omitempty"`
	Options Method `yaml:"options" json:"options,omitempty"`
	Head    Method `yaml:"head" json:"head,omitempty"`

	// A list of traits to apply to all methods declared (implicitly or
	// explicitly) for this resource. Individual methods can override this
	// declaration.
	Is []*Trait `yaml:"is" json:"is,omitempty"`

	// The resource type that this resource inherits.
	Type Unimplement `yaml:"type" json:"type,omitempty"`

	// The security schemes that apply to all methods declared (implicitly or
	// explicitly) for this resource.
	SecuredBy Unimplement `yaml:"securedBy" json:"securedBy,omitempty"`

	// Detailed information about any URI parameters of this resource.
	URIParameters Unimplement `yaml:"uriParameters" json:"uriParameters,omitempty"`

	// A nested resource, which is identified as any node whose name begins
	// with a slash ("/"), and is therefore treated as a relative URI.
	Resources Resources `yaml:",regexp:/.*" json:"resources,omitempty"`
}

// PostProcess for fill some field from RootDocument default config
func (t *Resource) PostProcess(conf PostProcessConfig) (err error) {
	if t == nil {
		return
	}
	if err = t.Annotations.PostProcess(conf); err != nil {
		return
	}
	if err = t.Get.PostProcess(conf); err != nil {
		return
	}
	if err = t.Patch.PostProcess(conf); err != nil {
		return
	}
	if err = t.Put.PostProcess(conf); err != nil {
		return
	}
	if err = t.Post.PostProcess(conf); err != nil {
		return
	}
	if err = t.Delete.PostProcess(conf); err != nil {
		return
	}
	if err = t.Options.PostProcess(conf); err != nil {
		return
	}
	if err = t.Head.PostProcess(conf); err != nil {
		return
	}
	if err = t.Type.PostProcess(conf); err != nil {
		return
	}
	if err = t.SecuredBy.PostProcess(conf); err != nil {
		return
	}
	if err = t.URIParameters.PostProcess(conf); err != nil {
		return
	}
	if err = t.Resources.PostProcess(conf); err != nil {
		return
	}
	return
}

// IsEmpty return true if it is empty
func (t Resource) IsEmpty() bool {
	return t.DisplayName == "" &&
		t.Description == "" &&
		t.Annotations.IsEmpty() &&
		t.Get.IsEmpty() &&
		t.Patch.IsEmpty() &&
		t.Put.IsEmpty() &&
		t.Post.IsEmpty() &&
		t.Delete.IsEmpty() &&
		t.Options.IsEmpty() &&
		t.Head.IsEmpty() &&
		len(t.Is) < 1 &&
		t.Type.IsEmpty() &&
		t.SecuredBy.IsEmpty() &&
		t.URIParameters.IsEmpty() &&
		t.Resources.IsEmpty()
}
