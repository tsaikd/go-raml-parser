package parser

// Traits map of Trait
type Traits map[string]*Trait

// PostProcess for fill some field from RootDocument default config
func (t *Traits) PostProcess(conf PostProcessConfig) (err error) {
	if t == nil {
		return
	}
	for _, trait := range *t {
		if err = trait.PostProcess(conf); err != nil {
			return
		}
	}
	return
}

// Trait like a method, can provide method-level nodes such as description,
// headers, query string parameters, and responses. Methods that use one or
// more traits inherit nodes of those traits. A resource and resource type
// can also use, and thus inherit from, one or more traits, which then apply
// to all methods of the resource and resource type. Traits are related to
// methods through a mixing pattern.
type Trait struct {
	String string
	Method
	TraitExtra
}

// UnmarshalYAML implement yaml unmarshaler
// a Property which MIGHT be a simple string or a map[string]interface{}
func (t *Trait) UnmarshalYAML(unmarshaler func(interface{}) error) (err error) {
	if err = unmarshaler(&t.String); err == nil {
		return
	}
	if !isErrorYAMLIntoString(err) {
		return
	}

	if err = unmarshaler(&t.Method); err != nil {
		return
	}
	if err = unmarshaler(&t.TraitExtra); err != nil {
		return
	}
	return
}

// PostProcess for fill some field from RootDocument default config
func (t *Trait) PostProcess(conf PostProcessConfig) (err error) {
	if t == nil {
		return
	}
	if err = t.Method.PostProcess(conf); err != nil {
		return
	}
	if err = t.TraitExtra.PostProcess(conf); err != nil {
		return
	}
	return
}

// TraitExtra contain fields no in Method
type TraitExtra struct {
	// The OPTIONAL usage node of a resource type or trait provides
	// instructions about how and when the resource type or trait should
	// be used. Documentation generators MUST describe this node in terms
	// of the characteristics of the resource and method, respectively.
	// However, the resources and methods MUST NOT inherit the usage node.
	// Neither resources nor methods allow a node named usage.
	Usage string `yaml:"usage" json:"usage,omitempty"`

	// The full resource URI relative to the baseUri if there is one.
	ResourcePath string `yaml:"resourcePath" json:"resourcePath,omitempty"`

	// The rightmost of the non-URI-parameter-containing path fragments.
	ResourcePathName string `yaml:"resourcePathName" json:"resourcePathName,omitempty"`

	// The name of the method
	MethodName string `yaml:"methodName" json:"methodName,omitempty"`
}

// PostProcess for fill some field from RootDocument default config
func (t *TraitExtra) PostProcess(conf PostProcessConfig) (err error) {
	return
}
