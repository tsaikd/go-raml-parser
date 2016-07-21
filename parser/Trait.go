package parser

// Traits map of Trait
type Traits map[string]*Trait

// IsEmpty return true if it is empty
func (t Traits) IsEmpty() bool {
	for _, elem := range t {
		if elem != nil {
			if !elem.IsEmpty() {
				return false
			}
		}
	}
	return true
}

// IsTraits slice of Trait
type IsTraits []*Trait

// IsEmpty return true if it is empty
func (t IsTraits) IsEmpty() bool {
	for _, elem := range t {
		if elem != nil {
			if !elem.IsEmpty() {
				return false
			}
		}
	}
	return true
}

// Trait wrap TraitRAML because TraitRAML may be a string for using trait
type Trait struct {
	String string `json:",omitempty"`

	TraitRAML
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
	if err = unmarshaler(&t.TraitRAML); err != nil {
		return
	}
	return
}

// IsEmpty return true if it is empty
func (t Trait) IsEmpty() bool {
	return t.String == "" &&
		t.TraitRAML.IsEmpty()
}

var _ fillTrait = &Trait{}

func (t *Trait) fillTrait(library Library) (err error) {
	if t == nil {
		return
	}

	name := t.String
	if name == "" {
		return
	}

	trait, err := library.GetTrait(name)
	if err != nil {
		return
	}

	*t = trait
	t.String = name

	return
}

var _ checkAnnotation = Trait{}

func (t Trait) checkAnnotation(conf PostProcessConfig) (err error) {
	return t.Annotations.checkAnnotationTargetLocation(TargetLocationTrait)
}

var _ checkUnusedTrait = Trait{}

func (t Trait) checkUnusedTrait(conf PostProcessConfig) (err error) {
	if t.String == "" {
		return
	}
	traitUsage := conf.TraitUsage()
	delete(traitUsage, conf.Library().Prefix()+t.String)
	return
}

// TraitRAML like a method, can provide method-level nodes such as description,
// headers, query string parameters, and responses. Methods that use one or
// more traits inherit nodes of those traits. A resource and resource type
// can also use, and thus inherit from, one or more traits, which then apply
// to all methods of the resource and resource type. Traits are related to
// methods through a mixing pattern.
type TraitRAML struct {
	Method

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

// IsEmpty return true if it is empty
func (t TraitRAML) IsEmpty() bool {
	return t.Method.IsEmpty() &&
		t.Usage == "" &&
		t.ResourcePath == "" &&
		t.ResourcePathName == "" &&
		t.MethodName == ""
}
